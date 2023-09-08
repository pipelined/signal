package signal

import "golang.org/x/exp/constraints"

// Float is a digital signal represented with floating-point values.
type Integer[T constraints.Integer] struct {
	buffer[T]
	bitDepth
}

// Slice the buffer with respect to channels.
func (b *Integer[T]) Slice(start, end int) GenSig[T] {
	start = b.BufferIndex(0, start)
	end = b.BufferIndex(0, end)
	return &Integer[T]{
		buffer: buffer[T]{
			channels: b.channels,
			data:     b.data[start:end],
			wrapFn:   b.buffer.wrapFn,
		},
		bitDepth: b.bitDepth,
	}
}
