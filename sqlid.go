package gomodel

import "sync"

type sqlStore struct {
	sqls []func(Executor) string
	sync.RWMutex
}

var store = sqlStore{
	sqls: make([]func(Executor) string, 0, 32),
}

func SqlById(executor Executor, id uint64) string {
	var creator func(Executor) string
	store.RLock()
	if id < uint64(len(store.sqls)) {
		creator = store.sqls[id]
	}
	store.RUnlock()
	if creator == nil {
		return ""
	}
	return creator(executor)
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
