package gomodel

import "database/sql"

type (
	Executor interface {
		Driver() Driver
		Table(model Model) *Table
		Prepare(sql string) (*sql.Stmt, error)

		Insert(model Model, fields uint64, resType ResultType) (int64, error)
		ArgsInsert(model Model, fields uint64, resType ResultType, args ...interface{}) (int64, error)

		Update(model Model, fields, whereFields uint64) (int64, error)
		ArgsUpdate(model Model, fields, whereFields uint64, args ...interface{}) (int64, error)

		Delete(model Model, whereFields uint64) (int64, error)
		ArgsDelete(model Model, whereFields uint64, args ...interface{}) (int64, error)

		One(model Model, fields, whereFields uint64) error
		ArgsOne(model Model, fields, whereFields uint64, args []interface{}, ptrs ...interface{}) error

		Limit(store Store, model Model, fields, whereFields uint64, start, count int64) error
		ArgsLimit(store Store, model Model, fields, whereFields uint64, args ...interface{}) error

		All(store Store, model Model, fields, whereFields uint64) error
		ArgsAll(store Store, model Model, fields, whereFields uint64, args ...interface{}) error

		Count(model Model, whereFields uint64) (count int64, err error)
		ArgsCount(model Model, whereFields uint64, args ...interface{}) (count int64, err error)

		IncrBy(model Model, field, whereFields uint64, count int) (int64, error)
		ArgsIncrBy(model Model, field, whereFields uint64, args ...interface{}) (int64, error)

		Exists(model Model, field, whereFields uint64) (bool, error)
		ArgsExists(model Model, field, whereFields uint64, args ...interface{}) (bool, error)

		ExecUpdate(sql string, args ...interface{}) (int64, error)
		Exec(sql string, resType ResultType, args ...interface{}) (int64, error)

		ExecById(sqlid uint64, resTyp ResultType, args ...interface{}) (int64, error)
		UpdateById(sqlid uint64, args ...interface{}) (int64, error)
		QueryById(sqlid uint64, args ...interface{}) Scanner
	}

	ResultType int
)

const (
	RES_NO   ResultType = iota // don't resolve sql.Result
	RES_ID                     // Result.LastInsertID
	RES_ROWS                   // Result.RowsAffected
)

// Update always returl the count of affected rows
func Update(stmt Stmt, err error, args ...interface{}) (int64, error) {
	return Exec(stmt, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func Exec(stmt Stmt, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(args...)
	return ResolveResult(res, err, typ)
}

// Query execute the query stmt, error stored in Scanner
func Query(stmt Stmt, err error, args ...interface{}) Scanner {
	if err != nil {
		return Scanner{Error: err}
	}
	rows, err := stmt.Query(args...)
	return Scanner{
		Error: err,
		Rows:  rows,
		Stmt:  stmt,
	}
}

// Update always returl the count of affected rows
func CloseUpdate(stmt Stmt, err error, args ...interface{}) (int64, error) {
	return CloseExec(stmt, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func CloseExec(stmt Stmt, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err == nil {
		defer stmt.Close()
	}

	return Exec(stmt, err, typ, args...)
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
