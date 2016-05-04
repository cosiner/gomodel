package gomodel

type SQLType uint64

const (
	MAX_NUMFIELDS = 30
)

const (
	// These are several predefined sql types
	INSERT SQLType = iota << (MAX_NUMFIELDS * 2)
	DELETE
	UPDATE
	INCRBY
	LIMIT
	ONE
	ALL
	COUNT
	EXISTS
)
