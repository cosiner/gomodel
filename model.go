package gomodel

import "github.com/cosiner/gohper/bitset"

type (
	// Model represent a database model mapping to a table
	Model interface {
		// Table return database table name that the model mapped
		Table() string
		// Vals store values of fields to given slice, the slice has length,
		// just setup value by index. The value order MUST same as fields defined
		// in strcuture
		Vals(fields uint, vals []interface{})
		// Ptrs is similar to Vals, but for field pointers
		Ptrs(fields uint, ptrs []interface{})
	}
)

var (
	// NumFIelds return fields count
	NumFields = bitset.BitCountUint
)

// FieldVals will extract values of fields from model, and concat with remains
// arguments
func FieldVals(v Model, fields uint, args ...interface{}) []interface{} {
	c, l := NumFields(fields), len(args)
	vals := make([]interface{}, c+l)
	v.Vals(fields, vals)

	for l = l - 1; l >= 0; l-- {
		vals[c+l] = args[l]
	}

	return vals
}

// FieldPtrs is similar to FieldVals, but for field pointers
func FieldPtrs(v Model, fields uint, args ...interface{}) []interface{} {
	c, l := NumFields(fields), len(args)
	ptrs := make([]interface{}, c+l)
	v.Ptrs(fields, ptrs)

	for l = l - 1; l >= 0; l-- {
		ptrs[c+l] = args[l]
	}

	return ptrs
}
