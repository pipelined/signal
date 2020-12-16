package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 19:42:20.042445 +0100 CET m=+0.017504438

type (
	unsignedChannel struct {
		buffer  Unsigned
		channel int
	}
)

// Append panics.
func (c unsignedChannel) Append(s Unsigned) {
	panic("appending signal to the single channel")
}

// AppendSample panics.
func (c unsignedChannel) AppendSample(s uint64) {
	panic("appending sample to the single channel")
}

// BufferIndex returns sample index in the channel of signal buffer.
func (c unsignedChannel) BufferIndex(channel int, index int) int {
	return c.channel * index
}

// Channels always returns 1.
func (c unsignedChannel) Channels() int {
	return 1
}

// Cap returns capacity of the channel.
func (c unsignedChannel) Cap() int {
	return c.buffer.Capacity()
}

// Capacity returns capacity of the channel.
func (c unsignedChannel) Capacity() int {
	return c.buffer.Capacity()
}

// Len returns length of the channel.
func (c unsignedChannel) Len() int {
	return c.buffer.Length()
}

// Length returns length of the channel.
func (c unsignedChannel) Length() int {
	return c.buffer.Length()
}

// Channel panics.
func (c unsignedChannel) Channel(channel int) Unsigned {
	panic("slicing channel of the channel")
}

// Sample returns signal value for provided channel and index.
func (c unsignedChannel) Sample(index int) uint64 {
	return c.buffer.Sample(index * c.channel)
}

// SetSample sets sample value for provided index.
func (c unsignedChannel) SetSample(index int, s uint64) {
	c.buffer.SetSample(c.buffer.BufferIndex(c.channel, index), s)
}

// Free panics.
func (c unsignedChannel) Free(*PoolAllocator) {
	panic("freeing single channel of the buffer")
}

// Slice slices buffer with respect to channels.
func (c unsignedChannel) Slice(start, end int) Unsigned {
	return unsignedChannel{
		buffer:  c.buffer.Slice(start, end),
		channel: c.channel,
	}
}

func (c unsignedChannel) BitDepth() BitDepth {
	return c.buffer.BitDepth()
}
