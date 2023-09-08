package signal

import (
	"math"

	"golang.org/x/exp/constraints"
)

type buffer[T SignalTypes] struct {
	channels
	data   []T
	wrapFn wrapSampleFunc[T]
}

// wrapSampleFunc is needed to wrap the provided sample value.
type wrapSampleFunc[T SignalTypes] func(T) T

func wrapFloatSample[T constraints.Float](v T) T {
	return v
}

func allocate[T SignalTypes](a Allocator) buffer[T] {
	return buffer[T]{
		data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (b *buffer[T]) AppendSample(v T) {
	if len(b.data) == cap(b.data) {
		return
	}
	b.data = append(b.data, b.wrapFn(v))
}

// SetSample sets sample value for provided index.
func (b *buffer[T]) SetSample(i int, v T) {
	b.data[i] = b.wrapFn(v)
}

// Capacity returns capacity of a single channel.
func (b *buffer[T]) Capacity() int {
	if b.channels == 0 {
		return 0
	}
	return cap(b.data) / int(b.channels)
}

// Length returns length of a single channel.
func (b *buffer[T]) Length() int {
	if b.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(b.data)) / float64(b.channels)))
}

// Cap returns capacity of whole buffer.
func (b *buffer[T]) Cap() int {
	return cap(b.data)
}

// Len returns length of whole buffer.
func (b *buffer[T]) Len() int {
	return len(b.data)
}

// Sample returns signal value for provided sample index.
func (b *buffer[T]) Sample(i int) T {
	return b.data[i]
}

// Append appends [0:Length] samples from src to current buffer.
// Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (dst *buffer[D]) Append(src GenSig[D]) {
	mustSameChannels(dst.Channels(), src.Channels())
	offset := dst.Len()
	if dst.Cap() < dst.Len()+src.Len() {
		dst.data = append(dst.data, make([]D, src.Len())...)
	} else {
		dst.data = dst.data[:dst.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		dst.SetSample(i+offset, D(dst.wrapFn(src.Sample(i))))
	}
	alignCapacity(&dst.data, dst.Channels(), dst.Cap())
}

// Channel returns a single channel of signal buffer. It behaves exactly as
// Floating buffer, but Append and AppendSample cause panic.
func (b *buffer[T]) Channel(c int) C[T] {
	return C[T]{
		buffer:  *b,
		channel: c,
	}
}

func (b *buffer[T]) clear() {
	for i := range b.data {
		b.data[i] = 0
	}
}