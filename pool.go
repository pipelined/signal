package signal

import (
	"sync"
)

// PoolAllocator allows to decrease a number of allocations at runtime.
// Internally it relies on sync.Pool to manage objects in memory.
type PoolAllocator[T SignalTypes] struct {
	pool  *sync.Pool
	alloc Allocator
}

// PoolAlloc returns new PoolAllocator.
func PoolAlloc[T SignalTypes](a Allocator) PoolAllocator[T] {
	return PoolAllocator[T]{
		alloc: a,
		pool: &sync.Pool{
			New: func() any {
				return Alloc[T](a)
			},
		},
	}
}

func (p *PoolAllocator[T]) Get() *Buffer[T] {
	return p.pool.Get().(*Buffer[T])
}

func (p *PoolAllocator[T]) Put(b *Buffer[T]) {
	mustSame(p.alloc.Capacity*p.alloc.Channels, b.Cap(), diffCapacity)
	b.clear()
	p.pool.Put(b)
}
