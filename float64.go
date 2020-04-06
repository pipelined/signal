package signal

// Float64 is a sequential float64 floating-point signal.
type Float64 struct {
	buffer [][]float64
	capacity
	channels
	*length
}

func (a Allocator) Float64() Float64 {
	buffer := make([][]float64, a.Channels)
	for i := range buffer {
		buffer[i] = make([]float64, a.Capacity)
	}
	return Float64{
		buffer:   buffer,
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
		length:   &length{},
	}
}

func (f Float64) WriteFloat64(floats [][]float64) int {
	mustSameChannels(f.Channels(), len(floats))
	l := f.Length()
	n := 0
	for channel := range f.buffer {
		if copied := copy(f.buffer[channel][l:], floats[channel]); copied > n {
			n = copied
		}
	}
	f.setLength(l + n)
	return n
}

func (f Float64) Append(src Float64) Float64 {
	mustSameChannels(f.Channels(), src.Channels())
	panic("not implemented")
}

func (f Float64) Sample(channel, pos int) float64 {
	return f.buffer[channel][pos]
}

func (f Float64) setSample(channel, pos int, value float64) {
	f.buffer[channel][pos] = value
}