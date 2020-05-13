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

func ReadFloat64(src Floating, dst []float64) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = src.Sample(pos)
	}
}

func ReadStripedFloat64(src Floating, dst [][]float64) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = src.Sample(src.ChannelPos(channel, pos))
		}
	}
}

// WriteFloat64 writes values from provided slice into the buffer.
func WriteFloat64(src []float64, dst Floating) Floating {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(src[pos])
	}
	return dst
}

// WriteStripedFloat64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of enclosing slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length.
func WriteStripedFloat64(src [][]float64, dst Floating) Floating {
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
