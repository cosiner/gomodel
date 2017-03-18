package driver

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cosiner/gomodel/utils"
)

type MySQL string

func NewMySQL(name string) MySQL {
	return MySQL(name)
}

func (m MySQL) String() string {
	return string(m)
}

func init() {
	Register("mysql", NewMySQL("mysql"))
}

func (MySQL) DSN(host, port, username, password, dbname string, cfg map[string]string) string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)

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

func (MySQL) Prepare(sql string) string {
	return sql
}

func (MySQL) SQLLimit() string {
	return "LIMIT ?, ?"
}

func (MySQL) ParamLimit(offset, count int) (int, int) {
	return offset, count
}

func (MySQL) PrimaryKey() string {
	return "PRIMARY"
}

func (MySQL) DuplicateKey(err error) string {
	if err == nil {
		return ""
	}

	// Duplicate entry ... for key 'keyname'
	const DUPLICATE = "Duplicate"
	const FORKEY = "for key"

	s := err.Error()
	index := strings.Index(s, DUPLICATE)
	if index < 0 {
		return ""
	}

	s = s[index+len(DUPLICATE):]
	index = strings.Index(s, FORKEY) + len(FORKEY)
	if index < 0 {
		return ""
	}

	s, _ = utils.TrimQuote(s[index:])
	return s
}

func (MySQL) ForeignKey(err error) string {
	if err == nil {
		return ""
	}

	// FOREIGN KEY (`keyname`)
	const FOREIGNKEY = "FOREIGN KEY "

	s := err.Error()
	index := strings.Index(s, FOREIGNKEY)
	if index < 0 {
		return ""
	}
	s = s[index+len(FOREIGNKEY):]
	index = strings.IndexByte(s, ')')
	if index < 0 {
		return ""
	}
	s = s[:index+1]

	return strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
}
