package gomodel

import (
	"database/sql"
	"testing"

	"github.com/cosiner/gohper/lib/errors"
	"github.com/cosiner/gohper/lib/test"
)

func TestErrDuplicateKey(t *testing.T) {
	tt := test.Wrap(t)
	err := errors.Err(`Duplicate entry '14d1b6c34a001-1648e0754a001' for key 'PRIMARY'`) // for combined primary key
	ErrForDuplicateKey(err, func(key string) error {
		tt.Eq(key, "PRIMARY")
		return nil
	})
}

func TestErrForeignKey(t *testing.T) {
	tt := test.Wrap(t)
	err := errors.Err("CONSTRAINT `article_vote_ibfk_1` FOREIGN KEY (`article_id`) REFERENCES `article` (`article_id`)")
	ErrForForeignKey(err, func(key string) error {
		tt.Eq(key, "article_id")
		return nil
	})
}

func TestErrNoRows(t *testing.T) {
	tt := test.Wrap(t)
	err := errors.Err("Error No Rows")

	tt.Eq(err, ErrForNoRows(sql.ErrNoRows, err))
}
