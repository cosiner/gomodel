// Package dberrs help processing database errors
package dberrs

import (
	"database/sql"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gomodel"
)

// Only tested for mysql
const (
	NonError = errors.Err("non error")
)

type (
	KeyParser func(err error) (key string)
)

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

	if k == key || key == "" {
		if newErr == NonError {
			return nil
		}
		return newErr
	}

	panic("unexpected key: " + k + ", expect: " + key)
}

func DuplicateKeyFunc(exec gomodel.Executor, err error, keyfunc func(key string) error) error {
	dk := exec.Driver().DuplicateKey
	return parseKeyFunc(err, keyfunc, dk)
}

func DuplicateKeyError(exec gomodel.Executor, err error, key string, newErr error) error {
	dk := exec.Driver().DuplicateKey
	return parseKeyError(err, key, newErr, dk)
}

func ForeignKeyFunc(exec gomodel.Executor, err error, keyfunc func(key string) error) error {
	fk := exec.Driver().ForeignKey
	return parseKeyFunc(err, keyfunc, fk)
}

func ForeignKeyError(exec gomodel.Executor, err error, key string, newErr error) error {
	fk := exec.Driver().ForeignKey
	return parseKeyError(err, key, newErr, fk)
}

func DuplicatePrimaryKeyError(exec gomodel.Executor, err error, newErr error) error {
	pk := exec.Driver().PrimaryKey()
	return parseKeyError(err, pk, newErr, exec.Driver().DuplicateKey)
}

func PrimaryKey(exec gomodel.Executor) string {
	return exec.Driver().PrimaryKey()
}

func AllowNoRows(err error) error {
	return NoRows(err, nil)
}

func NoRows(err, newErr error) error {
	if err == sql.ErrNoRows {
		return newErr
	}

	return err
}

func AllowNonExists(exist bool, err error) error {
	return NonExists(exist, err, nil)
}

func NonExists(exist bool, err, newErr error) error {
	if err == nil && !exist {
		return newErr
	}

	return err
}

func AllowNoAffects(c int64, err error) error {
	return NoAffects(c, err, nil)
}

func NoAffects(c int64, err, newErr error) error {
	if err == nil && c == 0 {
		return newErr
	}

	return err
}

func AllowHasEffects(c int64, err error) error {
	return HasAffects(c, err, nil)
}

func HasAffects(c int64, err, newErr error) error {
	if err == nil && c != 0 {
		return newErr
	}

	return err
}
