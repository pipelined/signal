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

func (s Float64) setLength(l int) Floating {
	s.length = length(l)
	return s
}

// Sample returns signal value for provided channel and position.
func (s Float64) Sample(channel, pos int) float64 {
	return s.buffer[interPos(s.Channels(), channel, pos)]
}

func (s Float64) setSample(channel, pos int, value float64) {
	s.buffer[interPos(s.Channels(), channel, pos)] = value
}

// Append appends data from src buffer to the end of the buffer.
// The result buffer has capacity and length equal to sum of lengths.
func (s Float64) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	newLen := s.Length() + src.Length()
	if s.Capacity() < newLen {
		s.buffer = append(s.buffer, make([]float64, (newLen-s.Capacity())*s.Channels())...)
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
