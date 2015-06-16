package gomodel

import "database/sql"

type (
	Tx struct {
		*sql.Tx
		db *DB
	}
)

func (tx Tx) Insert(v Model, fields uint, typ ResultType) (int64, error) {
	return tx.ArgsInsert(v, fields, typ, FieldVals(v, fields)...)
}

func (tx Tx) ArgsInsert(v Model, fields uint, typ ResultType, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(v).PrepareInsert(tx.Tx, fields)

	return StmtExec(stmt, err, typ, args...)
}

func (tx Tx) Update(v Model, fields, whereFields uint) (int64, error) {
	c1, c2 := FieldCount(fields), FieldCount(whereFields)
	args := make([]interface{}, c1+c2)
	v.Vals(fields, args)
	v.Vals(whereFields, args[c1:])

	return tx.ArgsUpdate(v, fields, whereFields, args...)
}

func (tx Tx) ArgsUpdate(v Model, fields, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(v).PrepareUpdate(tx.Tx, fields, whereFields)

	return StmtUpdate(stmt, err, args...)
}

func (tx Tx) Delete(v Model, whereFields uint) (int64, error) {
	return tx.ArgsDelete(v, whereFields, FieldVals(v, whereFields)...)
}

func (tx Tx) ArgsDelete(v Model, whereFields uint, args ...interface{}) (int64, error) {
	stmt, err := tx.db.Table(v).PrepareDelete(tx.Tx, whereFields)

	return StmtUpdate(stmt, err, args...)
}

// One select one row from database
func (tx Tx) One(v Model, fields, whereFields uint) error {
	stmt, err := tx.db.Table(v).PrepareOne(tx.Tx, fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, FieldVals(v, whereFields)...)

	return scanner.One(rows, FieldPtrs(v, fields)...)
}

func (tx Tx) Limit(s Store, v Model, fields, whereFields uint, start, count int) error {
	c := FieldCount(whereFields)
	args := make([]interface{}, c+2)
	v.Vals(whereFields, args)
	args[c], args[c+1] = start, count

	return tx.ArgsLimit(s, v, fields, whereFields, args...)
}

func (tx Tx) ArgsLimit(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := tx.db.Table(v).PrepareLimit(tx.Tx, fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

	return scanner.Limit(rows, s, args[len(args)-1].(int))
}

func (tx Tx) All(s Store, v Model, fields, whereFields uint) error {
	return tx.ArgsAll(s, v, fields, whereFields, FieldVals(v, whereFields)...)
}

// ArgsAll select all rows, the last two argument must be "start" and "count"
func (tx Tx) ArgsAll(s Store, v Model, fields, whereFields uint, args ...interface{}) error {
	stmt, err := tx.db.Table(v).PrepareAll(tx.Tx, fields, whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

	return scanner.All(rows, s, tx.db.InitialModels)
}

// Count return count of rows for model, arguments was extracted from Model
func (tx Tx) Count(v Model, whereFields uint) (count int64, err error) {
	return tx.ArgsCount(v, whereFields, FieldVals(v, whereFields)...)
}

//Args Count return count of rows for model use custome arguments
func (tx Tx) ArgsCount(v Model, whereFields uint,
	args ...interface{}) (count int64, err error) {
	stmt, err := tx.db.Table(v).PrepareCount(tx.Tx, whereFields)
	scanner, rows := StmtQuery(stmt, err, args...)

	err = scanner.One(rows, &count)

	return
}

// ExecUpdate execute a update operation, return resolved result
func (tx Tx) ExecUpdate(s string, needId bool, args ...interface{}) (int64, error) {
	return tx.Exec(s, RES_ROWS, args...)
}

// Exec execute a update operation, return resolved result
func (tx Tx) Exec(s string, typ ResultType, args ...interface{}) (int64, error) {
	res, err := tx.Tx.Exec(s, args...)

	return ResolveResult(res, err, typ)
}
