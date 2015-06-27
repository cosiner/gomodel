package gomodel

import "database/sql"

type (
	Executor interface {
		Exec(...interface{}) (sql.Result, error)
		Query(...interface{}) (*sql.Rows, error)
		Close() error
	}

	ResultType int
)

const (
	RES_NO ResultType = iota
	RES_ID
	RES_ROWS
)

// Update always returl the count of affected rows
func Update(exec Executor, err error, args ...interface{}) (int64, error) {
	return Exec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func Exec(exec Executor, err error, typ ResultType, args ...interface{}) (int64, error) {
	if err != nil {
		return 0, err
	}

	res, err := exec.Exec(args...)
	return ResolveResult(res, err, typ)
}

// Query execute the query stmt, error stored in Scanner
func Query(exec Executor, err error, args ...interface{}) Scanner {
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
func CloseUpdate(exec Executor, err error, args ...interface{}) (int64, error) {
	return CloseExec(exec, err, RES_ROWS, args...)
}

// Exec execute stmt with given arguments and resolve the result if error is nil
func CloseExec(exec Executor, err error, typ ResultType, args ...interface{}) (int64, error) {
	defer exec.Close()
	c, err := Exec(exec, err, typ, args...)

	return c, err
}

// Query execute the query stmt, error stored in Scanner
func CloseQuery(exec Executor, err error, args ...interface{}) Scanner {
	defer exec.Close()
	sc := Query(exec, err, args...)

	return sc
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
