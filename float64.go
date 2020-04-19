package signal

// Float64 is a sequential float64 floating-point signal.
type Float64 struct {
	buffer []float64
	capacity
	channels
	length
}

// Float64 allocates new sequential float64 signal buffer.
func (a Allocator) Float64() Float64 {
	return Float64{
		buffer:   make([]float64, a.Channels*a.Capacity),
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
	}
}

func (f Float64) setLength(l int) Floating {
	f.length = length(l)
	return f
}

// Sample returns signal value for provided channel and position.
func (f Float64) Sample(channel, pos int) float64 {
	return f.buffer[interPos(f.Channels(), channel, pos)]
}

func (f Float64) setSample(channel, pos int, value float64) {
	f.buffer[interPos(f.Channels(), channel, pos)] = value
}

// Append appends data from src buffer to the end of the buffer.
// The result buffer has capacity and length equal to sum of lengths.
func (f Float64) Append(src Floating) Floating {
	mustSameChannels(f.Channels(), src.Channels())
	newLen := f.Length() + src.Length()
	if f.Capacity() < newLen {
		f.buffer = append(f.buffer, make([]float64, (newLen-f.Capacity())*f.Channels())...)
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
