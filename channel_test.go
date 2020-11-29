package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestChannel(t *testing.T) {
	testAppendFloating := func() func(*testing.T) {
		result := signal.Allocator{
			Channels: 3,
			Capacity: 3,
		}.Float64()
		c := result.Channel(1)
		for i := 1; i < 4; i++ {
			c.AppendSample(float64(i))
		}
		return testOk(
			result,
			expected{
				length:   3,
				capacity: 3,
				data: [][]float64{
					{0, 0, 0},
					{1, 2, 3},
					{0, 0, 0},
				},
			},
		)
	}

	testSetFloating := func() func(*testing.T) {
		result := signal.Allocator{
			Channels: 3,
			Length:   3,
			Capacity: 3,
		}.Float64()
		c := result.Channel(1)
		for i := 0; i < c.Len(); i++ {
			c.SetSample(i, float64(i+1))
		}
		return testOk(
			result,
			expected{
				length:   3,
				capacity: 3,
				data: [][]float64{
					{0, 0, 0},
					{1, 2, 3},
					{0, 0, 0},
				},
			},
		)
	}

	t.Run("append floating", testAppendFloating())
	t.Run("set floating", testSetFloating())
}
