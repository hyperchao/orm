package tag

import (
	"iter"
	"reflect"
)

var (
	_ Value[any]  = (*value[any])(nil)
	_ Values[any] = (*values[any])(nil)
)

type meta[T any] struct {
	Name    string
	Attrs   T
	Indices []int
	Type    reflect.Type
}

type Value[T any] interface {
	Name() string
	Attrs() T
	Type() reflect.Type
	Interface() any
	Addr() any
	Value() reflect.Value
	Set(any any)
}

type Values[T any] interface {
	Contains(name string) bool
	Get(name string) Value[T]
	Iter() iter.Seq2[string, Value[T]]
	Len() int
}

type values[T any] struct {
	metas map[string]*meta[T]
	value reflect.Value
}

func (v *values[T]) Contains(name string) (ok bool) {
	_, ok = v.metas[name]
	return
}

func (v *values[T]) Get(field string) Value[T] {
	meta, ok := v.metas[field]
	if !ok {
		return nil
	}

	val := &value[T]{
		meta:      meta,
		rootValue: v.value,
	}
	return val
}

func (v *values[T]) Iter() iter.Seq2[string, Value[T]] {
	return func(yield func(string, Value[T]) bool) {
		for name := range v.metas {
			value := v.Get(name)
			if !yield(name, value) {
				return
			}
		}
	}
}

func (v *values[T]) Len() int {
	return len(v.metas)
}

type value[T any] struct {
	meta      *meta[T]
	rootValue reflect.Value // root rootValue
}

func (v *value[T]) Name() string {
	return v.meta.Name
}

func (v *value[T]) Attrs() T {
	return v.meta.Attrs
}

func (v *value[T]) Type() reflect.Type {
	return v.meta.Type
}

func (v *value[T]) Interface() any {
	return v.Value().Interface()
}

func (v *value[T]) Addr() any {
	return v.Value().Addr().Interface()
}

func (v *value[T]) Value() reflect.Value {
	return fieldByIndex(v.rootValue, v.meta.Indices)
}

func (v *value[T]) Set(val any) {
	switch v.meta.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.setInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.setUInt(val)
	case reflect.Float32, reflect.Float64:
		v.setFloat(val)
	case reflect.String:
		v.setString(val)
	case reflect.Bool:
		v.setBool(val)
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) setInt(val any) {
	switch val := val.(type) {
	case int:
		v.Value().SetInt(int64(val))
	case int8:
		v.Value().SetInt(int64(val))
	case int16:
		v.Value().SetInt(int64(val))
	case int32:
		v.Value().SetInt(int64(val))
	case int64:
		v.Value().SetInt(val)
	case uint:
		v.Value().SetInt(int64(val))
	case uint8:
		v.Value().SetInt(int64(val))
	case uint16:
		v.Value().SetInt(int64(val))
	case uint32:
		v.Value().SetInt(int64(val))
	case uint64:
		v.Value().SetInt(int64(val))
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) setUInt(val any) {
	switch val := val.(type) {
	case int:
		v.Value().SetUint(uint64(val))
	case int8:
		v.Value().SetUint(uint64(val))
	case int16:
		v.Value().SetUint(uint64(val))
	case int32:
		v.Value().SetUint(uint64(val))
	case int64:
		v.Value().SetUint(uint64(val))
	case uint:
		v.Value().SetUint(uint64(val))
	case uint8:
		v.Value().SetUint(uint64(val))
	case uint16:
		v.Value().SetUint(uint64(val))
	case uint32:
		v.Value().SetUint(uint64(val))
	case uint64:
		v.Value().SetUint(val)
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) setFloat(val any) {
	switch val := val.(type) {
	case float32:
		v.Value().SetFloat(float64(val))
	case float64:
		v.Value().SetFloat(val)
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) setString(val any) {
	switch val := val.(type) {
	case string:
		v.Value().SetString(val)
	case []byte:
		v.Value().SetString(string(val))
	case []rune:
		v.Value().SetString(string(val))
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) setBool(val any) {
	switch val := val.(type) {
	case bool:
		v.Value().SetBool(val)
	default:
		v.fallbackSet(val)
	}
}

func (v *value[T]) fallbackSet(val any) {
	rv := reflect.ValueOf(val)

	if !rv.IsValid() {
		// set nil
		v.Value().SetZero()
		return
	}

	if !rv.Type().AssignableTo(v.meta.Type) {
		if rv.CanConvert(v.meta.Type) {
			rv = rv.Convert(v.meta.Type)
		}
	}
	// always set, let it panic if set failed
	v.Value().Set(rv)
}
