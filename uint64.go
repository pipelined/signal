package signal

// Uint64 is uint64 unsigned fixed signal.
type Uint64 struct {
	buffer []uint64
	capacity
	channels
	bitDepth
	length
}

// Uint64 allocates new uint64 signal buffer.
func (a Allocator) Uint64(bd BitDepth) Uint64 {
	return Uint64{
		buffer:   make([]uint64, a.Channels*a.Capacity),
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
		bitDepth: bd.cap(BitDepth64),
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (s Uint64) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Sample returns signal value for provided channel and position.
func (s Uint64) Sample(channel, pos int) uint64 {
	return s.buffer[interPos(s.Channels(), channel, pos)]
}

func (s Uint64) setSample(channel, pos int, val uint64) {
	s.buffer[interPos(s.Channels(), channel, pos)] = val
}

func (s Uint64) setLength(l int) Unsigned {
	s.length = length(l)
	return s
}

func (s Uint64) Slice(start, end int) Unsigned {
	s.buffer = s.buffer[interPos(s.Channels(), 0, start):]
	s.capacity = capacity(s.Capacity() - start)
	return s.setLength(end - start)
}

// Append appends [0:Length] data from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// additional memory will be allocated and result buffer will have capacity and
// length equal to sum of lengths.
func (s Uint64) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	newLen := s.Length() + src.Length()
	if s.Capacity() < newLen {
		s.buffer = append(s.buffer, make([]uint64, (newLen-s.Capacity())*s.Channels())...)
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
