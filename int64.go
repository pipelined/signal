package signal

// Int64 is int64 signed fixed signal.
type Int64 struct {
	buffer []int64
	channels
	bitDepth
}

// Int64 allocates new sequential int64 signal buffer.
func (a Allocator) Int64(bd BitDepth) Int64 {
	return Int64{
		buffer:   make([]int64, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: bd.cap(BitDepth64),
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (s Int64) MaxBitDepth() BitDepth {
	return BitDepth64
}

func (s Int64) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

func (s Int64) Length() int {
	return len(s.buffer) / int(s.channels)
}

func (s Int64) Cap() int {
	return cap(s.buffer)
}

func (s Int64) Len() int {
	return len(s.buffer)
}

func (s Int64) Reset() Signed {
	return s.Slice(0, 0)
}

func (s Int64) AppendSample(value int64) Signed {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, s.BitDepth().SignedValue(value))
	return s
}

// Sample returns signal value for provided channel and position.
func (s Int64) Sample(pos int) int64 {
	return s.buffer[pos]
}

func (s Int64) SetSample(pos int, value int64) {
	s.buffer[pos] = s.BitDepth().SignedValue(value)
}

func (s Int64) Slice(start, end int) Signed {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Append appends [0:Length] data from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// additional memory will be allocated and result buffer will have capacity and
// length equal to sum of lengths.
func (s Int64) Append(src Signed) Signed {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// if capacity is not enough, then:
		// * extend buffer to cap;
		// * allocate and append buffer with length of source capacity;
		// * slice it to current data length;
		s.buffer = append(s.buffer[:s.Cap()], make([]int64, src.Cap())...)[:s.Len()]
	}
	result := Signed(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// WriteInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteInt64(s Signed, buf []int64) Signed {
	length := min(s.Cap()-s.Len(), len(buf))
	for pos := 0; pos < length; pos++ {
		s = s.AppendSample(buf[pos])
	}
	return s
}

// WriteStripedInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt64(s Signed, buf [][]int64) Signed {
	mustSameChannels(s.Channels(), len(buf))
	var length int
	for i := range buf {
		if len(buf[i]) > length {
			length = len(buf[i])
		}
	}
	length = min(length, s.Capacity()-s.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < s.Channels(); channel++ {
			if pos < len(buf[channel]) {
				s = s.AppendSample(buf[channel][pos])
			} else {
				s = s.AppendSample(0)
			}
		}
	}
	return s
}

func ReadInt64(s Signed, buf []int64) {
	length := min(s.Len(), len(buf))
	for pos := 0; pos < length; pos++ {
		buf[pos] = s.Sample(pos)
	}
}

func ReadStripedInt64(s Signed, buf [][]int64) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Length() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = s.Sample(s.ChannelPos(channel, pos))
		}
	}
}
