package signal

import (
	"strconv"
)

type (
	// Allocator defines allocation parameters for signal buffers.
	Allocator struct {
		Channels int
		Length   int
		Capacity int
	}
)

// Alloc acllocates signal buffers based on provided type parameter. Type parameter also determines
// bit depth of the buffer, ie int8 will be 8 bit depth and int64 is a 64 bit depth.
func Alloc[T SignalTypes](a Allocator) *Buffer[T] {
	return &Buffer[T]{
		data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
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
