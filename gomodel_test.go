package gomodel

import (
	"testing"

	"github.com/cosiner/gohper/strings2"
	"github.com/cosiner/gohper/testing2"
)

func TestCols(t *testing.T) {
	var cols Cols = &cols{cols: []string{"id", "age", "name"}}

	testing2.Eq(t, 3, cols.Length())
	testing2.
		Expect("?,?,?").Arg(cols.OnlyParam()).
		Expect("id,age,name").Arg(cols.String()).
		Expect("id=?,age=?,name=?").Arg(cols.Paramed()).
		Run(t, strings2.RemoveSpace)

	cols = singleCol("id")
	testing2.Eq(t, 1, cols.Length())
	testing2.
		Expect("?").Arg(cols.OnlyParam()).
		Expect("id").Arg(cols.String()).
		Expect("id=?").Arg(cols.Paramed()).
		Run(t, strings2.RemoveSpace)

	cols = _emptyCols
	testing2.Eq(t, 0, cols.Length())
	testing2.
		Expect("").Arg(cols.OnlyParam()).
		Expect("").Arg(cols.String()).
		Expect("").Arg(cols.Paramed()).
		Run(t, strings2.RemoveSpace)
}

func TestInterface(t *testing.T) {
	var _ Executor = &DB{}
	var _ Executor = Tx{}
}
