package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestAlloc(t *testing.T) {
	var (
		_ *signal.Buffer[float64] = signal.Alloc[float64](signal.Allocator{})
		_ *signal.Buffer[float32] = signal.Alloc[float32](signal.Allocator{})
		_ *signal.Buffer[int64]   = signal.Alloc[int64](signal.Allocator{})
	)
}
