package conv

import (
	"math"
	"math/bits"
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToBool(t *testing.T) {
	var x bool = true
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst bool
	}{
		{x, false, x},
		{&x, false, x},
		{1, true, false}, // other kind
	}

	for _, test := range tests {
		var dst bool
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToInt8(t *testing.T) {
	var x int8 = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst int8
	}{
		{int8(x), false, x},
		{&x, false, x},
		{int(x), true, 0},
		{int16(x), true, 0},
		{int32(x), true, 0},
		{int64(x), true, 0},
		{uint(x), true, 0}, // other kind
	}

	for _, test := range tests {
		var dst int8
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToInt16(t *testing.T) {
	var x int16 = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst int16
	}{
		{int8(x), false, x},
		{int16(x), false, x},
		{&x, false, x},
		{int(x), true, 0},
		{int32(x), true, 0},
		{int64(x), true, 0},
		{uint(x), true, 0}, // other kind
	}

	for _, test := range tests {
		var dst int16
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToInt32(t *testing.T) {
	var x int32 = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst int32
	}{
		{&x, false, x},
		{int(x), false, x},
		{int8(x), false, x},
		{int16(x), false, x},
		{int32(x), false, x},
		{int64(x), true, 0},
		{uint(x), true, 0}, // other kind
	}
	if bits.UintSize == 64 {
		tests[1].expectedErr = true
		tests[1].expectedDst = 0
	}

	for _, test := range tests {
		var dst int32
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToInt64(t *testing.T) {
	var x int64 = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst int64
	}{
		{&x, false, x},
		{int(x), false, x},
		{int8(x), false, x},
		{int16(x), false, x},
		{int32(x), false, x},
		{int64(x), false, x},
		{uint(x), true, 0}, // other kind
	}

	for _, test := range tests {
		var dst int64
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToInt(t *testing.T) {
	var x int = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst int
	}{
		{&x, false, x},
		{int64(x), false, x},
		{int(x), false, x},
		{int8(x), false, x},
		{int16(x), false, x},
		{int32(x), false, x},
		{uint(x), true, 0}, // other kind
	}
	if bits.UintSize == 32 {
		tests[1].expectedErr = true
		tests[1].expectedDst = 0
	}

	for _, test := range tests {
		var dst int
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}
}

func TestToMap_single(t *testing.T) {
	var x int = 1
	tests := []struct {
		src         interface{}
		expectedErr bool
		expectedDst map[int]interface{}
	}{
		{&x, false, map[int]interface{}{0: int(1)}},
		{true, false, map[int]interface{}{0: true}},
		{int(1), false, map[int]interface{}{0: int(1)}},
		{int8(1), false, map[int]interface{}{0: int8(1)}},
		{int16(1), false, map[int]interface{}{0: int16(1)}},
		{int32(1), false, map[int]interface{}{0: int32(1)}},
		{int64(1), false, map[int]interface{}{0: int64(1)}},
		{uint(1), false, map[int]interface{}{0: uint(1)}},
		{uint8(1), false, map[int]interface{}{0: uint8(1)}},
		{uint16(1), false, map[int]interface{}{0: uint16(1)}},
		{uint32(1), false, map[int]interface{}{0: uint32(1)}},
		{uint64(1), false, map[int]interface{}{0: uint64(1)}},
		{float32(1), false, map[int]interface{}{0: float32(1)}},
		{float64(1), false, map[int]interface{}{0: float64(1)}},
		{complex64(complex(1, 0)), false, map[int]interface{}{0: complex64(complex(1, 0))}},
		{complex128(complex(1, 0)), false, map[int]interface{}{0: complex128(complex(1, 0))}},
	}

	for _, test := range tests {
		var dst map[int]interface{}
		err := To(test.src, &dst)
		if test.expectedErr {
			assert.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
		assert.Equal(t, test.expectedDst, dst)
	}

}

func TestToMap_composite(t *testing.T) {
	srcMap := map[int]int{1: 1, 2: 2}
	var dst map[int]int
	err := To(srcMap, &dst)
	require.Nil(t, err)
	assert.Equal(t, srcMap, dst)

	srcSlice := []int{1, 2}
	var dst1 map[int]int
	expected1 := map[int]int{0: 1, 1: 2}
	err = To(srcSlice, &dst1)
	require.Nil(t, err)
	assert.Equal(t, expected1, dst1)

	srcArray := [2]int{1, 2}
	var dst2 map[int]int
	expected2 := map[int]int{0: 1, 1: 2}
	err = To(srcArray, &dst2)
	require.Nil(t, err)
	assert.Equal(t, expected2, dst2)

	srcStruct := struct {
		a, b int
	}{1, 2}
	var dst3 map[string]int
	expected3 := map[string]int{"a": 1, "b": 2}
	err = To(srcStruct, &dst3)
	require.Nil(t, err)
	assert.Equal(t, expected3, dst3)
}

type testStringer struct{}
type testStringStruct struct{}

func (s testStringer) String() string {
	return "testStringer.String"
}

func TestWeakToString(t *testing.T) {
	var strger testStringer
	tests := []struct {
		src      interface{}
		expected string
	}{
		{strger, strger.String()},
		{&strger, strger.String()},
		{true, "true"},
		{false, "false"},
		{int(-1), "-1"},
		{uint(1), "1"},
		{1.25, "1.25"},
		{"string", "string"},
		{1.0 / (1 << 20), "9.5367431640625E-07"},
		{[]int{1, 2}, "[1 2]"},
		{complex(1, 1), "(1+1i)"},
		{complex(1.0/(1<<20), 1.25), "(9.5367431640625E-07+1.25i)"},
		{uintptr(0x123), "0x123"},
		{func() {}, "<func() Value>"},
	}

	str := ""
	dst := reflect.ValueOf(&str).Elem()
	for _, test := range tests {
		src := reflect.ValueOf(test.src)
		err := weakToString(src, dst)
		require.Nil(t, err)
		require.Equal(t, test.expected, dst.String())
	}

	// unsafe pointer
	src := reflect.ValueOf(unsafe.Pointer(&str))
	err := weakToString(src, dst)
	require.Nil(t, err)
	require.Regexp(t, "0x[0-9a-fA-F]+", dst.String())
}

func TestWeakToBool(t *testing.T) {
	x := 1
	succTests := []struct {
		src      interface{}
		expected bool
	}{
		{&x, true},

		{true, true},
		{false, false},

		{int(1), true},
		{int(0), false},

		{uint(1), true},
		{uint(0), false},

		{float32(1), true},
		{float64(0), false},

		{complex(0, 1), true},
		{complex(0, 0), false},

		{"1", true},
		{"t", true},
		{"T", true},
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"0", false},
		{"f", false},
		{"F", false},
		{"false", false},
		{"FALSE", false},
		{"False", false},
	}

	failureTests := []struct {
		src      interface{}
		expected error
	}{
		{[]int{}, &CannotConvError{reflect.Slice, reflect.Bool}},
	}

	for _, test := range succTests {
		var dst bool
		err := weakTo(test.src, &dst)
		require.Nil(t, err)
		assert.Equal(t, test.expected, dst)
	}

	for _, test := range failureTests {
		var dst bool
		err := WeakTo(test.src, &dst)
		require.Equal(t, test.expected, err)
		assert.Equal(t, false, dst)
	}
}

func TestWeakToTimeDuration(t *testing.T) {
	x := "2s"
	succTests := []struct {
		src      interface{}
		expected time.Duration
	}{
		{&x, 2 * time.Second},
		{"2s", 2 * time.Second},
		{2 * time.Nanosecond, 2 * time.Nanosecond},
	}

	for _, test := range succTests {
		var dst time.Duration
		err := WeakTo(test.src, &dst)
		require.Nil(t, err)
		assert.Equal(t, test.expected, dst)
	}

}

type TestTime time.Time

func TestWeakToTimeTime(t *testing.T) {
	x := time.Now()
	dt := x.Format(TimeLayout)
	succTests := []struct {
		src      interface{}
		expected time.Time
	}{
		{dt, x},
		{&x, x},
		{x, x},
		{TestTime(x), x},
	}

	for _, test := range succTests {
		var dst time.Time
		err := WeakTo(test.src, &dst)
		require.Nil(t, err)
		assert.Equal(t, test.expected, dst)
	}

}

func TestIsOverflowInt(t *testing.T) {
	succTests := []struct {
		src interface{}
	}{
		{int64(1)},
		{uint64(1)},
		{float64(1)},
		{complex(1, 0)},
	}

	failureTests := []struct {
		src interface{}
	}{
		{uint64(math.MaxInt64 + 1)},
		{float64(math.MaxInt64 + 1)},
		{float64(math.MaxInt64)},
		{float64(1.2)},
		{complex(0, 1)},
		{complex(1.2, 1)},
	}

	dsts := []reflect.Value{
		reflect.ValueOf(int(0)),
		reflect.ValueOf(int8(0)),
		reflect.ValueOf(int16(0)),
		reflect.ValueOf(int32(0)),
		reflect.ValueOf(int64(0)),
	}
	for _, test := range succTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.False(t, isOverflowInt(src, dst))
		}
	}

	for _, test := range failureTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test)
			assert.True(t, isOverflowInt(src, dst))
		}
	}
}

