package driver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/bytes2"
	"github.com/cosiner/gohper/strings2"
)

type Postgres string

func (p Postgres) String() string {
	return string(p)
}

func (Postgres) DSN(host, port, username, password, dbname string, cfg map[string]string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s",
		username,
		password,
		host,
		port,
		dbname,
		strings2.JoinKVs(cfg, "=", " "),
	)
}

func (Postgres) Prepare(pool bytes2.Pool, sql string) string {
	buf := pool.Get(len(sql), false)

	index := 1
	for i, l := 0, len(sql); i < l; i++ {
		if c := sql[i]; c != '?' {
			buf = append(buf, c)
		} else {
			buf = append(buf, '$')
			buf = append(buf, strconv.Itoa(index)...)
			index++
		}
	}
	sql = string(buf)
	pool.Put(buf)
	return sql
}

func (Postgres) SQLLimit() string {
	return "LIMIT ? OFFSET ?"
}

func (Postgres) ParamLimit(start, count int64) (int64, int64) {
	return count, start
}

func (Postgres) PrimaryKey() string {
	return "PRIMARY"
}

func (Postgres) DuplicateKey(err error) string {
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

func (Postgres) ForeignKey(err error) string {
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
