package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/gohper/defval"
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/goutil/ast"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
	"github.com/cosiner/gohper/sortedmap"
	"github.com/cosiner/gohper/strings2"
)

var (
	outfile      string
	tmplfile     string
	copyTmpl     bool
	useCamelCase bool
)

func init() {
	flag.StringVar(&outfile, "o", "", "outtput file, default model_gen.go")
	flag.StringVar(&tmplfile, "t", "", "template file, first find in current directory, else use default file")

	// make it true to enable default CamelCase
	flag.BoolVar(&useCamelCase, "cc", false, "use CamelCase of constants")
	flag.BoolVar(&copyTmpl, "cp", false, "copy tmpl file to default path")

	flag.Parse()
	defval.String(&outfile, "model_gen.go")
	if tmplfile == "" {
		tmplfile = TmplName
		if !file.IsExist(tmplfile) {
			tmplfile = defTmplPath
		}
	}
}

const TmplName = "model.tmpl"

// change this if need
var defTmplPath = filepath.Join(path2.Home(), ".config", "go", TmplName)

func main() {
	if copyTmpl {
		errors.Fatal(file.Copy(defTmplPath, TmplName))
		return
	}
	args := flag.Args()
	if len(args) == 0 {
		return
	}

	mv := make(modelVisitor)
	if len(args) == 1 {
		errors.Fatalln(mv.parseDir(args[0]))
	} else {
		errors.Fatalln(mv.parseFiles(args...))
	}

	if len(mv) == 0 {
		return
	}

	errors.Fatal(
		file.OpenOrCreate(outfile, false, func(fd *os.File) error {
			t, err := template.ParseFiles(tmplfile)
			if err != nil {
				return err
			}

			return t.Execute(fd, mv.buildModelFields())
		}),
	)
}

type Table struct {
	Name   string
	Fields sortedmap.Map
}

type modelVisitor map[string]*Table

// add an model and it's field to parse result
func (mv modelVisitor) add(model, table, field, col string) {
	if table == "" {
		table = strings2.ToSnake(model)
	}

	if col == "" {
		col = strings2.ToSnake(field)
	}

	t, has := mv[model]
	if !has {
		t = &Table{Name: table}
		mv[model] = t
	}

	t.Fields.Set(field, col)
}

// needParse check whether a model should be parsed
// unexporeted model don't parse
// if visitor's model list is not empty, only parse model exist in list
// otherwise parse all
func (mv modelVisitor) needParse(model string) bool {
	return goutil.IsExported(model)
}

// parse ast tree to find exported struct and it's fields
func (mv modelVisitor) parseFiles(files ...string) error {
	for _, file := range files {
		err := mv.parseFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mv modelVisitor) parseDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		return mv.parseFile(path)
	})
}

func (mv modelVisitor) parseFile(file string) error {
	return ast.Parser{
		Struct: func(a *ast.Attrs) (err error) {
			if !mv.needParse(a.TypeName) {
				err = ast.TYPE_END
			} else if table := a.S.Tag.Get("table"); table == "-" {
				err = ast.TYPE_END
			} else if col := a.S.Tag.Get("column"); col != "-" {
				mv.add(a.TypeName, table, a.S.Field, col)
			}

			return
		},
	}.ParseFile(file)
}

// buildModelFields build model map from parse result
func (mv modelVisitor) buildModelFields() map[*Model][]*Field {
	names := make(map[*Model][]*Field, len(mv))

	for model, table := range mv {
		m := NewModel(model, table.Name)

		for field := range table.Fields.Indexes {
			names[m] = append(names[m], NewField(m, field))
		}
	}

	return names
}