func TestIsOverflowUint(t *testing.T) {
	succTests := []struct {
		src interface{}
	}{
		{int64(1)},
		{uint64(1)},
		{float64(1)},
		{complex(1, 0)},
	}

	failureTests := []struct {
		src interface{}
	}{
		{int64(-1)},
		{float64(math.MaxUint64 + 1)},
		{float64(math.MaxUint64)},
		{float64(1.2)},
		{float64(-1)},
		{complex(0, 1)},
		{complex(1.2, 0)},
		{complex(-1, 0)},
	}

	dsts := []reflect.Value{
		reflect.ValueOf(uint(0)),
		reflect.ValueOf(uint8(0)),
		reflect.ValueOf(uint16(0)),
		reflect.ValueOf(uint32(0)),
		reflect.ValueOf(uint64(0)),
	}
	for _, test := range succTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.False(t, isOverflowUint(src, dst))
		}
	}

	for _, test := range failureTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.True(t, isOverflowUint(src, dst))
		}
	}
}

func TestIsOverflowFloat(t *testing.T) {
	succTests := []struct {
		src interface{}
	}{
		{int64(1)},
		{uint64(1)},
		{float64(1)},
		{float64(0.25)},
		{complex(1, 0)},

		{int64(-1)},
		{float64(-1)},
		{float64(-0.25)},
		{complex(-1, 0)},
		{complex(-0.25, 0)},
	}

	failureTests := []struct {
		src interface{}
	}{
		{int64(math.MaxInt64)},
		{uint64(math.MaxUint64)},
		{complex(0, 1)},
	}

	dsts := []reflect.Value{
		reflect.ValueOf(float32(0)),
		reflect.ValueOf(float64(0)),
	}
	for _, test := range succTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.False(t, isOverflowFloat(src, dst))
		}
	}

	for _, test := range failureTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.True(t, isOverflowFloat(src, dst))
		}
	}
}

