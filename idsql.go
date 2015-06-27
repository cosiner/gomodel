package gomodel

import "sync/atomic"

var (
	currId uint64
)

type IdSql struct {
	ID  uint64
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
		ID:  atomic.AddUint64(&currId, 1),
		SQL: sql,
	}
}
