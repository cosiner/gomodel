package driver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosiner/gohper/strings2"
	"github.com/lib/pq"
)

type Postgres string

func (p Postgres) String() string {
	return string(p)
}

func init() {
	Register("postgres", Postgres("postgres"))
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

func (Postgres) Prepare(sql string) string {
	sql = strings.ToLower(sql)
	replacer := strings.NewReplacer(
		"from dual", "",
		"limit ?, ?", "limit ? offset ?",
		"limit ?,?", "limit ? offset ?",
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

func (Postgres) ParamLimit(start, count int64) (int64, int64) {
	return count, start
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
	if !is {
		return ""
	}

	if e.Code != errCode {
		return ""
	}

	// Key (`keyname`)=(`keyvalue`) already exists
	const KEY = "Key "
	detail := e.Detail
	if strings.HasPrefix(detail, KEY) {
		detail = detail[len(KEY):]
	} else {
		return ""
	}
	i := strings.IndexByte(detail, ')')
	if len(detail) < 3 || detail[0] != '(' || i <= 0 {
		return ""
	}

	key := strings.TrimSpace(detail[1:i])
	if strings.Contains(key, ",") {
		return "PRIMARY"
	}
	return key
}
