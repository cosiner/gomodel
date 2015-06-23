package main

import (
	"os"
	"path/filepath"

	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/goutil/ast"
	"github.com/cosiner/gohper/sortedmap"
	"github.com/cosiner/gohper/strings2"
)

type Table struct {
	Name   string
	Fields sortedmap.Map
}

type Visitor map[string]*Table

// add an model and it's field to parse result
func (v Visitor) add(model, table, field, col string) {
	if table == "" {
		table = strings2.ToSnake(model)
	}

	if col == "" {
		col = strings2.ToSnake(field)
	}

	t, has := v[model]
	if !has {
		t = &Table{Name: table}
		v[model] = t
	}

	t.Fields.Set(field, col)
}

// parse ast tree to find exported struct and it's fields
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
	return ast.Parser{
		Struct: func(a *ast.Attrs) (err error) {
			if !goutil.IsExported(a.TypeName) {
				err = ast.TYPE_END
			} else if table := a.S.Tag.Get("table"); table == "-" {
				err = ast.TYPE_END
			} else if col := a.S.Tag.Get("column"); col != "-" {
				v.add(a.TypeName, table, a.S.Field, col)
			}

			return
		},
	}.ParseFile(file)
}

// buildModelFields build model map from parse result
func (v Visitor) buildModelFields() map[*Model][]*Field {
	names := make(map[*Model][]*Field, len(v))

	for model, table := range v {
		m := NewModel(model, table.Name)
		fields := table.Fields
		for field, index := range fields.Indexes {
			names[m] = append(names[m], NewField(m, field, fields.Elements[index].(string)))
		}
	}

	return names
}
