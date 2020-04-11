package signal_test

import (
	"math"
	"reflect"
	"testing"

	"pipelined.dev/signal"
)

var (
	_ signal.Signed   = signal.Allocator{}.Int64Interleaved(signal.MaxBitDepth)
	_ signal.Signed   = signal.Allocator{}.Int64(signal.MaxBitDepth)
	_ signal.Floating = signal.Allocator{}.Float64()
)

func TestInt64Interleaved(t *testing.T) {
	alloc := signal.Allocator{Channels: 3, Capacity: 2}
	buf := alloc.Int64Interleaved(signal.MaxBitDepth)
	buf.WriteInt64([]int64{1, 2, 3})
	buf.WriteInt64([]int64{4, 5, 6, 7, 8, 9})
}

func TestInt64InterleavedAsFloat64(t *testing.T) {
	tests := map[string]struct {
		ints []int64
		signal.BitDepth
		props    signal.Allocator
		expected [][]float64
	}{
		"8 bits": {
			ints: []int64{math.MaxInt8, math.MaxInt8 * 2},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth8,
			expected: [][]float64{
				{1},
				{2},
			},
		},
		"16 bits": {
			ints: []int64{math.MaxInt16, math.MaxInt16 * 2},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth16,
			expected: [][]float64{
				{1},
				{2},
			},
		},
		"32 bits": {
			ints: []int64{math.MaxInt32, math.MaxInt32 * 2},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth32,
			expected: [][]float64{
				{1},
				{2},
			},
		},
		"24 bits": {
			ints: []int64{1<<23 - 1, (1<<23 - 1) * 2},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth24,
			expected: [][]float64{
				{1},
				{2},
			},
		},
	}

	for _, test := range tests {
		ints := test.props.Int64Interleaved(test.BitDepth)
		ints.WriteInt64(test.ints)

		result := signal.Allocator{
			Channels: test.props.Channels,
			Capacity: test.props.Capacity,
		}.Float64()
		signal.SignedAsFloating(ints, result)

		expected := signal.Allocator{
			Channels: len(test.expected),
			Capacity: len(test.expected[0]),
		}.Float64()
		expected.WriteFloat64(test.expected)
		assertEqual(t, "slices", result, expected)
	}
}

func TestWriteInt(t *testing.T) {
	testInt64 := func(s signal.Int64, expectedLen int, expected [][]int64, ints ...[][]int) func(t *testing.T) {
		return func(t *testing.T) {
			for i := range ints {
				s.WriteInt(ints[i])
			}
			assertEqual(t, "length", s.Length(), expectedLen)
			assertEqual(t, "slices", s.Data(), expected)
		}
	}

	t.Run("int64 full buffer", testInt64(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		10,
		[][]int64{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		[][]int{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	))
	t.Run("int64 short buffer", testInt64(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		3,
		[][]int64{
			{1, 2, 3, 0, 0, 0, 0, 0, 0, 0},
			{11, 12, 13, 0, 0, 0, 0, 0, 0, 0},
		},
		[][]int{
			{1, 2, 3},
			{11, 12, 13},
		},
	))
	t.Run("int64 multiple short buffers", testInt64(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		6,
		[][]int64{
			{1, 2, 3, 4, 5, 6, 0, 0, 0, 0},
			{11, 12, 13, 14, 15, 16, 0, 0, 0, 0},
		},
		[][]int{
			{1, 2, 3},
			{11, 12, 13},
		},
		[][]int{
			{4, 5, 6},
			{14, 15, 16},
		},
	))
	t.Run("int64 long buffer", testInt64(
		signal.Allocator{Capacity: 3, Channels: 2}.Int64(signal.MaxBitDepth),
		3,
		[][]int64{
			{1, 2, 3},
			{11, 12, 13},
		},
		[][]int{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
	))
	t.Run("int64 8-bits overflow", testInt64(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth8),
		1,
		[][]int64{
			{math.MaxInt8},
			{math.MinInt8},
		},
		[][]int{
			{math.MaxInt32},
			{math.MinInt32},
		},
	))
	t.Run("int64 16-bits overflow", testInt64(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth16),
		1,
		[][]int64{
			{math.MaxInt16},
			{math.MinInt16},
		},
		[][]int{
			{math.MaxInt64},
			{math.MinInt64},
		},
	))
}

func TestAppendFloat64(t *testing.T) {
	testOk := func(s signal.Float64, expected [][]float64, slices ...[][]float64) func(*testing.T) {
		return func(t *testing.T) {
			for _, slice := range slices {
				src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Float64()
				src.WriteFloat64(slice)
				s = s.Append(src)
			}
			assertEqual(t, "slices", s.Data(), expected)
		}
	}
	testPanic := func(s signal.Float64, expected [][]float64, slice [][]float64) func(*testing.T) {
		return func(t *testing.T) {
			src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Float64()
			src.WriteFloat64(slice)
			assertPanic(t, func() {
				s = s.Append(src)
			})
		}
	}

	t.Run("single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][]float64{
			{1, 2},
			{1, 2},
		},
		[][]float64{
			{1, 2},
			{1, 2},
		},
	))
	t.Run("multiple slices", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][]float64{
			{1, 2, 3, 4},
			{1, 2, 3, 4},
		},
		[][]float64{
			{1, 2},
			{1, 2},
		},
		[][]float64{
			{3, 4},
			{3, 4},
		},
	))
	t.Run("different channels slice", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][]float64{
			{1, 2},
			{1, 2},
		},
		[][]float64{
			{1, 2},
		},
	))
}

