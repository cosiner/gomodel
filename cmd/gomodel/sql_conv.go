package main

import (
	"bytes"

	"github.com/cosiner/gohper/bytes2"
	"github.com/cosiner/gohper/errors"
)

const (
	INIT = iota
	PARSING_MODEL
	PARSING_FIELD
)

func (v Visitor) modelTable(modelbuf *bytes.Buffer, table **Table) error {
	model := modelbuf.String()
	*table = v[model]
	if *table == nil {
		return errors.Newf("model %s isn't registered", model)
	}

	return nil
}

func (v Visitor) writeModel(sqlbuf, modelbuf *bytes.Buffer) error {
	var table *Table
	err := v.modelTable(modelbuf, &table)
	if err == nil {
		sqlbuf.WriteString(table.Name)
	}

	return err
}

func (v Visitor) writeField(table *Table, withModel bool, sqlbuf, modelbuf, fieldbuf *bytes.Buffer) error {
	field := fieldbuf.String()
	col := table.Fields.Get(field)
	if col == nil {
		return errors.Newf("field %s of model %s not found", field, modelbuf.String())
	}

	if withModel {
		sqlbuf.WriteString(table.Name)
		sqlbuf.WriteByte('.')
	}
	sqlbuf.WriteString(col.(string))

	return nil
}

func (v Visitor) conv(sql string) (s string, err error) {
	state := INIT
	sqlbuf := bytes2.NewBuffer(len(sql))
	modelbuf := bytes2.NewBuffer(8)
	fieldbuf := bytes2.NewBuffer(8)

	var table *Table
	var withModel bool
	for i := range sql {
		c := sql[i]
		switch state {
		case INIT:
			if c == '{' {
				state = PARSING_MODEL
				withModel = false
				modelbuf.Reset()
				fieldbuf.Reset()
			} else {
				sqlbuf.WriteByte(c)
			}

		case PARSING_MODEL:
			switch c {
			case '}':
				if err = v.writeModel(sqlbuf, modelbuf); err != nil {
					return
				}

				state = INIT
			case '.':
				withModel = true
				fallthrough
			case ':':
				if err = v.modelTable(modelbuf, &table); err != nil {
					return
				}

				state = PARSING_FIELD
			default:
				modelbuf.WriteByte(c)
			}
		case PARSING_FIELD:
			if c == ',' || c == '}' {
				if err = v.writeField(table, withModel, sqlbuf, modelbuf, fieldbuf); err != nil {
					return
				}

				fieldbuf.Reset()

				if c == '}' {
					state = INIT
				} else {
					sqlbuf.WriteByte(c)
				}
			} else if c == ' ' {
				sqlbuf.WriteByte(c)
			} else {
				fieldbuf.WriteByte(c)
			}
		}
	}

	return sqlbuf.String(), nil
}
