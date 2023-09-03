package signal

import (
	"golang.org/x/exp/constraints"
)

// Float is a digital signal represented with floating-point values.
type Float[T constraints.Float] struct {
	buffer[T]
}

// Slice the buffer with respect to channels.
func (b *Float[T]) Slice(start, end int) GenSig[T] {
	start = b.BufferIndex(0, start)
	end = b.BufferIndex(0, end)
	return &Float[T]{
		buffer: buffer[T]{
			channels: b.channels,
			data:     b.data[start:end],
		},
	}
}

// ReadFloat reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func Read[S, D SignalTypes](src GenSig[S], dst []D) int {
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
func ReadStriped[S, D SignalTypes](src GenSig[S], dst [][]D) (read int) {
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
func Write[S, D SignalTypes](src []S, dst GenSig[D]) int {
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
func WriteStriped[S, D SignalTypes](src [][]S, dst GenSig[D]) (written int) {
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
