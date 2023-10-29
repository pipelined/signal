package signal

type C[T SignalTypes] struct {
	Buffer  *Buffer[T]
	channel int
}

// BufferIndex returns sample index in the channel of signal Buffer.
func (c C[T]) BufferIndex(channel int, index int) int {
	return c.channel * index
}

// Channels always returns 1.
func (c C[T]) Channels() int {
	return 1
}

// Capacity returns capacity of the channel.
func (c C[T]) Capacity() int {
	return c.Buffer.Capacity()
}

// Length returns length of the channel.
func (c C[T]) Length() int {
	return c.Buffer.Length()
}

// Sample returns signal value for provided channel and index.
func (c C[T]) Sample(index int) T {
	return c.Buffer.Sample(index * c.channel)
}

// SetSample sets sample value for provided index.
func (c C[T]) SetSample(index int, s T) {
	c.Buffer.SetSample(c.Buffer.BufferIndex(c.channel, index), s)
}
