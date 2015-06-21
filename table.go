package gomodel

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/cosiner/gohper/goutil"
	"github.com/cosiner/gohper/reflect2"
	"github.com/cosiner/gohper/strings2"
)

type (
	// SQLBuilder build sql for model using fields and where fields
	SQLBuilder func(fields, whereFields uint) string

	// Table represent information of type
	// contains field count, table name, field names, field offsets
	//
	// because count sql and limit select sql is both stored in the same cache,
	// you should not use a empty fields for limit select, that will conflict with
	// count sql and get the wrong sql statement.
	Table struct {
		Name string
		Num  uint
		Cache

		columns        []string
		prefix         string // Name + "."
		colsCache      map[uint]Cols
		typedColsCache map[uint]Cols
	}
)

const (
	// _FIELD_TAG is tag name of database column
	_FIELD_TAG = "column"
)

// FieldsExcp create fieldset except given fields
func FieldsExcp(numField uint, fields uint) uint {
	return (1<<numField - 1) & (^fields)
}

// FieldsIdentity create signature from fields
func FieldsIdentity(numField uint, fields, whereFields uint) uint {
	return fields<<numField | whereFields
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Stmt(p Preparer, typ, fields, whereFields uint, build SQLBuilder) (*sql.Stmt, error) {
	id := FieldsIdentity(t.Num, fields, whereFields)

	sql_, stmt, err := t.Cache.GetStmt(p, typ, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(fields, whereFields)
		sqlPrinter.Print(false, sql_)

		return t.Cache.SetStmt(p, typ, id, sql_)
	}

	sqlPrinter.Print(true, sql_)

	return stmt, nil
}

func (t *Table) StmtInsert(p Preparer, fields uint) (*sql.Stmt, error) {
	return t.Stmt(p, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) StmtUpdate(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) StmtDelete(p Preparer, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) StmtLimit(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, SELECT_LIMIT, fields, whereFields, t.SQLLimit)
}

func (t *Table) StmtOne(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, SELECT_ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) StmtAll(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, SELECT_ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) StmtCount(p Preparer, whereFields uint) (*sql.Stmt, error) {
	return t.Stmt(p, SELECT_LIMIT, 0, whereFields, t.SQLCount)
}

// Stmt get sql from cache container, if cache not exist, then create new
func (t *Table) Prepare(p Preparer, typ, fields, whereFields uint, build SQLBuilder) (*sql.Stmt, error) {
	id := FieldsIdentity(t.Num, fields, whereFields)

	sql_, stmt, err := t.Cache.PrepareSQL(p, typ, id)
	if err != nil {
		return nil, err
	}

	if stmt == nil {
		sql_ = build(fields, whereFields)
		t.Cache.SetSQL(typ, id, sql_)
		sqlPrinter.Print(false, sql_)

		stmt, err = p.Prepare(sql_)
	} else {
		sqlPrinter.Print(true, sql_)
	}

	return stmt, err
}

func (t *Table) PrepareInsert(p Preparer, fields uint) (*sql.Stmt, error) {
	return t.Prepare(p, INSERT, fields, 0, t.SQLInsert)
}

func (t *Table) PrepareUpdate(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, UPDATE, fields, whereFields, t.SQLUpdate)
}

func (t *Table) PrepareDelete(p Preparer, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, DELETE, 0, whereFields, t.SQLDelete)
}

func (t *Table) PrepareLimit(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, SELECT_LIMIT, fields, whereFields, t.SQLLimit)
}

func (t *Table) PrepareOne(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, SELECT_ONE, fields, whereFields, t.SQLOne)
}

func (t *Table) PrepareAll(p Preparer, fields, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, SELECT_ALL, fields, whereFields, t.SQLAll)
}

func (t *Table) PrepareCount(p Preparer, whereFields uint) (*sql.Stmt, error) {
	return t.Prepare(p, SELECT_LIMIT, 0, whereFields, t.SQLCount)
}

