package main

import (
	"strings"

	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/strings2"
)

type Model struct {
	Name       string // struct's normal name
	Self       string
	Unexported string
	Table      string
}

func NewModel(name, table string) *Model {
	return &Model{
		Name:       name,
		Self:       strings2.ToLowerAbridge(name),
		Unexported: goutil.ToUnexported(name),
		Upper:      strings.ToUpper(name),
		Table:      table,
	}
}

type Field struct {
	Name  string // field's normal name
	Const string // field's const name is in STRUCTNAME_FIELDNAME case
}

func NewField(model *Model, field string) *Field {
	f := &Field{
		Name: field,
	}

	if useCamelCase {
		f.Const = model.Name + field
	} else {
		f.Const = model.Upper + "_" + strings.ToUpper(field)
	}

	return f
}
