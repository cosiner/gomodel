package driver

import (
	"fmt"
	"strings"

	"github.com/cosiner/gohper/strings2"
)

type MySQL string

func (m MySQL) String() string {
	return string(m)
}

func init() {
	Register("mysql", MySQL("mysql"))
}

func (MySQL) DSN(host, port, username, password, dbname string, cfg map[string]string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		username,
		password,
		host,
		port,
		dbname,
		strings2.JoinPairs(cfg, "=", "&"),
	)
}

func (MySQL) Prepare(sql string) string {
	return sql
}

func (MySQL) SQLLimit() string {
	return "LIMIT ?, ?"
}

func (MySQL) ParamLimit(start, count int64) (int64, int64) {
	return start, count
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

	s, _ = strings2.TrimQuote(s[index:])
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

	s, match := strings2.TrimWrap(s, "(", ")", true)
	if !match {
		return ""
	}
	s, match = strings2.TrimWrap(s, "`", "`", true)
	if !match {
		return ""
	}
	return s
}
