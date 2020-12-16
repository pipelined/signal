package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 20:01:52.987382 +0100 CET m=+0.013275530

import (
	"testing"

	"pipelined.dev/signal"
)

func TestFloat32(t *testing.T) {
	t.Run("float32", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Float32()
		signal.WriteStripedFloat32(
			[][]float32{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		result := signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Float32()
		result.Append(input.Slice(1, 3))
		return testOk(
			result,
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
