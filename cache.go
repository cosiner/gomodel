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

	// Cache store all the sql, statement by sql type and id
	// typically, the sql id of predefied sql type is
	// fields << numField( of Model) | whereFields,
	// it's used for a single model
	//
	// if custom is necessary, call cache.ExtendType(cache.Types()+1) to make
	// a new type, the sql id is bring your owns, also you can still use the standard
	// FieldIdentity(fields, whereFields) if possible
	Cache struct {
		cache []map[uint]cacheItem // [type]map[id]{sql, stmt}
	}
)

const (
	// These are five predefined sql types
	INSERT uint = iota
	DELETE
	UPDATE
	SELECT_LIMIT
	SELECT_ONE
	SELECT_ALL

	defaultTypeEnd
)

var (
	// Types defines the default sql types count, it's default applied to all
	// models.
	// Change it before register any models.
	Types = defaultTypeEnd
)

// NewCache create a common sql and statement cacher with given types count
// this will make no parameter checks
//
// the DB instance and each Table already embed a Cache, typically, it's not
// necessary to call this
func NewCache(sqlTypes uint) Cache {
	c := Cache{
		cache: make([]map[uint]cacheItem, sqlTypes),
	}

	for i := uint(0); i < sqlTypes; i++ {
		c.cache[i] = make(map[uint]cacheItem)
	}

	return c
}

// ExtendType typically used to extend types of Cache, but it also can be used
// to shrink the cacher, return value will be the new types count you passed in
//
// Example:
// //a.go
// newType1 := c.ExtendType(c.Types()+1)
// //b.go
// newType2 := c.ExtendType(c.Types()+1)
func (c *Cache) ExtendType(sqlType uint) uint {
	if l := uint(len(c.cache)); sqlType > l {
		cache := make([]map[uint]cacheItem, sqlType)
		copy(cache, c.cache)

		for ; l < sqlType; l++ {
			cache[l] = make(map[uint]cacheItem)
		}
		c.cache = cache
	} else {
		c.cache = c.cache[:sqlType]
	}

	return sqlType - 1
}

// Types return the sql types count of current Cache
func (c *Cache) Types() uint {
	return uint(len(c.cache))
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c *Cache) StmtById(prepare Preparer, sqlType uint, idsql IdSql) (*sql.Stmt, error) {
	if item, has := c.cache[sqlType][idsql.ID]; has {
		sqlPrinter.Print(true, item.sql)

		return item.stmt, nil
	}

	sql_ := idsql.SQL
	sqlPrinter.Print(false, sql_)

	stmt, err := prepare.Prepare(sql_)
	if err != nil {
		return nil, err
	}

	c.cache[sqlType][idsql.ID] = cacheItem{sql: sql_, stmt: stmt}

	return stmt, nil
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c *Cache) GetStmt(prepare Preparer, sqlType, sqlId uint) (string, *sql.Stmt, error) {
	item, has := c.cache[sqlType][sqlId]
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
func (c *Cache) SetStmt(prepare Preparer, sqlType uint, sqlId uint, sql string) (*sql.Stmt, error) {
	stmt, err := prepare.Prepare(sql)
	if err != nil {
		return nil, err
	}

	c.cache[sqlType][sqlId] = cacheItem{
		sql:  sql,
		stmt: stmt,
	}

	return stmt, nil
}

func (c *Cache) PrepareById(prepare Preparer, sqlType uint, idsql IdSql) (*sql.Stmt, error) {
	item, has := c.cache[sqlType][idsql.ID]
	if !has {
		item.sql = idsql.SQL
		c.cache[sqlType][idsql.ID] = item
	}

	sqlPrinter.Print(has, item.sql)

	stmt, err := prepare.Prepare(item.sql)
	return stmt, err
}

func (c *Cache) PrepareSQL(prepare Preparer, sqlType, sqlId uint) (string, *sql.Stmt, error) {
	item, has := c.cache[sqlType][sqlId]
	if !has {
		return "", nil, nil
	}

	stmt, err := prepare.Prepare(item.sql)
	return item.sql, stmt, err
}

func (c *Cache) SetSQL(sqlType, sqlId uint, sql string) {
	c.cache[sqlType][sqlId] = cacheItem{
		sql: sql,
	}
}
