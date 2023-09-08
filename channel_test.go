package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func testChannel[T signal.SignalTypes](result signal.GenSig[T]) func(*testing.T) {
	channel := 1
	c := result.Channel(channel).Slice(0, 2)
	for i := 0; i < c.Len(); i++ {
		c.SetSample(i, T(i+1))
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

func TestChannel(t *testing.T) {
	t.Run("Floating channel", testChannel(signal.AllocFloat[float64](signal.Allocator{
		Channels: 3,
		Length:   3,
		Capacity: 3,
	})))

	t.Run("Integer channel", testChannel(signal.AllocInteger[int64](signal.Allocator{
		Channels: 3,
		Length:   3,
		Capacity: 3,
	}, signal.BitDepth32)))
}
