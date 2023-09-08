package signal

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type (
	GetFloatFunc[T constraints.Float]     func() *Float[T]
	ReleaseFloatFunc[T constraints.Float] func(*Float[T])

	PFloatAllocator[T constraints.Float] struct {
		Get     GetFloatFunc[T]
		Release ReleaseFloatFunc[T]
	}

	floatPool[T constraints.Float] struct {
		p *sync.Pool
	}
)

type (
	GetIntegerFunc[T constraints.Integer]     func(BitDepth) *Integer[T]
	ReleaseIntegerFunc[T constraints.Integer] func(*Integer[T])

	PIntegerAllocator[T constraints.Integer] struct {
		Get     GetIntegerFunc[T]
		Release ReleaseIntegerFunc[T]
	}

	integerPool[T constraints.Integer] struct {
		p *sync.Pool
	}
)

func (p *floatPool[T]) Get() *Float[T] {
	return p.p.Get().(*Float[T])
}

func (p *floatPool[T]) Put(t *Float[T]) {
	p.p.Put(t)
}

func (p *integerPool[T]) Get() *Integer[T] {
	return p.p.Get().(*Integer[T])
}

func (p *integerPool[T]) Put(t *Integer[T]) {
	p.p.Put(t)
}

func GetFloatPool[T constraints.Float](a Allocator) PFloatAllocator[T] {
	pool := &floatPool[T]{
		p: &sync.Pool{
			New: func() any {
				return &Float[T]{
					buffer: buffer[T]{
						data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
						channels: channels(a.Channels),
						wrapFn:   wrapFloatSample[T],
					},
				}
			},
		},
	}
	return PFloatAllocator[T]{
		Get: func() *Float[T] {
			return pool.Get()
		},
		Release: func(f *Float[T]) {
			mustSameCapacity(a.Capacity*a.Channels, f.Cap())
			f.buffer.clear()
			pool.Put(f)
		},
	}
}

func GetIntegerPool[T constraints.Integer](a Allocator) PIntegerAllocator[T] {
	pool := &integerPool[T]{
		p: &sync.Pool{
			New: func() any {
				return &Integer[T]{
					buffer: buffer[T]{
						data:     make([]T, a.Channels*a.Length, a.Channels*a.Capacity),
						channels: channels(a.Channels),
					},
				}
			},
		},
	}
	return PIntegerAllocator[T]{
		Get: func(bd BitDepth) *Integer[T] {
			b := pool.Get()
			b.bitDepth = limitBitDepth[T](bd)
			b.buffer.wrapFn = wrapIntegerSample[T](bd)
			return b
		},
		Release: func(f *Integer[T]) {
			mustSameCapacity(a.Capacity*a.Channels, f.Cap())
			f.buffer.clear()
			pool.Put(f)
		},
	}
}
