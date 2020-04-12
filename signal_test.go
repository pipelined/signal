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

// func TestInt64Interleaved(t *testing.T) {
// 	alloc := signal.Allocator{Channels: 3, Capacity: 2}
// 	buf := alloc.Int64Interleaved(signal.MaxBitDepth)
// 	buf.WriteInt64([]int64{1, 2, 3})
// 	buf.WriteInt64([]int64{4, 5, 6, 7, 8, 9})
// }

func TestInt64InterleavedAsFloat64(t *testing.T) {
	tests := map[string]struct {
		ints []int64
		signal.BitDepth
		props    signal.Allocator
		expected [][]float64
	}{
		"8 bits": {
			ints: []int64{math.MaxInt8, math.MaxInt8 + 1},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth8,
			expected: [][]float64{
				{1},
				{1},
			},
		},
		"16 bits": {
			ints: []int64{math.MaxInt16, math.MaxInt16 + 1},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth16,
			expected: [][]float64{
				{1},
				{1},
			},
		},
		"32 bits": {
			ints: []int64{math.MaxInt32, math.MaxInt32 + 1},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth32,
			expected: [][]float64{
				{1},
				{1},
			},
		},
		"24 bits": {
			ints: []int64{1<<23 - 1, (1<<23 - 1) + 1},
			props: signal.Allocator{
				Channels: 2,
				Capacity: 1,
			},
			BitDepth: signal.BitDepth24,
			expected: [][]float64{
				{1},
				{1},
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

func TestWrite(t *testing.T) {
	type expected struct {
		length int
		data   interface{}
	}
	testOk := func(writer signal.Signal, data interface{}, ex expected) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			switch w := writer.(type) {
			case signal.Int64:
				switch d := data.(type) {
				case [][]int:
					w.WriteInt(d)
				case [][]int64:
					w.WriteInt64(d)
				}
				assertEqual(t, "slices", w.Data(), ex.data)
			case signal.Int64Interleaved:
				switch d := data.(type) {
				case []int:
					w.WriteInt(d)
				case []int64:
					w.WriteInt64(d)
				}
				assertEqual(t, "slices", w.Data(), ex.data)
			case signal.Float64:
				d := data.([][]float64)
				w.WriteFloat64(d)
				assertEqual(t, "slices", w.Data(), ex.data)
			default:
				t.Fatalf("unsupported write type %T", writer)
			}
			assertEqual(t, "length", writer.Length(), ex.length)
		}
	}

	t.Run("int64 int full buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 int short buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int{
			{1, 2, 3},
			{11, 12, 13},
		},
		expected{
			length: 3,
			data: [][]int64{
				{1, 2, 3, 0, 0, 0, 0, 0, 0, 0},
				{11, 12, 13, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	))
	t.Run("int64 int long buffer", testOk(
		signal.Allocator{Capacity: 3, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 int 8-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth8),
		[][]int{
			{math.MaxInt32},
			{math.MinInt32},
		},
		expected{
			length: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 int 16-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth16),
		[][]int{
			{math.MaxInt64},
			{math.MinInt64},
		},
		expected{
			length: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("int64 int64 full buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 int64 short buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 3},
			{11, 12, 13},
		},
		expected{
			length: 3,
			data: [][]int64{
				{1, 2, 3, 0, 0, 0, 0, 0, 0, 0},
				{11, 12, 13, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	))
	t.Run("int64 int64 long buffer", testOk(
		signal.Allocator{Capacity: 3, Channels: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 int64 8-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth8),
		[][]int64{
			{math.MaxInt32},
			{math.MinInt32},
		},
		expected{
			length: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 int64 16-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64(signal.BitDepth16),
		[][]int64{
			{math.MaxInt64},
			{math.MinInt64},
		},
		expected{
			length: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("int64interleaved int full buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		expected{
			length: 10,
			data:   []int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		},
	))
	t.Run("int64interleaved int short buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int{1, 11, 2, 12, 3, 13},
		expected{
			length: 3,
			data:   []int64{1, 11, 2, 12, 3, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	))
	t.Run("int64interleaved int long buffer", testOk(
		signal.Allocator{Capacity: 3, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		expected{
			length: 3,
			data:   []int64{1, 11, 2, 12, 3, 13},
		},
	))
	t.Run("int64interleaved int 8-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64Interleaved(signal.BitDepth8),
		[]int{math.MaxInt32, math.MinInt32},
		expected{
			length: 1,
			data:   []int64{math.MaxInt8, math.MinInt8},
		},
	))
	t.Run("int64interleaved int 16-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64Interleaved(signal.BitDepth16),
		[]int{math.MaxInt64, math.MinInt64},
		expected{
			length: 1,
			data:   []int64{math.MaxInt16, math.MinInt16},
		},
	))
	t.Run("int64interleaved int full buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		expected{
			length: 10,
			data:   []int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		},
	))
	t.Run("int64interleaved int short buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int64{1, 11, 2, 12, 3, 13},
		expected{
			length: 3,
			data:   []int64{1, 11, 2, 12, 3, 13, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	))
	t.Run("int64interleaved int long buffer", testOk(
		signal.Allocator{Capacity: 3, Channels: 2}.Int64Interleaved(signal.MaxBitDepth),
		[]int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20},
		expected{
			length: 3,
			data:   []int64{1, 11, 2, 12, 3, 13},
		},
	))
	t.Run("int64interleaved int 8-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64Interleaved(signal.BitDepth8),
		[]int64{math.MaxInt32, math.MinInt32},
		expected{
			length: 1,
			data:   []int64{math.MaxInt8, math.MinInt8},
		},
	))
	t.Run("int64interleaved int 16-bits overflow", testOk(
		signal.Allocator{Capacity: 1, Channels: 2}.Int64Interleaved(signal.BitDepth16),
		[]int64{math.MaxInt64, math.MinInt64},
		expected{
			length: 1,
			data:   []int64{math.MaxInt16, math.MinInt16},
		},
	))
	t.Run("float64 full buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Float64(),
		[][]float64{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 10,
			data: [][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("float64 short buffer", testOk(
		signal.Allocator{Capacity: 10, Channels: 2}.Float64(),
		[][]float64{
			{1, 2, 3},
			{11, 12, 13},
		},
		expected{
			length: 3,
			data: [][]float64{
				{1, 2, 3, 0, 0, 0, 0, 0, 0, 0},
				{11, 12, 13, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	))
	t.Run("float64 long buffer", testOk(
		signal.Allocator{Capacity: 3, Channels: 2}.Float64(),
		[][]float64{
			{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		},
		expected{
			length: 3,
			data: [][]float64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
}

func TestAppend(t *testing.T) {
	type expected struct {
		length   int
		capacity int
		data     interface{}
	}
	testOk := func(appender signal.Signal, data interface{}, ex expected) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			var result signal.Signal
			switch a := appender.(type) {
			case signal.Int64:
				d := data.([][][]int64)
				for _, slice := range d {
					src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Int64(signal.MaxBitDepth)
					src.WriteInt64(slice)
					a = a.Append(src)
				}
				assertEqual(t, "slices", a.Data(), ex.data)
				result = a
			case signal.Int64Interleaved:
				d := data.([][]int64)
				for _, slice := range d {
					src := signal.Allocator{Channels: a.Channels(), Capacity: len(slice)}.Int64Interleaved(signal.MaxBitDepth)
					src.WriteInt64(slice)
					a = a.Append(src)
				}
				assertEqual(t, "slices", a.Data(), ex.data)
				result = a
			case signal.Float64:
				d := data.([][][]float64)
				for _, slice := range d {
					src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Float64()
					src.WriteFloat64(slice)
					a = a.Append(src)
				}
				assertEqual(t, "slices", a.Data(), ex.data)
				result = a
			default:
				t.Fatalf("unsupported append type %T", appender)
			}
			assertEqual(t, "length", result.Length(), ex.length)
			assertEqual(t, "capacity", result.Capacity(), ex.capacity)
		}
	}
	testPanic := func(appender signal.Signal, data signal.Signal) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			switch a := appender.(type) {
			case signal.Int64:
				d := data.(signal.Int64)
				assertPanic(t, func() {
					a.Append(d)
				})
			case signal.Int64Interleaved:
				d := data.(signal.Int64Interleaved)
				assertPanic(t, func() {
					a.Append(d)
				})
			case signal.Float64:
				d := data.(signal.Float64)
				assertPanic(t, func() {
					a.Append(d)
				})
			default:
				t.Fatalf("unsupported append panic type %T", appender)
			}
		}
	}

	t.Run("int64 single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		[][][]int64{
			{
				{1, 2},
				{1, 2},
			},
		},
		expected{
			capacity: 2,
			length:   2,
			data: [][]int64{
				{1, 2},
				{1, 2},
			},
		},
	))
	t.Run("int64 multiple slices", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		[][][]int64{
			{
				{1, 2},
				{1, 2},
			},
			{
				{3, 4},
				{3, 4},
			},
		},
		expected{
			capacity: 4,
			length:   4,
			data: [][]int64{
				{1, 2, 3, 4},
				{1, 2, 3, 4},
			},
		},
	))
	t.Run("int64 different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		signal.Allocator{Channels: 1, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
	t.Run("int64 different bit depth", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.BitDepth8),
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
	t.Run("int64interleaved single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64Interleaved(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 11, 12},
		},
		expected{
			capacity: 2,
			length:   2,
			data:     []int64{1, 2, 11, 12},
		},
	))
	t.Run("int64interleaved multiple slices", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64Interleaved(signal.MaxBitDepth),
		[][]int64{
			{1, 2, 11, 12},
			{3, 4, 13, 14},
		},
		expected{
			capacity: 4,
			length:   4,
			data:     []int64{1, 2, 11, 12, 3, 4, 13, 14},
		},
	))
	t.Run("int64interleaved different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64Interleaved(signal.MaxBitDepth),
		signal.Allocator{Channels: 1, Capacity: 2}.Int64Interleaved(signal.MaxBitDepth),
	))
	t.Run("int64interleaved different bit depth", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64Interleaved(signal.BitDepth8),
		signal.Allocator{Channels: 2, Capacity: 2}.Int64Interleaved(signal.MaxBitDepth),
	))
	t.Run("float64 single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][][]float64{
			{
				{1, 2},
				{11, 12},
			},
		},
		expected{
			length:   2,
			capacity: 2,
			data: [][]float64{
				{1, 2},
				{11, 12},
			},
		},
	))
	t.Run("float64 multiple slices", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][][]float64{
			{
				{1, 2},
				{11, 12},
			}, {
				{3, 4},
				{13, 14},
			},
		},
		expected{
			length:   4,
			capacity: 4,
			data: [][]float64{
				{1, 2, 3, 4},
				{11, 12, 13, 14},
			},
		},
	))
	t.Run("float64 different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		signal.Allocator{Channels: 1, Capacity: 2}.Float64(),
	))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}
