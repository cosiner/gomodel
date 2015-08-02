package dbutils

import (
	"reflect"

	"github.com/cosiner/gohper/slices"
	"github.com/cosiner/gomodel"
)

func QueryCountById(runner gomodel.CommonRunner, sqlid uint64, args ...interface{}) (int64, error) {
	var count int64
	sc := runner.QueryById(sqlid, args...)
	err := sc.One(&count)
	return count, err
}

func CondArgs(whereFields, condField, index uint64, args ...interface{}) (uint64, []interface{}) {
	var usable bool
	switch val := args[index].(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		usable = val != 0
	case string:
		usable = val != ""
	case bool:
		usable = val
	default:
		panic("unsupported type: " + reflect.TypeOf(args[index]).String())
	}

	if usable {
		whereFields |= condField
	} else {
		args = slices.RemoveElement(args, int(index))
	}
	return whereFields, args
}
