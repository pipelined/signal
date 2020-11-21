package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-11-21 19:22:12.866661 +0100 CET m=+0.018783912

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint32(t *testing.T) {
	t.Run("uint32", func() func(t *testing.T) {
		input := signal.Allocator{
			Channels: 3,
			Capacity: 3,
			Length:   3,
		}.Uint32(signal.BitDepth32)
		signal.WriteStripedUint32(
			[][]uint32{
				{},
				{1, 2, 3},
				{11, 12, 13, 14},
			},
			input,
		)
		result := signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Uint32(signal.BitDepth32)
		result.Append(input.Slice(1, 3))
		return testOk(
			result,
			expected{
				length:   2,
				capacity: 2,
				data: [][]uint32{
					{0, 0},
					{2, 3},
					{12, 13},
				},
			},
		)
	}())
}
