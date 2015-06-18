package gomodel

import "sync/atomic"

var (
	currId uint32
)

type IdSql struct {
	ID  uint
	SQL func() string
}

func (i *IdSql) String() string {
	return i.SQL()
}

func NewIdSql(create func() string) *IdSql {
	return &IdSql{
		ID:  uint(atomic.AddUint32(&currId, 1)),
		SQL: create,
	}
}
