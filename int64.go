package signal

// Int64 is int64 signed fixed signal.
type Int64 struct {
	buffer []int64
	capacity
	channels
	bitDepth
	length
}

// Int64 allocates new sequential int64 signal buffer.
func (a Allocator) Int64(bd BitDepth) Int64 {
	return Int64{
		buffer:   make([]int64, a.Capacity*a.Channels),
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
		bitDepth: bd.cap(BitDepth64),
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (s Int64) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Sample returns signal value for provided channel and position.
func (s Int64) Sample(channel, pos int) int64 {
	return s.buffer[interPos(s.Channels(), channel, pos)]
}

func (s Int64) setSample(channel, pos int, val int64) {
	s.buffer[interPos(s.Channels(), channel, pos)] = val
}

func (s Int64) setLength(l int) Signed {
	s.length = length(l)
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
	newLen := s.Length() + src.Length()
	if s.Capacity() < newLen {
		s.buffer = append(s.buffer, make([]int64, (newLen-s.Capacity())*s.Channels())...)
		s.capacity = capacity(newLen)
	}
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < src.Length(); pos++ {
			s.setSample(channel, pos+s.Length(), src.Sample(channel, pos))
		}
	}
	s.length = length(newLen)
	return s
}
