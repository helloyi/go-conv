package conv

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/maltegrosse/go-bytesize"
)

func toBool(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		dst.SetBool(src.Bool())
		return nil

	case reflect.Interface, reflect.Ptr:
		return toBool(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}
}

func toInt(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Type().Size() < src.Type().Size() {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetInt(src.Int())
		return nil

	case reflect.Interface, reflect.Ptr:
		return toInt(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}
}

func toUint(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if dst.Type().Size() < src.Type().Size() {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetUint(src.Uint())
		return nil

	case reflect.Interface, reflect.Ptr:
		return toUint(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}
}

func toFloat(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Int:
		if src.Type().Size() == 8 /* int64 */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		// else in32
		if dst.Kind() == reflect.Float32 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Int())) /* float64 fraction 52bit */

	case reflect.Int8, reflect.Int16:
		dst.SetFloat(float64(src.Int()))

	case reflect.Int32:
		if dst.Kind() == reflect.Float32 /* fraction 23bit */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Int())) /* float64 fraction 52bit */

	case reflect.Uint:
		if src.Type().Size() == 8 /* uint64 */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		// else uin32
		if dst.Kind() == reflect.Float32 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Uint())) /* float64 fraction 52bit */

	case reflect.Uint8, reflect.Uint16:
		dst.SetFloat(float64(src.Uint()))

	case reflect.Uint32:
		if dst.Kind() == reflect.Float32 /* fraction 23bit */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Uint())) /* float64 fraction 52bit */

	case reflect.Float32, reflect.Float64:
		if dst.Type().Size() < src.Type().Size() {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetFloat(src.Float())

	case reflect.Interface, reflect.Ptr:
		return toFloat(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toComplex(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Int:
		if src.Type().Size() == 8 /* int64 */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		// else in32
		if dst.Kind() == reflect.Complex64 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(complex(float64(src.Int()), 0))

	case reflect.Int8, reflect.Int16:
		dst.SetComplex(complex(float64(src.Int()), 0))

	case reflect.Int32:
		if dst.Kind() == reflect.Complex64 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(complex(float64(src.Int()), 0))

	case reflect.Uint:
		if src.Type().Size() == 8 /* uint64 */ {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		// else uin32
		if dst.Kind() == reflect.Complex64 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(complex(float64(src.Uint()), 0))

	case reflect.Uint8, reflect.Uint16:
		dst.SetComplex(complex(float64(src.Uint()), 0))

	case reflect.Uint32:
		if dst.Kind() == reflect.Complex64 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(complex(float64(src.Uint()), 0))

	case reflect.Float32:
		dst.SetComplex(complex(src.Float(), 0))

	case reflect.Float64:
		if dst.Kind() == reflect.Complex64 {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(complex(src.Float(), 0))

	case reflect.Complex64, reflect.Complex128:
		if dst.Type().Size() < src.Type().Size() {
			return &CannotConvError{src.Kind(), dst.Kind()}
		}
		dst.SetComplex(src.Complex())

	case reflect.Interface, reflect.Ptr:
		return toComplex(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toArray(src, dst reflect.Value) error {
	return toArray0(src, dst, to0)
}

func toArray0(src, dst reflect.Value, convTo func(src, dst reflect.Value) error) error {
	switch src.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		fallthrough
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		if dst.Cap() <= 0 {
			break // TODO
		}

		dstElem := dst.Index(0)
		if err := convTo(src, dstElem); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < src.Len(); i++ {
			if i >= dst.Cap() {
				break // TODO
			}

			srcElem := src.Index(i)
			dstElem := dst.Index(i)
			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
		}

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if i >= dst.Cap() {
				break // TODO
			}

			srcElem := src.Field(i)
			dstElem := dst.Index(i)

			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
		}

	case reflect.Interface, reflect.Ptr:
		return toArray0(indirect(src), dst, convTo)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toInterface(src, dst reflect.Value) error {
	dst.Set(src)
	return nil
}

func toMap(src, dst reflect.Value) error {
	return toMap0(src, dst, to0)
}

func toMap0(src, dst reflect.Value, convTo func(src, dst reflect.Value) error) error {
	switch src.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		fallthrough
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dst.Type(), 1))
		}
		dst.SetMapIndex(reflect.ValueOf(nil), src)

	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dst.Type(), src.Len()))
		}

		iter := src.MapRange()
		for iter.Next() {
			key := iter.Key()

			srcElem := iter.Value()
			dstElem := dst.MapIndex(key)
			if dstElem.Kind() == reflect.Invalid {
				dstElem = reflect.New(dst.Type().Elem())
				dstElem = dstElem.Elem()
			}

			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
			dst.SetMapIndex(key, dstElem)
		}

	case reflect.Slice, reflect.Array:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dst.Type(), src.Len()))
		}

		for i := 0; i < src.Len(); i++ {
			key := reflect.ValueOf(i) // TODO: check map key type

			srcElem := src.Index(i)
			dstElem := dst.MapIndex(key)
			if dstElem.Kind() == reflect.Invalid {
				dstElem = reflect.New(dst.Type().Elem())
				dstElem = dstElem.Elem()
			}

			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
			dst.SetMapIndex(key, dstElem)
		}

	case reflect.Struct:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dst.Type(), src.NumField()))
		}

		for i := 0; i < src.NumField(); i++ {
			name := src.Type().Field(i).Name
			key := reflect.ValueOf(reflect.TypeOf(name)) // TODO: check map key type

			srcField := src.Field(i)
			dstElem := dst.MapIndex(key)
			if dstElem.Kind() == reflect.Invalid {
				dstElem = reflect.New(dst.Type().Elem())
				dstElem = dstElem.Elem()
			}

			if err := convTo(srcField, dstElem); err != nil {
				return err
			}
			dst.SetMapIndex(key, dstElem)
		}

	case reflect.Interface, reflect.Ptr:
		return toMap0(indirect(src), dst, convTo)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toPtr(src, dst reflect.Value) error {
	realdst := dst
	if dst.IsNil() {
		realdst = reflect.New(dst.Type().Elem())
	}
	if err := to0(src, realdst.Elem()); err != nil {
		return err
	}
	dst.Set(realdst)

	return nil
}

