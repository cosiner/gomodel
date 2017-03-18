package dberrs

import (
	"errors"
	"testing"

	"github.com/cosiner/gohper/testing2"
	"github.com/cosiner/gomodel/driver"
)

func TestError(t *testing.T) {
	tt := testing2.Wrap(t)
	mysql := driver.MySQL("mysql")

	err := errors.New(`Duplicate entry '14d1b6c34a001-1648e0754a001' for key 'PRIMARY'`) // for combined primary key
	tt.Eq(mysql.PrimaryKey(), mysql.DuplicateKey(err))

	err = errors.New("CONSTRAINT `article_vote_ibfk_1` FOREIGN KEY (`article_id`) REFERENCES `article` (`article_id`)")
	tt.Eq("article_id", mysql.ForeignKey(err))
}
