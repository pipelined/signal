package signal_test

import (
	"testing"

	"golang.org/x/exp/constraints"
	"pipelined.dev/signal"
)

type expectedAllocation struct {
	channels int
	length   int
	capacity int
	signal.BitDepth
}

func TestPool(t *testing.T) {
	testOk := func(t *testing.T, allocs int, channels, length, capacity int) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			alloc := signal.Allocator{channels, length, capacity}
			for i := 0; i < allocs; i++ {
				// floating
				fp := signal.GetPool[float64](alloc)
				fp.Put(testFloat(t, channels, length, capacity, fp.Get()))
				// signed
				ip := signal.GetPool[int32](alloc)
				// ip.Release(testInteger(t, channels, length, capacity, ip.Get(signal.BitDepth8), signal.BitDepth8))
				ip.Put(testInteger(t, channels, length, capacity, ip.Get(), signal.BitDepth32))
			}
		}
	}

	t.Run("empty allocs",
		testOk(t, 10, 0, 0, 0),
	)
	t.Run("10 allocs",
		testOk(t, 10, 1, 0, 512),
	)
	t.Run("10 allocs length",
		testOk(t, 10, 2, 256, 512),
	)
	t.Run("100 allocs",
		testOk(t, 100, 100, 0, 512),
	)
}

func testFloat[T constraints.Float](t *testing.T, channels, length, capacity int, s *signal.Buffer[T]) *signal.Buffer[T] {
	t.Helper()
	assertAllocation(
		t,
		s,
		expectedAllocation{
			channels: channels,
			length:   length,
			capacity: capacity,
		})
	s.AppendSample(1)
	return s
}

func testInteger[T constraints.Signed](t *testing.T, channels, length, capacity int, s *signal.Buffer[T], mbd signal.BitDepth) *signal.Buffer[T] {
	t.Helper()
	assertAllocation(
		t,
		s,
		expectedAllocation{
			channels: channels,
			length:   length,
			capacity: capacity,
			BitDepth: mbd,
		})
	s.AppendSample(1)
	if mbd != s.BitDepth() {
		t.Fatalf("Invalid Buffer bit depth: %v expected: %v", s.BitDepth(), mbd)
	}
	return s
}

func assertAllocation[T signal.SignalTypes](t *testing.T, s *signal.Buffer[T], e expectedAllocation) {
	t.Helper()
	if e.channels != s.Channels() {
		t.Fatalf("Invalid number of channels: %v expected: %v", s.Channels(), e.channels)
	}
	if e.length != s.Length() {
		t.Fatalf("Invalid Buffer length: %v expected: %v", s.Length(), e.length)
	}
	if e.capacity != s.Capacity() {
		t.Fatalf("Invalid Buffer capacity: %v expected: %v", s.Capacity(), e.capacity)
	}
}
