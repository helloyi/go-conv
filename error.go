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

func (e *CannotConvError) Error() string {
	return fmt.Sprintf("cannot convert %s to %s", e.srcKind, e.dstKind)
}
