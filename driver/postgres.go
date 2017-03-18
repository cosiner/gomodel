package driver

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosiner/gomodel/utils"
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
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "user=%s password=%s host=%s port=%s dbname=%s", username, password, host, port, dbname)
	for k, v := range cfg {
		buf.WriteByte(' ')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	return buf.String()
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
	const PGERR_DUPLICATE string = "23505"
	return p.pgKey(PGERR_DUPLICATE, err)
}

func (p Postgres) ForeignKey(err error) string {
	const PGERR_FOREIGN string = "23503"
	return p.pgKey(PGERR_FOREIGN, err)
}

type PGError interface {
	Get(k byte) (v string)
}

func (p Postgres) pgKey(errCode string, err error) string {
	if err == nil {
		return ""
	}
	e, is := err.(PGError)
	if !is || e.Get('C') != errCode {
		return ""
	}

	// Key (`keyname`)=(`keyvalue`) already exists
	detail := e.Get('D')
	i := strings.IndexByte(detail, '(')
	if i >= 0 {
		detail = detail[i+1:]
		i = strings.IndexByte(detail, ')')
		if i >= 0 {
			detail = detail[:i]
		}
	}
	detail, _ = utils.TrimQuote(detail)

	if strings.Contains(detail, ",") {
		return "PRIMARY"
	}
	return detail
}
