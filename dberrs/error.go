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
	NonError = errors.Err("non error")
)

type (
	Driverer interface {
		Driver() string
	}

	StringDriver string

	KeyParser func(err error) (key string)

	DriverInfo struct {
		ForeignKeyParser   KeyParser
		DuplicateKeyParser KeyParser
		PrimaryKeyWord     string
	}
)

func (s StringDriver) Driver() string {
	return string(s)
}

var (
	driverInfos = map[string]DriverInfo{
		"mysql": DriverInfo{
			ForeignKeyParser:   mysqlForeignKey,
			DuplicateKeyParser: mysqlDuplicateKey,
			PrimaryKeyWord:     "PRIMARY",
		},
	}
)

func RegisterDriverInfo(driver Driverer, info DriverInfo) {
	if _, has := driverInfos[driver.Driver()]; has {
		panic("driver info for " + driver.Driver() + " already registered")
	}
	if info.ForeignKeyParser == nil || info.DuplicateKeyParser == nil || info.PrimaryKeyWord == "" {
		panic("some info for driver " + driver.Driver() + " is lacked ")
	}

	driverInfos[driver.Driver()] = info
}

func mysqlDuplicateKey(err error) string {
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

func mysqlForeignKey(err error) string {
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

func DuplicateKeyFunc(driver Driverer, err error, keyfunc func(key string) error) error {
	duplicateKey := driverInfos[driver.Driver()].DuplicateKeyParser
	return parseKeyFunc(err, keyfunc, duplicateKey)
}

func DuplicateKeyError(driver Driverer, err error, key string, newErr error) error {
	duplicateKey := driverInfos[driver.Driver()].DuplicateKeyParser
	return parseKeyError(err, key, newErr, duplicateKey)
}

func ForeignKeyFunc(driver Driverer, err error, keyfunc func(key string) error) error {
	foreignKey := driverInfos[driver.Driver()].ForeignKeyParser
	return parseKeyFunc(err, keyfunc, foreignKey)
}

func ForeignKeyError(driver Driverer, err error, key string, newErr error) error {
	foreignKey := driverInfos[driver.Driver()].ForeignKeyParser
	return parseKeyError(err, key, newErr, foreignKey)
}

func DuplicatePrimaryKeyError(driver Driverer, err error, newErr error) error {
	info := driverInfos[driver.Driver()]
	return parseKeyError(err, info.PrimaryKeyWord, newErr, info.DuplicateKeyParser)
}

func PrimaryKey(driver Driverer) string {
	return driverInfos[driver.Driver()].PrimaryKeyWord
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

func HasAffects(c int64, err, newErr error) error {
	if err == nil && c != 0 {
		return newErr
	}

	return err
}
