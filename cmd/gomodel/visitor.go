package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/gohper/ds/sortedmap"
	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/utils/ast"
	"github.com/cosiner/gohper/utils/pair"
)

type Table struct {
	Name    string
	Nocache string
	Fields  sortedmap.Map

	initialed bool
}

type Visitor struct {
	Models map[string]*Table // [modelname]modeltable
	SQLs   map[string]string // [sqlid]sqlstring
}

func newVisitor() Visitor {
	return Visitor{
		Models: make(map[string]*Table),
		SQLs:   make(map[string]string),
	}
}

// add an model and it's field to parse result
func (v Visitor) add(model, table, field, col string) {
	if table == "" {
		table = strings2.ToSnake(model)
	}

	if col == "" {
		col = strings2.ToSnake(field)
	}

	t, has := v.Models[model]
	if !has {
		t = &Table{Name: table, Fields: sortedmap.New(), initialed: true}
		v.Models[model] = t
	} else if t.Fields.Indexes == nil {
		t.Name = table
		t.Fields = sortedmap.New()
	}

	t.Fields.Set(field, col)
}

func (v Visitor) addNocahe(model, nocache string) {
	t, has := v.Models[model]
	if !has {
		t = &Table{}
		v.Models[model] = t
	}

	t.Nocache = nocache
}

func (v Visitor) parseFiles(files ...string) error {
	for _, file := range files {
		err := v.parseFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (v Visitor) parseDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return v.parseFile(path)
	})
}

func (v Visitor) parseFile(file string) error {
	parser := ast.Parser{
		Struct: func(a *ast.Attrs) error {
			if !goutil.IsExported(a.TypeName) {
				return ast.TYPE_END
			}

			var table string
			if table = a.S.Tag.Get("table"); table == "-" {
				return ast.TYPE_END
			}

			if nocache := a.S.Tag.Get("nocache"); nocache != "" {
				v.addNocahe(a.TypeName, nocache)
			}

			if col := a.S.Tag.Get("column"); col != "-" {
				v.add(a.TypeName, table, a.S.Field, col)
			}

			return nil
		},

		ParseDoc: true,
		Func: func(a *ast.Attrs) (err error) {
			v.extractSQLs(a.TypeDoc)

			return nil
		},
	}
	return parser.ParseFile(file)
}

// buildModelFields build model map from parse result
func (v Visitor) buildModelFields() map[*Model][]*Field {
	names := make(map[*Model][]*Field, len(v.Models))

	for model, table := range v.Models {
		m := NewModel(model, table.Name, table.Nocache)
		fields := table.Fields

		for _, field := range fields.Values {
			names[m] = append(names[m], NewField(field.Key, field.Value.(string)))
		}
	}

	return names
}

func (v Visitor) extractSQLs(docs []string) {
	const (
		INIT = iota
		PARSING

		GOMODEL = "//gomodel "
	)

	var (
		sqls  []string
		name  string
		state = INIT
	)

	for _, d := range docs {
		if state == PARSING {
			d = d[len("//"):]

			if strings.HasPrefix(d, "]") {
				v.SQLs[name] = strings.Join(sqls, " ")
				state = INIT
			} else {
				sqls = append(sqls, strings.TrimSpace(d))
			}

		} else if strings.HasPrefix(d, GOMODEL) {
			d = d[len(GOMODEL):]
			p := pair.Parse(d, "=").Trim()

			if !strings.HasSuffix(p.Value, "[") {
				v.SQLs[p.Key] = p.Value
			} else {
				name = p.Key
				sqls = sqls[:0]

				state = PARSING
			}
		}
	}
}
