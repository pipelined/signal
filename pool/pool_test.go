package pool_test

import (
	"testing"

	"pipelined.dev/signal"
	"pipelined.dev/signal/pool"
)

type expected struct {
	capacity int
	channels int
}

func TestPool(t *testing.T) {
	testSigned := func(t *testing.T, allocs int, alloc signal.Allocator) func(t *testing.T) {
		return func(t *testing.T) {
			p := pool.Signed(
				alloc.Int64,
				signal.BitDepth24,
			)
			e := expected{capacity: alloc.Capacity, channels: alloc.Channels}
			for i := 0; i < allocs; i++ {
				b := p.Get()
				assertResult(t, b, e)
				b.AppendSample(1)
				p.Put(b)
			}
		}
	}
	testUnsigned := func(t *testing.T, allocs int, alloc signal.Allocator) func(t *testing.T) {
		return func(t *testing.T) {
			p := pool.Unsigned(
				alloc.Uint64,
				signal.BitDepth24,
			)
			e := expected{capacity: alloc.Capacity, channels: alloc.Channels}
			for i := 0; i < allocs; i++ {
				b := p.Get()
				assertResult(t, b, e)
				b.AppendSample(1)
				p.Put(b)
			}
		}
	}
	testFloating := func(t *testing.T, allocs int, alloc signal.Allocator) func(t *testing.T) {
		return func(t *testing.T) {
			p := pool.Floating(
				alloc.Float64,
			)
			e := expected{capacity: alloc.Capacity, channels: alloc.Channels}
			for i := 0; i < allocs; i++ {
				b := p.Get()
				assertResult(t, b, e)
				b.AppendSample(1)
				p.Put(b)
			}
		}
	}

	t.Run("signed 10 allocs",
		testSigned(t, 10, signal.Allocator{
			Channels: 1,
			Capacity: 512,
		}),
	)
	t.Run("signed 100 allocs",
		testSigned(t, 100, signal.Allocator{
			Channels: 100,
			Capacity: 512,
		}),
	)
	t.Run("signed 10 allocs",
		testUnsigned(t, 10, signal.Allocator{
			Channels: 1,
			Capacity: 512,
		}),
	)
	t.Run("signed 100 allocs",
		testUnsigned(t, 100, signal.Allocator{
			Channels: 100,
			Capacity: 512,
		}),
	)
	t.Run("signed 10 allocs",
		testFloating(t, 10, signal.Allocator{
			Channels: 1,
			Capacity: 512,
		}),
	)
	t.Run("signed 100 allocs",
		testFloating(t, 100, signal.Allocator{
			Channels: 100,
			Capacity: 512,
		}),
	)
}

func assertResult(t *testing.T, s signal.Signal, e expected) {
	t.Helper()
	if e.channels != s.Channels() {
		t.Fatalf("Invalid number of channels: %v expected: %v", s.Channels(), e.channels)
	}
	if e.capacity != s.Capacity() {
		t.Fatalf("Invalid buffer capacity: %v expected: %v", s.Capacity(), e.capacity)
	}
	if s.Length() != 0 {
		t.Fatalf("non-zero buffer length: %v", s.Length())
	}
}
