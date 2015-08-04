package dbutils

import (
	"reflect"

	"github.com/cosiner/gohper/errors"
	"github.com/cosiner/gohper/slices"
	"github.com/cosiner/gomodel"
	"github.com/cosiner/gomodel/dberrs"
)

func QueryCountById(runner gomodel.CommonRunner, sqlid uint64, args ...interface{}) (int64, error) {
	var count int64
	sc := runner.QueryById(sqlid, args...)
	err := sc.One(&count)
	return count, err
}

type CondOption struct {
	OtherFields    uint64
	CondFields     []uint64
	CondArgIndexes []int
	Args           []interface{}
}

func (opt CondOption) CondArgs() (uint64, []interface{}) {
	var (
		skipped     int
		otherFields = opt.OtherFields
		args        = opt.Args
		indexes     = opt.CondArgIndexes
	)

	for i, condField := range opt.CondFields {
		var usable bool
		index := indexes[i] - skipped
		switch val := opt.Args[index].(type) {
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
			otherFields |= condField
		} else {
			args = slices.RemoveElement(args, index)
			skipped++
		}
	}
	return otherFields, args
}

func CondArgs(otherFields uint64, condFields []uint64, condArgIndex []int, args []interface{}) (uint64, []interface{}) {
	opt := (CondOption{
		OtherFields:    otherFields,
		CondFields:     condFields,
		CondArgIndexes: condArgIndex,
		Args:           args,
	})
	return opt.CondArgs()
}

func CheckFieldForIncrBy(field, fields uint64, count int64) {
	if gomodel.NumFields(field) == 0 || field&fields == 0 {
		panic(errors.Newf("unexpected field type %d", field))
	}

	if count != -1 && count != 1 {
		panic(errors.Newf("unexpected field incrby count %d, must be -1 or 1", count))
	}
}

type IncrByFunc func(runner gomodel.CommonRunner, field uint64, count int64, whereArg string) error

func FuncForIncrByFieldCount(defaultRunner gomodel.CommonRunner, model gomodel.Model, fields, whereField uint64, noAffectsError error) IncrByFunc {
	if gomodel.NumFields(whereField) != 1 || gomodel.NumFields(fields) == 0 {
		panic(errors.Newf("unexpected field count of fields %d and whereField %d", fields, whereField))
	}

	return func(runner gomodel.CommonRunner, field uint64, count int64, whereArg string) error {
		CheckFieldForIncrBy(field, fields, count)
		if runner == nil {
			runner = defaultRunner
		}
		c, err := runner.ArgsIncrBy(model, field, whereField, count, whereArg)
		return dberrs.NoAffects(c, err, noAffectsError)
	}

}
