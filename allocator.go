package signal

import (
	"fmt"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Allocator provides allocation of various signal buffers.
type Allocator struct {
	Channels int
	Length   int
	Capacity int
}

// AllocFloat allocates a new float signal buffer.
func AllocFloat[T constraints.Float](a Allocator) *Float[T] {
	return &Float[T]{
		buffer: buffer[T]{data: make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
			channels: channels(a.Channels),
		},
	}
}

// AllocInt allocates a new integer signal buffer. If bd exceeds maximum
// bit debth for a given type, function will panic.
func AllocInt[T constraints.Integer](a Allocator, bd BitDepth) *Integer[T] {
	return &Integer[T]{
		buffer: buffer[T]{
			data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
			channels: channels(a.Channels),
		},
		bitDepth: limitBitDepth(bd, maxBitDebth[T]()),
	}
}

// defaultBitDepth limits bit depth value to max and returns max if it is 0.
func limitBitDepth(b, max BitDepth) bitDepth {
	if b == 0 {
		return bitDepth(max)
	}
	if b > max {
		panic(fmt.Sprintf("maximum bit debth: %v got: %v", max, b))
	}
	return bitDepth(b)
}

// maxBitDebth returns a maximum bit debth for a given type, ie. 64 bits for int64 and uint64.
func maxBitDebth[T constraints.Integer]() BitDepth {
	return BitDepth(unsafe.Sizeof(new(T)) * 8)
}
