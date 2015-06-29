package gomodel

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/reflect2"
	"github.com/cosiner/gohper/slices"
	"github.com/cosiner/gohper/strings2"
)

type (
	// SQLBuilder build sql for model using fields and where fields
	SQLBuilder func(fields, whereFields uint64) string

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

		columns        []string
		prefix         string // Name + "."
		colsCache      map[uint64]Cols
		typedColsCache map[uint64]Cols
	}
)

// FieldsIdentity create signature from fields
func FieldsIdentity(sqlType SQLType, numField, fields, whereFields uint64) uint64 {
	return fields<<numField | whereFields | uint64(sqlType)
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Stmt(prepare Preparer, sqlType SQLType, fields, whereFields uint64, build SQLBuilder) (*sql.Stmt, error) {
	id := FieldsIdentity(sqlType, t.NumFields, fields, whereFields)

	sql_, stmt, err := t.cache.GetStmt(prepare, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(fields, whereFields)
		sqlPrinter.Print(false, sql_)

		return t.cache.SetStmt(prepare, id, sql_)
	}

	sqlPrinter.Print(true, sql_)

	return stmt, nil
}

func (t *Table) StmtInsert(prepare Preparer, fields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) StmtUpdate(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) StmtDelete(prepare Preparer, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) StmtLimit(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, LIMIT, fields, whereFields, t.SQLLimit)
}

func (t *Table) StmtOne(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) StmtAll(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) StmtCount(prepare Preparer, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, LIMIT, 0, whereFields, t.SQLCount)
}

func (t *Table) StmtIncrBy(prepare Preparer, field, whereFields uint64) (*sql.Stmt, error) {
	return t.Stmt(prepare, INCRBY, field, whereFields, t.SQLIncrBy)
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Prepare(prepare Preparer, sqlType SQLType, fields, whereFields uint64, build SQLBuilder) (*sql.Stmt, error) {
	id := FieldsIdentity(sqlType, t.NumFields, fields, whereFields)

	sql_, stmt, err := t.cache.PrepareSQL(prepare, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(fields, whereFields)
		t.cache.SetSQL(id, sql_)
		sqlPrinter.Print(false, sql_)

		stmt, err = prepare.Prepare(sql_)
	} else {
		sqlPrinter.Print(true, sql_)
	}

	return stmt, err
}

func (t *Table) PrepareInsert(prepare Preparer, fields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) PrepareUpdate(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) PrepareDelete(prepare Preparer, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) PrepareLimit(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, LIMIT, fields, whereFields, t.SQLLimit)
}

func (t *Table) PrepareOne(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) PrepareAll(prepare Preparer, fields, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) PrepareCount(prepare Preparer, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, LIMIT, 0, whereFields, t.SQLCount)
}

func (t *Table) PrepareIncrBy(prepare Preparer, field, whereFields uint64) (*sql.Stmt, error) {
	return t.Prepare(prepare, INCRBY, field, whereFields, t.SQLIncrBy)
}

// InsertSQL create insert sql for given fields
func (t *Table) SQLInsert(fields, _ uint64) string {
	cols := t.Cols(fields)

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		t.Name,
		cols.String(),
		cols.OnlyParam())
}

// UpdateSQL create update sql for given fields
func (t *Table) SQLUpdate(fields, whereFields uint64) string {
	return fmt.Sprintf("UPDATE %s SET %s %s",
		t.Name,
		t.Cols(fields).Paramed(),
		t.Where(whereFields))
}

// DeleteSQL create delete sql for given fields
func (t *Table) SQLDelete(_, whereFields uint64) string {
	return fmt.Sprintf("DELETE FROM %s %s", t.Name, t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLLimit(fields, whereFields uint64) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT ?, ?",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLOne(fields, whereFields uint64) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT 1",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// AllSQL create select sql for given fields
func (t *Table) SQLAll(fields, whereFields uint64) string {
	return fmt.Sprintf("SELECT %s FROM %s %s",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// SQLForCount create select count sql
func (t *Table) SQLCount(_, whereFields uint64) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s %s",
		t.Name,
		t.Where(whereFields))
}

// SQLIncrBy create sql for increase/decrease field value
func (t *Table) SQLIncrBy(field, whereFields uint64) string {
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
	cols := t.colsCache[fields]
	if cols == nil {
		cols = t.cols(fields, "")
		t.colsCache[fields] = cols
	}

	return cols
}

// Col return column name of field
func (t *Table) Col(field uint64) string {
	return t.Cols(field).String()
}

// TabCols return column names for given fields with type's table name as prefix
// like table.column
func (t *Table) TabCols(fields uint64) Cols {
	cols := t.typedColsCache[fields]
	if cols == nil {
		cols = t.cols(fields, t.prefix)
		t.typedColsCache[fields] = cols
	}

	return cols
}

// Col return column name of field
func (t *Table) TabCol(field uint64) string {
	return t.TabCols(field).String()
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

const (
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG = "column"
)

// parseModel will first use field tag as column name, the tag key is 'column',
// if no tag specified, use field name's camel_case, disable a field by put 'notcol'
// in field tag
func parseModel(v Model, db *DB) *Table {
	if c, is := v.(Columner); is {
		return newTable(v.Table(), c.Columns())
	}

	typ := reflect2.IndirectType(v)
	num := typ.NumField()

	cols := make([]string, 0)
	for i := 0; i < num; i++ {
		field := typ.Field(i)
		col := field.Name
		// Exported + !(anonymous && structure)
		if goutil.IsExported(col) &&
			!(field.Anonymous &&
				field.Type.Kind() == reflect.Struct) {

			tagName := field.Tag.Get(_FIELD_TAG)
			if tagName == "-" {
				continue
			}
			if tagName != "" {
				col = tagName
			}

			cols = append(cols, strings2.ToSnake(col))
		}
	}

	if len(cols) > MAX_NUMFIELDS {
		panic(fmt.Sprint("can't register model with fields count over ", MAX_NUMFIELDS))
	}

	return newTable(v.Table(), slices.FitCapToLenForString(cols))
}

func newTable(table string, cols []string) *Table {
	return &Table{
		NumFields: uint64(len(cols)),
		Name:      table,
		cache:     newCache(),

		columns:        cols,
		prefix:         table + ".",
		colsCache:      make(map[uint64]Cols),
		typedColsCache: make(map[uint64]Cols),
	}
}
