package tag

import (
	"reflect"
	"strings"
	"sync"
)

type ParseFunc[T any] func(string) (string, T)

type Parser[T any] struct {
	cache sync.Map
	fun   ParseFunc[T]
}

func NewParser[T any](f ParseFunc[T]) *Parser[T] {
	return &Parser[T]{
		fun:   f,
		cache: sync.Map{},
	}
}

func (p *Parser[T]) Parse(tagName string, val any) Values[T] {
	var metas map[string]*meta[T]

	rt := indirectT(reflect.TypeOf(val))
	ret, ok := p.cache.Load(rt)
	if ok {
		metas = ret.(map[string]*meta[T])
	} else {
		metas = make(map[string]*meta[T])
		p.traverse(tagName, rt, nil, metas, make(map[reflect.Type]struct{}))
		p.cache.Store(rt, metas)
	}
	return &values[T]{
		metas: metas,
		value: reflect.ValueOf(val),
	}
}

func (p *Parser[T]) traverse(
	tagName string,
	rt reflect.Type,
	path []int,
	metas map[string]*meta[T],
	visitedTypes map[reflect.Type]struct{}) {

	if rt.Kind() != reflect.Struct {
		return
	}
	if _, visited := visitedTypes[rt]; visited {
		return
	}
	visitedTypes[rt] = struct{}{}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			// skip unexported field
			continue
		}

		tagVal := field.Tag.Get(tagName)
		tagVal = strings.TrimSpace(tagVal)
		if tagVal != "" {
			name, attrs := p.fun(tagVal)
			indices := make([]int, len(path), len(path)+1)
			copy(indices, path)
			indices = append(indices, i)
			metas[name] = &meta[T]{
				Name:    name,
				Attrs:   attrs,
				Indices: indices,
				Type:    field.Type,
			}
		} else if isStructOrIndirectToStruct(field.Type) {
			p.traverse(tagName, indirectT(field.Type), append(path, i), metas, visitedTypes)
		}
	}
}

func isStructOrIndirectToStruct(r reflect.Type) bool {
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	return r.Kind() == reflect.Struct
}

func indirectT(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
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

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	if len(index) == 0 {
		return v
	}
	v = indirect(v)
	return fieldByIndex(v.Field(index[0]), index[1:])
}
