package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
	"github.com/cosiner/gohper/terminal/color"
)

var (
	outfile  string
	tmplfile string
	copyTmpl bool

	parseModel bool
	parseSQL   bool

	useAst bool
)

func init() {
	flag.StringVar(&outfile, "o", "model_gen.go", "output file, default model_gen.go")
	flag.StringVar(&tmplfile, "t", Tmplfile, "template file, search current directory firstly, use default file $HOME/.config/go/"+Tmplfile+" if not found")
	flag.BoolVar(&copyTmpl, "cp", false, "copy tmpl file to default path")
	flag.BoolVar(&useAst, "ast", true, "parse sql ast")

	flag.BoolVar(&parseModel, "model", false, "generate model functions")
	flag.BoolVar(&parseSQL, "sql", false, "generate sqls")
	flag.Parse()

	if !file.IsExist(tmplfile) {
		tmplfile = defTmplPath
	}
}

const Tmplfile = "model.tmpl"

// change this if need
var defTmplPath = filepath.Join(path2.Home(), ".config", "go", Tmplfile)

func main() {
	if copyTmpl {
		errors.Fatal(file.Copy(defTmplPath, Tmplfile))
		return
	}
	args := flag.Args()
	if len(args) == 0 {
		return
	}

	if !parseSQL && !parseModel {
		return
	}

	v := newVisitor()
	if len(args) == 1 {
		errors.Fatalln(v.parseDir(args[0]))
	} else {
		errors.Fatalln(v.parseFiles(args...))
	}

	if len(v.Models) == 0 {
		return
	}

	var result struct {
		Models map[*Model][]*Field
		SQLs   map[string]string
	}

	if parseSQL {
		var err error
		for name, sql := range v.SQLs {
			if useAst {
				sql, err = v.astConv(sql)
			} else {
				sql, err = v.conv(sql)
			}
			if err != nil {
				color.LightRed.Errorf("%s: %s\n", name, err)
			} else {
				v.SQLs[name] = sql
			}
		}

		result.SQLs = v.SQLs
	}

	if parseModel {
		result.Models = v.buildModelFields()
	}

	errors.Fatal(
		file.OpenOrCreate(outfile, false, func(fd *os.File) error {
			t, err := template.ParseFiles(tmplfile)
			if err != nil {
				return err
			}

			return t.Execute(fd, result)
		}),
	)
}
