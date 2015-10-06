package gomodel

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/cosiner/gohper/reflect2"
	"github.com/cosiner/gohper/slices"
	"github.com/cosiner/gohper/strings2"
)

type (
	// SQLBuilder build sql for model using fields and where fields
	SQLBuilder func(driver Driver, fields, whereFields uint64) string

	// Table represent information of type
	// contains field count, table name, field names, field offsets
	//
	// because count sql and limit select sql is both stored in the same cache,
	// you should not use a empty fields for limit select, that will conflict with
	// count sql and get the wrong sql statement.
	Table struct {
		Name      string
		NumFields uint64
		cache     cache

		columns   []string
		prefix    string // Name + "."
		colsCache map[uint64]Cols
	}
)

// FieldsIdentity create signature from fields
func FieldsIdentity(sqlType SQLType, numField, fields, whereFields uint64) uint64 {
	return fields<<numField | whereFields | uint64(sqlType)
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Stmt(exec Executor, sqlType SQLType, fields, whereFields uint64, build SQLBuilder) (Stmt, error) {
	id := FieldsIdentity(sqlType, t.NumFields, fields, whereFields)

	sql_, stmt, err := t.cache.GetStmt(exec, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(exec.Driver(), fields, whereFields)
		sqlPrinter.Print(false, sql_)

		stmt, err =  t.cache.SetStmt(exec, id, sql_)
		if err != nil {
			return nil, err
		}
	}

	sqlPrinter.Print(true, sql_)

	return WrapStmt(STMT_NOPCLOSE, stmt, nil)
}

func (t *Table) StmtInsert(exec Executor, fields uint64) (Stmt, error) {
	return t.Stmt(exec, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) StmtUpdate(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) StmtDelete(exec Executor, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) StmtLimit(exec Executor, fields, whereFields uint64) (Stmt, error) {
	stmt, err := t.Stmt(exec, LIMIT, fields, whereFields, t.SQLLimit)
	return stmt, err
}

func (t *Table) StmtOne(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) StmtAll(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) StmtCount(exec Executor, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, LIMIT, 0, whereFields, t.SQLCount)
}

func (t *Table) StmtIncrBy(exec Executor, field, whereFields uint64) (Stmt, error) {
	return t.Stmt(exec, INCRBY, field, whereFields, t.SQLIncrBy)
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Prepare(exec Executor, sqlType SQLType, fields, whereFields uint64, build SQLBuilder) (Stmt, error) {
	id := FieldsIdentity(sqlType, t.NumFields, fields, whereFields)

	sql_, stmt, err := t.cache.PrepareSQL(exec, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(exec.Driver(), fields, whereFields)
		t.cache.SetSQL(id, sql_)
		sqlPrinter.Print(false, sql_)

		stmt, err = exec.Prepare(sql_)
	} else {
		sqlPrinter.Print(true, sql_)
	}

	return WrapStmt(STMT_CLOSEABLE, stmt, err)
}

func (t *Table) PrepareInsert(exec Executor, fields uint64) (Stmt, error) {
	return t.Prepare(exec, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) PrepareUpdate(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) PrepareDelete(exec Executor, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) PrepareLimit(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, LIMIT, fields, whereFields, t.SQLLimit)
}

func (t *Table) PrepareOne(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) PrepareAll(exec Executor, fields, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) PrepareCount(exec Executor, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, LIMIT, 0, whereFields, t.SQLCount)
}

func (t *Table) PrepareIncrBy(exec Executor, field, whereFields uint64) (Stmt, error) {
	return t.Prepare(exec, INCRBY, field, whereFields, t.SQLIncrBy)
}

// InsertSQL create insert sql for given fields
func (t *Table) SQLInsert(_ Driver, fields, _ uint64) string {
	cols := t.Cols(fields)

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		t.Name,
		cols.String(),
		cols.OnlyParam())
}

// UpdateSQL create update sql for given fields
func (t *Table) SQLUpdate(_ Driver, fields, whereFields uint64) string {
	return fmt.Sprintf("UPDATE %s SET %s %s",
		t.Name,
		t.Cols(fields).Paramed(),
		t.Where(whereFields))
}

