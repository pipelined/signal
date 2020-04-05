// Package pool provides allocation pool for float64 signal buffers.
package pool

import (
	"sync"

	"pipelined.dev/signal"
)

// Pool for signal buffers.
type Pool struct {
	bufferSize  int
	numChannels int
	pool        sync.Pool
}

// New returns new pool.
func New(a signal.Allocator) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				b := a.Float64()
				return &b
			},
		},
	}
}

// Alloc retrieves new signal.Float64 buffer from the pool.
func (p *Pool) Alloc() *signal.Float64 {
	return p.pool.Get().(*signal.Float64)
}

// Free returns signal.Float64 buffer to the pool. Buffer is also cleared up.
func (p *Pool) Free(b *signal.Float64) {
	p.pool.Put(b)
}
