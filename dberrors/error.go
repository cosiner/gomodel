// Package dberrors help processing database errors
package dberrors

import (
	"database/sql"
	"strings"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/strings2"
)

// Only tested for mysql
const (
	// PRIMARY_KEY
	PRIMARY_KEY = "PRIMARY"
	NonError    = errors.Err("non error")
)

type (
	KeyError struct {
		Key   string
		Error error
	}
)

func duplicateKey(err error) string {
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

func DuplicateKeyFunc(err error, keyfunc func(key string) error) error {
	if key := duplicateKey(err); key != "" {
		if e := keyfunc(key); e != nil {
			err = e
		} else if e == NonError {
			err = nil
		}
	}
	return err
}

func DuplicateKeyError(err error, keyError KeyError) error {
	if key := duplicateKey(err); key != "" {
		if key == keyError.Key {
			return keyError.Error
		}
		panic("unexpected duplicate key " + keyError.Key)
	}

	return err
}

func foreignKey(err error) string {
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

func ForeignKey(err error, keyfunc func(key string) error) error {
	if key := foreignKey(err); key != "" {
		if e := keyfunc(key); e != nil {
			err = e
		} else if e == NonError {
			err = nil
		}
	}
	return err
}

func NoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		return newErr
	} else if err == NonError {
		return nil
	}

	return err
}

func NoAffects(c int64, err, newErr error) error {
	if err == nil && c == 0 {
		return newErr
	} else if err == NonError {
		return nil
	}

	return err
}
