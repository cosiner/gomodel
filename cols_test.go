package gomodel

import (
	"testing"

	"github.com/cosiner/gohper/lib/test"
	"github.com/cosiner/gohper/lib/types"
)

func TestCols(t *testing.T) {
	tt := test.Wrap(t)
	var cols Cols = &cols{cols: []string{"id", "age", "name"}}
	tt.Eq(3, cols.Length())
	tt.Eq("?,?,?", types.RemoveSpace(cols.OnlyParam()))
	tt.Eq("id,age,name", types.RemoveSpace(cols.String()))
	tt.Eq("id=?,age=?,name=?", types.RemoveSpace(cols.Paramed()))

	cols = singleCol("id")
	tt.Eq(1, cols.Length())
	tt.Eq("?", types.RemoveSpace(cols.OnlyParam()))
	tt.Eq("id", types.RemoveSpace(cols.String()))
	tt.Eq("id=?", types.RemoveSpace(cols.Paramed()))

	cols = zeroCols
	tt.Eq(0, cols.Length())
	tt.Eq("", types.RemoveSpace(cols.OnlyParam()))
	tt.Eq("", types.RemoveSpace(cols.String()))
	tt.Eq("", types.RemoveSpace(cols.Paramed()))
}
