package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestChannel(t *testing.T) {
	testFloating := func() func(*testing.T) {
		result := signal.AllocFloat[float64](signal.Allocator{
			Channels: 3,
			Length:   3,
			Capacity: 3,
		})
		channel := 1
		c := result.Channel(channel).Slice(0, 2)
		for i := 0; i < c.Len(); i++ {
			c.SetSample(i, float64(i+1))
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

	t.Run("Floating channel", testFloating())
}
