package gomodel

import "database/sql"

type (
	// Store defines the interface to store data from databqase rows
	Store interface {
		// Make will be called twice, first to allocate initial data space, second to specified
		// the final row count
		Make(size int)
		// If index is greater than initial size of 'Make'(only occured in
		// Scanner.All), Store should allocate new space,
		// return false if want to stop scanning
		Ptrs(index int, ptrs []interface{}) bool
	}

	// Scanner scan database rows to data Store when Error is nil, if the Rows is
	// empty, sql.ErrNoRows was returned, the Rows will always be be Closed
	Scanner struct {
		Error error
		Rows  *sql.Rows
	}
)

func _rowCount(c int) int {
	const DEFAULT_ROW_COUNT = 10

	if c >= 0 {
		return c
	}

	return DEFAULT_ROW_COUNT
}

const (
	_SCAN_ALL   = true
	_SCAN_LIMIT = !_SCAN_ALL
)

func (sc Scanner) multiple(s Store, count int, scanType bool) error {
	if sc.Error != nil {
		return sc.Error
	}

	var (
		index int
		ptrs  []interface{}
		err   error
	)

	rows := sc.Rows
	for rows.Next() && (index < count || scanType == _SCAN_ALL) {
		if index == 0 {
			cols, _ := rows.Columns()
			s.Make(count)
			ptrs = make([]interface{}, len(cols))
		}

		if !s.Ptrs(index, ptrs) {
			break
		}

		if err = rows.Scan(ptrs...); err != nil {
			_ = rows.Close()

			return err
		}
		index++
	}

	if index == 0 {
		err = sql.ErrNoRows
	} else {
		s.Make(index)
	}
	_ = rows.Close()

	return err
}

func (sc Scanner) All(s Store, initsize int) error {
	return sc.multiple(s, _rowCount(initsize), _SCAN_ALL)
}

func (sc Scanner) Limit(s Store, rowCount int) error {
	return sc.multiple(s, _rowCount(rowCount), _SCAN_LIMIT)
}

func (sc Scanner) One(ptrs ...interface{}) error {
	if sc.Error != nil {
		return sc.Error
	}

	var err error

	rows := sc.Rows
	if rows.Next() {
		err = rows.Scan(ptrs...)
	} else {
		err = sql.ErrNoRows
	}
	_ = rows.Close()

	return err
}

type (
	StringStore struct {
		Values []string
	}

	IntStore struct {
		Values []int
	}

	BoolStore struct {
		Values []bool
	}
)

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

func (s *BoolStore) Make(size int) {
	if s.Values == nil {
		s.Values = make([]bool, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *BoolStore) Ptrs(index int, ptrs []interface{}) bool {
	if len := len(s.Values); index == len {
		values := make([]bool, 2*len)
		copy(values, s.Values)
		s.Values = values
	}

	ptrs[0] = &s.Values[index]
	return true
}
