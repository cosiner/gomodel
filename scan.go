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
