// Package dberrs help processing database errors
package dberrs

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

type KeyParser func(err error) (key string)

var duplicateKey = func(err error) string {
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

var foreignKey = func(err error) string {
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

	index += len(FOREIGNKEY) + 2
	s = s[index:]
	return s[:strings.IndexByte(s, ')')-1]
}

// SetupKeyParsers change the default key parser if new parser is non-nill
func SetupKeyParsers(duplicate, foreign KeyParser) {
	if duplicate != nil {
		duplicateKey = duplicate
	}

	if foreign != nil {
		foreignKey = foreign
	}
}

func parseKeyFunc(err error, keyfunc func(key string) error, getKey KeyParser) error {
	key := getKey(err)
	if key == "" {
		return err
	}

	if e := keyfunc(key); e == NonError {
		err = nil
	} else if e != nil {
		err = e
	}

	return err
}

func parseKeyError(err error, key string, newErr error, getKey KeyParser) error {
	k := getKey(err)
	if k == "" {
		return err
	}

	if k == key {
		return newErr
	}

	panic("unexpected key: " + k + ", expect: " + key)
}

func DuplicateKeyFunc(err error, keyfunc func(key string) error) error {
	return parseKeyFunc(err, keyfunc, duplicateKey)
}

func DuplicateKeyError(err error, key string, newErr error) error {
	return parseKeyError(err, key, newErr, duplicateKey)
}

func ForeignKeyFunc(err error, keyfunc func(key string) error) error {
	return parseKeyFunc(err, keyfunc, foreignKey)
}

func ForeignKeyError(err error, key string, newErr error) error {
	return parseKeyError(err, key, newErr, foreignKey)
}

func NoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		return newErr
	}

	return err
}

func NoAffects(c int64, err, newErr error) error {
	if err == nil && c == 0 {
		return newErr
	}

	return err
}
