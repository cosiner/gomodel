package gomodel

import "sync"

var (
	// InitialSQLCount is the initial capacity of sql storage,
	// it should bee changed before NewDB
	InitialSQLCount uint64 = 256
)

type Tabler interface {
	Table(model Model) *Table
}

type sqlStore struct {
	sqls []func(Tabler) string
	sync.Mutex
}

var store sqlStore

func initSqlStore() {
	if store.sqls == nil {
		store.sqls = make([]func(Tabler) string, 0, InitialSQLCount)
	}
}

func sqlById(tabler Tabler, id uint64) string {
	return store.sqls[id](tabler)
}

func NewSqlId(create func(Tabler) string) (id uint64) {
	store.Lock()

	id = uint64(len(store.sqls))
	store.sqls = append(store.sqls, create)

	store.Unlock()
	return
}
