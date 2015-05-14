package gomodel

import (
	"github.com/cosiner/gohper/strings2"
)

type (
	Cols interface {
		// String return columns string join with ",",
		// result like "foo, bar"
		String() string

		// Paramed return columns string joind with "=?,", last "," was trimed,
		// result like "foo=?, bar=?"
		Paramed() string

		// OnlyParam return columns placeholdered string,
		// each column was replaced with "?"
		// result like "?, ?, ?, ?", count of "?" is colums length
		OnlyParam() string

		// Join append suffix string to each columns then join them with the seperator
		Join(suffix, sep string) string

		// Length return columns count
		Length() int
	}

	// cols used for columns more than two
	cols struct {
		cols        []string
		str         string
		paramed     string
		onlyParamed string
	}

	// singleCol means only one column
	singleCol string

	// nilCOls means there is no columns
	nilCols string
)

var (
	zeroCols Cols = nilCols("")
)

func (c *cols) String() string {
	if c.str == "" {
		c.str = strings2.SuffixJoin(c.cols, "", ",")
	}

	return c.str
}

func (c *cols) Paramed() string {
	if c.paramed == "" {
		c.paramed = strings2.SuffixJoin(c.cols, "=?", ",")
	}

	return c.paramed
}

func (c *cols) OnlyParam() string {
	if c.onlyParamed == "" {
		c.onlyParamed = strings2.RepeatJoin("?", ",", len(c.cols))
	}

	return c.onlyParamed
}

func (c *cols) Join(suffix, sep string) string {
	return strings2.SuffixJoin(c.cols, suffix, sep)
}

func (c *cols) Length() int {
	return len(c.cols)
}

func (c singleCol) String() string {
	return string(c)
}

func (c singleCol) Paramed() string {
	return string(c) + "=?"
}

func (c singleCol) OnlyParam() string {
	return "?"
}

func (c singleCol) Join(suffix, _ string) string {
	return string(c) + suffix
}

func (c singleCol) Length() int {
	return 1
}

func (nilCols) String() string {
	return ""
}

func (nilCols) Paramed() string {
	return ""
}

func (nilCols) OnlyParam() string {
	return ""
}

func (nilCols) Join(_, _ string) string {
	return ""
}

func (nilCols) Length() int {
	return 0
}
