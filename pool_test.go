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
	testOk := func(t *testing.T, allocs int, p *signal.Pool) func(t *testing.T) {
		return func(t *testing.T) {
			for i := 0; i < allocs; i++ {
				// floating
				p.PutFloat32(testFloating(t, p, p.GetFloat32()))
				p.PutFloat64(testFloating(t, p, p.GetFloat64()))
				// signed
				p.PutInt8(testSigned(t, p, p.GetInt8(signal.MaxBitDepth), signal.BitDepth8))
				p.PutInt16(testSigned(t, p, p.GetInt16(signal.MaxBitDepth), signal.BitDepth16))
				p.PutInt32(testSigned(t, p, p.GetInt32(signal.MaxBitDepth), signal.BitDepth32))
				p.PutInt64(testSigned(t, p, p.GetInt64(signal.MaxBitDepth), signal.BitDepth64))
				// unsigned
				p.PutUint8(testUnsigned(t, p, p.GetUint8(signal.MaxBitDepth), signal.BitDepth8))
				p.PutUint16(testUnsigned(t, p, p.GetUint16(signal.MaxBitDepth), signal.BitDepth16))
				p.PutUint32(testUnsigned(t, p, p.GetUint32(signal.MaxBitDepth), signal.BitDepth32))
				p.PutUint64(testUnsigned(t, p, p.GetUint64(signal.MaxBitDepth), signal.BitDepth64))
			}
		}
	}

	t.Run("nil pool",
		testOk(t, 10, nil),
	)
	t.Run("empty allocs",
		testOk(t, 10, signal.Allocator{}.Pool()),
	)
	t.Run("10 allocs",
		testOk(t, 10, signal.Allocator{
			Channels: 1,
			Capacity: 512,
		}.Pool()),
	)
	t.Run("10 allocs length",
		testOk(t, 10, signal.Allocator{
			Channels: 2,
			Length:   256,
			Capacity: 512,
		}.Pool()),
	)
	t.Run("100 allocs",
		testOk(t, 100, signal.Allocator{
			Channels: 100,
			Capacity: 512,
		}.Pool()),
	)
}

func testSigned(t *testing.T, p *signal.Pool, s signal.Signed, mbd signal.BitDepth) signal.Signed {
	t.Helper()
	if s == nil {
		return s
	}
	a := p.Allocator()
	assertAllocation(
		t,
		s,
		expectedAllocation{
			channels: a.Channels,
			length:   a.Length,
			capacity: a.Capacity,
			BitDepth: mbd,
		})
	s.AppendSample(1)
	return s
}

func testUnsigned(t *testing.T, p *signal.Pool, s signal.Unsigned, mbd signal.BitDepth) signal.Unsigned {
	t.Helper()
	if s == nil {
		return s
	}
	a := p.Allocator()
	assertAllocation(
		t,
		s,
		expectedAllocation{
			channels: a.Channels,
			length:   a.Length,
			capacity: a.Capacity,
			BitDepth: mbd,
		})
	s.AppendSample(1)
	return s
}

func testFloating(t *testing.T, p *signal.Pool, s signal.Floating) signal.Floating {
	t.Helper()
	if s == nil {
		return s
	}
	a := p.Allocator()
	assertAllocation(
		t,
		s,
		expectedAllocation{
			channels: a.Channels,
			length:   a.Length,
			capacity: a.Capacity,
		})
	s.AppendSample(1)
	return s
}

func assertAllocation(t *testing.T, s signal.Signal, e expectedAllocation) {
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
