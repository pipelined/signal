package signal

// Float64 is a sequential float64 floating-point signal.
type Float64 struct {
	buffer [][]float64
	capacity
	channels
	*length
}

// Float64 allocates new sequential float64 signal buffer.
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

// Sample returns signal value for provided channel and position.
func (f Float64) Sample(channel, pos int) float64 {
	return f.buffer[channel][pos]
}

// Data returns underlying signal buffer.
func (f Float64) Data() [][]float64 {
	return f.buffer
}

func (f Float64) setSample(channel, pos int, value float64) {
	f.buffer[channel][pos] = value
}

// WriteFloat64 writes values from provided slice into buffer.
// Provided slice must have the exactly same number of channels.
// Length is updated with longest nested slice length.
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

// Append appends data from src buffer to the end of the buffer.
// The result buffer has capacity and length equal to sum of lengths.
func (f Float64) Append(src Float64) Float64 {
	mustSameChannels(f.Channels(), src.Channels())
	l := f.Length() + src.Length()
	result := make([][]float64, f.Channels())
	for channel := range result {
		result[channel] = append(f.buffer[channel][:f.Length()], src.buffer[channel][:src.Length()]...)
	}
	return Float64{
		buffer:   result,
		capacity: capacity(l),
		channels: f.channels,
		length:   &length{value: l},
	}
}
