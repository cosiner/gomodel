package driver

import (
	"github.com/cosiner/gomodel"
)

var drivers = make(map[string]gomodel.Driver)

func Register(name string, driver gomodel.Driver) bool {
	_, has := drivers[name]
	if !has {
		drivers[name] = driver
	}
	return !has
}

func init() {
	Register("mysql", MySQL("mysql"))
	Register("postgres", Postgres("postgres"))
}

func Get(name string) gomodel.Driver {
	return drivers[name]
}
