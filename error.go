package gomodel

import (
	"database/sql"
	"strings"

	"github.com/cosiner/gohper/lib/types"
)

// Only tested for mysql

// PRIMARY_KEY for combined foreign key
const PRIMARY_KEY = "PRIMARY"

func ErrForDuplicateKey(err error, newErrFunc func(key string) error) error {
	// Duplicate entry ... for key 'keyname'
	const duplicate = "Duplicate"
	const forKey = "for key"
	if err != nil {
		s := err.Error()
		index := strings.Index(s, duplicate)
		if index >= 0 {
			s = s[index+len(duplicate):]
			if index = strings.Index(s, forKey) + len(forKey); index >= 0 {
				s, _ = types.TrimQuote(s[index:])
				if e := newErrFunc(s); e != nil {
					err = e
				}
			}
		}
	}
	return err
}

func ErrForForeignKey(err error, newErrFunc func(key string) error) error {
	// FOREIGN KEY (`keyname`)
	const foreign = "FOREIGN KEY "
	if err != nil {
		s := err.Error()
		index := strings.Index(s, foreign)
		if index > 0 {
			index += len(foreign) + 2
			s = s[index:]
			next := strings.IndexByte(s, ')') - 1
			if e := newErrFunc(s[:next]); e != nil {
				err = e
			}
		}
	}
	return err
}

func ErrForNoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		err = newErr
	}
	return err
}
