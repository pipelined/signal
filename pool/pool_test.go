package pool_test

import (
	"testing"

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
		p := pool.New(test.numChannels, test.bufferSize)
		for i := 0; i < test.allocs; i++ {
			b := p.Alloc()
			if test.numChannels != b.NumChannels() {
				t.Fatalf("Invalid number of channels: %v expected: %v", b.NumChannels(), test.numChannels)
			}
			if test.bufferSize != b.Size() {
				t.Fatalf("Invalid buffer size: %v expected: %v", b.Size(), test.bufferSize)
			}
			p.Free(b)
		}
	}
}
