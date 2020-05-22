package signal

import "sync"

type Pool struct {
	allocator Allocator
	i8        sync.Pool
	i16       sync.Pool
	i32       sync.Pool
	i64       sync.Pool
	u8        sync.Pool
	u16       sync.Pool
	u32       sync.Pool
	u64       sync.Pool
	f32       sync.Pool
	f64       sync.Pool
}

func (a Allocator) Pool() *Pool {
	return &Pool{
		allocator: a,
		i8:        signedPool(a.Int8, BitDepth8),
		i16:       signedPool(a.Int16, BitDepth16),
		i32:       signedPool(a.Int32, BitDepth32),
		i64:       signedPool(a.Int64, BitDepth64),
		u8:        unsignedPool(a.Uint8, BitDepth8),
		u16:       unsignedPool(a.Uint16, BitDepth16),
		u32:       unsignedPool(a.Uint32, BitDepth32),
		u64:       unsignedPool(a.Uint64, BitDepth64),
		f32:       floatingPool(a.Float32),
		f64:       floatingPool(a.Float64),
	}
}

func signedPool(alloc func(BitDepth) Signed, mbd BitDepth) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return alloc(mbd)
		},
	}
}

func unsignedPool(alloc func(BitDepth) Unsigned, mbd BitDepth) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return alloc(mbd)
		},
	}
}

func floatingPool(alloc func() Floating) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return alloc()
		},
	}
}

func (p *Pool) Allocator() Allocator {
	if p != nil {
		return p.allocator
	}
	return Allocator{}
}
