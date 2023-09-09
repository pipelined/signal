package signal

import (
	"math"
)

// Buffer is a buffer that contains digital signal of given type.
// Each type is associated with certain bit depth, ie: int8 is 8 bits, float32 is 32 bits.
type Buffer[T SignalTypes] struct {
	channels
	data []T
	bitDepth
}

// Slice the Buffer with respect to channels.
func (b *Buffer[T]) Slice(start, end int) *Buffer[T] {
	start = b.BufferIndex(0, start)
	end = b.BufferIndex(0, end)
	return &Buffer[T]{
		channels: b.channels,
		data:     b.data[start:end],
		bitDepth: b.bitDepth,
	}
}

// AppendSample appends sample at the end of the Buffer.
// Sample is not appended if Buffer capacity is reached.
func (b *Buffer[T]) AppendSample(v T) {
	if len(b.data) == cap(b.data) {
		return
	}
	b.data = append(b.data, v)
}

// SetSample sets sample value for provided index.
func (b *Buffer[T]) SetSample(i int, v T) {
	b.data[i] = v
}

// Capacity returns capacity of a single channel.
func (b *Buffer[T]) Capacity() int {
	if b.channels == 0 {
		return 0
	}
	return cap(b.data) / int(b.channels)
}

// Length returns length of a single channel.
func (b *Buffer[T]) Length() int {
	if b.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(b.data)) / float64(b.channels)))
}

// Cap returns capacity of whole Buffer.
func (b *Buffer[T]) Cap() int {
	return cap(b.data)
}

// Len returns length of whole Buffer.
func (b *Buffer[T]) Len() int {
	return len(b.data)
}

// Sample returns signal value for provided sample index.
func (b *Buffer[T]) Sample(i int) T {
	return b.data[i]
}

// Append appends [0:Length] samples from src to current Buffer.
// Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (dst *Buffer[D]) Append(src *Buffer[D]) {
	mustSameChannels(dst.Channels(), src.Channels())
	offset := dst.Len()
	if dst.Cap() < dst.Len()+src.Len() {
		dst.data = append(dst.data, make([]D, src.Len())...)
	} else {
		dst.data = dst.data[:dst.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		dst.SetSample(i+offset, src.Sample(i))
	}
	alignCapacity(&dst.data, dst.Channels(), dst.Cap())
}

// Channel returns a single channel of signal Buffer. It behaves exactly as
// Floating Buffer, but Append and AppendSample cause panic.
func (b *Buffer[T]) Channel(c int) C[T] {
	return C[T]{
		Buffer:  *b,
		channel: c,
	}
}

func (b *Buffer[T]) clear() {
	for i := range b.data {
		b.data[i] = 0
	}
}
