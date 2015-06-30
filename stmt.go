package gomodel

import (
	"database/sql"
)

type Stmt interface {
	Exec(...interface{}) (sql.Result, error)
	Query(...interface{}) (*sql.Rows, error)
	QueryRow(...interface{}) *sql.Row
	Close() error
}

type NopCloseStmt struct {
	*sql.Stmt
}

func (s NopCloseStmt) Close() error {
	return nil
}

const (
	STMT_CLOSEABLE = true
	STMT_NOPCLOSE  = false
)

func WrapStmt(closeable bool, stmt *sql.Stmt, err error) (Stmt, error) {
	if err != nil {
		return nil, err
	}

	if closeable {
		if stmt == nil {
			return nil, err
		}

		return stmt, err
	}

	return NopCloseStmt{stmt}, nil
}
