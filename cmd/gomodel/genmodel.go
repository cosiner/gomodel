package main

import (
	"flag"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/gohper/defval"
	encfile "github.com/cosiner/gohper/encoding2/file"
	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/terminal/color"
)

var (
	outfile    string
	tmplfile   string
	copyTmpl   bool
	sqlmapping string

	useAst bool
)

func init() {
	flag.StringVar(&outfile, "o", "", "outtput file, default model_gen.go")
	flag.StringVar(&tmplfile, "t", "", "template file, first find in current directory, else use default file")
	flag.BoolVar(&copyTmpl, "cp", false, "copy tmpl file to default path")
	flag.StringVar(&sqlmapping, "m", "", "sql mapping file to convert")
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

	mv := make(modelVisitor)
	if len(args) == 1 {
		errors.Fatalln(mv.parseDir(args[0]))
	} else {
		errors.Fatalln(mv.parseFiles(args...))
	}

	if len(mv) == 0 {
		return
	}

	type Result struct {
		Models map[*Model][]*Field
		SQLs   map[string]string
	}

	var result Result

	if sqlmapping != "" {
		mapping, err := encfile.ReadProperties(sqlmapping)
		errors.Fatalln(err)

		errors.Fatal(
			file.OpenOrCreate(outfile, false, func(fd *os.File) error {
				for name, sql := range mapping {
					sql, err := strings2.TrimQuote(sql)
					if err != nil {
						color.LightRed.Errorln(err)
						continue
					}
					if useAst {
						sql, err = mv.astConv(sql)
					} else {
						sql, err = mv.conv(sql)
					}
					if err != nil {
						color.LightRed.Errorf("%s: %s", name, err)
					} else {
						mapping[name] = sql
					}
				}

				return nil
			}),
		)

		result.SQLs = mapping
	}

	errors.Fatal(
		file.OpenOrCreate(outfile, false, func(fd *os.File) error {
			t, err := template.ParseFiles(tmplfile)
			if err != nil {
				return err
			}
			result.Models = mv.buildModelFields()

			return t.Execute(fd, result)
		}),
	)
}
