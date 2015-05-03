package gomodel

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/testing2"
)

func TestCols(t *testing.T) {
	tt := testing2.Wrap(t)
	var cols Cols = &cols{cols: []string{"id", "age", "name"}}
	tt.Eq(3, cols.Length())
	tt.Eq("?,?,?", strings2.RemoveSpace(cols.OnlyParam()))
	tt.Eq("id,age,name", strings2.RemoveSpace(cols.String()))
	tt.Eq("id=?,age=?,name=?", strings2.RemoveSpace(cols.Paramed()))

	cols = singleCol("id")
	tt.Eq(1, cols.Length())
	tt.Eq("?", strings2.RemoveSpace(cols.OnlyParam()))
	tt.Eq("id", strings2.RemoveSpace(cols.String()))
	tt.Eq("id=?", strings2.RemoveSpace(cols.Paramed()))

	cols = zeroCols
	tt.Eq(0, cols.Length())
	tt.Eq("", strings2.RemoveSpace(cols.OnlyParam()))
	tt.Eq("", strings2.RemoveSpace(cols.String()))
	tt.Eq("", strings2.RemoveSpace(cols.Paramed()))
}