func TestIsOverflowComplex(t *testing.T) {
	succTests := []struct {
		src interface{}
	}{
		{int64(1)},
		{uint64(1)},
		{float64(1)},
		{float64(0.25)},
		{complex(1, 0)},
		{complex(1, 1)},

		{int64(-1)},
		{float64(-1)},
		{float64(-0.25)},
		{complex(-1, 0)},
		{complex(-0.25, 0)},
	}

	failureTests := []struct {
		src interface{}
	}{
		{int64(math.MaxInt64)},
		{uint64(math.MaxUint64)},
	}

	dsts := []reflect.Value{
		reflect.ValueOf(complex64(complex(0, 0))),
		reflect.ValueOf(complex(0, 0)),
	}

	for _, test := range succTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.False(t, isOverflowComplex(src, dst))
		}
	}

	for _, test := range failureTests {
		for _, dst := range dsts {
			src := reflect.ValueOf(test.src)
			assert.True(t, isOverflowComplex(src, dst))
		}
	}
}

func TestWeakToInt(t *testing.T) {
	x := 1
	succTests := []struct {
		src      interface{}
		expected int64
	}{
		{&x, 1},

		{true, 1},
		{false, 0},

		{int(1), 1},
		{int8(1), 1},
		{int16(1), 1},
		{int32(1), 1},
		{int64(1), 1},

		{uint(1), 1},
		{uint8(1), 1},
		{uint16(1), 1},
		{uint32(1), 1},
		{uint64(1), 1},

		{float32(1), 1},
		{float64(1), 1},

		{complex64(complex(1, 0)), 1},
		{complex(1, 0), 1},

		{"1", 1},
		{"0", 0},
	}

	y := int16(math.MaxInt8 + 1)
	failureTests := []struct {
		src      interface{}
		expected error
	}{
		{[]int{}, &CannotConvError{reflect.Slice, reflect.Int8}},
		{int16(math.MinInt8 - 1), &OverflowError{int16(math.MinInt8 - 1), reflect.Int16, reflect.Int8}},
		{uint16(math.MaxInt8 + 1), &OverflowError{uint16(math.MaxInt8 + 1), reflect.Uint16, reflect.Int8}},
		{float32(math.MinInt8 - 1), &OverflowError{float32(math.MinInt8 - 1), reflect.Float32, reflect.Int8}},
		{complex(math.MinInt8-1, 0), &OverflowError{complex(math.MinInt8-1, 0), reflect.Complex128, reflect.Int8}},
		{strconv.FormatInt(math.MinInt8-1, 10), &OverflowError{strconv.FormatInt(math.MinInt8-1, 10), reflect.String, reflect.Int8}},
		{&y, &OverflowError{y, reflect.Int16, reflect.Int8}},
	}

	var (
		i   int
		i8  int8
		i16 int16
		i32 int32
		i64 int64
	)
	dsts := []reflect.Value{
		reflect.ValueOf(&i).Elem(),
		reflect.ValueOf(&i8).Elem(),
		reflect.ValueOf(&i16).Elem(),
		reflect.ValueOf(&i32).Elem(),
		reflect.ValueOf(&i64).Elem(),
	}

	for _, test := range succTests {
		src := reflect.ValueOf(test.src)
		for _, dst := range dsts {
			err := weakToInt(src, dst)
			require.Nilf(t, err, "src=%s, dst=%s", src.Kind(), dst.Kind())
			assert.EqualValues(t, test.expected, dst.Interface())
		}
	}

	for _, test := range failureTests {
		src := reflect.ValueOf(test.src)
		var i8 int8
		dst := reflect.ValueOf(&i8).Elem()
		err := weakToInt(src, dst)
		require.Zero(t, dst.Interface())
		assert.EqualValuesf(t, test.expected, err, "%s", err.Error())
	}
}