func TestAppendInt64(t *testing.T) {
	testInt64 := func(s signal.Int64, expected [][]int64, slices ...[][]int64) func(*testing.T) {
		return func(t *testing.T) {
			for _, slice := range slices {
				src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Int64(signal.MaxBitDepth)
				src.WriteInt64(slice)
				s = s.Append(src)
			}
			assertEqual(t, "slices", s.Data(), expected)
		}
	}
	t.Run("single slice", testInt64(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2},
			{1, 2},
		},
		[][]int64{
			{1, 2},
			{1, 2},
		},
	))
	t.Run("multiple slices", testInt64(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 3, 4},
			{1, 2, 3, 4},
		},
		[][]int64{
			{1, 2},
			{1, 2},
		},
		[][]int64{
			{3, 4},
			{3, 4},
		},
	))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\ninvalid %v value: %+v \nexpected: %+v", name, result, expected)
	}
}

func assertPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

// func assertFloatingBuffers(t *testing.T, result, expected signal.Floating) {
// 	assertBuffers(t, result, expected)

// 	for channel := 0; channel < expected.Channels(); channel++ {
// 		for pos := 0; pos < expected.Length(); pos++ {
// 			if es, rs := expected.Sample(channel, pos), result.Sample(channel, pos); es != rs {
// 				t.Fatalf("invalid value: %v expected %v", rs, es)
// 			}
// 		}
// 	}
// }

// func assertSignedBuffers(t *testing.T, result, expected signal.Signed) {
// 	assertBuffers(t, result, expected)
// 	for channel := 0; channel < expected.Channels(); channel++ {
// 		for pos := 0; pos < expected.Length(); pos++ {
// 			if es, rs := expected.Sample(channel, pos), result.Sample(channel, pos); es != rs {
// 				t.Fatalf("invalid value: %v expected %v", rs, es)
// 			}
// 		}
// 	}
// }

// func assertUnsignedBuffers(t *testing.T, result, expected signal.Unsigned) {
// 	assertBuffers(t, result, expected)
// 	for channel := 0; channel < expected.Channels(); channel++ {
// 		for pos := 0; pos < expected.Length(); pos++ {
// 			if es, rs := expected.Sample(channel, pos), result.Sample(channel, pos); es != rs {
// 				t.Fatalf("invalid value: %v expected %v", rs, es)
// 			}
// 		}
// 	}
// }

// func assertBuffers(t *testing.T, result, expected signal.Signal) {
// 	if expected.Channels() != result.Channels() {
// 		t.Fatalf("invalid num channels: %v expeced %v", result.Channels(), expected.Channels())
// 	}
// 	if expected.Size() != result.Size() {
// 		t.Fatalf("invalid buffer size: %v expeced %v", result.Size(), expected.Size())
// 	}
// }
