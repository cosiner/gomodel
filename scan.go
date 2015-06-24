package gomodel

import "database/sql"

const DEFAULT_ROW_COUNT = 10

type Store interface {
	// Make will be called twice, first to allocate data space, second to specified
	// the row count
	Make(size int)
	// If index is greater than size of 'Make',
	// Store should allocate new space, return false if don't need remains rows
	Ptrs(index int, ptrs []interface{}) bool
}

type StringStore struct {
	Values []string
}

func _rowCount(c int) int {
	if c >= 0 {
		return c
	}

	return DEFAULT_ROW_COUNT
}

func (s *StringStore) Make(size int) {
	if s.Values == nil {
		s.Values = make([]string, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *StringStore) Ptrs(index int, ptrs []interface{}) bool {
	if len := len(s.Values); index == len {
		values := make([]string, 2*len)
		copy(values, s.Values)
		s.Values = values
	}

	ptrs[0] = &s.Values[index]
	return true
}

type IntStore struct {
	Values []int
}

func (s *IntStore) Make(size int) {
	if s.Values == nil {
		s.Values = make([]int, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *IntStore) Ptrs(index int, ptrs []interface{}) bool {
	if len := len(s.Values); index == len {
		values := make([]int, 2*len)
		copy(values, s.Values)
		s.Values = values
	}

	ptrs[0] = &s.Values[index]
	return true
}

type Scanner struct {
	Error error
}

func (sc Scanner) multiple(rows *sql.Rows, s Store, count int, all bool) error {
	if sc.Error != nil {
		return sc.Error
	}

	var (
		index int
		ptrs  []interface{}
		err   error
	)
	for rows.Next() && (index < count || all) {
		if index == 0 {
			cols, _ := rows.Columns()
			s.Make(count)
			ptrs = make([]interface{}, len(cols))
		}

		if !s.Ptrs(index, ptrs) {
			break
		}

		if err = rows.Scan(ptrs...); err != nil {
			rows.Close()

			return err
		}
		index++
	}

	if index == 0 {
		err = sql.ErrNoRows
	} else {
		s.Make(index)
	}
	rows.Close()

	return err
}

func (sc Scanner) All(rows *sql.Rows, s Store, initsize int) error {
	return sc.multiple(rows, s, _rowCount(initsize), true)
}

// Limit scan rows to scanner
// if rows has elements, scanner's Make method will be called to allocate space,
// the size will be rowCount, and fields count will get from rows.Columns().
// if therre is no rows, sql.ErrNoRows was returned.
//
// it's mostly designed for that customed search
func (sc Scanner) Limit(rows *sql.Rows, s Store, rowCount int) error {
	return sc.multiple(rows, s, _rowCount(rowCount), false)
}

// One scan once then close rows, if no data, sql.ErrNoRows was returned
func (sc Scanner) One(rows *sql.Rows, ptrs ...interface{}) error {
	if sc.Error != nil {
		return sc.Error
	}

	var err error
	if rows.Next() {
		err = rows.Scan(ptrs...)
	} else {
		err = sql.ErrNoRows
	}
	rows.Close()

	return err
}
