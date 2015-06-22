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
func NewCache(types uint) Cache {
	c := Cache{
		cache: make([]map[uint]cacheItem, types),
	}

	for i := uint(0); i < types; i++ {
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
func (c *Cache) ExtendType(typ uint) uint {
	if l := uint(len(c.cache)); typ > l {
		cache := make([]map[uint]cacheItem, typ)
		copy(cache, c.cache)

		for ; l < typ; l++ {
			cache[l] = make(map[uint]cacheItem)
		}
		c.cache = cache
	} else {
		c.cache = c.cache[:typ]
	}

	return typ - 1
}

// Types return the sql types count of current Cache
func (c *Cache) Types() uint {
	return uint(len(c.cache))
}

// StmtById search a prepared statement for given sql type by id, if not found,
// create with the creator, and prepared the sql to a statement, cache it, then
// return
func (c *Cache) StmtById(p Preparer, typ uint, is IdSql) (*sql.Stmt, error) {
	if item, has := c.cache[typ][is.ID]; has {
		sqlPrinter.Print(true, item.sql)

		return item.stmt, nil
	}

	sql_ := is.SQL
	sqlPrinter.Print(false, sql_)

	stmt, err := p.Prepare(sql_)
	if err != nil {
		return nil, err
	}

	c.cache[typ][is.ID] = cacheItem{sql: sql_, stmt: stmt}

	return stmt, nil
}

// GetStmt get sql and statement from cacher, if not found, "" and nil was returned
func (c *Cache) GetStmt(p Preparer, typ, id uint) (string, *sql.Stmt, error) {
	item, has := c.cache[typ][id]
	if !has {
		return "", nil, nil
	}

	var err error
	if item.stmt == nil {
		item.stmt, err = p.Prepare(item.sql)
	}

	return item.sql, item.stmt, err
}

// SetStmt prepare a sql to statement, cache then return it
func (c *Cache) SetStmt(p Preparer, typ uint, id uint, sql string) (*sql.Stmt, error) {
	stmt, err := p.Prepare(sql)
	if err != nil {
		return nil, err
	}

	c.cache[typ][id] = cacheItem{
		sql:  sql,
		stmt: stmt,
	}

	return stmt, nil
}

func (c *Cache) PrepareById(p Preparer, typ uint, is IdSql) (*sql.Stmt, error) {
	item, has := c.cache[typ][is.ID]
	if !has {
		item.sql = is.SQL
		c.cache[typ][is.ID] = item
	}

	sqlPrinter.Print(has, item.sql)

	stmt, err := p.Prepare(item.sql)
	return stmt, err
}

func (c *Cache) PrepareSQL(p Preparer, typ, id uint) (string, *sql.Stmt, error) {
	item, has := c.cache[typ][id]
	if !has {
		return "", nil, nil
	}

	stmt, err := p.Prepare(item.sql)
	return item.sql, stmt, err
}

func (c *Cache) SetSQL(typ, id uint, sql string) {
	c.cache[typ][id] = cacheItem{
		sql: sql,
	}
}
