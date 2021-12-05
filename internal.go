package conv

import (
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
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
		key := reflect.Zero(dst.Type().Key())
		dstElem := mapIndex(dst, key)
		if err := convTo(src, dstElem); err != nil {
			return err
		}
		dst.SetMapIndex(key, dstElem)

	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMapWithSize(dst.Type(), src.Len()))
		}

		iter := src.MapRange()
		for iter.Next() {
			key := iter.Key()

			srcElem := iter.Value()
			dstElem := mapIndex(dst, key)
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
			dstElem := mapIndex(dst, key)

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
			key := reflect.ValueOf(name) // TODO: check map key type

			srcField := src.Field(i)
			dstElem := mapIndex(dst, key)

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

func toNetHardwareAddr(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		haddr, err := net.ParseMAC(s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(haddr))
		return nil

	case reflect.Interface, reflect.Ptr:
		return toNetHardwareAddr(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
		// TODO: toBytes(src, dst)
	}
}

func toNetIP(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		ip := net.ParseIP(s)
		if len(ip) == 0 {
			return errors.New("invalid ip")
		}
		dst.Set(reflect.ValueOf(ip))
		return nil

	case reflect.Interface, reflect.Ptr:
		return toNetIP(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
		// TODO: toBytes(src, dst)
	}
}

func toNetURL(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		url, err := url.Parse(s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(*url))
		return nil

	case reflect.Interface, reflect.Ptr:
		return toNetURL(indirect(src), dst)

	default:
		return toStruct(src, dst)
	}
}

func toMailAddress(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		addr, err := mail.ParseAddress(s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(*addr))
		return nil

	case reflect.Interface, reflect.Ptr:
		return toMailAddress(indirect(src), dst)

	default:
		return toStruct(src, dst)
	}
}

func toRegexpRegexp(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.String:
		s := src.String()
		r, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		dst.Set(reflect.ValueOf(*r))
		// TODO: POSIX regexp
		return nil

	case reflect.Interface, reflect.Ptr:
		return toRegexpRegexp(indirect(src), dst)

	default:
		return toStruct(src, dst)
	}
}

func weakToBool(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		dst.SetBool(src.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dst.SetBool(src.Int() != 0)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		dst.SetBool(src.Uint() != 0)

	case reflect.Float32, reflect.Float64:
		dst.SetBool(src.Float() != 0)

	case reflect.Complex64, reflect.Complex128:
		dst.SetBool(src.Complex() != complex(0, 0))

	case reflect.String:
		// "1", "t", "T", "true", "TRUE", "True":
		// "0", "f", "F", "false", "FALSE", "False":
		b, err := strconv.ParseBool(src.String())
		if err != nil {
			return err
		}
		dst.SetBool(b)

	case reflect.Interface, reflect.Ptr:
		return weakToBool(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func weakToInt(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		if src.Bool() {
			dst.SetInt(1)
		} else {
			dst.SetInt(0)
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if isOverflowInt(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetInt(src.Int())

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if isOverflowInt(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetInt(int64(src.Uint()))

	case reflect.Float32, reflect.Float64:
		if isOverflowInt(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetInt(int64(src.Float()))

	case reflect.Complex64, reflect.Complex128:
		if isOverflowInt(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetInt(int64(real(src.Complex())))

	case reflect.String:
		i, err := strconv.ParseInt(src.String(), 10, 64)
		if err != nil {
			return err
		}
		if dst.OverflowInt(i) {
			return &OverflowError{src.String(), src.Kind(), dst.Kind()}
		}
		dst.SetInt(i)

	case reflect.Interface, reflect.Ptr:
		return weakToInt(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func weakToUint(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		if src.Bool() {
			dst.SetUint(1)
		} else {
			dst.SetUint(0)
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if isOverflowUint(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetUint(src.Uint())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if isOverflowUint(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		dst.SetUint(uint64(src.Int()))

	case reflect.Float32, reflect.Float64:
		if isOverflowUint(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		f := src.Float()
		dst.SetUint(uint64(f))

	case reflect.Complex64, reflect.Complex128:
		if isOverflowUint(src, dst) {
			return &OverflowError{src.Interface(), src.Kind(), dst.Kind()}
		}
		c := src.Complex()
		dst.SetUint(uint64(real(c)))

	case reflect.String:
		u, err := strconv.ParseUint(src.String(), 10, 64)
		if err != nil {
			return err
		}
		if dst.OverflowUint(u) {
			return &OverflowError{src.String(), src.Kind(), dst.Kind()}
		}
		dst.SetUint(u)

	case reflect.Interface, reflect.Ptr:
		return weakToUint(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}

func weakToFloat(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		if src.Bool() {
			dst.SetFloat(1)
		} else {
			dst.SetFloat(0)
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if isOverflowFloat(src, dst) {
			return &OverflowError{src.Int(), src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Int()))

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if isOverflowFloat(src, dst) {
			return &OverflowError{src.Uint(), src.Kind(), dst.Kind()}
		}
		dst.SetFloat(float64(src.Uint()))

	case reflect.Float32, reflect.Float64:
		if isOverflowFloat(src, dst) {
			return &OverflowError{src.Float(), src.Kind(), dst.Kind()}
		}
		f := src.Float()
		dst.SetFloat(f)

	case reflect.Complex64, reflect.Complex128:
		if isOverflowFloat(src, dst) {
			return &OverflowError{src.Complex(), src.Kind(), dst.Kind()}
		}
		c := src.Complex()
		dst.SetFloat(real(c))

	case reflect.String:
		f, err := strconv.ParseFloat(src.String(), 64)
		if err != nil {
			return err
		}
		if dst.OverflowFloat(f) {
			return &OverflowError{src.String(), src.Kind(), dst.Kind()}
		}
		dst.SetFloat(f)

	case reflect.Interface, reflect.Ptr:
		return weakToFloat(indirect(src), dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}

	return nil
}
