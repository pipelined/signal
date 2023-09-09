package signal

import (
	"strconv"
)

type (
	Allocator struct {
		Channels int
		Length   int
		Capacity int
	}
)

func Alloc[T SignalTypes](a Allocator) *Buffer[T] {
	return &Buffer[T]{
		buffer:   allocate[T](a),
		bitDepth: bitDepth(getBitDepth[T]()),
	}
}

// maxBitDebth returns a maximum bit debth for a given type, ie. 64 bits for int64 and uint64.
func getBitDepth[T SignalTypes]() BitDepth {
	switch any(new(T)).(type) {
	case *int8, *uint8:
		return BitDepth8
	case *int16, *uint16:
		return BitDepth16
	case *int32, *uint32, *float32:
		return BitDepth32
	default:
		return strconv.IntSize
	}
}
