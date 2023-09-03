package signal

type C[T SignalTypes] struct {
	buffer  buffer[T]
	channel int
}

// BufferIndex returns sample index in the channel of signal buffer.
func (c C[T]) BufferIndex(channel int, index int) int {
	return c.channel * index
}

// Channels always returns 1.
func (c C[T]) Channels() int {
	return 1
}

// Cap returns capacity of the channel.
func (c C[T]) Cap() int {
	return c.buffer.Capacity()
}

// Capacity returns capacity of the channel.
func (c C[T]) Capacity() int {
	return c.buffer.Capacity()
}

// Len returns length of the channel.
func (c C[T]) Len() int {
	return c.buffer.Length()
}

// Length returns length of the channel.
func (c C[T]) Length() int {
	return c.buffer.Length()
}

// Sample returns signal value for provided channel and index.
func (c C[T]) Sample(index int) T {
	return c.buffer.Sample(index * c.channel)
}

// SetSample sets sample value for provided index.
func (c C[T]) SetSample(index int, s T) {
	c.buffer.SetSample(c.buffer.BufferIndex(c.channel, index), s)
}

// Slice slices buffer with respect to channels.
func (c C[T]) Slice(start, end int) C[T] {
	return C[T]{
		buffer: buffer[T]{
			data: c.buffer.data[start:end],
		},
		channel: c.channel,
	}
}
