package dbutils

import (
	"testing"

	"github.com/cosiner/gohper/testing2"
)

func TestCondArgs(t *testing.T) {
	tt := testing2.Wrap(t)

	fields, args := CondArgs(0xf00, 0x2, 1, "1", "2", "3")
	tt.True(fields == (0xf00 | 0x2))
	tt.DeepEq([]interface{}{"1", "2", "3"}, args)

	fields, args = CondArgs(0xf00, 0x2, 1, "1", "", "3")
	tt.True(fields == (0xf00))
	tt.DeepEq([]interface{}{"1", "3"}, args)

	fields, args = CondArgs(0xf00, 0x2, 1, 1, 2, 3)
	tt.True(fields == (0xf00 | 0x2))
	tt.DeepEq([]interface{}{1, 2, 3}, args)

	fields, args = CondArgs(0xf00, 0x2, 1, 1, 0, 3)
	tt.True(fields == (0xf00))
	tt.DeepEq([]interface{}{1, 3}, args)

	fields, args = CondArgs(0xf00, 0x2, 1, true, true, true)
	tt.True(fields == (0xf00 | 0x2))
	tt.DeepEq([]interface{}{true, true, true}, args)

	fields, args = CondArgs(0xf00, 0x2, 1, true, false, true)
	tt.True(fields == (0xf00))
	tt.DeepEq([]interface{}{true, true}, args)

	defer tt.Recover()
	CondArgs(0xf00, 0x2, 1, true, []byte{}, true)
}
