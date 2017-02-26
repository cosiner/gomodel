package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
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

func packageName(outfile string) string {
	abs, err := filepath.Abs(outfile)
	if err != nil {
		return ""
	}
	return filepath.Base(filepath.Dir(abs))
}

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
				color.Red.Errorf("%s: %s\n", name, err)
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
		file.Open(outfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, func(fd *os.File) error {
			t, err := template.ParseFiles(tmplfile)
			if err != nil {
				return err
			}
			err = t.Execute(fd, result)
			if err != nil {
				return err
			}

			pkg := packageName(outfile)
			if pkg == "" {
				return nil
			}
			fd.Seek(0, os.SEEK_SET)

			content, err := ioutil.ReadAll(fd)
			if err != nil {
				return err
			}
			fset := token.NewFileSet()
			ast, err := parser.ParseFile(fset, outfile, bytes.NewReader(content), parser.ParseComments)
			if err != nil {
				return err
			}
			if ast.Name.Name == pkg {
				return nil
			}
			ast.Name.Name = pkg
			fd.Truncate(0)
			return format.Node(fd, fset, ast)
		}),
	)
}
