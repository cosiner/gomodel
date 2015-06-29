package gomodel

import "database/sql"

type (
	// cacheItem keeps the sql and prepared statement of it
	cacheItem struct {
		sql  string
		stmt *sql.Stmt
	}

	Preparer interface {
		Tabler
		Prepare(sql string) (*sql.Stmt, error)
	}

	cache map[uint64]cacheItem // map[id]{sql, stmt}
)

func newCache() cache {
	return make(cache)
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c cache) StmtById(prepare Preparer, sqlid uint64) (*sql.Stmt, error) {
	if item, has := c[sqlid]; has {
		sqlPrinter.Print(true, item.sql)

		return item.stmt, nil
	}

	sql_ := sqlById(prepare, sqlid)
	sqlPrinter.Print(false, sql_)

	stmt, err := prepare.Prepare(sql_)
	if err != nil {
		return nil, err
	}

	c[sqlid] = cacheItem{sql: sql_, stmt: stmt}

	return stmt, nil
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c cache) GetStmt(prepare Preparer, sqlid uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		return "", nil, nil
	}

	var err error
	if item.stmt == nil {
		item.stmt, err = prepare.Prepare(item.sql)
	}

	return item.sql, item.stmt, err
}

// SetStmt prepare a sql to statement, cache then return it
func (c cache) SetStmt(prepare Preparer, sqlid uint64, sql string) (*sql.Stmt, error) {
	stmt, err := prepare.Prepare(sql)
	if err != nil {
		return nil, err
	}

	c[sqlid] = cacheItem{
		sql:  sql,
		stmt: stmt,
	}

	return stmt, nil
}

func (c cache) PrepareById(prepare Preparer, sqlid uint64) (*sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		item.sql = sqlById(prepare, sqlid)
		c[sqlid] = item
	}

	sqlPrinter.Print(has, item.sql)

	stmt, err := prepare.Prepare(item.sql)
	return stmt, err
}

func (c cache) PrepareSQL(prepare Preparer, sqlid uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlid]
	if !has {
		return "", nil, nil
	}

	stmt, err := prepare.Prepare(item.sql)
	return item.sql, stmt, err
}

func (c cache) SetSQL(sqlid uint64, sql string) {
	c[sqlid] = cacheItem{
		sql: sql,
	}
}
