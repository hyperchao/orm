package orm

import (
	"reflect"
	"strings"
	"sync"
)

var cache sync.Map

type empty struct{}

func (empty) Scan(any) error {
	return nil
}

func isStructOrIndirectToStruct(r reflect.Type) bool {
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	return r.Kind() == reflect.Struct
}

// indirect returns the item at the end of pointer indirection.
// when meet a nil pointer, create a zero value for it to point to
func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}
	return v
}

func indirectT(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// extract a slice of interfaces from struct for sql.Rows.Scan to use
func getColumnPlaceholder(conf *config, v reflect.Value, columns []string) []any {
	if len(columns) == 0 {
		return nil
	}

	r := make([]interface{}, len(columns))
	mapping := getColumnIndexMapping(conf, v.Type())
	for i, column := range columns {
		indexes := mapping[column]
		if len(indexes) > 0 {
			f := fieldByIndex(v, indexes)
			r[i] = f.Addr().Interface()
		} else {
			r[i] = empty{}
		}
	}

	return r
}

func getColumnIndexMapping(conf *config, t reflect.Type) map[string][]int {
	t = indirectT(t)
	ret, ok := cache.Load(t)
	if ok {
		return ret.(map[string][]int)
	}
	mapping := make(map[string][]int)
	visitedTypes := make(map[reflect.Type]struct{})
	walk(conf, t, nil, mapping, visitedTypes)
	cache.Store(t, mapping)
	return mapping
}

func walk(conf *config, t reflect.Type, path []int, mapping map[string][]int, visitedTypes map[reflect.Type]struct{}) {
	if t.Kind() != reflect.Struct {
		return
	}
	if _, visited := visitedTypes[t]; visited {
		return
	}
	visitedTypes[t] = struct{}{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			// skip unexported field
			continue
		}
		fieldName := f.Tag.Get(conf.tagName)
		fieldName = strings.TrimSpace(fieldName)
		if fieldName != "" {
			_, ok := mapping[fieldName]
			if !ok {
				indexes := make([]int, len(path), len(path)+1)
				copy(indexes, path)
				indexes = append(indexes, i)
				mapping[fieldName] = indexes
			}
		} else if isStructOrIndirectToStruct(f.Type) {
			walk(conf, indirectT(f.Type), append(path, i), mapping, visitedTypes)
		}
	}
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	if len(index) == 0 {
		return v
	}
	v = indirect(v)
	return fieldByIndex(v.Field(index[0]), index[1:])
}
