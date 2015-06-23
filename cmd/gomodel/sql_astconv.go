package main

import (
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/sqlparser"
)

type Section struct {
	Columns map[string]map[string]*sqlparser.ColName
	Tables  map[string]struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}
}

func newSection() *Section {
	return &Section{
		Columns: make(map[string]map[string]*sqlparser.ColName),
		Tables: make(map[string]struct {
			Table    *sqlparser.TableName
			Subquery *Section
		}),
	}
}

func (s *Section) tableColumns(table string) map[string]*sqlparser.ColName {
	cols, has := s.Columns[table]
	if !has {
		cols = make(map[string]*sqlparser.ColName)
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
		cols[colName] = col
	}
}

func (s *Section) AddTable(as string, tab *sqlparser.TableName) string {
	if as == "" {
		as = string(tab.Name)
	}
	s.Tables[as] = struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}{Table: tab}

	return as
}

func (s *Section) AddSubquery(as string) *Section {
	newsec := newSection()
	s.Tables[as] = struct {
		Table    *sqlparser.TableName
		Subquery *Section
	}{Subquery: newsec}

	return newsec
}

func (mv modelVisitor) astConv(sql string) (string, error) {
	node, err := sqlparser.Parse(sql, true)
	if err != nil {
		return "", err
	}

	s := newSection()
	s.Inspect(node)

	if err = s.replace(mv); err != nil {
		return "", err
	}
	buf := sqlparser.NewTrackedBuffer(nil)
	node.Format(buf)
	return buf.String(), nil
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

func (s *Section) model(mv modelVisitor, tab *sqlparser.TableName) (*Table, error) {
	tabname := string(tab.Name)
	model := mv[tabname]
	if model == nil {
		return nil, errors.Newf("model %s hasn't been registered", tabname)
	}

	return model, nil
}

func (s *Section) replace(mv modelVisitor) error {
	for tabalias, cols := range s.Columns {
		tabnode, has := s.Tables[tabalias]
		if !has {
			return errors.Newf("table alias %s not found in sql", tabalias)
		}

		if tabnode.Table == nil {
			// don't replace columns in subquery section for alias
			continue
		}

		model, err := s.model(mv, tabnode.Table)
		if err != nil {
			return err
		}

		for field, col := range cols {
			col.Name = []byte(model.Fields.DefGet(field, "").(string))
		}
	}

	for _, tabnode := range s.Tables {
		if tabnode.Table != nil {
			if string(tabnode.Table.Name) == "DUAL" {
				continue
			}

			model, err := s.model(mv, tabnode.Table)
			if err != nil {
				return err
			}

			tabnode.Table.Name = []byte(model.Name)
		} else {
			if err := tabnode.Subquery.replace(mv); err != nil {
				return err
			}
		}
	}

	return nil
}
