package signal_test

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 19:42:20.042445 +0100 CET m=+0.017504438

import (
	"testing"

	"pipelined.dev/signal"
)

func TestUnsignedChannel(t *testing.T) {
	testUnsigned := func() func(*testing.T) {
		result := signal.Allocator{
			Channels: 3,
			Length:   3,
			Capacity: 3,
		}.Uint64(signal.BitDepth32)
		channel := 1
		c := result.Channel(channel).Slice(0, 2)
		for i := 0; i < c.Len(); i++ {
			c.SetSample(i, uint64(i+1))
		}
		return func(t *testing.T) {
			assertEqual(t, "channels", c.Channels(), 1)
			assertEqual(t, "length", c.Len(), c.Length())
			assertEqual(t, "capacity", c.Cap(), c.Capacity())
			for i := 0; i < c.Cap(); i++ {
				assertEqual(t, "index", c.Sample(i), result.Sample(c.BufferIndex(channel, i)))
			}
		}
	}
	testPanic := func() func(*testing.T) {
		result := signal.Allocator{
			Channels: 3,
			Length:   3,
			Capacity: 3,
		}.Uint64(signal.BitDepth32)
		c := result.Channel(1)
		return func(t *testing.T) {
			assertPanic(t, func() {
				c.Append(nil)
			})
			assertPanic(t, func() {
				c.AppendSample(0)
			})
			assertPanic(t, func() {
				c.Channel(0)
			})
			assertPanic(t, func() {
				c.Free(nil)
			})
		}
	}

	t.Run("Unsigned channel", testUnsigned())
	t.Run("panic Unsigned channel", testPanic())
}