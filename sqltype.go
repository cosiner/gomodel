package gomodel

type SQLType uint64

const (
	MAX_NUMFIELDS = 30
)

const (
	// These are five predefined sql types
	INSERT SQLType = iota + 1<<(MAX_NUMFIELDS*2)
	DELETE
	UPDATE
	INCRBY
	LIMIT
	ONE
	ALL
)
