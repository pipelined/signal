package signal_test

import (
	"testing"

	"golang.org/x/exp/constraints"
	"pipelined.dev/signal"
)

func TestFloat(t *testing.T) {
	t.Run("float64", testGenericFloat[float64]())
	t.Run("float32", testGenericFloat[float32]())
}

type expectedGeneric[T constraints.Float] struct {
	length   int
	capacity int
	data     [][]T
}

func testGenericFloat[T constraints.Float]() func(t *testing.T) {
	input := signal.Alloc[T](signal.Allocator{
		Channels: 3,
		Capacity: 3,
		Length:   3,
	})
	signal.WriteStriped(
		[][]T{
			{},
			{1, 2, 3},
			{11, 12, 13, 14},
		},
		input,
	)
	r := signal.Alloc[T](signal.Allocator{
		Channels: 3,
		Capacity: 2,
	})
	r.Append(input.Slice(1, 3))
	ex := expectedGeneric[T]{
		length:   2,
		capacity: 2,
		data: [][]T{
			{0, 0},
			{2, 3},
			{12, 13},
		},
	}
	return func(t *testing.T) {
		t.Helper()
		if ex.capacity != 0 {
			assertEqual(t, "capacity", r.Capacity(), ex.capacity)
		}
		if ex.length != 0 {
			assertEqual(t, "length", r.Length(), ex.length)
		}
		assertEqual(t, "slices", resultGeneric[T](r), ex.data)
	}
}

func resultGeneric[T constraints.Float](src *signal.Buffer[T]) [][]T {
	result := make([][]T, src.Channels())
	for i := range result {
		result[i] = make([]T, src.Length())
	}
	signal.ReadStriped(src, result)
	return result
}

func TestSlice(t *testing.T) {
	alloc := signal.Allocator{
		Channels: 3,
		Capacity: 3,
		Length:   3,
	}
	t.Run("floating", func() func(t *testing.T) {
		input := signal.Alloc[float64](alloc)
		signal.WriteStriped[float64, float64](
			[][]float64{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk[float64](
			input.Slice(1, 3),
			expected{
				length:   2,
				capacity: 2,
				data: [][]float64{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())

	t.Run("slice same size", func() func(t *testing.T) {
		input := signal.Alloc[float64](alloc)
		signal.WriteStriped[float64, float64](
			[][]float64{
				{},
				{1, 2},
				{11, 12, 13},
			},
			input,
		)
		return testOk[float64](
			input.Slice(0, 3),
			expected{
				length:   3,
				capacity: 3,
				data: [][]float64{
					{0, 0, 0},
					{1, 2, 0},
					{11, 12, 13},
				},
			},
		)
	}())
}
