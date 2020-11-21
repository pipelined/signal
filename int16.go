package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-11-21 19:22:12.858605 +0100 CET m=+0.010727510

import (
	"math"
	"sync"
)

// Int16 is int16 Signed fixed-point signal.
type Int16 struct {
	buffer []int16
	channels
	bitDepth
}

// Int16 allocates a new sequential int16 signal buffer.
func (a Allocator) Int16(bd BitDepth) Signed {
	return &Int16{
		buffer:   make([]int16, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth16),
	}
}

// GetInt16 selects a new sequential int16 signal buffer.
// from the pool.
func (p *PoolAllocator) GetInt16(bd BitDepth) Signed {
	s := p.i16.Get().(*Int16)
	s.channels = channels(p.Channels)
	s.buffer = s.buffer[:p.Length*p.Channels]
	s.bitDepth = limitBitDepth(bd, BitDepth16)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s *Int16) AppendSample(value int64) {
	if len(s.buffer) == cap(s.buffer) {
		return
	}
	s.buffer = append(s.buffer, int16(s.BitDepth().SignedValue(value)))
}

// SetSample sets sample value for provided index.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Int16) SetSample(i int, value int64) {
	s.buffer[i] = int16(s.BitDepth().SignedValue(value))
}

// PutInt16 places signal buffer back to the pool. If a type of
// provided buffer isn't Int16 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (s *Int16) Free(p *PoolAllocator) {
	mustSameCapacity(s.Cap(), p.Channels*p.Capacity)
	for i := range s.buffer {
		s.buffer[i] = 0
	}
	p.i16.Put(s)
}

func int16pool(size int) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &Int16{
				buffer: make([]int16, 0, size),
			}
		},
	}
}

// Capacity returns capacity of a single channel.
func (s *Int16) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s *Int16) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s *Int16) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s *Int16) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and index.
func (s *Int16) Sample(i int) int64 {
	return int64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Signed buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (s *Int16) Append(src Signed) {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	offset := s.Len()
	if s.Cap() < s.Len()+src.Len() {
		s.buffer = append(s.buffer, make([]int16, src.Len())...)
	} else {
		s.buffer = s.buffer[:s.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		s.SetSample(i+offset, src.Sample(i))
	}
	alignCapacity(&s.buffer, s.Channels(), s.Cap())
}

// Slice slices buffer with respect to channels.
func (s *Int16) Slice(start, end int) Signed {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	return &Int16{
		channels: s.channels,
		buffer:   s.buffer[start:end],
		bitDepth: s.bitDepth,
	}
}

// ReadInt16 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadInt16(src Signed, dst []int16) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = int16(BitDepth16.SignedValue(src.Sample(i)))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedInt16 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedInt16(src Signed, dst [][]int16) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = int16(BitDepth16.SignedValue(src.Sample(src.BufferIndex(c, i))))
		}
	}
	return
}

// WriteInt16 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteInt16(src []int16, dst Signed) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, int64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedInt16 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedInt16(src [][]int16, dst Signed) (written int) {
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
				dst.SetSample(dst.BufferIndex(c, i), int64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}
