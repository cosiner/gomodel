package driver

import (
	"bytes"

	"github.com/cosiner/gomodel"
)

type SQLite3 string

func NewSQLite3(name string) gomodel.Driver {
	return SQLite3(name)
}

func (m SQLite3) String() string {
	return string(m)
}

func init() {
	Register("sqlite3", NewMySQL("sqlite3"))
}

func (SQLite3) DSN(_, _, _, _, path string, cfg map[string]string) string {
	var buf bytes.Buffer
	buf.WriteString(path)

	var i int
	for k, v := range cfg {
		if i == 0 {
			buf.WriteByte('?')
		} else {
			buf.WriteByte('&')
		}
		buf.WriteByte(' ')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	return buf.String()
}

func (SQLite3) Prepare(sql string) string {
	return sql
}

func (SQLite3) SQLLimit() string {
	return "LIMIT ?, ?"
}

func (SQLite3) ParamLimit(offset, count int) (int, int) {
	return offset, count
}

func (SQLite3) PrimaryKey() string {
	return ""
}

func (SQLite3) DuplicateKey(err error) string {
	return ""
}

func (SQLite3) ForeignKey(err error) string {
	return ""
}
