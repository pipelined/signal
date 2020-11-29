package signal

type (
	floatingChannel struct {
		buf     Floating
		channel int
	}
)

// TODO
func (c floatingChannel) Append(s Floating) {
}

// TODO
func (c floatingChannel) AppendSample(s float64) {
	for i := 0; i < c.buf.Channels(); i++ {
		if c.channel == i {
			c.buf.AppendSample(s)
			continue
		}
		c.buf.AppendSample(0)
	}
}

func (c floatingChannel) BufferIndex(channel int, pos int) int {
	return c.channel * pos
}

func (c floatingChannel) Channels() int {
	return 1
}

func (c floatingChannel) Cap() int {
	return c.buf.Capacity()
}

func (c floatingChannel) Capacity() int {
	return c.buf.Capacity()
}

func (c floatingChannel) Len() int {
	return c.buf.Length()
}

func (c floatingChannel) Length() int {
	return c.buf.Length()
}

func (c floatingChannel) Channel(channel int) Floating {
	return floatingChannel{
		buf:     c.buf,
		channel: channel,
	}
}

func (c floatingChannel) Sample(pos int) float64 {
	return c.buf.Sample(pos * c.channel)
}

func (c floatingChannel) SetSample(pos int, s float64) {
	c.buf.SetSample(c.buf.BufferIndex(c.channel, pos), s)
}

func (c floatingChannel) Free(*PoolAllocator) {
	panic("freeing single channel of the buffer")
}

// TODO
func (c floatingChannel) Slice(start, end int) Floating {
	return floatingChannel{
		buf:     c.buf.Slice(start, end),
		channel: c.channel,
	}
}
