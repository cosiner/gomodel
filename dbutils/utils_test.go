package dbutils

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestCondArgs(t *testing.T) {
	tt := testing2.Wrap(t)

	const (
		A uint64 = 1 << iota
		B
		C
		D
	)

	fields, args := CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", "2", "3"})
	tt.True(fields == A|B)
	tt.DeepEq([]interface{}{"1", "2", "3"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", "", "3"})
	tt.True(fields == B)
	tt.DeepEq([]interface{}{"1", "3"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", "", ""})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{0, 1}, []interface{}{"1", "", ""})
	tt.True(fields == A)
	tt.DeepEq([]interface{}{"1", ""}, args)

	fields, args = CondArgs(0, nil, nil, []interface{}{"1", "", ""})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1", "", ""}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", "2", "3"})
	tt.True(fields == A|B)
	tt.DeepEq([]interface{}{"1", "2", "3"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", false, "3"})
	tt.True(fields == B)
	tt.DeepEq([]interface{}{"1", "3"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", false, false})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{0, 1}, []interface{}{"1", false, false})
	tt.True(fields == A)
	tt.DeepEq([]interface{}{"1", false}, args)

	fields, args = CondArgs(0, nil, nil, []interface{}{"1", false, false})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1", false, false}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", 0, "3"})
	tt.True(fields == B)
	tt.DeepEq([]interface{}{"1", "3"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{1, 2}, []interface{}{"1", 0, 0})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1"}, args)

	fields, args = CondArgs(0, []uint64{A, B}, []int{0, 1}, []interface{}{"1", 0, 0})
	tt.True(fields == A)
	tt.DeepEq([]interface{}{"1", 0}, args)

	fields, args = CondArgs(0, nil, nil, []interface{}{"1", 0, 0})
	tt.True(fields == 0)
	tt.DeepEq([]interface{}{"1", 0, 0}, args)

	defer tt.Recover()
	CondArgs(0, []uint64{A, B}, []int{0, 1}, []interface{}{"1", []byte{}, 0})
}
