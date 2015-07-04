package gomodel

import "database/sql"

type (
	// Store defines the interface to store data from databqase rows
	Store interface {
		// Init will be called twice, first to allocate initial data space, second to specified
		// the final row count
		// Init initial the data store with size rows, if size is not enough,
		// Realloc will be called
		Init(size int)

		// Final indicate the last found rows
		Final(size int)

		// Ptrs should store pointers of data store at given index to the ptr parameter
		Ptrs(index int, ptrs []interface{})

		// Realloc will occurred if the initial size is not enough, only occured
		// when call the All method of Scanner.
		// The return value shold be the new size of Store.
		// If don't want to continue, just return a non-positive number.
		Realloc(currSize int) (latest int)
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
	defer rows.Close()

	for rows.Next() && (index < count || scanType == _SCAN_ALL) {
		if index == 0 {
			cols, _ := rows.Columns()
			s.Init(count)
			ptrs = make([]interface{}, len(cols))
		}

		if index == count {
			if count = s.Realloc(count); count <= 0 {
				break // don't continue
			}
		}

		s.Ptrs(index, ptrs)

		if err = rows.Scan(ptrs...); err != nil {
			return err
		}
		index++
	}

	if index == 0 {
		err = sql.ErrNoRows
	} else {
		s.Final(index)
	}

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

	rows := sc.Rows
	defer rows.Close()

	var err error
	if rows.Next() {
		err = rows.Scan(ptrs...)
	} else {
		err = sql.ErrNoRows
	}

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

func (s *StringStore) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]string, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *StringStore) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *StringStore) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *StringStore) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]string, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of StringStore")
}

func (s *StringStore) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}

func (s *IntStore) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]int, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *IntStore) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *IntStore) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *IntStore) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]int, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of IntStore")
}

func (s *IntStore) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}
func (s *BoolStore) Init(size int) {
	if cap(s.Values) < size {
		s.Values = make([]bool, size)
	} else {
		s.Values = s.Values[:size]
	}
}

func (s *BoolStore) Final(size int) {
	s.Values = s.Values[:size]
}

func (s *BoolStore) Ptrs(index int, ptrs []interface{}) {
	ptrs[0] = &s.Values[index]
}

func (s *BoolStore) Realloc(count int) int {
	if c := cap(s.Values); c == count {
		values := make([]bool, 2*c)
		copy(values, s.Values)
		s.Values = values

		return 2 * c
	} else if c > count {
		s.Values = s.Values[:c]

		return c
	}

	panic("unexpected capacity of BoolStore")
}

func (s *BoolStore) Clear() {
	if s.Values != nil {
		s.Values = s.Values[:0]
	}
}
