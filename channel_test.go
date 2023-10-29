package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func testChannel[T signal.SignalTypes](result *signal.Buffer[T]) func(*testing.T) {
	channel := 1
	c := result.Channel(channel)
	for i := 0; i < c.Length(); i++ {
		c.SetSample(i, T(i+1))
	}
	return func(t *testing.T) {
		assertEqual(t, "channels", c.Channels(), 1)
		assertEqual(t, "length", result.Length(), c.Length())
		assertEqual(t, "capacity", result.Capacity(), c.Capacity())
		for i := 0; i < c.Capacity(); i++ {
			assertEqual(t, "index", c.Sample(i), result.Sample(c.BufferIndex(channel, i)))
		}
	}
}

func TestChannel(t *testing.T) {
	t.Run("Floating channel", testChannel(signal.Alloc[float64](signal.Allocator{
		Channels: 3,
		Length:   3,
		Capacity: 3,
	})))

	t.Run("Integer channel", testChannel(signal.Alloc[int64](signal.Allocator{
		Channels: 3,
		Length:   3,
		Capacity: 3,
	})))
}
