package signal_test

import (
	"math"
	"testing"

	"pipelined.dev/signal"
)

func TestSlice(t *testing.T) {
	alloc := signal.Allocator{
		Channels: 3,
		Capacity: 3,
		Length:   3,
	}
	t.Run("floating", func() func(t *testing.T) {
		input := alloc.Float64()
		signal.WriteStripedFloat64(
			[][]float64{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk(
			signal.Slice(input, 1, 3),
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

	t.Run("signed", func() func(t *testing.T) {
		input := alloc.Int64(signal.MaxBitDepth)
		signal.WriteStripedInt64(
			[][]int64{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk(
			signal.Slice(input, 1, 3),
			expected{
				length:   2,
				capacity: 2,
				data: [][]int64{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
	t.Run("unsigned", func() func(t *testing.T) {
		input := alloc.Uint64(signal.MaxBitDepth)
		signal.WriteStripedUint64(
			[][]uint64{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk(
			signal.Slice(input, 1, 3),
			expected{
				length:   2,
				capacity: 2,
				data: [][]uint64{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
}

func TestConversions(t *testing.T) {
	t.Skip()
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}

	// floating buf
	floating := alloc.Float64()
	floats := [][]float64{
		{-1, 0, 1},
	}
	signal.WriteStripedFloat64(floats, floating)

	// signed buf
	signed := alloc.Int64(signal.MaxBitDepth)
	ints := [][]int64{
		{math.MinInt64, 0, math.MaxInt64},
	}
	signal.WriteStripedInt64(ints, signed)

	// unsigned buf
	unsigned := alloc.Uint64(signal.MaxBitDepth)
	uints := [][]uint64{
		{0, math.MaxInt64 + 1, math.MaxUint64},
	}
	signal.WriteStripedUint64(uints, unsigned)

	t.Run("floating", func() func(*testing.T) {
		output := alloc.Float64()
		return func(t *testing.T) {
			signal.AsFloating(floating, output)
			assertEqual(t, "floating ", result(output), floats)
			signal.AsFloating(signed, output)
			assertEqual(t, "signed", result(output), floats)
			signal.AsFloating(unsigned, output)
			assertEqual(t, "unsigned", result(output), floats)
		}
	}())
	t.Run("signed", func() func(*testing.T) {
		output := alloc.Int64(signal.MaxBitDepth)
		return func(t *testing.T) {
			signal.AsSigned(floating, output)
			assertEqual(t, "floating ", result(output), ints)
			signal.AsSigned(signed, output)
			assertEqual(t, "signed", result(output), ints)
			signal.AsSigned(unsigned, output)
			assertEqual(t, "unsigned", result(output), ints)
		}
	}())
	t.Run("unsigned", func() func(*testing.T) {
		output := alloc.Uint64(signal.MaxBitDepth)
		return func(t *testing.T) {
			signal.AsUnsigned(floating, output)
			assertEqual(t, "floating ", result(output), uints)
			signal.AsUnsigned(signed, output)
			assertEqual(t, "signed", result(output), uints)
			signal.AsUnsigned(unsigned, output)
			assertEqual(t, "unsigned", result(output), uints)
		}
	}())
}