func toSlice(src, dst reflect.Value) error {
	return toSlice0(src, dst, to0)
}

func toSlice0(src, dst reflect.Value, convTo func(src, dst reflect.Value) error) error {
	switch src.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		fallthrough
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		if dst.IsNil() {
			dst.Set(reflect.MakeSlice(dst.Type(), 1, 1))
		}

		dstElem := sliceIndex(dst, 0)
		if err := convTo(src, dstElem); err != nil {
			return err
		}

	case reflect.Array, reflect.Slice:
		if dst.IsNil() {
			dst.Set(reflect.MakeSlice(dst.Type(), src.Len(), src.Len()))
		}

		for i := 0; i < src.Len(); i++ {
			srcElem := src.Index(i)
			dstElem := sliceIndex(dst, i)

			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
		}

	case reflect.Struct:
		if dst.IsNil() {
			dst.Set(reflect.MakeSlice(dst.Type(), src.NumField(), src.NumField()))
		}

		for i := 0; i < src.NumField(); i++ {
			srcElem := src.Field(i)
			dstElem := sliceIndex(dst, i)

			if err := convTo(srcElem, dstElem); err != nil {
				return err
			}
		}

	case reflect.Interface, reflect.Ptr:
		return toSlice0(indirect(src), dst, convTo)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
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

func toString(src, dst reflect.Value) error {
	stringer := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	if src.Type().Implements(stringer) {
		string, _ := src.Type().MethodByName("String")
		s := string.Func.Call(nil)[0].String()
		dst.SetString(s)
		return nil
	}

	switch src.Kind() {
	// TODO: bytes to string
	case reflect.String:
		dst.SetString(src.String())

	case reflect.Interface, reflect.Ptr:
		toString(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toStruct(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		fallthrough
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		fallthrough
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		fallthrough
	case reflect.Float32, reflect.Float64:
		fallthrough
	case reflect.Complex64, reflect.Complex128:
		if dst.NumField() == 0 {
			return nil
		}
		dstElem := dst.Field(0)
		if err := to0(src, dstElem); err != nil {
			return err
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < src.Len(); i++ {
			if i >= dst.NumField() {
				break
			}

			srcElem := src.Index(i)
			dstElem := dst.Field(i)

			if dstElem.Kind() == reflect.Invalid {
				continue // TODO
			}

			if err := to0(srcElem, dstElem); err != nil {
				return err
			}
		}

	case reflect.Map:
		iter := src.MapRange()
		for iter.Next() {
			srcKey := iter.Key().String()

			isMatchCase := true
			if srcKey == strings.ToLower(srcKey) {
				isMatchCase = false
			}

			var dstField reflect.Value
			if isMatchCase {
				dstField = dst.FieldByName(srcKey)
			} else {
				dstField = dst.FieldByNameFunc(func(fn string) bool {
					return strings.ToLower(fn) == srcKey
				})
			}

			if dstField.Kind() == reflect.Invalid { // not exist
				continue
			}

			if err := to0(iter.Value(), dstField); err != nil {
				return err
			}
		}

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			srcField := src.Field(i)

			srcFname := src.Type().Field(i).Name
			dstField := dst.FieldByName(srcFname)

			if dstField.Kind() == reflect.Invalid { // not exist field
				continue
			}

			if err := to0(srcField, dstField); err != nil {
				return err
			}
		}

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func toTimeDuration(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		str := src.String()
		dur, err := time.ParseDuration(str)
		if err != nil {
			return err
		}
		dst.SetInt(int64(dur))

	case reflect.Interface, reflect.Ptr:
		return toTimeDuration(indirect(src), dst)

	default:
		return toInt(src, dst)
	}

	return nil
}

func toTimeTime(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		t, err := time.Parse(TimeLayout, s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(t))

	case reflect.Interface, reflect.Ptr:
		return toTimeTime(indirect(src), dst)

	default:
		return toStruct(src, dst)
	}

	return nil
}

func toByteSize(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		bs, err := bytesize.Parse(s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(ByteSize(bs)))

	case reflect.Interface, reflect.Ptr:
		return toByteSize(indirect(src), dst)

	default:
		return toStruct(src, dst)
	}

	return nil
}
