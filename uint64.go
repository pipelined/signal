package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-05-29 23:12:40.772476 +0200 CEST m=+0.019803586

import "math"

// Uint64 is uint64 Unsigned fixed-point signal.
type Uint64 struct {
	buffer []uint64
	channels
	bitDepth
}

// Uint64 allocates a new sequential uint64 signal buffer.
func (a Allocator) Uint64(bd BitDepth) Unsigned {
	return Uint64{
		buffer:   make([]uint64, a.Channels*a.Length, a.Channels*a.Capacity),
		channels: channels(a.Channels),
		bitDepth: limitBitDepth(bd, BitDepth64),
	}
}

func (s Uint64) setBitDepth(bd BitDepth) Unsigned {
	s.bitDepth = limitBitDepth(bd, BitDepth64)
	return s
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Uint64) AppendSample(value uint64) Unsigned {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, uint64(s.BitDepth().UnsignedValue(value)))
	return s
}

// SetSample sets sample value for provided position.
// Sample values are capped by maximum value of the buffer bit depth.
func (s Uint64) SetSample(pos int, value uint64) {
	s.buffer[pos] = uint64(s.BitDepth().UnsignedValue(value))
}

// GetUint64 selects a new sequential uint64 signal buffer.
// from the pool.
func (p *Pool) GetUint64(bd BitDepth) Unsigned {
	if p == nil {
		return nil
	}
	return p.u64.Get().(Unsigned).setBitDepth(bd)
}

// PutUint64 places signal buffer back to the pool. If a type of
// provided buffer isn't Uint64 or its capacity doesn't equal
// allocator capacity, the function will panic.
func (p *Pool) PutUint64(s Unsigned) {
	if p == nil {
		return
	}
	if _, ok := s.(Uint64); !ok {
		panic("pool put uint64 invalid type")
	}
	mustSameCapacity(s.Capacity(), p.allocator.Capacity)
	p.u64.Put(s.Slice(0, p.allocator.Length))
}

// Capacity returns capacity of a single channel.
func (s Uint64) Capacity() int {
	if s.channels == 0 {
		return 0
	}
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Uint64) Length() int {
	if s.channels == 0 {
		return 0
	}
	return int(math.Ceil(float64(len(s.buffer)) / float64(s.channels)))
}

// Cap returns capacity of whole buffer.
func (s Uint64) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Uint64) Len() int {
	return len(s.buffer)
}

// Sample returns signal value for provided channel and position.
func (s Uint64) Sample(pos int) uint64 {
	return uint64(s.buffer[pos])
}

// Append appends [0:Length] samples from src to current buffer and returns new
// Unsigned buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s Uint64) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]uint64, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Unsigned(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// Slice slices buffer with respect to channels.
func (s Uint64) Slice(start, end int) Unsigned {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// ReadUint64 reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadUint64(src Unsigned, dst []uint64) int {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = uint64(BitDepth64.UnsignedValue(src.Sample(pos)))
	}
	return chanLen(length, src.Channels())
}

// ReadStripedUint64 reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedUint64(src Unsigned, dst [][]uint64) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		length := min(len(dst[channel]), src.Length())
		if length > read {
			read = length
		}
		for pos := 0; pos < length; pos++ {
			dst[channel][pos] = uint64(BitDepth64.UnsignedValue(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
	return
}

// WriteUint64 writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteUint64(src []uint64, dst Unsigned) int {
	length := min(dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst.SetSample(pos, uint64(src[pos]))
	}
	return chanLen(length, dst.Channels())
}

// WriteStripedUint64 writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedUint64(src [][]uint64, dst Unsigned) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for channel := 0; channel < dst.Channels(); channel++ {
		for pos := 0; pos < written; pos++ {
			if pos < len(src[channel]) {
				dst.SetSample(dst.ChannelPos(channel, pos), uint64(src[channel][pos]))
			} else {
				dst.SetSample(dst.ChannelPos(channel, pos), 0)
			}
		}
	}
	return
}
