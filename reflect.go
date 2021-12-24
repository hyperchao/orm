package orm

import (
	"reflect"
	"strings"
	"sync"
)

const tagName = "orm"

var cache sync.Map

// extract a slice of interfaces from struct for sql.Rows.Scan to use
func getColumnPlaceholder(v reflect.Value, columns []string) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	r := make([]interface{}, len(columns))
	mapping := getColumnIndexMapping(v.Type())
	for i, cloumn := range columns {
		indexes := mapping[cloumn]
		if len(indexes) > 0 {
			f := v.FieldByIndex(indexes)
			r[i] = f.Addr().Interface()
		} else {
			r[i] = empty{}
		}
	}

	return r
}

func getColumnIndexMapping(t reflect.Type) map[string][]int {
	ret, ok := cache.Load(t)
	if ok {
		return ret.(map[string][]int)
	}
	mapping := make(map[string][]int)
	walk(t, nil, mapping)
	cache.Store(t, mapping)
	return mapping
}

func walk(t reflect.Type, path []int, mapping map[string][]int) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			// skip unexported field
			continue
		}
		fieldName := f.Tag.Get(tagName)
		fieldName = strings.TrimSpace(fieldName)
		if fieldName != "" {
			_, ok := mapping[fieldName]
			if !ok {
				indexes := make([]int, len(path), len(path)+1)
				copy(indexes, path)
				indexes = append(indexes, i)
				mapping[fieldName] = indexes
			}
		}
		if f.Type.Kind() == reflect.Struct {
			walk(f.Type, append(path, i), mapping)
		}
	}
}

// func assertKind(t reflect.Type, k reflect.Kind) {
// if t.Kind() != k {
// panic("expect kind: %d")
// }
// }
