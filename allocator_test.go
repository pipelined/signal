package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func TestAlloc(t *testing.T) {
	var (
		_ signal.GenSig[float64] = signal.AllocFloat[float64](signal.Allocator{})
		_ signal.GenSig[float32] = signal.AllocFloat[float32](signal.Allocator{})
		_ signal.GenSig[int64]   = signal.AllocInteger[int64](signal.Allocator{}, signal.BitDepth16)
	)
}

func TestBitDepth(t *testing.T) {
	assertPanic(t, func() {
		signal.AllocInteger[int16](signal.Allocator{1, 1, 1}, signal.MaxBitDepth)
	})
}