// InsertSQL create insert sql for given fields
func (t *Table) SQLInsert(fields, _ uint) string {
	cols := t.Cols(fields)

	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)",
		t.Name,
		cols.String(),
		cols.OnlyParam())
}

// UpdateSQL create update sql for given fields
func (t *Table) SQLUpdate(fields, whereFields uint) string {
	return fmt.Sprintf("UPDATE %s SET %s %s",
		t.Name,
		t.Cols(fields).Paramed(),
		t.Where(whereFields))
}

// DeleteSQL create delete sql for given fields
func (t *Table) SQLDelete(_, whereFields uint) string {
	return fmt.Sprintf("DELETE FROM %s %s", t.Name, t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLLimit(fields, whereFields uint) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT ?, ?",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// LimitSQL create select sql for given fields
func (t *Table) SQLOne(fields, whereFields uint) string {
	return fmt.Sprintf("SELECT %s FROM %s %s LIMIT 1",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// AllSQL create select sql for given fields
func (t *Table) SQLAll(fields, whereFields uint) string {
	return fmt.Sprintf("SELECT %s FROM %s %s",
		t.Cols(fields),
		t.Name,
		t.Where(whereFields))
}

// SQLForCount create select count sql
func (t *Table) SQLCount(_, whereFields uint) string {
	return fmt.Sprintf("SELECT COUNT(*) FROM %s %s",
		t.Name,
		t.Where(whereFields))
}

// Where create where clause for given fields, the 'WHERE' word is included
func (t *Table) Where(fields uint) string {
	cols := t.Cols(fields)
	if cols.Length() == 0 {
		return ""
	}

	return "WHERE " + cols.Join("=?", " AND ")
}

// TabWhere is slimilar with Where, but prepend a table name for each column
func (t *Table) TabWhere(fields uint) string {
	cols := t.TabCols(fields)
	if cols.Length() != 0 {
		return ""
	}

	return "WHERE " + cols.Join("=?", " AND ")
}

// Cols return column names for given fields
// if fields is only one, return single column
// else return column slice
func (t *Table) Cols(fields uint) Cols {
	cols := t.colsCache[fields]
	if cols == nil {
		cols = t.cols(fields, "")
		t.colsCache[fields] = cols
	}

	return cols
}

// Col return column name of field
func (t *Table) Col(field uint) string {
	return t.Cols(field).String()
}

// TabCols return column names for given fields with type's table name as prefix
// like table.column
func (t *Table) TabCols(fields uint) Cols {
	cols := t.typedColsCache[fields]
	if cols == nil {
		cols = t.cols(fields, t.prefix)
		t.typedColsCache[fields] = cols
	}

	return cols
}

// Col return column name of field
func (t *Table) TabCol(field uint) string {
	return t.TabCols(field).String()
}

// cols get fields names, each field prepend a prefix string
// if fields count is 0, type emptyCols was returned
// if fields count is 1, type singleCols was returned
// otherwise, type cols was returned
func (t *Table) cols(fields uint, prefix string) Cols {
	fieldNames := t.columns
	if colCount := FieldCount(fields); colCount > 1 {
		names := make([]string, colCount)
		var index int
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				names[index] = prefix + fieldNames[i]
				index++
			}
		}

		return &cols{cols: names}
	} else if colCount == 1 {
		for i, l := uint(0), uint(len(fieldNames)); i < l; i++ {
			if (1<<i)&fields != 0 {
				return singleCol(prefix + fieldNames[i])
			}
		}
	}

	return _emptyCols
}

// parse will first use field tag as column name, the tag key is 'column',
// if no tag specified, use field name's camel_case, disable a field by put 'notcol'
// in field tag
func parse(v Model, db *DB) *Table {
	typ := reflect2.IndirectType(v)
	num := typ.NumField()
	cols := make([]string, 0, num)

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

	return &Table{
		Num:   uint(num),
		Name:  v.Table(),
		Cache: NewCache(Types),

		columns:        cols,
		prefix:         v.Table() + ".",
		colsCache:      make(map[uint]Cols),
		typedColsCache: make(map[uint]Cols),
	}
}
