package pool_test

import (
	"testing"

	"pipelined.dev/signal"
	"pipelined.dev/signal/pool"
)

func TestPool(t *testing.T) {
	tests := []struct {
		numChannels int
		bufferSize  int
		allocs      int
	}{
		{
			numChannels: 1,
			bufferSize:  512,
			allocs:      10,
		},
		{
			numChannels: 100,
			bufferSize:  1024,
			allocs:      1000,
		},
	}
	for _, test := range tests {
		p := pool.New(signal.Allocator{Channels: test.numChannels, Capacity: test.bufferSize})
		for i := 0; i < test.allocs; i++ {
			b := p.Alloc()
			if test.numChannels != b.Channels() {
				t.Fatalf("Invalid number of channels: %v expected: %v", b.Channels(), test.numChannels)
			}
			if test.bufferSize != b.Capacity() {
				t.Fatalf("Invalid buffer size: %v expected: %v", b.Capacity(), test.bufferSize)
			}
			p.Free(b)
		}
	}
}
