package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-25 04:56:22.154712 +0200 CEST m=+0.014133941

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint32(t *testing.T) {
	t.Run("uint32", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Uint32(signal.BitDepth32).
			Append(signal.WriteStripedUint32(
				[][]uint32{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Uint32(signal.BitDepth32)),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]uint32{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}