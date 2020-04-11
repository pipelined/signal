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

func TestWriteInt64(t *testing.T) {
	type expected struct {
		length int
		data   [][]int64
	}
	testOk := func(s signal.Int64, ints [][]int, ex expected) func(t *testing.T) {
		return func(t *testing.T) {
			s.WriteInt(ints)
			assertEqual(t, "length", s.Length(), ex.length)
			assertEqual(t, "slices", s.Data(), ex.data)
		}
	}

	t.Run("full buffer", testOk(
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
	t.Run("short buffer", testOk(
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
	t.Run("long buffer", testOk(
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
	t.Run("8-bits overflow", testOk(
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
	t.Run("16-bits overflow", testOk(
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
}

func TestAppendFloat64(t *testing.T) {
	type expected struct {
		length   int
		capacity int
		data     [][]float64
	}
	testOk := func(s signal.Float64, slices [][][]float64, ex expected) func(*testing.T) {
		return func(t *testing.T) {
			for _, slice := range slices {
				src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Float64()
				src.WriteFloat64(slice)
				s = s.Append(src)
			}
			assertEqual(t, "slices", s.Data(), ex.data)
			assertEqual(t, "length", s.Length(), ex.length)
			assertEqual(t, "capacity", s.Capacity(), ex.capacity)
		}
	}
	testPanic := func(s signal.Float64, slice [][]float64) func(*testing.T) {
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
	t.Run("multiple slices", testOk(
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
	t.Run("different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		[][]float64{
			{1, 2},
		},
	))
}

func TestAppendInt64(t *testing.T) {
	testOk := func(s signal.Int64, expected [][]int64, slices ...[][]int64) func(*testing.T) {
		return func(t *testing.T) {
			for _, slice := range slices {
				src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Int64(signal.MaxBitDepth)
				src.WriteInt64(slice)
				s = s.Append(src)
			}
			assertEqual(t, "slices", s.Data(), expected)
		}
	}
	testPanic := func(s signal.Int64, slice [][]int64) func(*testing.T) {
		return func(t *testing.T) {
			src := signal.Allocator{Channels: len(slice), Capacity: len(slice[0])}.Int64(signal.MaxBitDepth)
			src.WriteInt64(slice)
			assertPanic(t, func() {
				s = s.Append(src)
			})
		}
	}

	t.Run("single slice", testOk(
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
	t.Run("multiple slices", testOk(
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
	t.Run("different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		[][]int64{
			{1, 2},
		},
	))
	t.Run("different bit depth", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.BitDepth8),
		[][]int64{
			{1, 2},
			{1, 2},
		},
	))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\ngot: \t%+v \nexpected: \t%+v", name, result, expected)
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
