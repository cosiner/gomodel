package gomodel

import (
	"strings"

	"github.com/cosiner/gohper/strings2"
)

// Only tested for mysql

// PRIMARY_KEY for composite foreign key
const PRIMARY_KEY = "PRIMARY"

type err struct{}

var Error = err{}

func (err) DuplicateKey(err error) string {
	if err == nil {
		return ""
	}

	// Duplicate entry ... for key 'keyname'
	const duplicate = "Duplicate"
	const forKey = "for key"

	s := err.Error()
	index := strings.Index(s, duplicate)
	if index >= 0 {
		s = s[index+len(duplicate):]
		if index = strings.Index(s, forKey) + len(forKey); index >= 0 {
			s, _ = strings2.TrimQuote(s[index:])
			return s
		}
	}

	return ""
}

func (err) ForeignKey(err error) string {
	if err == nil {
		return ""
	}

	// FOREIGN KEY (`keyname`)
	const foreign = "FOREIGN KEY "

	s := err.Error()
	index := strings.Index(s, foreign)
	if index > 0 {
		index += len(foreign) + 2
		s = s[index:]
		return s[:strings.IndexByte(s, ')')-1]
	}

	return ""
}
