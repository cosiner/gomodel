package dberrors

import (
	"testing"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/testing2"
)

func TestError(t *testing.T) {
	tt := testing2.Wrap(t)
	err := errors.Err(`Duplicate entry '14d1b6c34a001-1648e0754a001' for key 'PRIMARY'`) // for combined primary key
	tt.Eq(PRIMARY_KEY, Error.DuplicateKey(err))

	err = errors.Err("CONSTRAINT `article_vote_ibfk_1` FOREIGN KEY (`article_id`) REFERENCES `article` (`article_id`)")
	tt.Eq("article_id", Error.ForeignKey(err))
}
