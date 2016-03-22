package gomodel

import (
	"strings"

	"github.com/cosiner/gohper/slices"
)

type (
	Cols interface {
		// Cols return colum names
		Names() []string
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

	// MultipleCols used for columns more than two
	MultipleCols struct {
		Cols        []string
		str         string
		paramed     string
		onlyParamed string
	}

	// singleCol means only one column
	SingleCol string

	// emptyCols means there is no columns
	emptyCols string

	ColNum int
)

const _emptyCols = emptyCols("")

func (c *MultipleCols) Names() []string {
	return c.Cols
}

func (c *MultipleCols) String() string {
	if c.str == "" {
		c.str = strings.Join(c.Cols, ",")
	}

	return c.str
}

func (c *MultipleCols) Paramed() string {
	if c.paramed == "" {
		c.paramed = slices.Strings(c.Cols).Join("=?", ",")
	}

	return c.paramed
}

func (c *MultipleCols) OnlyParam() string {
	if c.onlyParamed == "" {
		c.onlyParamed = OnlyParamed(len(c.Cols))
	}

	return c.onlyParamed
}

func (c *MultipleCols) Join(suffix, sep string) string {
	return slices.Strings(c.Cols).Join(suffix, sep)
}

func (c *MultipleCols) Length() int {
	return len(c.Cols)
}

func (c SingleCol) Names() []string {
	return []string{string(c)}
}

func (c SingleCol) String() string {
	return string(c)
}

func (c SingleCol) Paramed() string {
	return string(c) + "=?"
}

func (c SingleCol) OnlyParam() string {
	return "?"
}

func (c SingleCol) Join(suffix, _ string) string {
	return string(c) + suffix
}

func (c SingleCol) Length() int {
	return 1
}

func (emptyCols) Names() []string {
	return nil
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

func OnlyParamed(n int) string {
	switch n {
	case 0:
		return ""
	case 1:
		return "?"
	case 2:
		return "?,?"
	case 3:
		return "?,?,?"
	case 4:
		return "?,?,?,?"
	case 5:
		return "?,?,?,?,?"
	case 6:
		return "?,?,?,?,?,?"
	case 7:
		return "?,?,?,?,?,?,?"
	case 8:
		return "?,?,?,?,?,?,?,?"
	case 9:
		return "?,?,?,?,?,?,?,?,?"
	case 10:
		return "?,?,?,?,?,?,?,?,?,?"
	}
	return slices.MakeStrings("?", int(n)).Join("", ",")
}
