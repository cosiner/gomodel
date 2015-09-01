package gomodel

import (
	"github.com/cosiner/gohper/bytes2"
)

type Driver interface {
	String() string
	DSN(host, port, username, password, dbname string, cfg map[string]string) string
	// Prepare should replace the standard placeholder '?' with driver specific placeholder,
	// for postgresql it's '$n'
	Prepare(buf bytes2.Pool, sql string) string
	SQLLimit() string
	ParamLimit(start, count int64) (int64, int64)
	PrimaryKey() string
	DuplicateKey(err error) string
	ForeignKey(err error) string
}