// DeleteSQL create delete sql for given fields
func (t *Table) SQLDelete(_ Driver, _, whereFields uint64) string {
	return fmt.Sprintf("DELETE FROM %s %s", t.Name, t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLLimit(driver Driver, fields, whereFields uint64) string {

	return fmt.Sprintf("SELECT %s FROM %s %s "+driver.SQLLimit(),
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLOne(_ Driver, fields, whereFields uint64) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT 1",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// AllSQL create select sql for given fields
func (t *Table) SQLAll(_ Driver, fields, whereFields uint64) string {
	return fmt.Sprintf("SELECT %s FROM %s %s",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// SQLForCount create select count sql
func (t *Table) SQLCount(_ Driver, _, whereFields uint64) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s %s",
		t.Name,
		t.Where(whereFields))
}

// SQLIncrBy create sql for increase/decrease field value
func (t *Table) SQLIncrBy(_ Driver, field, whereFields uint64) string {
	if n := NumFields(field); n != 1 {
		panic("IncrBy only allow update one field, but got " + strconv.Itoa(n))
	}

	col := t.Col(field)

	return fmt.Sprintf("UPDATE %s SET %s=%s+? %s",
		t.Name,
		col,
		col,
		t.Where(whereFields),
	)
}

// Where create where clause for given fields, the 'WHERE' word is included
func (t *Table) Where(fields uint64) string {
	cols := t.Cols(fields)
	if cols.Length() == 0 {
		return ""
	}

	return "WHERE " + cols.Join("=?", " AND ")
}

// TabWhere is slimilar with Where, but prepend a table name for each column
func (t *Table) TabWhere(fields uint64) string {
	cols := t.TabCols(fields)
	if cols.Length() != 0 {
		return ""
	}

	return "WHERE " + cols.Join("=?", " AND ")
}

// Cols return column names for given fields
// if fields is only one, return single column
// else return column slice
func (t *Table) Cols(fields uint64) Cols {
	return t.colsByType(_COLS, fields)
}

// Col return column name of field
func (t *Table) Col(field uint64) string {
	return t.Cols(field).String()
}

// TabCols return column names for given fields with type's table name as prefix
// like table.column
func (t *Table) TabCols(fields uint64) Cols {
	return t.colsByType(_TAB_COLS, fields)
}

// Col return column name of field
func (t *Table) TabCol(field uint64) string {
	return t.TabCols(field).String()
}

const (
	_COLS = iota + 1<<(2*MAX_NUMFIELDS)
	_TAB_COLS
)

func (t *Table) colsByType(typ, fields uint64) Cols {
	cols := t.colsCache[typ|fields]
	if cols == nil {
		cols = t.cols(fields, "")
		t.colsCache[typ|fields] = cols
	}

	return cols
}

// cols get fields names, each field prepend a prefix string
// if fields count is 0, type emptyCols was returned
// if fields count is 1, type singleCols was returned
// otherwise, type cols was returned
func (t *Table) cols(fields uint64, prefix string) Cols {
	fieldNames := t.columns
	if colCount := NumFields(fields); colCount > 1 {
		names := make([]string, colCount)
		var index int
		for i, l := uint64(0), uint64(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				names[index] = prefix + fieldNames[i]
				index++
			}
		}

		return &cols{cols: names}
	} else if colCount == 1 {
		for i, l := uint64(0), uint64(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				return singleCol(prefix + fieldNames[i])
			}
		}
	}

	return _emptyCols
}

// parseModel will first use field tag as column name, the tag key is 'column',
// if no tag specified, use field name's camel_case, disable a field or model
// by set '-' as field tag value
func parseModel(v Model, db *DB) *Table {
	var nocache bool
	if nc, is := v.(Nocacher); is {
		nocache = nc.Nocache()
	}

	if c, is := v.(Columner); is {
		return newTable(v.Table(), c.Columns(), nocache)
	}

	typ := reflect2.IndirectType(v)
	num := typ.NumField()

	cols := make([]string, 0)

	for i := 0; i < num; i++ {
		field := typ.Field(i)
		col := strings2.ToSnake(field.Name)

		b, err := strconv.ParseBool(field.Tag.Get("nocache"))
		if err == nil {
			nocache = b
		}
		if nocache {
			break
		}

		if !field.Anonymous ||
			field.Type.Kind() != reflect.Struct {

			colTag := field.Tag.Get("column")
			if colTag == "-" {
				continue
			}
			if colTag != "" {
				col = colTag
			}

			cols = append(cols, col)
		}
	}

	return newTable(
		v.Table(),
		slices.FitCapToLenString(cols),
		nocache,
	)
}

// newTable create Table for a Model with the table name and columns, if nocache,
// it will not allocate cache memory
func newTable(table string, cols []string, nocache bool) *Table {
	if len(cols) > MAX_NUMFIELDS {
		panic(fmt.Sprint("can't register model with fields count over ", MAX_NUMFIELDS))
	}

	t := &Table{
		NumFields: uint64(len(cols)),
		Name:      table,
	}

	if !nocache {
		t.prefix = table + "."
		t.columns = cols
		t.cache = newCache()
		t.colsCache = make(map[uint64]Cols)
	}

	return t
}
