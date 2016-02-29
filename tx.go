package gomodel

import "database/sql"

type (
	Tx struct {
		*sql.Tx
		db        *DB
		isSuccess bool
	}
)

func newTx(tx *sql.Tx, db *DB) *Tx {
	return &Tx{
		Tx:        tx,
		db:        db,
		isSuccess: true,
	}
}

func (tx *Tx) Driver() Driver {
	return tx.db.Driver()
}

func (tx *Tx) Table(model Model) *Table {
	return tx.db.Table(model)
}

func (tx *Tx) Insert(model Model, fields uint64, resType ResultType) (int64, error) {
	return tx.ArgsInsert(model, fields, resType, FieldVals(model, fields)...)
}

func (tx *Tx) ArgsInsert(model Model, fields uint64, resType ResultType, args ...interface{}) (int64, error) {
	stmt, err := tx.Table(model).PrepareInsert(tx, fields)

	return CloseExec(stmt, err, resType, args...)
}

func (tx *Tx) Update(model Model, fields, whereFields uint64) (int64, error) {
	c1, c2 := NumFields(fields), NumFields(whereFields)
	args := make([]interface{}, c1+c2)
	model.Vals(fields, args)
	model.Vals(whereFields, args[c1:])

	return tx.ArgsUpdate(model, fields, whereFields, args...)
}

func (tx *Tx) ArgsUpdate(model Model, fields, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := tx.Table(model).PrepareUpdate(tx, fields, whereFields)

	return CloseUpdate(stmt, err, args...)
}

func (tx *Tx) Delete(model Model, whereFields uint64) (int64, error) {
	return tx.ArgsDelete(model, whereFields, FieldVals(model, whereFields)...)
}

func (tx *Tx) ArgsDelete(model Model, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := tx.Table(model).PrepareDelete(tx, whereFields)

	return CloseUpdate(stmt, err, args...)
}

// One select one row from database
func (tx *Tx) One(model Model, fields, whereFields uint64) error {
	return tx.ArgsOne(model, fields, whereFields, FieldVals(model, whereFields))
}

func (tx *Tx) ArgsOne(model Model, fields, whereFields uint64, args []interface{}, ptrs ...interface{}) error {
	stmt, err := tx.Table(model).PrepareOne(tx, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	if len(ptrs) == 0 {
		ptrs = FieldPtrs(model, fields)
	}
	return scanner.One(ptrs...)
}

func (tx *Tx) Limit(store Store, model Model, fields, whereFields uint64, start, count int64) error {
	args := FieldVals(model, whereFields, start, count)

	return tx.ArgsLimit(store, model, fields, whereFields, args...)
}

// The last two arguments must be "start" and "count" of limition with type "int"
func (tx *Tx) ArgsLimit(store Store, model Model, fields, whereFields uint64, args ...interface{}) error {
	stmt, err := tx.Table(model).PrepareLimit(tx, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	return scanner.Limit(store, args[len(args)-1].(int))
}

func (tx *Tx) All(store Store, model Model, fields, whereFields uint64) error {
	return tx.ArgsAll(store, model, fields, whereFields, FieldVals(model, whereFields)...)
}

// ArgsAll select all rows, the last two argument must be "start" and "count"
func (tx *Tx) ArgsAll(store Store, model Model, fields, whereFields uint64, args ...interface{}) error {
	stmt, err := tx.Table(model).PrepareAll(tx, fields, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()

	return scanner.All(store, tx.db.InitialModels)
}

// Count return count of rows for model, arguments was extracted from Model
func (tx *Tx) Count(model Model, whereFields uint64) (count int64, err error) {
	return tx.ArgsCount(model, whereFields, FieldVals(model, whereFields)...)
}

// ArgsCount return count of rows for model use custome arguments
func (tx *Tx) ArgsCount(model Model, whereFields uint64,
	args ...interface{}) (count int64, err error) {
	stmt, err := tx.Table(model).PrepareCount(tx, whereFields)
	scanner := Query(stmt, err, args...)
	defer scanner.Close()
	err = scanner.One(&count)

	return
}

func (tx *Tx) IncrBy(model Model, field, whereFields uint64, count int) (int64, error) {
	args := make([]interface{}, NumFields(whereFields)+1)
	args[0] = count
	model.Vals(whereFields, args[1:])

	return tx.ArgsIncrBy(model, field, whereFields, args...)
}

func (tx *Tx) ArgsIncrBy(model Model, field, whereFields uint64, args ...interface{}) (int64, error) {
	stmt, err := tx.Table(model).PrepareIncrBy(tx, field, whereFields)

	return CloseUpdate(stmt, err, args...)
}

// ExecUpdate execute a update operation, return resolved result
func (tx *Tx) ExecUpdate(sql string, args ...interface{}) (int64, error) {
	return tx.Exec(sql, RES_ROWS, args...)
}

// Exec execute a update operation, return resolved result
func (tx *Tx) Exec(sql string, resType ResultType, args ...interface{}) (int64, error) {
	res, err := tx.Tx.Exec(sql, args...)

	return ResolveResult(res, err, resType)
}

// ExecById execute a update operation, return rows affected
func (tx *Tx) ExecById(sqlid uint64, resType ResultType, args ...interface{}) (int64, error) {
	stmt, err := tx.PrepareById(sqlid)

	return CloseExec(stmt, err, resType, args...)
}

// UpdateById execute a update operation, return resolved result
func (tx *Tx) UpdateById(sqlid uint64, args ...interface{}) (int64, error) {
	return tx.ExecById(sqlid, RES_ROWS, args...)
}

func (tx *Tx) QueryById(sqlid uint64, args ...interface{}) Scanner {
	stmt, err := tx.PrepareById(sqlid)

	return Query(stmt, err, args...)
}

// Done check if error is nil then commit transaction, otherwise rollback.
// Done should be called only once, otherwise it will panic.
// Done should be called in deferred function to avoid uncommitted/unrollbacked
// transaction caused by panic.
//
// Example:
//  func operation() (err error) {
//  	tx, err := db.Begin()
//  	if err != nil {
//  		return err
//  	}
//  	defer tx.Close()
//
//  	// do something
//
//		tx.Success(true)
//  	return // err must be a named return value, otherwise, error in deferred
//  	       isSuccessfunction will be lost
//  }
func (tx *Tx) Close() error {
	if tx.isSuccess {
		err := tx.Commit()
		if err == sql.ErrTxDone {
			panic("commit transaction twice")
		}
		return err
	}

	err := tx.Rollback()
	if err == sql.ErrTxDone {
		panic("rollback a committed/rollbacked transaction")
	}
	return err
}

func (tx *Tx) Success(success bool) {
	tx.isSuccess = tx.isSuccess && success
}

func (tx *Tx) PrepareById(sqlid uint64) (Stmt, error) {
	stmt, err := tx.db.cache.PrepareById(tx, sqlid)

	return WrapStmt(STMT_CLOSEABLE, stmt, err)
}
