package main

import (
	"io/ioutil"
	"text/template"
)

func parseTemplate(fname string) (*template.Template, error) {
	var content string
	if fname != "" {
		c, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, err
		}
		content = string(c)
	} else {
		content = genTemplate
	}
	return template.New("").Parse(content)
}

var genTemplate = `
package gomodel

import (
    "github.com/cosiner/gomodel"
)

{{range $model, $fields := .Models}}
{{$normal := $model.Name}}
{{$unexport := $model.Unexported}}
{{$upper := $model.Upper}}
{{$self := $model.Self}}
{{$recv := (printf "(%s *%s)" $self $normal)}}
const (
    {{range $index, $field := $fields}}{{$upper}}_{{$field.Upper}} {{if eq $index  0}} uint64 = 1 << iota {{end}}
    {{end}}
    {{$normal}}FieldEnd = iota
    {{$normal}}FieldsAll = 1 << {{$normal}}FieldEnd-1
    {{range $index, $field := $fields}}{{$normal}}FieldsExcp{{$field.Name}} = {{$normal}}FieldsAll &(^{{$upper}}_{{$field.Upper}})
    {{end}}

    {{$normal}}Table = "{{$model.Table}}"
    {{range $index, $field := $fields}}{{$normal}}{{$field.Name}}Col = "{{$field.Column}}"
    {{end}}
)

var (
    {{$normal}}Instance = new({{$normal}})
)

func {{$recv}} Table() string {
    return {{$normal}}Table
}

{{if $model.Nocache}}
func {{$recv}} Nocache() bool {
    return {{$model.Nocache}}
}
{{end}}

func {{$recv}} Columns() []string {
    return []string{
    {{range $index, $field:=$fields}}{{$normal}}{{$field.Name}}Col,{{end}}
    }
}

func {{$recv}} Vals(fields uint64, vals []interface{}) {
    if fields != 0 {
    if fields == {{$normal}}FieldsAll {
        {{range $index, $field:=$fields}}vals[{{$index}}]={{$self}}.{{$field.Name}}
        {{end}}
    } else {
       index := 0
    {{range $fields}} if fields&{{$upper}}_{{.Upper}} != 0 {
        vals[index] = {{$self}}.{{.Name}}
        index++
        }
    {{end}}  }
    }
}

func {{$recv}} Ptrs(fields uint64, ptrs []interface{}) {
    if fields != 0 {
        if fields == {{$normal}}FieldsAll {
        {{range $index, $field:=$fields}}ptrs[{{$index}}]=&({{$self}}.{{$field.Name}})
        {{end}}
         } else {
        index := 0
        {{range $fields}} if fields&{{$upper}}_{{.Upper}} != 0 {
            ptrs[index] = &({{$self}}.{{.Name}})
            index++
        }
    {{end}}}
    }
}

func {{$recv}} TxDo(exec gomodel.Executor, do func(*gomodel.Tx, *{{$normal}}) error) error {
    var (
        tx *gomodel.Tx
        err error
    )
    switch r := exec.(type) {
    case *gomodel.Tx:
        tx = r
    case *gomodel.DB:
        tx, err = r.Begin()
    	if err != nil {
    		return err
    	}
    	defer tx.Close()
    default:
        panic("unexpected underlay type of gomodel.Executor")
    }

    err = do(tx, {{$self}})
    tx.Success(err == nil)
    return err
}

type (
    {{$normal}}Store struct {
        Values []{{$normal}}
        Fields uint64
    }
)

func (s *{{$normal}}Store) Init(size int) {
    if cap(s.Values) < size {
        s.Values = make([]{{$normal}}, size)
    } else {
        s.Values = s.Values[:size]
    }
}

func (s *{{$normal}}Store) Final(size int) {
    s.Values = s.Values[:size]
}

func (s *{{$normal}}Store) Ptrs(index int, ptrs []interface{}) {
    s.Values[index].Ptrs(s.Fields, ptrs)
}

func (s *{{$normal}}Store) Realloc(count int) int {
    if c := cap(s.Values); c == count {
        values := make([]{{$normal}}, 2*c)
        copy(values, s.Values)
        s.Values = values

        return 2 * c
    } else if c > count {
        s.Values = s.Values[:c]

        return c
    }

    panic("unexpected capacity of {{$normal}}Store")
}

func (a *{{$normal}}Store) Clear() {
    if a.Values != nil {
        a.Values = a.Values[:0]
    }
}
{{end}}

{{ $length := len .SQLs }} {{ if gt $length 0 }}
// Generated SQL
var (
{{range $name, $sql := .SQLs}}
    {{$name}} = gomodel.NewSqlId(func(gomodel.Executor) string {
        return "{{$sql}}"
    })
{{end}}
)
{{ end }}

`
