package gomodel

import "sync/atomic"

var (
	currId uint32
)

type IdSql struct {
	ID  uint
	SQL string
}

func (i IdSql) String() string {
	return i.SQL
}

func NewIdSql(create interface{}) IdSql {
	var sql string
	switch c := create.(type) {
	case string:
		sql = c
	case func() string:
		sql = c()
	default:
		panic("invalid sql type")
	}

	return IdSql{
		ID:  uint(atomic.AddUint32(&currId, 1)),
		SQL: sql,
	}
}
