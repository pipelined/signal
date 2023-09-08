package signal

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

type (
	Allocator struct {
		Channels int
		Length   int
		Capacity int
	}
)

// AllocFloat allocates a new float signal buffer.
func AllocFloat[T constraints.Float](a Allocator) *Float[T] {
	buffer := allocate[T](a)
	buffer.wrapFn = wrapFloatSample[T]
	return &Float[T]{
		buffer: buffer,
	}
}

// AllocInteger allocates a new integer signal buffer. If bd exceeds maximum
// bit debth for a given type, function will panic.
func AllocInteger[T constraints.Integer](a Allocator, bd BitDepth) *Integer[T] {
	buffer := allocate[T](a)
	buffer.wrapFn = wrapIntegerSample[T](bd)
	return &Integer[T]{
		buffer:   buffer,
		bitDepth: limitBitDepth[T](bd),
	}
}

// defaultBitDepth limits bit depth value to max and returns max if it is 0.
func limitBitDepth[T constraints.Integer](bd BitDepth) bitDepth {
	max := maxBitDebth[T]()
	if bd == 0 {
		return bitDepth(max)
	}
	if bd > max {
		panic(fmt.Sprintf("maximum bit debth: %v got: %v", max, bd))
	}
	return bitDepth(bd)
}

// maxBitDebth returns a maximum bit debth for a given type, ie. 64 bits for int64 and uint64.
func maxBitDebth[T constraints.Integer]() BitDepth {
	switch any(new(T)).(type) {
	case *int8:
		return BitDepth8
	case *int16:
		return BitDepth16
	case *int32:
		return BitDepth32
	default:
		return strconv.IntSize
	}
}
