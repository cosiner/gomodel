package gomodel

import (
	"github.com/cosiner/gohper/ds/bitset"
)

type (
	// Model represent a database model mapping to a table
	Model interface {
		// Table return database table name that the model mapped
		Table() string
		// Vals store values of fields to given slice, the slice has length,
		// just setup value by index. The value order MUST same as fields defined
		// in strcuture
		Vals(fields uint64, vals []interface{})
		// Ptrs is similar to Vals, but for field pointers
		Ptrs(fields uint64, ptrs []interface{})
	}

	// Columner is a optional interface for Model, if Model implements this interface,
	// it's no need to parse Model info with reflection
	Columner interface {
		// Columns return all column names for this Model
		Columns() []string
	}

	// Nocacher is a optional interface for Model, if Model implements this interface,
	// and NoCache method return true, it will not allocate memory to store
	// sql, stmt, columns for this Model, all sqls, stmts must be stored in DB instance.
	//
	// If Nocache, the only methods to get Stmt are DB.StmtById and Tx.PrepareById,
	// implements it only when you actually know what are you do.
	Nocacher interface {
		Nocache() bool
	}
)

var (
	// NumFIelds return fields count
	NumFields = bitset.BitCount
)

func Fields(fields ...uint64) uint64 {
	var f uint64
	for _, field := range fields {
		f |= field
	}
	return f
}

// NumFieldsExcp create fieldset except given fields
func NumFieldsExcp(numField uint64, fields ...uint64) uint64 {
	return AllFieldsExcp(1<<numField-1, fields...)
}

// AllFieldsExcp create fieldset except given fields
func AllFieldsExcp(allFields uint64, fields ...uint64) uint64 {
	return allFields & ^Fields(fields...)
}

func AllFieldsOrSome(allFields uint64, fields ...uint64) uint64 {
	f := Fields(fields...)
	if f != 0 {
		return f
	}
	return allFields
}

// FieldVals will extract values of fields from model, and concat with remains
// arguments
func FieldVals(model Model, fields uint64, args ...interface{}) []interface{} {
	c, l := NumFields(fields), len(args)
	vals := make([]interface{}, c+l)
	model.Vals(fields, vals)

	for l = l - 1; l >= 0; l-- {
		vals[c+l] = args[l]
	}

	return vals
}

// FieldPtrs is similar to FieldVals, but for field pointers.
// FieldPtrs only used for query operations such as One, Limit, All.
func FieldPtrs(model Model, fields uint64, args ...interface{}) []interface{} {
	c, l := NumFields(fields), len(args)
	ptrs := make([]interface{}, c+l)
	model.Ptrs(fields, ptrs)

	for l = l - 1; l >= 0; l-- {
		ptrs[c+l] = args[l]
	}

	return ptrs
}
