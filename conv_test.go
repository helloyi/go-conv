package conv

import (
	"math/bits"
	"testing"

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