func TestWeakToUint(t *testing.T) {
	x := 1
	succTests := []struct {
		src      interface{}
		expected uint64
	}{
		{&x, 1},

		{true, 1},
		{false, 0},

		{int(1), 1},
		{int8(1), 1},
		{int16(1), 1},
		{int32(1), 1},
		{int64(1), 1},

		{uint(1), 1},
		{uint8(1), 1},
		{uint16(1), 1},
		{uint32(1), 1},
		{uint64(1), 1},

		{float32(1), 1},
		{float64(1), 1},

		{complex64(complex(1, 0)), 1},
		{complex(1, 0), 1},

		{"1", 1},
		{"0", 0},
	}

	failureTests := []struct {
		src      interface{}
		expected error
	}{
		{[]int{}, &CannotConvError{reflect.Slice, reflect.Uint8}},
		{int16(-1), &OverflowError{int16(-1), reflect.Int16, reflect.Uint8}},
		{uint16(math.MaxUint16), &OverflowError{uint16(math.MaxUint16), reflect.Uint16, reflect.Uint8}},
		{float32(-1), &OverflowError{float32(-1), reflect.Float32, reflect.Uint8}},
		{complex(-1, 0), &OverflowError{complex(-1, 0), reflect.Complex128, reflect.Uint8}},
	}

	var (
		u   uint
		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
	)
	dsts := []reflect.Value{
		reflect.ValueOf(&u).Elem(),
		reflect.ValueOf(&u8).Elem(),
		reflect.ValueOf(&u16).Elem(),
		reflect.ValueOf(&u32).Elem(),
		reflect.ValueOf(&u64).Elem(),
	}

	for _, test := range succTests {
		src := reflect.ValueOf(test.src)
		for _, dst := range dsts {
			err := weakToUint(src, dst)
			require.Nilf(t, err, "src=%s, dst=%s", src.Kind(), dst.Kind())
			assert.EqualValues(t, test.expected, dst.Interface())
		}
	}

	for _, test := range failureTests {
		src := reflect.ValueOf(test.src)
		u8 := uint8(0)
		dst := reflect.ValueOf(&u8).Elem()
		err := weakToUint(src, dst)
		require.Zero(t, dst.Interface())
		assert.EqualValuesf(t, test.expected, err, "%s", err.Error())
	}
}

