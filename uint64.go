package signal

// Uint64 is uint64 signed fixed signal.
type Uint64 struct {
	buffer []uint64
	channels
	bitDepth
}

// Uint64 allocates new sequential uint64 signal buffer.
func (a Allocator) Uint64(bd BitDepth) Unsigned {
	return Uint64{
		buffer:   make([]uint64, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: defaultBitDepth(bd, BitDepth64),
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (s Uint64) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Capacity returns capacity of a single channel.
func (s Uint64) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s Uint64) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s Uint64) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s Uint64) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s Uint64) AppendSample(value uint64) Unsigned {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, s.BitDepth().UnsignedValue(value))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Uint64) Sample(pos int) uint64 {
	return s.buffer[pos]
}

// SetSample sets sample value for provided position.
func (s Uint64) SetSample(pos int, value uint64) {
	s.buffer[pos] = s.BitDepth().UnsignedValue(value)
}

// Slice slices buffer with respect to channels.
func (s Uint64) Slice(start, end int) Unsigned {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s Uint64) Reset() Unsigned {
	return s.Slice(0, 0)
}

// Append appends data from src to current buffer and returns new
// Unsigned buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// additional memory will be allocated and result buffer will have capacity and
// length equal to sum of lengths.
func (s Uint64) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// if capacity is not enough, then:
		// * extend buffer to cap;
		// * allocate and append buffer with length of source capacity;
		// * slice it to current data length;
		s.buffer = append(s.buffer[:s.Cap()], make([]uint64, src.Cap())...)[:s.Len()]
	}
	result := Unsigned(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// ReadUint64 reads values from the buffer into provided slice.
func ReadUint64(src Unsigned, dst []uint64) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = src.Sample(pos)
	}
}

// ReadStripedUint64 reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStripedUint64(src Unsigned, dst [][]uint64) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = src.Sample(src.ChannelPos(channel, pos))
		}
	}
}

// WriteUint64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteUint64(src []uint64, dst Unsigned) Unsigned {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(src[pos])
	}
	return dst
}

// WriteStripedUint64 appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended. Sample values are capped by maximum value
// of the buffer bit depth.
func WriteStripedUint64(src [][]uint64, dst Unsigned) Unsigned {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(src[channel][pos])
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}
