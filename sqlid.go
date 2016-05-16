package gomodel

import "sync"

var (
	// InitialSQLCount is the initial capacity of sql storage,
	// it should bee changed before NewDB
	InitialSQLCount uint64 = 256
)

const (
	_InitialSQLBufsize = 256
)

type sqlStore struct {
	sqls []func(Executor) string
	sync.Mutex
}

var store sqlStore

func initSqlStore() {
	if store.sqls == nil {
		store.sqls = make([]func(Executor) string, 0, InitialSQLCount)
	}
}

func sqlById(executor Executor, id uint64) string {
	return store.sqls[id](executor)
}

// NewSqlId create an id for this sql creator used in methods like XXXById
func NewSqlId(create func(Executor) string) (id uint64) {
	store.Lock()

	id = uint64(len(store.sqls))
	store.sqls = append(store.sqls, create)

	store.Unlock()
	return
}

type SqlIdKeeper struct {
	ids map[string]string
	mu  sync.RWMutex
}

func NewSqlIdKeeper() *SqlIdKeeper {
	return &SqlIdKeeper{
		ids: make(map[string]string),
	}
}

func (s *SqlIdKeeper) Get(key string) (string, bool) {
	s.mu.RLock()
	val, has := s.ids[key]
	s.mu.RUnlock()
	return val, has
}

func (s *SqlIdKeeper) Set(key, val string) {
	s.mu.Lock()
	s.ids[key] = val
	s.mu.Unlock()
}
