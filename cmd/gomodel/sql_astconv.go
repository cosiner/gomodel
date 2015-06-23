package main

import (
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/sqlparser"
)

type Section struct {
	Columns map[string]map[string][]*sqlparser.ColName
	Tables  map[string][]struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}
}

func newSection() *Section {
	return &Section{
		Columns: make(map[string]map[string][]*sqlparser.ColName),
		Tables: make(map[string][]struct {
			Table    *sqlparser.TableName
			Subquery *Section
		}),
	}
}

func (s *Section) tableColumns(table string) map[string][]*sqlparser.ColName {
	cols, has := s.Columns[table]
	if !has {
		cols = make(map[string][]*sqlparser.ColName)
		s.Columns[table] = cols
	}

	return cols
}

func (s *Section) AddColumn(table string, col *sqlparser.ColName) {
	tabName := string(col.Qualifier)
	if tabName == "" {
		tabName = table
	}

	cols := s.tableColumns(tabName)
	if colName := string(col.Name); colName != "?" {
		cols[colName] = append(cols[colName], col)
	}
}

func (s *Section) AddTable(as string, tab *sqlparser.TableName) string {
	if as == "" {
		as = string(tab.Name)
	}
	s.Tables[as] = append(s.Tables[as], struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}{
		Table: tab,
	})

	return as
}

func (s *Section) AddSubquery(as string) *Section {
	newsec := newSection()
	s.Tables[as] = append(s.Tables[as], struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}{
		Subquery: newsec,
	})

	return newsec
}

func (s *Section) Inspect(node sqlparser.SQLNode) {
	var table string
	sqlparser.Inspect(node, func(node sqlparser.SQLNode) bool {
		switch node := node.(type) {
		case *sqlparser.TableName:
			table = s.AddTable("", node)
		case *sqlparser.AliasedTableExpr:
			switch tab := node.Expr.(type) {
			case *sqlparser.TableName:
				table = s.AddTable(string(node.As), tab)
			case *sqlparser.Subquery:
				table = string(node.As)
				s.AddSubquery(table).Inspect(tab)
			}

			return false // stop inspect subnodes
		case *sqlparser.ColName:
			s.AddColumn(table, node)
		}

		return true
	})
}

func (s *Section) modelTable2(v Visitor, tab *sqlparser.TableName) (*Table, error) {
	tabname := string(tab.Name)
	model := v.Models[tabname]
	if model == nil {
		return nil, errors.Newf("model %s hasn't been registered", tabname)
	}

	return model, nil
}

func (s *Section) replace(v Visitor) error {
	for tabalias, cols := range s.Columns {
		tabnodes, has := s.Tables[tabalias]
		if !has {
			return errors.Newf("table alias %s not found in sql", tabalias)
		}

		tabnode := tabnodes[0]
		if tabnode.Table == nil {
			// don't replace columns of subquery section because of alias
			continue
		}

		model, err := s.modelTable2(v, tabnode.Table)
		if err != nil {
			return err
		}

		for field, cols2 := range cols {
			for _, col := range cols2 {
				col.Name = []byte(model.Fields.DefGet(field, "").(string))
			}
		}
	}

	for _, tabnodes := range s.Tables {
		for _, tabnode := range tabnodes {
			if tabnode.Table != nil {
				if string(tabnode.Table.Name) == "DUAL" {
					continue
				}

				model, err := s.modelTable2(v, tabnode.Table)
				if err != nil {
					return err
				}

				tabnode.Table.Name = []byte(model.Name)
			} else {
				if err := tabnode.Subquery.replace(v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (v Visitor) astConv(sql string) (string, error) {
	node, err := sqlparser.Parse(sql, true)
	if err != nil {
		return "", err
	}

	s := newSection()
	s.Inspect(node)

	if err = s.replace(v); err != nil {
		return "", err
	}
	buf := sqlparser.NewTrackedBuffer(nil)
	node.Format(buf)
	return buf.String(), nil
}
