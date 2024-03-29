//package conv
package conv

import (
	"errors"
	"reflect"

	"github.com/maltegrosse/go-bytesize"
)

// TimeLayout default time layout for convert To time.Time
var TimeLayout = "Mon Jan 2 15:04:05 -0700 MST 2006"

// ByteSize type of byte size
type ByteSize bytesize.ByteSize

// To convert to src to dst
func To(src, dst interface{}) error {
	return to(src, dst)
}

// WeakTo convert to src to dst (weak type convert)
func WeakTo(src, dst interface{}) error {
	return weakTo(src, dst)
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
	if !dst.CanSet() {
		return &CannotSetError{}
	}

	switch dst.Type().PkgPath() {
	case "time":
		switch dst.Type().Name() {
		case "Duration":
			return toTimeDuration(src, dst)
		case "Time":
			return toTimeTime(src, dst)
		}
	case "net":
		switch dst.Type().Name() {
		case "IP":
			return toNetIP(src, dst)
		case "HardwareAddr":
			return toNetHardwareAddr(src, dst)
		}
	case "net/url":
		if dst.Type().Name() == "URL" {
			return toNetURL(src, dst)
		}
	case "net/mail":
		if dst.Type().Name() == "Address" {
			return toMailAddress(src, dst)
		}
	case "regexp":
		if dst.Type().Name() == "Regexp" {
			return toRegexpRegexp(src, dst)
		}
	case "github.com/helloyi/go-conv":
		if dst.Type().Name() == "ByteSize" {
			return toByteSize(src, dst)
		}
	}

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

func weakTo(src, dst interface{}) error {
	dstv := reflect.ValueOf(dst)
	if dstv.Kind() != reflect.Ptr {
		return errors.New("non-pointer of dst")
	}
	srcv := reflect.ValueOf(src)

	return weakTo0(srcv, dstv.Elem())
}

func weakTo0(src, dst reflect.Value) error {
	if !dst.CanSet() {
		return &CannotSetError{}
	}

	switch dst.Type().PkgPath() {
	case "time":
		switch dst.Type().Name() {
		case "Duration":
			return weakToTimeDuration(src, dst)
		case "Time":
			return weakToTimeTime(src, dst)
		}
	case "net":
		switch dst.Type().Name() {
		case "IP":
			return weakToNetIP(src, dst)
		case "HardwareAddr":
			return weakToNetHardwareAddr(src, dst)
		}
	case "net/url":
		if dst.Type().Name() == "URL" {
			return weakToNetURL(src, dst)
		}
	case "net/mail":
		if dst.Type().Name() == "Address" {
			return weakToMailAddress(src, dst)
		}
	case "regexp":
		if dst.Type().Name() == "Regexp" {
			return weakToRegexpRegexp(src, dst)
		}
	case "github.com/helloyi/go-conv":
		if dst.Type().Name() == "ByteSize" {
			return weakToByteSize(src, dst)
		}
	}

	switch dst.Kind() {
	case reflect.Bool:
		return weakToBool(src, dst)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return weakToInt(src, dst)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return weakToUint(src, dst)

	case reflect.Float32, reflect.Float64:
		return weakToFloat(src, dst)

	case reflect.Complex64, reflect.Complex128:
		return weakToComplex(src, dst)

	case reflect.String:
		return weakToString(src, dst)

	case reflect.Array:
		return weakToArray(src, dst)

	case reflect.Interface:
		return weakToInterface(src, dst)

	case reflect.Map:
		return weakToMap(src, dst)

	case reflect.Ptr:
		return weakToPtr(src, dst)

	case reflect.Slice:
		return weakToSlice(src, dst)

	case reflect.Struct:
		return weakToStruct(src, dst)

	default:
		return &CannotConvError{src.Kind(), dst.Kind()}
	}
}
