package conv

import (
	"reflect"
)

func indirect(v reflect.Value) reflect.Value {
	for {
		if v.Kind() != reflect.Interface && v.Kind() != reflect.Ptr {
			return v
		}
		v = v.Elem()
	}
}

func mapIndex(m, key reflect.Value) reflect.Value {
	val := m.MapIndex(key)
	if val.Kind() != reflect.Invalid {
		return val
	}
	val = reflect.New(m.Type().Elem())
	return val.Elem()
}

func sliceIndex(slice reflect.Value, i int) (elem reflect.Value) {
	if i >= slice.Len() {
		elem = reflect.New(slice.Type().Elem())
		newSlice := reflect.Append(slice, slice.Elem())
		if slice.UnsafeAddr() != newSlice.UnsafeAddr() {
			slice.Set(newSlice)
		}
	}
	return slice.Index(i)
}
