package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-25 04:56:22.153274 +0200 CEST m=+0.012696062

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUint16(t *testing.T) {
	t.Run("uint16", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Uint16(signal.BitDepth16).
			Append(signal.WriteStripedUint16(
				[][]uint16{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Uint16(signal.BitDepth16)),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]uint16{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}