package gomodel

import "log"

type SQLPrinter func(string, ...interface{})

var sqlPrinter SQLPrinter = func(string, ...interface{}) {}

func (p SQLPrinter) Print(fromcache bool, sql string) {
	p("Cached: %t, SQL: %s\n", fromcache, sql)
}

// SQLPrint enable sql print for each operation
func SQLPrint(enable bool, printer func(formart string, v ...interface{})) {
	if !enable {
		return
	}

	sqlPrinter = printer
	if sqlPrinter == nil {
		sqlPrinter = log.Printf
	}
}
