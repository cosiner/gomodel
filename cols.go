package gomodel

import "github.com/cosiner/gohper/slices"

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

	// emptyCols means there is no columns
	emptyCols string
)

const _emptyCols = emptyCols("")

func (c *cols) String() string {
	if c.str == "" {
		c.str = slices.Strings(c.cols).Join("", ",")
	}

	return c.str
}

func (c *cols) Paramed() string {
	if c.paramed == "" {
		c.paramed = slices.Strings(c.cols).Join("=?", ",")
	}

	return c.paramed
}

func (c *cols) OnlyParam() string {
	if c.onlyParamed == "" {
		c.onlyParamed = slices.MakeStrings("?", len(c.cols)).Join("", ",")
	}

	return c.onlyParamed
}

func (c *cols) Join(suffix, sep string) string {
	return slices.Strings(c.cols).Join(suffix, sep)
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

func (emptyCols) String() string {
	return ""
}

func (emptyCols) Paramed() string {
	return ""
}

func (emptyCols) OnlyParam() string {
	return ""
}

func (emptyCols) Join(_, _ string) string {
	return ""
}

func (emptyCols) Length() int {
	return 0
}
