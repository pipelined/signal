package signal

type (
	// Int64 is a sequential int64 signal.
	Int64 struct {
		buffer [][]int64
		capacity
		channels
		bitDepth
		*length
	}

	// Int64Interleaved is an interleaved int64 signal.
	Int64Interleaved struct {
		buffer []int64
		capacity
		channels
		bitDepth
		*length
	}
)

// Int64 allocates new sequential int64 signal buffer.
func (a Allocator) Int64(bd BitDepth) Int64 {
	if bd == 0 || bd == MaxBitDepth {
		bd = BitDepth64
	}
	buffer := make([][]int64, a.Channels)
	for i := range buffer {
		buffer[i] = make([]int64, a.Capacity)
	}
	return Int64{
		buffer:   buffer,
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
		bitDepth: bitDepth(bd),
		length:   &length{},
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (f Int64) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Sample returns signal value for provided channel and position.
func (f Int64) Sample(channel, pos int) int64 {
	return f.buffer[channel][pos]
}

// Data returns underlying signal buffer.
func (f Int64) Data() [][]int64 {
	return f.buffer
}

func (f Int64) setSample(channel, pos int, val int64) {
	f.buffer[channel][pos] = val
}

// WriteInt64 writes values from provided slice into buffer.
// Provided slice must have the exactly same number of channels.
// Length is updated with longest nested slice length.
func (f Int64) WriteInt64(ints [][]int64) int {
	mustSameChannels(f.Channels(), len(ints))
	var (
		length = f.Length()
		copied = 0
	)
	for channel := range f.buffer {
		c := 0
		for pos := length; pos < f.Capacity() && c < len(ints[channel]); pos, c = pos+1, c+1 {
			f.buffer[channel][pos] = f.BitDepth().SignedValue(ints[channel][c])
		}
		if copied < c {
			copied = c
		}
	}
	f.setLength(length + copied)
	return copied
}

// WriteInt writes values from provided slice into buffer.
// Provided slice must have the exactly same number of channels.
// Length is updated with longest nested slice length.
func (f Int64) WriteInt(ints [][]int) int {
	mustSameChannels(f.Channels(), len(ints))
	copied := 0
	for channel := range f.buffer {
		c := 0
		pos := f.Length()
		for pos < f.Capacity() && c < len(ints[channel]) {
			f.buffer[channel][pos] = f.BitDepth().SignedValue(int64(ints[channel][c]))
			pos++
			c++
		}
		if copied < c {
			copied = c
		}
	}
	f.setLength(f.Length() + copied)
	return copied
}

// Append appends [0:Length] data from src to current buffer and returns new buffer.
func (f Int64) Append(src Int64) Int64 {
	mustSameChannels(f.Channels(), f.Channels())
	mustSameBitDepth(f.BitDepth(), src.BitDepth())
	l := f.Length() + src.Length()
	result := make([][]int64, f.Channels())
	for channel := range result {
		result[channel] = append(f.buffer[channel][:f.Length()], src.buffer[channel][:src.Length()]...)
	}
	return Int64{
		buffer:   result,
		capacity: capacity(l),
		channels: f.channels,
		length:   &length{value: l},
		bitDepth: f.bitDepth,
	}
}

// Int64Interleaved returns new int64 interleaved buffer. If non-nill parameter
// is provided, the values will a. copied into result buffer. Result buffer will
// always have size provided a. properties.
func (a Allocator) Int64Interleaved(bd BitDepth) Int64Interleaved {
	if bd == 0 || bd == MaxBitDepth {
		bd = BitDepth64
	}
	return Int64Interleaved{
		buffer:   make([]int64, a.Capacity*a.Channels),
		capacity: capacity(a.Capacity),
		channels: channels(a.Channels),
		bitDepth: bitDepth(bd),
		length:   &length{},
	}
}

// MaxBitDepth returns maximal bit depth for this type.
func (f Int64Interleaved) MaxBitDepth() BitDepth {
	return BitDepth64
}

// Sample returns signal value for provided channel and position.
func (f Int64Interleaved) Sample(channel, pos int) int64 {
	return f.buffer[interPos(f.Channels(), channel, pos)]
}

// Data returns underlying signal buffer.
func (f Int64Interleaved) Data() []int64 {
	return f.buffer
}

func (f Int64Interleaved) setSample(channel, pos int, val int64) {
	f.buffer[interPos(f.Channels(), channel, pos)] = val
}

// WriteInt64 writes values from provided slice into buffer.
// Length is updated with slice length.
func (f Int64Interleaved) WriteInt64(ints []int64) int {
	bufLen := f.Length() * f.Channels()
	c := copy(f.buffer[bufLen:], ints)
	f.setLength(interLen(f.Channels(), bufLen+c))
	return c
}

// WriteInt writes values from provided slice into buffer.
// Length is updated with slice length.
func (f Int64Interleaved) WriteInt(ints []int) int {
	pos := 0
	c := 0
	for pos < f.Capacity() && c < len(ints) {
		f.buffer[pos] = int64(ints[c])
		pos++
		c++
	}
	f.setLength(interLen(f.Channels(), f.Length()*f.Channels()+c))
	return c
}

func (f Int64Interleaved) Append(src Int64Interleaved) Int64Interleaved {
	mustSameChannels(f.Channels(), src.Channels())
	mustSameBitDepth(f.BitDepth(), src.BitDepth())
	panic("not implemented")
}
