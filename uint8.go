package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 20:52:36.092991 +0100 CET m=+0.012766905

import (
	"math"
	"sync"
)

// u8 is uint8 Unsigned fixed-point signal.
type u8 struct {
	buffer []uint8
	channels
	bitDepth
}

// Uint8 allocates a new sequential uint8 signal buffer.
func (a Allocator) Uint8(bd BitDepth) Unsigned {
	return &u8{
		buffer:   make([]uint8, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth8),
	}
}

// Uint8 selects a new sequential uint8 signal buffer.
// from the pool.
func (p *PoolAllocator) Uint8(bd BitDepth) Unsigned {
	s := p.u8.Get().(*u8)
	s.channels = channels(p.Channels)
	s.buffer = s.buffer[:p.Length*p.Channels]
	s.bitDepth = limitBitDepth(bd, BitDepth8)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s *u8) AppendSample(value uint64) {
	if len(s.buffer) == cap(s.buffer) {
		return
	}
	s.buffer = append(s.buffer, uint8(s.BitDepth().UnsignedValue(value)))
}

// SetSample sets sample value for provided index.
// Sample values are capped by maximum value of the buffer bit depth.
func (s u8) SetSample(i int, value uint64) {
	s.buffer[i] = uint8(s.BitDepth().UnsignedValue(value))
}

// Putu8 places signal buffer back to the pool. If a type of
// provided buffer isn't u8 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (s *u8) Free(p *PoolAllocator) {
	mustSameCapacity(s.Cap(), p.Channels*p.Capacity)
	for i := range s.buffer {
		s.buffer[i] = 0
	}
	p.u8.Put(s)
}

func uint8pool(size int) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &u8{
				buffer: make([]uint8, 0, size),
			}
		},
	}
}

// Capacity returns capacity of a single channel.
func (s *u8) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s *u8) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s *u8) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s *u8) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided sample index.
func (s *u8) Sample(i int) uint64 {
	return uint64(s.buffer[i])
}

// Append appends [0:Length] samples from src to current buffer and returns
// new Unsigned buffer. Both buffers must have same number of channels and
// bit depth, otherwise function will panic.
func (s *u8) Append(src Unsigned) {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	offset := s.Len()
	if s.Cap() < s.Len()+src.Len() {
		s.buffer = append(s.buffer, make([]uint8, src.Len())...)
	} else {
		s.buffer = s.buffer[:s.Len()+src.Len()]
	}
	for i := 0; i < src.Len(); i++ {
		s.SetSample(i+offset, src.Sample(i))
	}
	alignCapacity(&s.buffer, s.Channels(), s.Cap())
}

// Slice slices buffer with respect to channels.
func (s *u8) Slice(start, end int) Unsigned {
	start = s.BufferIndex(0, start)
	end = s.BufferIndex(0, end)
	return &u8{
		channels: s.channels,
		buffer:   s.buffer[start:end],
		bitDepth: s.bitDepth,
	}
}

// Channel returns a single channel of signal buffer. It behaves exactly as
// Unsigned buffer, but Append and AppendSample cause panic.
func (s *u8) Channel(c int) Unsigned {
	return unsignedChannel{
		buffer:  s,
		channel: c,
	}
}

// ReadUint8 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadUint8(src Unsigned, dst []uint8) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = uint8(BitDepth8.UnsignedValue(src.Sample(i)))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedUint8 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedUint8(src Unsigned, dst [][]uint8) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = uint8(BitDepth8.UnsignedValue(src.Sample(src.BufferIndex(c, i))))
		}
	}
	return
}

// WriteUint8 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteUint8(src []uint8, dst Unsigned) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, uint64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedUint8 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedUint8(src [][]uint8, dst Unsigned) (written int) {
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
				dst.SetSample(dst.BufferIndex(c, i), uint64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}
