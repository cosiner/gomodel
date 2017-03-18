package main

import (
	"strings"

	"github.com/cosiner/gomodel/utils"
)

type Model struct {
	Name       string // struct's normal name
	Self       string
	Unexported string
	Upper      string
	Table      string
	Nocache    string
}

func NewModel(name, table, nocache string) *Model {
	return &Model{
		Name:       name,
		Self:       utils.ToLowerAbridgeCase(name),
		Unexported: utils.UnexportedName(name),
		Upper:      strings.ToUpper(name),
		Table:      table,
		Nocache:    nocache,
	}
}

type Field struct {
	Name   string // field's normal name
	Upper  string // field's const name is in STRUCTNAME_FIELDNAME case
	Column string
}

func NewField(field, col string) *Field {
	return &Field{
		Name:   field,
		Upper:  strings.ToUpper(field),
		Column: col,
	}
}
