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

func isOverflowInt(src, dst reflect.Value) bool {
	var x int64
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x = src.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u := src.Uint()
		if u != uint64(int64(u)) {
			return true
		}

		x = int64(u)

	case reflect.Float32, reflect.Float64:
		f := src.Float()
		if f != float64(int64(f)) {
			return true
		}

		x = int64(f)

	case reflect.Complex64, reflect.Complex128:
		c := src.Complex()
		if imag(c) != 0 {
			return true
		}

		f := real(c)
		if f != float64(int64(f)) {
			return true
		}

		x = int64(f)

	default:
		panic("invalid kind")
	}

	return dst.OverflowInt(x)
}

func isOverflowUint(src, dst reflect.Value) bool {
	var x uint64
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := src.Int()
		if i != int64(uint64(i)) {
			return true
		}

		x = uint64(i)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x = src.Uint()

	case reflect.Float32, reflect.Float64:
		f := src.Float()
		if f != float64(uint64(f)) {
			return true
		}

		x = uint64(f)

	case reflect.Complex64, reflect.Complex128:
		c := src.Complex()
		if imag(c) != 0 {
			return true
		}

		f := real(c)
		if f != float64(uint64(f)) {
			return true
		}

		x = uint64(f)

	default:
		panic("invalid kind")
	}

	return dst.OverflowUint(x)
}
