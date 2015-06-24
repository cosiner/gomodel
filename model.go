package gomodel

import "github.com/cosiner/gohper/bitset"

type (
	// Model represent a database model
	Model interface {
		Table() string
		// Vals store values of fields to given slice
		Vals(fields uint, vals []interface{})
		Ptrs(fields uint, ptrs []interface{})
	}
)

var (
	FieldCount = bitset.BitCountUint
)

func FieldVals(v Model, fields uint, args ...interface{}) []interface{} {
	c, l := FieldCount(fields), len(args)
	vals := make([]interface{}, c+l)
	v.Vals(fields, vals)

	for l = l - 1; l >= 0; l-- {
		vals[c+l] = args[l]
	}

	return vals
}

func FieldPtrs(v Model, fields uint, args ...interface{}) []interface{} {
	c, l := FieldCount(fields), len(args)
	ptrs := make([]interface{}, c+l)
	v.Ptrs(fields, ptrs)

	for l = l - 1; l >= 0; l-- {
		ptrs[c+l] = args[l]
	}

	return ptrs
}
