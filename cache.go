package gomodel

import "database/sql"

type (
	// cacheItem keeps the sql and prepared statement of it
	cacheItem struct {
		sql  string
		stmt *sql.Stmt
	}

	cache map[uint64]cacheItem // map[id]{sql, stmt}
)

func newCache() cache {
	return make(cache)
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c cache) StmtById(exec Executor, sqlid uint64) (*sql.Stmt, error) {
	if item, has := c[sqlid]; has {
		sqlPrinter.Print(true, item.sql)

		return item.stmt, nil
	}

	sql_ := sqlById(exec, sqlid)
	sql_ = exec.Driver().Prepare(sqlBufpool, sql_)
	sqlPrinter.Print(false, sql_)

	stmt, err := exec.Prepare(sql_)
	if err != nil {
		return nil, err
	}
	c[sqlid] = cacheItem{sql: sql_, stmt: stmt}

	return stmt, nil
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c cache) GetStmt(exec Executor, sqlid uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		return "", nil, nil
	}
	var err error
	if item.stmt == nil {
		item.sql = exec.Driver().Prepare(sqlBufpool, item.sql)
		item.stmt, err = exec.Prepare(item.sql)
	}

	return item.sql, item.stmt, err
}

// SetStmt exec a sql to statement, cache then return it
func (c cache) SetStmt(exec Executor, sqlid uint64, sql string) (*sql.Stmt, error) {
	sql = exec.Driver().Prepare(sqlBufpool, sql)
	stmt, err := exec.Prepare(sql)
	if err != nil {
		return nil, err
	}
	c[sqlid] = cacheItem{
		sql:  sql,
		stmt: stmt,
	}

	return stmt, nil
}

func (c cache) PrepareById(exec Executor, sqlid uint64) (*sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		item.sql = exec.Driver().Prepare(sqlBufpool, sqlById(exec, sqlid))
		c[sqlid] = item
	}
	sqlPrinter.Print(has, item.sql)

	stmt, err := exec.Prepare(item.sql)
	return stmt, err
}

func (c cache) PrepareSQL(exec Executor, sqlid uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		return "", nil, nil
	}

	item.sql = exec.Driver().Prepare(sqlBufpool, item.sql)
	stmt, err := exec.Prepare(item.sql)
	return item.sql, stmt, err
}

func (c cache) SetSQL(sqlid uint64, sql string) {
	c[sqlid] = cacheItem{
		sql: sql,
	}
}
