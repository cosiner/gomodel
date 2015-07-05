package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/os2/file"
	"github.com/cosiner/gohper/os2/path2"
	"github.com/cosiner/gohper/terminal/color"
	"github.com/cosiner/gohper/unsafe2"
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
	flag.StringVar(&outfile, "o", "model_gen.go", "output file")
	flag.StringVar(&tmplfile, "t", Tmplfile, "template file, search current directory firstly, use default file $HOME/.config/go/"+Tmplfile+" if not found")
	flag.BoolVar(&copyTmpl, "cp", false, "copy tmpl file to default path")
	flag.BoolVar(&useAst, "ast", true, "parse sql ast")

	flag.BoolVar(&parseModel, "model", false, "generate model functions")
	flag.BoolVar(&parseSQL, "sql", false, "generate sqls")

	flag.Usage = func() {
		fmt.Println("gomodel [OPTIONS] DIR|FILES...")
		flag.PrintDefaults()
	}
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
		flag.Usage()
		return
	}

	if !parseSQL && !parseModel {
		fmt.Println("neighter -model nor -sql are specified.")
		return
	}

	v := newVisitor()
	if len(args) == 1 {
		errors.Fatalln(v.parseDir(args[0]))
	} else {
		errors.Fatalln(v.parseFiles(args...))
	}

	if len(v.Models) == 0 {
		fmt.Println("no models found.")
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
		file.OpenOrCreate(outfile, false, func(fd *os.File) (err error) {
			t := template.New("tmpl").Funcs(map[string]interface{}{
				"HasPrefix": strings.HasPrefix,
				"Slice": func(s string, start, end int) string {
					if end == -1 {
						end = len(s)
					}
					return s[start:end]
				},
			})

			tfd, err := os.Open(tmplfile)
			if err != nil {
				return err
			}
			defer tfd.Close()

			content, err := ioutil.ReadAll(tfd)
			if err != nil {
				return err
			}

			t, err = t.Parse(unsafe2.String(content))
			if err != nil {
				return err
			}

			return t.Execute(fd, result)
		}),
	)
}
