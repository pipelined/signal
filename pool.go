package signal

import (
	"sync"
)

type (
	GetFunc[T SignalTypes] func() *Buffer[T]
	Release[T SignalTypes] func(*Buffer[T])

	PAllocator[T SignalTypes] struct {
		Get     GetFunc[T]
		Release Release[T]
	}

	pool[T SignalTypes] struct {
		p *sync.Pool
	}
)

func (p *pool[T]) Get() *Buffer[T] {
	return p.p.Get().(*Buffer[T])
}

func (p *pool[T]) Put(t *Buffer[T]) {
	p.p.Put(t)
}

func GetPool[T SignalTypes](a Allocator) PAllocator[T] {
	pool := &pool[T]{
		p: &sync.Pool{
			New: func() any {
				return &Buffer[T]{
					buffer: buffer[T]{
						data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
						channels: channels(a.Channels),
					},
					bitDepth: bitDepth(getBitDepth[T]()),
				}
			},
		},
	}
	return PAllocator[T]{
		Get: func() *Buffer[T] {
			return pool.Get()
		},
		Release: func(f *Buffer[T]) {
			mustSameCapacity(a.Capacity*a.Channels, f.Cap())
			f.buffer.clear()
			pool.Put(f)
		},
	}
}
