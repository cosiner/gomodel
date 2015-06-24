// Package database is a library help for interact with database by model
package gomodel

import "database/sql"

type (
	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		// driver string
		*sql.DB
		tables map[string]*Table
		Cache

		// initial models count for 'All'
		InitialModels int
	}
)

// Open create a database manager and connect to database server
func Open(driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db := NewDB()
	err := db.Connect(driver, dsn, maxIdle, maxOpen)

	return db, err
}

// New create a new db structure
func NewDB() *DB {
	return &DB{
		tables:        make(map[string]*Table),
		InitialModels: 10,
	}
}

// Connect to database server
func (db *DB) Connect(driver, dsn string, maxIdle, maxOpen int) error {
	db_, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	db_.SetMaxIdleConns(maxIdle)
	db_.SetMaxOpenConns(maxOpen)
	db.DB = db_
	db.Cache = NewCache(Types) // use global types count

	return nil
}

// register save table of model
func (db *DB) register(v Model, table string) *Table {
	t := parse(v, db)
	db.tables[table] = t

	return t
}

// Table return infomation of given model
// if table not exist, do parse and save it
func (db *DB) Table(v Model) *Table {
	table := v.Table()
	if t, has := db.tables[table]; has {
		return t
	}

	return db.register(v, table)
}

func (db *DB) Insert(v Model, fields uint, typ ResultType) (int64, error) {
	return db.ArgsInsert(v, fields, typ, FieldVals(v, fields)...)
}

func (db *DB) ArgsInsert(v Model, fields uint, typ ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtInsert(db.DB, fields)

	return Exec(stmt, err, typ, args...)
}

func (db *DB) Update(v Model, fields, whereFields uint) (int64, error) {
	c1, c2 := FieldCount(fields), FieldCount(whereFields)
	args := make([]interface{}, c1+c2)
	v.Vals(fields, args)
	v.Vals(whereFields, args[c1:])

	return db.ArgsUpdate(v, fields, whereFields, args...)
}

func (db *DB) ArgsUpdate(v Model, fields, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtUpdate(db.DB, fields, whereFields)

	return Update(stmt, err, args...)
}

func (db *DB) Delete(v Model, whereFields uint) (int64, error) {
	return db.ArgsDelete(v, whereFields, FieldVals(v, whereFields)...)
}

func (db *DB) ArgsDelete(v Model, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtDelete(db.DB, whereFields)

	return Update(stmt, err, args...)
}

// One select one row from database
func (db *DB) One(v Model, fields, whereFields uint) error {
	stmt, err := db.Table(v).StmtOne(db.DB, fields, whereFields)
	scanner, rows := Query(stmt, err, FieldVals(v, whereFields)...)

	return scanner.One(rows, FieldPtrs(v, fields)...)
}

func (db *DB) Limit(s Store, v Model, fields, whereFields uint, start, count int) error {
	args := FieldVals(v, whereFields, start, count)

	return db.ArgsLimit(s, v, fields, whereFields, args...)
}

func (db *DB) ArgsLimit(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := db.Table(v).StmtLimit(db.DB, fields, whereFields)
	scanner, rows := Query(stmt, err, args...)

	return scanner.Limit(rows, s, args[len(args)-1].(int))
}

func (db *DB) All(s Store, v Model, fields, whereFields uint) error {
	return db.ArgsAll(s, v, fields, whereFields, FieldVals(v, whereFields)...)
}

// ArgsAll select all rows, the last two argument must be "start" and "count"
func (db *DB) ArgsAll(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := db.Table(v).StmtAll(db.DB, fields, whereFields)
	scanner, rows := Query(stmt, err, args...)

	return scanner.All(rows, s, db.InitialModels)
}

// Count return count of rows for model, arguments was extracted from Model
func (db *DB) Count(v Model, whereFields uint) (count int64, err error) {
	return db.ArgsCount(v, whereFields, FieldVals(v, whereFields)...)
}

//Args Count return count of rows for model use custome arguments
func (db *DB) ArgsCount(v Model, whereFields uint,
	args ...interface{}) (count int64, err error) {
	t := db.Table(v)

	stmt, err := t.StmtCount(db.DB, whereFields)
	scanner, rows := Query(stmt, err, args...)

	err = scanner.One(rows, &count)

	return
}

// ExecUpdate execute a update operation, return resolved result
func (db *DB) ExecUpdate(s string, needId bool, args ...interface{}) (int64, error) {
	return db.Exec(s, RES_ROWS, args...)
}

// Exec execute a update operation, return resolved result
func (db *DB) Exec(s string, typ ResultType, args ...interface{}) (int64, error) {
	res, err := db.DB.Exec(s, args...)

	return ResolveResult(res, err, typ)
}

func (db *DB) ExecById(typ uint, is IdSql, resTyp ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.StmtById(db, typ, is)
	return Exec(stmt, err, resTyp, args...)
}

func (db *DB) UpdateById(typ uint, is IdSql, args ...interface{}) (int64, error) {
	return db.ExecById(typ, is, RES_ROWS, args...)
}

func (db *DB) QueryById(typ uint, is IdSql, args ...interface{}) (Scanner, *sql.Rows) {
	stmt, err := db.StmtById(db, typ, is)
	return Query(stmt, err, args...)
}

var emptyTX = Tx{}

func (db *DB) Begin() (Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return emptyTX, err
	}

	return Tx{
		Tx: tx,
		db: db,
	}, nil
}
