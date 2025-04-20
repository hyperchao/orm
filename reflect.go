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

// extract a slice of interfaces from struct for sql.Rows.Scan to use
func getColumnPlaceholder(conf *config, v reflect.Value, columns []string) []any {
	if len(columns) == 0 {
		return nil
	}

	r := make([]interface{}, len(columns))
	mapping := getColumnIndexMapping(conf, v)
	for i, column := range columns {
		indexes := mapping[column]
		if len(indexes) > 0 {
			f := v.FieldByIndex(indexes)
			r[i] = f.Addr().Interface()
		} else {
			r[i] = empty{}
		}
	}

	return r
}

func getColumnIndexMapping(conf *config, t reflect.Value) map[string][]int {
	ret, ok := cache.Load(t)
	if ok {
		return ret.(map[string][]int)
	}
	mapping := make(map[string][]int)
	walk(conf, t, nil, mapping)
	cache.Store(t, mapping)
	return mapping
}

func walk(conf *config, v reflect.Value, path []int, mapping map[string][]int) {
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Type().Field(i)
		if f.PkgPath != "" {
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
			walk(conf, indirect(v.Field(i)), append(path, i), mapping)
		}
	}
}

func isSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}
