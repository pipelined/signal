package pool

import (
	"sync"

	"pipelined.dev/signal"
)

type (
	// AllocSignedFunc allocates new signed buffer.
	AllocSignedFunc func(signal.BitDepth) signal.Signed
	// AllocUnsignedFunc allocates new unssigned buffer.
	AllocUnsignedFunc func(signal.BitDepth) signal.Unsigned
	// AllocFloatingFunc allocates new floating buffer.
	AllocFloatingFunc func() signal.Floating
)

type (
	// FloatingPool provides pool of floating buffers.
	FloatingPool struct {
		pool sync.Pool
	}

	// SignedPool provides pool of signed buffers.
	SignedPool struct {
		pool sync.Pool
	}

	// UnsignedPool provides pool of unsigned buffers.
	UnsignedPool struct {
		pool sync.Pool
	}
)

// Floating returns new FloatingPool.
func Floating(alloc AllocFloatingFunc) *FloatingPool {
	return &FloatingPool{
		pool: sync.Pool{
			New: func() interface{} {
				return alloc()
			},
		},
	}
}

// Get retrieves new floating buffer from the pool.
func (p *FloatingPool) Get() signal.Floating {
	return p.pool.Get().(signal.Floating)
}

// Put returns floating buffer to the pool. Buffer is also reset.
func (p *FloatingPool) Put(s signal.Floating) {
	p.pool.Put(s.Reset())
}

// Signed returns new SignedPool.
func Signed(alloc AllocSignedFunc, bd signal.BitDepth) *SignedPool {
	return &SignedPool{
		pool: sync.Pool{
			New: func() interface{} {
				return alloc(bd)
			},
		},
	}
}

// Get retrieves new signed buffer from the pool.
func (p *SignedPool) Get() signal.Signed {
	return p.pool.Get().(signal.Signed)
}

// Put returns signed buffer to the pool. Buffer is also reset.
func (p *SignedPool) Put(s signal.Signed) {
	p.pool.Put(s.Reset())
}

// Unsigned returns new UnsignedPool.
func Unsigned(alloc AllocUnsignedFunc, bd signal.BitDepth) *UnsignedPool {
	return &UnsignedPool{
		pool: sync.Pool{
			New: func() interface{} {
				return alloc(bd)
			},
		},
	}
}

// Get retrieves new unsigned buffer from the pool.
func (p *UnsignedPool) Get() signal.Unsigned {
	return p.pool.Get().(signal.Unsigned)
}

// Put returns unsigned buffer to the pool. Buffer is also reset.
func (p *UnsignedPool) Put(s signal.Unsigned) {
	p.pool.Put(s.Reset())
}
