package signal

type (
	// Int64 is a sequential int64 signal.
	Int64 struct {
		buffer []int64
		capacity
		channels
		bitDepth
		length
	}
)

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
func (f Int64) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Sample returns signal value for provided channel and position.
func (f Int64) Sample(channel, pos int) int64 {
	return f.buffer[interPos(f.Channels(), channel, pos)]
}

func (f Int64) setSample(channel, pos int, val int64) {
	f.buffer[interPos(f.Channels(), channel, pos)] = val
}

func (f Int64) setLength(l int) Signed {
	f.length = length(l)
	return f
}

// Append appends [0:Length] data from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// additional memory will be allocated and result buffer will have capacity and
// length equal to sum of lengths.
func (f Int64) Append(src Signed) Signed {
	mustSameChannels(f.Channels(), src.Channels())
	mustSameBitDepth(f.BitDepth(), src.BitDepth())
	newLen := f.Length() + src.Length()
	if f.Capacity() < newLen {
		f.buffer = append(f.buffer, make([]int64, (newLen-f.Capacity())*f.Channels())...)
		f.capacity = capacity(newLen)
	}
	for channel := 0; channel < f.Channels(); channel++ {
		for pos := 0; pos < src.Length(); pos++ {
			f.setSample(channel, pos+f.Length(), src.Sample(channel, pos))
		}
	}
	f.length = length(newLen)
	return f
}
