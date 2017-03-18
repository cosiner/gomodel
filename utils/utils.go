package utils

import (
	"fmt"
	"os/user"
	"strings"

	"os"

	"io"

	"unicode"
	"unicode/utf8"
)

func TruncCapToLen(ss []string) []string {
	l, cap := len(ss), cap(ss)
	if l == cap {
		return ss
	}
	ns := make([]string, l)
	copy(ns, ss)
	return ns
}

func ToSnakeCase(s string) string {
	num := len(s)
	need := false // need determin if it's necessery to add a '_'

	snake := make([]byte, 0, len(s)*2)
	for i := 0; i < num; i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
			if need {
				snake = append(snake, '_')
			}
		} else {
			// if previous is '_' or ' ',
			// there is no need to add extra '_' before
			need = (c != '_' && c != ' ')
		}

		snake = append(snake, c)
	}

	return string(snake)
}

func ToLowerAbridgeCase(str string) (s string) {
	if str == "" {
		return ""
	}
	first, _ := utf8.DecodeRuneInString(str)
	if first == utf8.RuneError {
		return ""
	}

	arbi := []rune{unicode.ToLower(first)}
	for _, r := range str {
		if unicode.IsUpper(r) {
			arbi = append(arbi, unicode.ToLower(r))
		}
	}
	return string(arbi)
}

func ConvToInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case uint:
		return int64(v), nil
	}
	return 0, fmt.Errorf("%v(%T) is not an integer", v, v)
}

func MakeStrings(ele string, size int) []string {
	s := make([]string, size)
	if ele != "" {
		for i := 0; i < size; i++ {
			s[i] = ele
		}
	}
	return s
}

func JoinStrings(s []string, suffix, sep string) string {
	if suffix == "" {
		return strings.Join(s, sep)
	}

	l := len(s)
	if l == 0 {
		return ""
	}

	buf := make([]byte, 0, l*2+len(suffix)*l+len(sep)*(l-1))
	for i := 0; i < l; i++ {
		buf = append(buf, s[i]...)
		buf = append(buf, suffix...)
		if i != l-1 {
			buf = append(buf, sep...)
		}
	}
	return string(buf)
}

func TrimQuote(str string) (string, bool) {
	str = strings.TrimSpace(str)
	l := len(str)
	if l == 0 {
		return "", true
	}

	if s, e := str[0], str[l-1]; s == '\'' || s == '"' || s == '`' || e == '\'' || e == '"' || e == '`' {
		if l != 1 && s == e {
			str = str[1 : l-1]
		} else {
			return "", false
		}
	}

	return str, true
}

func RemoveSliceItem(v []interface{}, index int) []interface{} {
	l := len(v)
	if index < 0 || index >= l {
		return v
	}

	if l >= 2*(index+1) {
		for ; index > 0; index-- {
			v[index] = v[index-1]
		}
		v = v[1:]
	} else {
		for ; index < l-1; index++ {
			v[index] = v[index+1]
		}
		v = v[:l-1]
	}

	return v
}
func UnexportedName(n string) string {
	if n == "" {
		return n
	}
	return strings.ToLower(n[:1]) + n[1:]
}

func FatalOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}

func UserHome() string {
	u, _ := user.Current()
	return u.HomeDir
}

func IsFileExist(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
}

func CopyFile(dst, src string) error {
	srcfd, err := os.OpenFile(src, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer srcfd.Close()
	dstfd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstfd.Close()

	_, err = io.Copy(dstfd, srcfd)
	return err
}
