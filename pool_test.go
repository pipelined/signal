package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

type expectedAllocation struct {
	channels int
	length   int
	capacity int
	signal.BitDepth
}

func TestPool(t *testing.T) {
	assertAllocation := func(t *testing.T, s signal.Signal, e expectedAllocation) {
		t.Helper()
		if e.channels != s.Channels() {
			t.Fatalf("Invalid number of channels: %v expected: %v", s.Channels(), e.channels)
		}
		if e.length != s.Length() {
			t.Fatalf("Invalid buffer length: %v expected: %v", s.Length(), e.length)
		}
		if e.capacity != s.Capacity() {
			t.Fatalf("Invalid buffer capacity: %v expected: %v", s.Capacity(), e.capacity)
		}
		if f, ok := s.(signal.Fixed); ok {
			if e.BitDepth != f.BitDepth() {
				t.Fatalf("Invalid buffer bit depth: %v expected: %v", f.BitDepth(), e.BitDepth)
			}
		}
	}
	testSigned := func(t *testing.T, channels, length, capacity int, s signal.Signed, mbd signal.BitDepth) signal.Signed {
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
		return s
	}
	testUnsigned := func(t *testing.T, channels, length, capacity int, s signal.Unsigned, mbd signal.BitDepth) signal.Unsigned {
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
		return s
	}
	testFloating := func(t *testing.T, channels, length, capacity int, s signal.Floating) signal.Floating {
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

	testOk := func(t *testing.T, allocs int, channels, length, capacity int) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			p := signal.GetPoolAllocator(channels, length, capacity)
			for i := 0; i < allocs; i++ {
				// floating
				testFloating(t, channels, length, capacity, p.GetFloat32()).Free(p)
				testFloating(t, channels, length, capacity, p.GetFloat64()).Free(p)
				// signed
				testSigned(t, channels, length, capacity, p.GetInt8(signal.MaxBitDepth), signal.BitDepth8).Free(p)
				testSigned(t, channels, length, capacity, p.GetInt16(signal.MaxBitDepth), signal.BitDepth16).Free(p)
				testSigned(t, channels, length, capacity, p.GetInt32(signal.MaxBitDepth), signal.BitDepth32).Free(p)
				testSigned(t, channels, length, capacity, p.GetInt64(signal.MaxBitDepth), signal.BitDepth64).Free(p)
				// unsigned
				testUnsigned(t, channels, length, capacity, p.GetUint8(signal.MaxBitDepth), signal.BitDepth8).Free(p)
				testUnsigned(t, channels, length, capacity, p.GetUint16(signal.MaxBitDepth), signal.BitDepth16).Free(p)
				testUnsigned(t, channels, length, capacity, p.GetUint32(signal.MaxBitDepth), signal.BitDepth32).Free(p)
				testUnsigned(t, channels, length, capacity, p.GetUint64(signal.MaxBitDepth), signal.BitDepth64).Free(p)
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

func TestGetPool(t *testing.T) {
	// a1 := signal.Allocator{
	// 	Channels: 10,
	// 	Capacity: 512,
	// }
	// a2 := signal.Allocator{
	// 	Channels: 10,
	// 	Length:   512,
	// 	Capacity: 512,
	// }

	p1 := signal.GetPoolAllocator(10, 0, 512)
	p2 := signal.GetPoolAllocator(10, 512, 512)
	if p1 != p2 {
		t.Fatal("p1 must be equal to p2")
	}

	signal.ClearPoolAllocatorCache()
	p3 := signal.GetPoolAllocator(10, 0, 512)
	if p1 == p3 {
		t.Fatal("p1 must not be equal to p3")
	}
}
