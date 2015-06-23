package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/gohper/defval"
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
	"github.com/cosiner/gohper/terminal/color"
)

var (
	outfile  string
	tmplfile string
	copyTmpl bool

	useAst bool
)

func init() {
	flag.StringVar(&outfile, "o", "", "outtput file, default model_gen.go")
	flag.StringVar(&tmplfile, "t", "", "template file, first find in current directory, else use default file")
	flag.BoolVar(&copyTmpl, "cp", false, "copy tmpl file to default path")
	flag.BoolVar(&useAst, "ast", true, "parse sql ast")
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

	v := newVisitor()
	if len(args) == 1 {
		errors.Fatalln(v.parseDir(args[0]))
	} else {
		errors.Fatalln(v.parseFiles(args...))
	}

	if len(v.Models) == 0 {
		return
	}

	errors.Fatal(
		file.OpenOrCreate(outfile, false, func(fd *os.File) error {
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

			return nil
		}),
	)

	errors.Fatal(
		file.OpenOrCreate(outfile, false, func(fd *os.File) error {
			t, err := template.ParseFiles(tmplfile)
			if err != nil {
				return err
			}

			return t.Execute(fd, struct {
				Models map[*Model][]*Field
				SQLs   map[string]string
			}{
				v.buildModelFields(),
				v.SQLs,
			})
		}),
	)
}
