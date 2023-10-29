package signal

import (
	"sync"
)

type (

	// PoolAllocator allows to decrease a number of allocations at runtime.
	// Internally it relies on sync.Pool to manage objects in memory.
	PoolAllocator[T SignalTypes] struct {
		Get getFunc[T]
		Put putFunc[T]
	}

	getFunc[T SignalTypes] func() *Buffer[T]
	putFunc[T SignalTypes] func(*Buffer[T])

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

func GetPool[T SignalTypes](a Allocator) PoolAllocator[T] {
	pool := &pool[T]{
		p: &sync.Pool{
			New: func() any {
				return Alloc[T](a)
			},
		},
	}
	return PoolAllocator[T]{
		Get: func() *Buffer[T] {
			return pool.Get()
		},
		Put: func(f *Buffer[T]) {
			mustSame(a.Capacity*a.Channels, f.Cap(), diffCapacity)
			f.clear()
			pool.Put(f)
		},
	}
}
