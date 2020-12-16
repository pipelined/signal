package signal

// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// 2020-12-16 20:52:36.096957 +0100 CET m=+0.016732471

type (
	floatingChannel struct {
		buffer  Floating
		channel int
	}
)

// Append panics.
func (c floatingChannel) Append(s Floating) {
	panic("appending signal to the single channel")
}

// AppendSample panics.
func (c floatingChannel) AppendSample(s float64) {
	panic("appending sample to the single channel")
}

// BufferIndex returns sample index in the channel of signal buffer.
func (c floatingChannel) BufferIndex(channel int, index int) int {
	return c.channel * index
}

// Channels always returns 1.
func (c floatingChannel) Channels() int {
	return 1
}

// Cap returns capacity of the channel.
func (c floatingChannel) Cap() int {
	return c.buffer.Capacity()
}

// Capacity returns capacity of the channel.
func (c floatingChannel) Capacity() int {
	return c.buffer.Capacity()
}

// Len returns length of the channel.
func (c floatingChannel) Len() int {
	return c.buffer.Length()
}

// Length returns length of the channel.
func (c floatingChannel) Length() int {
	return c.buffer.Length()
}

// Channel panics.
func (c floatingChannel) Channel(channel int) Floating {
	panic("slicing channel of the channel")
}

// Sample returns signal value for provided channel and index.
func (c floatingChannel) Sample(index int) float64 {
	return c.buffer.Sample(index * c.channel)
}

// SetSample sets sample value for provided index.
func (c floatingChannel) SetSample(index int, s float64) {
	c.buffer.SetSample(c.buffer.BufferIndex(c.channel, index), s)
}

// Free panics.
func (c floatingChannel) Free(*PoolAllocator) {
	panic("freeing single channel of the buffer")
}

// Slice slices buffer with respect to channels.
func (c floatingChannel) Slice(start, end int) Floating {
	return floatingChannel{
		buffer:  c.buffer.Slice(start, end),
		channel: c.channel,
	}
}
