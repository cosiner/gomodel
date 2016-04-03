// Package gomodel is a library help for interact with database efficiently
package gomodel

import (
	"database/sql"

	"github.com/cosiner/gohper/conv"
	"github.com/cosiner/gohper/errors"
)

type (
	// DB holds database connections, store all tables
	DB struct {
		*sql.DB
		driver Driver
		tables map[string]*Table
		cache  cache

		// initial models count for select 'All', default 20
		InitialModels int
	}
)

// Open create a database manager and connect to database server
func Open(driver Driver, dsn string, maxIdle, maxOpen int) (*DB, error) {
	db := NewDB()
	err := db.Connect(driver, dsn, maxIdle, maxOpen)

	return db, err
}

// NewDB create a new DB instance
func NewDB() *DB {
	initSqlStore()

	return &DB{
		tables:        make(map[string]*Table),
		InitialModels: 20,
	}
}

// Connect to database server
func (db *DB) Connect(driver Driver, dsn string, maxIdle, maxOpen int) error {
	db_, err := sql.Open(driver.String(), dsn)
	if err != nil {
		return err
	}
	db.driver = driver

	db_.SetMaxIdleConns(maxIdle)
	db_.SetMaxOpenConns(maxOpen)
	db.DB = db_
	db.cache = newCache()

	return nil
}

func (db *DB) Driver() Driver {
	return db.driver
}

// Table return infomation of given model
// if table not exist, do parse and save it
func (db *DB) Table(model Model) *Table {
	table := model.Table()
	t, has := db.tables[table]
	if !has {
		t = parseModel(model, db)
		db.tables[table] = t
	}

	return t
}

func (db *DB) Insert(model Model, fields uint64, resType ResultType) (int64, error) {
	return db.ArgsInsert(model, fields, resType, FieldVals(model, fields)...)
}

func (db *DB) ArgsInsert(model Model, fields uint64, resType ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.Table(model).StmtInsert(db, fields)

	return Exec(stmt, err, resType, args...)
}

func (db *DB) Update(model Model, fields, whereFields uint64) (int64, error) {
	c1, c2 := NumFields(fields), NumFields(whereFields)
	args := make([]interface{}, c1+c2)
	model.Vals(fields, args)
	model.Vals(whereFields, args[c1:])

	return db.ArgsUpdate(model, fields, whereFields, args...)
}

func (db *DB) ArgsUpdate(model Model, fields, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := db.Table(model).StmtUpdate(db, fields, whereFields)

	return Update(stmt, err, args...)
}

func (db *DB) Delete(model Model, whereFields uint64) (int64, error) {
	return db.ArgsDelete(model, whereFields, FieldVals(model, whereFields)...)
}

func (db *DB) ArgsDelete(model Model, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := db.Table(model).StmtDelete(db, whereFields)

	return Update(stmt, err, args...)
}

// One select one row from database
func (db *DB) One(model Model, fields, whereFields uint64) error {
	return db.ArgsOne(model, fields, whereFields, FieldVals(model, whereFields))
}

func (db *DB) ArgsOne(model Model, fields, whereFields uint64, args []interface{}, ptrs ...interface{}) error {
	stmt, err := db.Table(model).StmtOne(db, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	if len(ptrs) == 0 {
		ptrs = FieldPtrs(model, fields)
	}
	return scanner.One(ptrs...)
}

func (db *DB) Limit(store Store, model Model, fields, whereFields uint64, start, count int64) error {
	args := FieldVals(model, whereFields, start, count)

	return db.ArgsLimit(store, model, fields, whereFields, args...)
}

// The last two arguments must be "start" and "count" of limition with type "int"
func (db *DB) ArgsLimit(store Store, model Model, fields, whereFields uint64, args ...interface{}) error {
	argc := len(args)
	if argc < 2 {
		panic(errors.Newf("ArgsLimit need at least two parameters, but only got %d", argc))
	}
	offset, err := conv.IfaceToInt(args[argc-2])
	errors.Panicln(err)
	count, err := conv.IfaceToInt(args[argc-1])
	errors.Panicln(err)

	arg1, arg2 := db.driver.ParamLimit(offset, count)
	args[argc-2], args[argc-1] = arg1, arg2

	stmt, err := db.Table(model).StmtLimit(db, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	return scanner.Limit(store, int(count))
}

func (db *DB) All(store Store, model Model, fields, whereFields uint64) error {
	return db.ArgsAll(store, model, fields, whereFields, FieldVals(model, whereFields)...)
}

func (db *DB) ArgsAll(store Store, model Model, fields, whereFields uint64, args ...interface{}) error {
	stmt, err := db.Table(model).StmtAll(db, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	return scanner.All(store, db.InitialModels)
}

// Count return count of rows for model, arguments was extracted from Model
func (db *DB) Count(model Model, whereFields uint64) (count int64, err error) {
	return db.ArgsCount(model, whereFields, FieldVals(model, whereFields)...)
}

// ArgsCount return count of rows for model use custome arguments
func (db *DB) ArgsCount(model Model, whereFields uint64, args ...interface{}) (count int64, err error) {
	t := db.Table(model)

	stmt, err := t.StmtCount(db, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	err = scanner.One(&count)

	return
}

func (db *DB) IncrBy(model Model, field, whereFields uint64, count int) (int64, error) {
	args := make([]interface{}, NumFields(whereFields)+1)
	args[0] = count
	model.Vals(whereFields, args[1:])

	return db.ArgsIncrBy(model, field, whereFields, args...)
}

func (db *DB) ArgsIncrBy(model Model, field, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := db.Table(model).StmtIncrBy(db, field, whereFields)

	return Update(stmt, err, args...)
}

func (db *DB) ExecById(sqlid uint64, resTyp ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.StmtById(sqlid)

	return Exec(stmt, err, resTyp, args...)
}

func (db *DB) UpdateById(sqlid uint64, args ...interface{}) (int64, error) {
	return db.ExecById(sqlid, RES_ROWS, args...)
}

func (db *DB) QueryById(sqlid uint64, args ...interface{}) Scanner {
	stmt, err := db.StmtById(sqlid)

	return Query(stmt, err, args...)
}

func (db *DB) prepare(sql string) (Stmt, error) {
	sql = db.driver.Prepare(sql)
	sqlPrinter(sql)
	stmt, err := db.DB.Prepare(sql)
	return WrapStmt(true, stmt, err)
}

func (db *DB) Query(sql string, args ...interface{}) Scanner {
	stmt, err := db.prepare(sql)
	return Query(stmt, err, args...)
}

func (db *DB) Exec(sql string, resTyp ResultType, args ...interface{}) (int64, error) {
	stmt, err := db.prepare(sql)
	return CloseExec(stmt, err, resTyp, args...)
}

func (db *DB) ExecUpdate(sql string, args ...interface{}) (int64, error) {
	return db.Exec(sql, RES_ROWS, args...)
}

var emptyTX = &Tx{}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return emptyTX, err
	}

	return newTx(tx, db), nil
}

func (db *DB) TxDo(do func(*Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Close()

	err = do(tx)
	tx.Success(err == nil)
	return err
}

func (db *DB) StmtById(sqlid uint64) (Stmt, error) {
	stmt, err := db.cache.StmtById(db, sqlid)

	return WrapStmt(STMT_NOPCLOSE, stmt, err)
}
