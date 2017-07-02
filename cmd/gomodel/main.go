package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/flag"
	"github.com/cosiner/gomodel/utils"
)

type Flags struct {
	Out      string `names:"-o" default:"model_gen.go" usage:"output file"`
	Model    bool   `names:"-model" default:"true" usage:"implete Model for structure"`
	SQL      bool   `names:"-sql" default:"true" usage:"generate SQL"`
	Template string `names:"-t" default:"" usage:"template file, if not found using default"`
	Args     []string
}

func packageName(outfile string) string {
	abs, err := filepath.Abs(outfile)
	if err != nil {
		return ""
	}
	return filepath.Base(filepath.Dir(abs))
}

func main() {
	var flags Flags
	flag.NewFlagSet(flag.Flag{Arglist: "[FLAG]... FILE|DIR..."}).ParseStruct(&flags)

	files := flags.Args
	if len(files) == 0 {
		utils.FatalOnError(errors.New("no input files to parse."))
		return
	}
	tmpl, err := parseTemplate(flags.Template)
	if err != nil {
		utils.FatalOnError(fmt.Errorf("parse template failed: %s", err.Error()))
	}

	v := newVisitor()
	if len(files) == 1 {
		utils.FatalOnError(v.parseDir(files[0]))
	} else {
		utils.FatalOnError(v.parseFiles(files...))
	}
	if len(v.Models) == 0 {
		utils.FatalOnError(errors.New("no models found."))
	}

	var data templateData
	if flags.SQL {
		for name, sql := range v.SQLs {
			sql, err = v.astConv(sql)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
			} else {
				v.SQLs[name] = sql
			}
		}
		data.SQLs = v.SQLs
	}
	if flags.Model {
		data.Models = v.buildModelFields()
	}

	utils.FatalOnError(executeTemplate(flags.Out, tmpl, data))
}

type templateData struct {
	Models map[*Model][]*Field
	SQLs   map[string]string
}

func executeTemplate(outfile string, tmpl *template.Template, data templateData) error {
	outfd, err := os.OpenFile(outfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outfd.Close()
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return err
	}
	pkg := packageName(outfile)
	if pkg == "" {
		return nil
	}

	fset := token.NewFileSet()
	ast, err := parser.ParseFile(fset, outfile, &buf, parser.ParseComments)
	if err != nil {
		return err
	}
	if ast.Name.Name == pkg {
		return nil
	}
	ast.Name.Name = pkg
	outfd.Truncate(0)
	return format.Node(outfd, fset, ast)
}
