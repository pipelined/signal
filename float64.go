package signal

// Float64 is a sequential float64 floating-point signal.
type Float64 struct {
	buffer []float64
	channels
}

// Float64 allocates new sequential float64 signal buffer.
func (a Allocator) Float64() Float64 {
	return Float64{
		buffer:   make([]float64, 0, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

func (s Float64) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

func (s Float64) Length() int {
	return len(s.buffer) / int(s.channels)
}

func (s Float64) Cap() int {
	return cap(s.buffer)
}

func (s Float64) Len() int {
	return len(s.buffer)
}

func (s Float64) Reset() Floating {
	return s.Slice(0, 0)
}

func (s Float64) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, value)
	return s
}

// Sample returns signal value for provided channel and position.
func (s Float64) Sample(pos int) float64 {
	return s.buffer[pos]
}

func (s Float64) SetSample(pos int, value float64) {
	s.buffer[pos] = value
}

func (s Float64) Slice(start, end int) Floating {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Append appends data from src buffer to the end of the buffer.
func (s Float64) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// if capacity is not enough, then:
		// * extend buffer to cap;
		// * allocate and append buffer with length of source capacity;
		// * slice it to current data length;
		s.buffer = append(s.buffer[:s.Cap()], make([]float64, src.Cap())...)[:s.Len()]
	}
	result := Floating(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}
