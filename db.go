// Package database is a library help for interact with database by model
package gomodel

import (
	"database/sql"

	"github.com/cosiner/gohper/bitset"
)

type (
	// Model represent a database model
	Model interface {
		Table() string
		// Vals store values of fields to given slice
		Vals(fields uint, vals []interface{})
		Ptrs(fields uint, ptrs []interface{})
	}

	// DB holds database connection, all typeinfos, and sql cache
	DB struct {
		// driver string
		*sql.DB
		tables map[string]*Table
		Cacher

		ModelCount int
	}

	ResultType int
)

const (
	RES_NO ResultType = iota
	RES_ID
	RES_ROWS
)

var (
	FieldCount = bitset.BitCountUint
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
		tables:     make(map[string]*Table),
		ModelCount: 10,
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
	db.Cacher = NewCacher(Types, db) // use global types count

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

func FieldVals(v Model, fields uint) []interface{} {
	args := make([]interface{}, FieldCount(fields))
	v.Vals(fields, args)

	return args
}

func FieldPtrs(v Model, fields uint) []interface{} {
	ptrs := make([]interface{}, FieldCount(fields))
	v.Ptrs(fields, ptrs)

	return ptrs
}

func (db *DB) Insert(v Model, fields uint, typ ResultType) (int64, error) {
	return db.ArgsInsert(v, fields, typ, FieldVals(v, fields))
}

func (db *DB) ArgsInsert(v Model, fields uint, typ ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtInsert(fields)

	return StmtExec(stmt, err, typ, args...)
}

func (db *DB) Update(v Model, fields, whereFields uint) (int64, error) {
	c1, c2 := FieldCount(fields), FieldCount(whereFields)
	args := make([]interface{}, c1+c2)
	v.Vals(fields, args)
	v.Vals(whereFields, args[c1:])

	return db.ArgsUpdate(v, fields, whereFields, args...)
}

func (db *DB) ArgsUpdate(v Model, fields, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtUpdate(fields, whereFields)

	return StmtUpdate(stmt, err, args...)
}

func (db *DB) Delete(v Model, whereFields uint) (int64, error) {
	return db.ArgsDelete(v, whereFields, FieldVals(v, whereFields))
}

func (db *DB) ArgsDelete(v Model, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := db.Table(v).StmtDelete(whereFields)

	return StmtUpdate(stmt, err, args...)
}

// One select one row from database
func (db *DB) One(v Model, fields, whereFields uint) error {
	stmt, err := db.Table(v).StmtOne(fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, FieldVals(v, whereFields)...)

	return scanner.One(rows, FieldPtrs(v, fields)...)
}

func (db *DB) Limit(s Store, v Model, fields, whereFields uint, start, count int) error {
	c := FieldCount(whereFields)
	args := make([]interface{}, c+2)
	v.Vals(whereFields, args)
	args[c], args[c+1] = start, count

	return db.ArgsLimit(s, v, fields, whereFields, args...)
}

func (db *DB) ArgsLimit(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := db.Table(v).StmtLimit(fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

	return scanner.Limit(rows, s, args[len(args)-1].(int))
}

func (db *DB) All(s Store, v Model, fields, whereFields uint) error {
	return db.ArgsAll(s, v, fields, whereFields, FieldVals(v, whereFields))
}

// ArgsAll select all rows, the last two argument must be "start" and "count"
func (db *DB) ArgsAll(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := db.Table(v).StmtAll(fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

	return scanner.All(rows, s, db.ModelCount)
}

// Count return count of rows for model, arguments was extracted from Model
func (db *DB) Count(v Model, whereFields uint) (count int64, err error) {
	return db.ArgsCount(v, whereFields, FieldVals(v, whereFields))
}

//Args Count return count of rows for model use custome arguments
func (db *DB) ArgsCount(v Model, whereFields uint,
	args ...interface{}) (count int64, err error) {
	t := db.Table(v)

	stmt, err := t.StmtCount(whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

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

// StmtUpdate always returl the count of affected rows
func StmtUpdate(stmt *sql.Stmt, err error, args ...interface{}) (int64, error) {
	return StmtExec(stmt, err, RES_ROWS, args...)
}

// StmtExec execute stmt with given arguments and resolve the result if error is nil
func StmtExec(stmt *sql.Stmt, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	return ResolveResult(res, err, typ)
}

// StmtQuery execute the query stmt, error stored in Scanner
func StmtQuery(stmt *sql.Stmt, err error, args ...interface{}) (Scanner, *sql.Rows) {
	if err != nil {
		return Scanner{err}, nil
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return Scanner{err}, nil
	}

	return normalScanner, rows
}

// ResolveResult resolve sql result, if need id, return last insert id
// else return affected rows count
func ResolveResult(res sql.Result, err error, typ ResultType) (int64, error) {
	if err != nil {
		return 0, err
	}

	switch typ {
	case RES_NO:
		return 0, nil
	case RES_ID:
		return res.LastInsertId()
	case RES_ROWS:
		return res.RowsAffected()
	default:
		panic("unexpected result type")
	}
}
