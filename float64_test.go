package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-09-11 21:34:57.813026 +0200 CEST m=+0.017265635

import (
	"testing"

	"pipelined.dev/signal"
)

func TestFloat64(t *testing.T) {
	t.Run("float64", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Float64()
		signal.WriteStripedFloat64(
			[][]float64{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		return testOk(
			signal.Allocator{
				Channels: 3,
				Capacity: 2,
			}.Float64().Append(input.Slice(1, 3)),
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
}