func TestWeakToFloat(t *testing.T) {
	x := 1
	succTests := []struct {
		src      interface{}
		expected float64
	}{
		{&x, 1},

		{true, 1},
		{false, 0},

		{int(1), 1},
		{int8(1), 1},
		{int16(1), 1},
		{int32(1), 1},
		{int64(1), 1},

		{uint(1), 1},
		{uint8(1), 1},
		{uint16(1), 1},
		{uint32(1), 1},
		{uint64(1), 1},

		{float32(1), 1},
		{float64(1), 1},

		{complex64(complex(1, 0)), 1},
		{complex(1, 0), 1},

		{"1", 1},
		{"0", 0},
	}

	failureTests := []struct {
		src      interface{}
		expected error
	}{
		{[]int{}, &CannotConvError{reflect.Slice, reflect.Float64}},
		{int64(math.MaxInt64), &OverflowError{int64(math.MaxInt64), reflect.Int64, reflect.Float64}},
		{uint64(math.MaxUint64), &OverflowError{uint64(math.MaxUint64), reflect.Uint64, reflect.Float64}},
	}

	var (
		f32 float32
		f64 float64
	)
	dsts := []reflect.Value{
		reflect.ValueOf(&f32).Elem(),
		reflect.ValueOf(&f64).Elem(),
	}

	for _, test := range succTests {
		src := reflect.ValueOf(test.src)
		for _, dst := range dsts {
			err := weakToFloat(src, dst)
			require.Nilf(t, err, "src=%s, dst=%s", src.Kind(), dst.Kind())
			assert.EqualValues(t, test.expected, dst.Interface())
		}
	}

	for _, test := range failureTests {
		src := reflect.ValueOf(test.src)
		f := float64(0)
		dst := reflect.ValueOf(&f).Elem()
		err := weakToFloat(src, dst)
		require.Zero(t, dst.Interface())
		assert.EqualValuesf(t, test.expected, err, "%s", err.Error())
	}
}

func TestWeakToComplex(t *testing.T) {
	x := 1
	succTests := []struct {
		src      interface{}
		expected complex128
	}{
		{&x, complex(1, 0)},

		{true, complex(1, 0)},
		{false, complex(0, 0)},

		{int(1), complex(1, 0)},
		{int8(1), complex(1, 0)},
		{int16(1), complex(1, 0)},
		{int32(1), complex(1, 0)},
		{int64(1), complex(1, 0)},

		{uint(1), complex(1, 0)},
		{uint8(1), complex(1, 0)},
		{uint16(1), complex(1, 0)},
		{uint32(1), complex(1, 0)},
		{uint64(1), complex(1, 0)},

		{float32(1), complex(1, 0)},
		{float64(1), complex(1, 0)},

		{complex64(complex(1, 1)), complex(1, 1)},
		{complex(1, 1), complex(1, 1)},

		{"1+1i", complex(1, 1)},
		{"0", complex(0, 0)},
	}

	failureTests := []struct {
		src      interface{}
		expected error
	}{
		{[]int{}, &CannotConvError{reflect.Slice, reflect.Complex128}},
		{int64(math.MaxInt64), &OverflowError{int64(math.MaxInt64), reflect.Int64, reflect.Complex128}},
		{uint64(math.MaxUint64), &OverflowError{uint64(math.MaxUint64), reflect.Uint64, reflect.Complex128}},
	}

	var (
		c64  complex64
		c128 complex128
	)
	dsts := []reflect.Value{
		reflect.ValueOf(&c64).Elem(),
		reflect.ValueOf(&c128).Elem(),
	}

	for _, test := range succTests {
		src := reflect.ValueOf(test.src)
		for _, dst := range dsts {
			err := weakToComplex(src, dst)
			require.Nilf(t, err, "src=%s, dst=%s", src.Kind(), dst.Kind())
			require.EqualValues(t, test.expected, dst.Interface())
		}
	}

	for _, test := range failureTests {
		src := reflect.ValueOf(test.src)
		c := complex(0, 0)
		dst := reflect.ValueOf(&c).Elem()
		err := weakToComplex(src, dst)
		require.Zero(t, dst.Interface())
		require.EqualValuesf(t, test.expected, err, "%s", err.Error())
	}
}

func BenchmarkToBool(b *testing.B) {
	var src, dst bool
	for i := 0; i < b.N; i++ {
		To(src, &dst)
	}
}

func BenchmarkToInt(b *testing.B) {
	var src, dst int = 1, 0
	for i := 0; i < b.N; i++ {
		To(src, &dst)
	}
}
