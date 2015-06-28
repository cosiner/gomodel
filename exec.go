package gomodel

import "database/sql"

type (
	QueryExecer interface {
		Exec(...interface{}) (sql.Result, error)
		Query(...interface{}) (*sql.Rows, error)
		Close() error
	}

	ResultType int
)

const (
	RES_NO   ResultType = iota // don't resolve sql.Result
	RES_ID                     // Result.LastInsertID
	RES_ROWS                   // Result.RowsAffected
)

// Update always returl the count of affected rows
func Update(exec QueryExecer, err error, args ...interface{}) (int64, error) {
	return Exec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func Exec(exec QueryExecer, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err != nil {
		return 0, err
	}

	res, err := exec.Exec(args...)
	return ResolveResult(res, err, typ)
}

// Query execute the query stmt, error stored in Scanner
func Query(exec QueryExecer, err error, args ...interface{}) Scanner {
	if err != nil {
		return Scanner{Error: err}
	}

	rows, err := exec.Query(args...)
	return Scanner{
		Error: err,
		Rows:  rows,
	}
}

// Update always returl the count of affected rows
func CloseUpdate(exec QueryExecer, err error, args ...interface{}) (int64, error) {
	return CloseExec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func CloseExec(exec QueryExecer, err error, typ ResultType, args ...interface{}) (int64, error) {
	defer exec.Close()

	return Exec(exec, err, typ, args...)
}

// Query execute the query stmt, error stored in Scanner
func CloseQuery(exec QueryExecer, err error, args ...interface{}) Scanner {
	defer exec.Close()

	return Query(exec, err, args...)
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
