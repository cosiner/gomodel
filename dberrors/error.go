// Package dberrors help processing database errors
package dberrors

import (
	"database/sql"
	"strings"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/strings2"
)

// Only tested for mysql

// PRIMARY_KEY for composite foreign key
const PRIMARY_KEY = "PRIMARY"
const NonError = errors.Err("non error")

type KeyFunc func(key string) error

func DuplicateKey(err error) string {
	if err == nil {
		return ""
	}

	// Duplicate entry ... for key 'keyname'
	const duplicate = "Duplicate"
	const forKey = "for key"

	s := err.Error()
	index := strings.Index(s, duplicate)
	if index >= 0 {
		s = s[index+len(duplicate):]
		if index = strings.Index(s, forKey) + len(forKey); index >= 0 {
			s, _ = strings2.TrimQuote(s[index:])
			return s
		}
	}

	return ""
}

func WrapDuplicateKey(err error, keyfunc KeyFunc) error {
	if key := DuplicateKey(err); key != "" {
		if e := keyfunc(key); e != nil {
			err = e
		} else if e == NonError {
			err = nil
		}
	}
	return err
}

func ForeignKey(err error) string {
	if err == nil {
		return ""
	}

	// FOREIGN KEY (`keyname`)
	const foreign = "FOREIGN KEY "

	s := err.Error()
	index := strings.Index(s, foreign)
	if index > 0 {
		index += len(foreign) + 2
		s = s[index:]
		return s[:strings.IndexByte(s, ')')-1]
	}

	return ""
}

func WrapForeignKey(err error, keyfunc KeyFunc) error {
	if key := ForeignKey(err); key != "" {
		if e := keyfunc(key); e != nil {
			err = e
		} else if e == NonError {
			err = nil
		}
	}
	return err
}

func WrapNoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		return newErr
	} else if err == NonError {
		return nil
	}

	return err
}

func WrapNoAffects(c int64, err, newErr error) error {
	if err == nil && c == 0 {
		return newErr
	} else if err == NonError {
		return nil
	}

	return err
}
