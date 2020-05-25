package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-25 04:56:22.160829 +0200 CEST m=+0.020251148

import (
	"testing"

	"pipelined.dev/signal"
)

func TestInt32(t *testing.T) {
	t.Run("int32", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.Int32(signal.BitDepth32).
			Append(signal.WriteStripedInt32(
				[][]int32{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.Int32(signal.BitDepth32)),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]int32{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}