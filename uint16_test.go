package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 20:01:52.991988 +0100 CET m=+0.017880682

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint16(t *testing.T) {
	t.Run("uint16", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Uint16(signal.BitDepth16)
		signal.WriteStripedUint16(
			[][]uint16{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		result := signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Uint16(signal.BitDepth16)
		result.Append(input.Slice(1, 3))
		return testOk(
			result,
			expected{
				length:   2,
				capacity: 2,
				data: [][]uint16{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
}
