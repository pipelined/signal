package signal

import (
	"math"

	"golang.org/x/exp/constraints"
)

// F is a digital signal represented with floating-point values.
type F[T constraints.Float] struct {
	channels
	buffer []T
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s *F[T]) AppendSample(v T) {
	if len(s.buffer) == cap(s.buffer) {
		return
	}
	s.buffer = append(s.buffer, v)
}

// SetSample sets sample value for provided index.
func (s *F[T]) SetSample(i int, v T) {
	s.buffer[i] = v
}

// Capacity returns capacity of a single channel.
func (s *F[T]) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s *F[T]) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s *F[T]) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s *F[T]) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided sample index.
func (s *F[T]) Sample(i int) T {
	return s.buffer[i]
}

// Slice slices buffer with respect to channels.
func (s *F[T]) Slice(start, end int) *F[T] {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	return &F[T]{
		channels: s.channels,
		buffer:   s.buffer[start:end],
	}
}

// ReadFloat reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadFloat[S, D constraints.Float](src *F[S], dst []D) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = D(src.Sample(i))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedFloat reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedFloat[S, D constraints.Float](src *F[S], dst [][]D) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = D(src.Sample(src.BufferIndex(c, i)))
		}
	}
	return
}

// WriteFloat writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteFloat[S, D constraints.Float](src []S, dst *F[D]) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, D(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedFloat64 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedFloat[S, D constraints.Float](src [][]S, dst *F[D]) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for c := 0; c < dst.Channels(); c++ {
		for i := 0; i < written; i++ {
			if i < len(src[c]) {
				dst.SetSample(dst.BufferIndex(c, i), D(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Floating buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func Append[S, D constraints.Float](src *F[S], dst *F[D]) {
	mustSameChannels(dst.Channels(), src.Channels())
	offset := dst.Len()
	if dst.Cap() < dst.Len()+src.Len() {
		dst.buffer = append(dst.buffer, make([]D, src.Len())...)
	} else {
		dst.buffer = dst.buffer[:dst.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		dst.SetSample(i+offset, D(src.Sample(i)))
	}
	alignCapacity(&dst.buffer, dst.Channels(), dst.Cap())
}
