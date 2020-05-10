package signal

// Uint64 is uint64 signed fixed signal.
type Uint64 struct {
	buffer []uint64
	channels
	bitDepth
}

// Uint64 allocates new sequential uint64 signal buffer.
func (a Allocator) Uint64(bd BitDepth) Uint64 {
	return Uint64{
		buffer:   make([]uint64, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: bd.cap(BitDepth64),
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (s Uint64) MaxBitDepth() BitDepth {
	return BitDepth64
}

func (s Uint64) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

func (s Uint64) Length() int {
	return len(s.buffer) / int(s.channels)
}

func (s Uint64) Cap() int {
	return cap(s.buffer)
}

func (s Uint64) Len() int {
	return len(s.buffer)
}

func (s Uint64) Reset() Unsigned {
	return s.Slice(0, 0)
}

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

func (s Uint64) SetSample(pos int, value uint64) {
	s.buffer[pos] = s.BitDepth().UnsignedValue(value)
}

func (s Uint64) Slice(start, end int) Unsigned {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
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

func ReadUint64(src Unsigned, dst []uint64) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = src.Sample(pos)
	}
}

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

// WriteStripedUint64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
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
