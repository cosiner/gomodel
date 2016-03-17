package driver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/strings2"
	"github.com/lib/pq"
)

type Postgres string

func NewPostgres(name string) Postgres {
	return Postgres(name)
}

func (p Postgres) String() string {
	return string(p)
}

func init() {
	Register("postgres", NewPostgres("postgres"))
}

func (Postgres) DSN(host, port, username, password, dbname string, cfg map[string]string) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s",
		username,
		password,
		host,
		port,
		dbname,
		strings2.JoinPairs(cfg, "=", " "),
	)
}

func (p Postgres) Prepare(sql string) string {
	sql = strings.ToLower(sql)
	l := p.SQLLimit()
	replacer := strings.NewReplacer(
		"from dual", "",
		"FROM DUAL", "",
		"limit ?, ?", l,
		"LIMIT ?, ?", l,
		"limit ?,?", l,
		"LIMIT ?,?", l,
	)
	sql = replacer.Replace(sql)

	buf := make([]byte, 0, len(sql))
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
	return string(buf)
}

func (Postgres) SQLLimit() string {
	return "LIMIT ? OFFSET ?"
}

func (Postgres) ParamLimit(offset, count int) (int, int) {
	return count, offset
}

func (Postgres) PrimaryKey() string {
	return "PRIMARY"
}

func (p Postgres) DuplicateKey(err error) string {
	const PGERR_DUPLICATE pq.ErrorCode = "23505"
	return p.pgKey(PGERR_DUPLICATE, err)
}

func (p Postgres) ForeignKey(err error) string {
	const PGERR_FOREIGN pq.ErrorCode = "23503"
	return p.pgKey(PGERR_FOREIGN, err)
}

func (p Postgres) pgKey(errCode pq.ErrorCode, err error) string {
	if err == nil {
		return ""
	}
	e, is := err.(*pq.Error)
	if !is || e.Code != errCode {
		return ""
	}

	if e.Code != errCode {
		return ""
	}

	// Key (`keyname`)=(`keyvalue`) already exists
	detail := e.Detail
	i := strings.IndexByte(detail, '(')
	if i >= 0 {
		detail = detail[i+1:]
		i = strings.IndexByte(detail, ')')
		if i >= 0 {
			detail = detail[:i]
		}
	}
	detail, _ = strings2.TrimQuote(detail)

	if strings.Contains(detail, ",") {
		return "PRIMARY"
	}
	return detail
}
