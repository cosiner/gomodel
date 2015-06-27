package gomodel

import "database/sql"

type (
	// cacheItem keeps the sql and prepared statement of it
	cacheItem struct {
		sql  string
		stmt *sql.Stmt
	}

	Preparer interface {
		Prepare(sql string) (*sql.Stmt, error)
	}

	cache map[uint64]cacheItem // [type]map[id]{sql, stmt}
)

func newCache() cache {
	return make(cache)
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c cache) StmtById(prepare Preparer, idsql IdSql) (*sql.Stmt, error) {
	if item, has := c[idsql.ID]; has {
		sqlPrinter.Print(true, item.sql)

		return item.stmt, nil
	}

	sql_ := idsql.SQL
	sqlPrinter.Print(false, sql_)

	stmt, err := prepare.Prepare(sql_)
	if err != nil {
		return nil, err
	}

	c[idsql.ID] = cacheItem{sql: sql_, stmt: stmt}

	return stmt, nil
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c cache) GetStmt(prepare Preparer, sqlId uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlId]
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
func (c cache) SetStmt(prepare Preparer, sqlId uint64, sql string) (*sql.Stmt, error) {
	stmt, err := prepare.Prepare(sql)
	if err != nil {
		return nil, err
	}

	c[sqlId] = cacheItem{
		sql:  sql,
		stmt: stmt,
	}

	return stmt, nil
}

func (c cache) PrepareById(prepare Preparer, idsql IdSql) (*sql.Stmt, error) {
	item, has := c[idsql.ID]
	if !has {
		item.sql = idsql.SQL
		c[idsql.ID] = item
	}

	sqlPrinter.Print(has, item.sql)

	stmt, err := prepare.Prepare(item.sql)
	return stmt, err
}

func (c cache) PrepareSQL(prepare Preparer, sqlId uint64) (string, *sql.Stmt, error) {
	item, has := c[sqlId]
	if !has {
		return "", nil, nil
	}

	stmt, err := prepare.Prepare(item.sql)
	return item.sql, stmt, err
}

func (c cache) SetSQL(sqlId uint64, sql string) {
	c[sqlId] = cacheItem{
		sql: sql,
	}
}
