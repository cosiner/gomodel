package dbutils

import (
	"fmt"
	"reflect"

	"github.com/cosiner/gomodel"
	"github.com/cosiner/gomodel/dberrs"
	"github.com/cosiner/gomodel/utils"
)

func QueryOneResultById(exec gomodel.Executor, sqlid uint64, ptr interface{}, args ...interface{}) error {
	sc := exec.QueryById(sqlid, args...)
	defer sc.Close()
	return sc.One(ptr)
}

func QueryCountById(exec gomodel.Executor, sqlid uint64, args ...interface{}) (int64, error) {
	var count int64
	err := QueryOneResultById(exec, sqlid, &count, args...)
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
			args = utils.RemoveSliceItem(args, index)
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

func CheckFieldForIncrBy(field, fields uint64, count int64) error {
	if gomodel.NumFields(field) == 0 || field&fields == 0 {
		return fmt.Errorf("unexpected field type %d", field)
	}

	if count != -1 && count != 1 {
		return fmt.Errorf("unexpected field incrby count %d, must be -1 or 1", count)
	}
}

type IncrByFunc func(exec gomodel.Executor, field uint64, whereArgs ...interface{}) error

func FuncForIncrByFieldCount(defaultRunner gomodel.Executor, model gomodel.Model, fields, whereFields uint64, noAffectsError error) IncrByFunc {
	if gomodel.NumFields(whereFields) == 0 || gomodel.NumFields(fields) == 0 {
		return fmt.Errorf("unexpected field count of fields %d and whereField %d", fields, whereFields)
	}

	return func(exec gomodel.Executor, field uint64, whereArgs ...interface{}) error {
		var count int64
		switch arg := whereArgs[0].(type) {
		case int:
			count = int64(arg)
		case int8:
			count = int64(arg)
		case int16:
			count = int64(arg)
		case int32:
			count = int64(arg)
		case int64:
			count = int64(arg)
		case uint8:
			count = int64(arg)
		case uint16:
			count = int64(arg)
		case uint32:
			count = int64(arg)
		case uint64:
			count = int64(arg)
		default:
			return fmt.Errorf("count %v must be an integer", arg)
		}
		err := CheckFieldForIncrBy(field, fields, count)
		if err != nil {
			return err
		}

		if exec == nil {
			exec = defaultRunner
		}
		c, err := exec.ArgsIncrBy(model, field, whereFields, whereArgs...)
		return dberrs.NoAffects(c, err, noAffectsError)
	}
}
