//package conv
package conv

import (
	"errors"
	"reflect"
)

// To convert to src to dst
func To(src, dst interface{}) error {
	return to(src, dst)
}

func to(src, dst interface{}) error {
	dstv := reflect.ValueOf(dst)
	if dstv.Kind() != reflect.Ptr {
		return errors.New("non-pointer of dst")
	}
	srcv := reflect.ValueOf(src)

	return to0(srcv, dstv.Elem())
}

func to0(src, dst reflect.Value) (err error) {
	switch dst.Kind() {
	case reflect.Bool:
		return toBool(src, dst)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toInt(src, dst)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return toUint(src, dst)

	case reflect.Float32, reflect.Float64:
		return toFloat(src, dst)

	case reflect.Complex64, reflect.Complex128:
		return toComplex(src, dst)

	case reflect.Array:
		return toArray(src, dst)

	case reflect.Interface:
		return toInterface(src, dst)

	case reflect.Map:
		return toMap(src, dst)

	case reflect.Ptr:
		return toPtr(src, dst)

	case reflect.Slice:
		return toSlice(src, dst)

	case reflect.String:
		return toString(src, dst)

	case reflect.Struct:
		return toStruct(src, dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}
}
