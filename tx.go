package gomodel

import "database/sql"

type (
	Tx struct {
		*sql.Tx
		db *DB
	}
)

func (tx Tx) Insert(model Model, fields uint, resType ResultType) (int64, error) {
	return tx.ArgsInsert(model, fields, resType, FieldVals(model, fields)...)
}

func (tx Tx) ArgsInsert(model Model, fields uint, resType ResultType, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(model).PrepareInsert(tx.Tx, fields)

	return CloseExec(stmt, err, resType, args...)
}

func (tx Tx) Update(model Model, fields, whereFields uint) (int64, error) {
	c1, c2 := NumFields(fields), NumFields(whereFields)
	args := make([]interface{}, c1+c2)
	model.Vals(fields, args)
	model.Vals(whereFields, args[c1:])

	return tx.ArgsUpdate(model, fields, whereFields, args...)
}

func (tx Tx) ArgsUpdate(model Model, fields, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(model).PrepareUpdate(tx.Tx, fields, whereFields)

	return CloseUpdate(stmt, err, args...)
}

func (tx Tx) Delete(model Model, whereFields uint) (int64, error) {
	return tx.ArgsDelete(model, whereFields, FieldVals(model, whereFields)...)
}

func (tx Tx) ArgsDelete(model Model, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(model).PrepareDelete(tx.Tx, whereFields)

	return CloseUpdate(stmt, err, args...)
}

// One select one row from database
// ArgsOne is unnecessary, just put result in model
func (tx Tx) One(model Model, fields, whereFields uint) error {
	stmt, err := tx.db.Table(model).PrepareOne(tx.Tx, fields, whereFields)
	scanner := CloseQuery(stmt, err, FieldVals(model, whereFields)...)

	return scanner.One(FieldPtrs(model, fields)...)
}

func (tx Tx) Limit(store Store, model Model, fields, whereFields uint, start, count int) error {
	c := NumFields(whereFields)
	args := make([]interface{}, c+2)
	model.Vals(whereFields, args)
	args[c], args[c+1] = start, count

	return tx.ArgsLimit(store, model, fields, whereFields, args...)
}

// The last two arguments must be "start" and "count" of limition with type "int"
func (tx Tx) ArgsLimit(store Store, model Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := tx.db.Table(model).PrepareLimit(tx.Tx, fields, whereFields)
	scanner := CloseQuery(stmt, err, args...)

	return scanner.Limit(store, args[len(args)-1].(int))
}

func (tx Tx) All(store Store, model Model, fields, whereFields uint) error {
	return tx.ArgsAll(store, model, fields, whereFields, FieldVals(model, whereFields)...)
}

// ArgsAll select all rows, the last two argument must be "start" and "count"
func (tx Tx) ArgsAll(store Store, model Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := tx.db.Table(model).PrepareAll(tx.Tx, fields, whereFields)
	scanner := CloseQuery(stmt, err, args...)

	return scanner.All(store, tx.db.InitialModels)
}

// Count return count of rows for model, arguments was extracted from Model
func (tx Tx) Count(model Model, whereFields uint) (count int64, err error) {
	return tx.ArgsCount(model, whereFields, FieldVals(model, whereFields)...)
}

// ArgsCount return count of rows for model use custome arguments
func (tx Tx) ArgsCount(model Model, whereFields uint,
	args ...interface{}) (count int64, err error) {
	stmt, err := tx.db.Table(model).PrepareCount(tx.Tx, whereFields)
	scanner := CloseQuery(stmt, err, args...)

	err = scanner.One(&count)

	return
}

// ExecUpdate execute a update operation, return resolved result
func (tx Tx) ExecUpdate(sql string, args ...interface{}) (int64, error) {
	return tx.Exec(sql, RES_ROWS, args...)
}

// Exec execute a update operation, return resolved result
func (tx Tx) Exec(sql string, resType ResultType, args ...interface{}) (int64, error) {
	res, err := tx.Tx.Exec(sql, args...)

	return ResolveResult(res, err, resType)
}

// Done check if error is nil then commit transaction, otherwise rollback, the error
// will be returned without change.
func (tx Tx) Done(err error) error {
	if err == nil {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()
	}

	return err
}
