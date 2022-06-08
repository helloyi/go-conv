package conv

import (
	"fmt"
	"reflect"
)

// CannotConvError ...
type CannotConvError struct {
	srcKind reflect.Kind
	dstKind reflect.Kind
}

// OverflowError ...
type OverflowError struct {
	num     interface{}
	srcKind reflect.Kind
	dstKind reflect.Kind
}

type CannotSetError struct{}

func (e *CannotConvError) Error() string {
	return fmt.Sprintf("cannot convert %s to %s", e.srcKind, e.dstKind)
}

func (e *OverflowError) Error() string {
	return fmt.Sprintf("convert %v to %s overflows", e.num, e.dstKind.String())
}

func (e *CannotSetError) Error() string {
	return fmt.Sprintf("cannot set")
}
