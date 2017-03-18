package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cosiner/flag"
	"github.com/cosiner/gomodel/utils"
)

type Flags struct {
	Out      string `names:"-o" default:"model_gen.go" usage:"output file"`
	Template string `names:"-t" default:"model.tmpl" usage:"template file, use default file $HOME/.config/go/model.tmpl if not found"`
	Copy     bool   `names:"-cp" usage:"copy tmpl file to default path"`
	Model    bool   `names:"-model" default:"true" usage:"generate model functions"`
	SQL      bool   `names:"-sql" default:"true" usage:"generate sqls"`
	Args     []string

	DefaultPath string `names:"-"`
}

func parseFlags() Flags {
	var flags Flags
	flags.DefaultPath = filepath.Join(utils.UserHome(), ".config", "go", "model.tmpl")
	flag.NewFlagSet(flag.Flag{
		Arglist: "[FLAG]... FILE|DIR...",
	}).ParseStruct(&flags)
	if !utils.IsFileExist(flags.Template) {
		flags.Template = flags.DefaultPath
	}
	return flags
}

func packageName(outfile string) string {
	abs, err := filepath.Abs(outfile)
	if err != nil {
		return ""
	}
	return filepath.Base(filepath.Dir(abs))
}

func main() {
	flags := parseFlags()
	if flags.Copy {
		utils.FatalOnError(utils.CopyFile(flags.DefaultPath, flags.Template))
		return
	}
	files := flags.Args
	if len(files) == 0 {
		utils.FatalOnError(errors.New("no input files to parse."))
		return
	}

	if !flags.Model && !flags.SQL {
		utils.FatalOnError(errors.New("neighter -model nor -sql are specified."))
		return
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
		var err error
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

	fd, err := os.OpenFile(flags.Out, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	utils.FatalOnError(err)
	defer fd.Close()

	utils.FatalOnError(executeTemplate(flags.Out, fd, flags.Template, data))
}

type templateData struct {
	Models map[*Model][]*Field
	SQLs   map[string]string
}

func executeTemplate(outfile string, outfd *os.File, tmpl string, data templateData) error {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return err
	}
	err = t.Execute(outfd, data)
	if err != nil {
		return err
	}

	pkg := packageName(outfile)
	if pkg == "" {
		return nil
	}
	outfd.Seek(0, os.SEEK_SET)

	content, err := ioutil.ReadAll(outfd)
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
	outfd.Truncate(0)
	return format.Node(outfd, fset, ast)
}
