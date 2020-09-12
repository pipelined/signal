package signal

import "sync"

var cache = struct {
	sync.Mutex
	pools map[int]*pool
}{
	pools: map[int]*pool{},
}

// PoolAllocator allows to decrease a number of allocations at runtime. Internally
// it relies on sync.PoolAllocator to manage objects in memory. It provides a pool
// per signal buffer type.
type PoolAllocator struct {
	Channels int
	Capacity int
	Length   int
	*pool
}

type pool struct {
	i8  sync.Pool
	i16 sync.Pool
	i32 sync.Pool
	i64 sync.Pool
	u8  sync.Pool
	u16 sync.Pool
	u32 sync.Pool
	u64 sync.Pool
	f32 sync.Pool
	f64 sync.Pool
}

// GetPoolAllocator returns pool for provided buffer dimensions. Pools are
// cached internally, so multiple calls with same dimentions will return
// the same pool instance.
func GetPoolAllocator(channels, length, capacity int) PoolAllocator {
	size := channels * capacity
	cache.Lock()
	defer cache.Unlock()

	if p, ok := cache.pools[size]; ok {
		return PoolAllocator{
			Length:   length,
			Channels: channels,
			Capacity: capacity,
			pool:     p,
		}
	}

	pool := pool{
		i8: sync.Pool{
			New: func() interface{} {
				return make([]int8, 0, size)
			},
		},
		i16: sync.Pool{
			New: func() interface{} {
				return make([]int16, 0, size)
			},
		},
		i32: sync.Pool{
			New: func() interface{} {
				return make([]int32, 0, size)
			},
		},
		i64: sync.Pool{
			New: func() interface{} {
				return make([]int64, 0, size)
			},
		},
		u8: sync.Pool{
			New: func() interface{} {
				return make([]uint8, 0, size)
			},
		},
		u16: sync.Pool{
			New: func() interface{} {
				return make([]uint16, 0, size)
			},
		},
		u32: sync.Pool{
			New: func() interface{} {
				return make([]uint32, 0, size)
			},
		},
		u64: sync.Pool{
			New: func() interface{} {
				return make([]uint64, 0, size)
			},
		},
		f32: sync.Pool{
			New: func() interface{} {
				return make([]float32, 0, size)
			},
		},
		f64: sync.Pool{
			New: func() interface{} {
				return make([]float64, 0, size)
			},
		},
	}

	cache.pools[size] = &pool
	return PoolAllocator{
		Channels: channels,
		Length:   length,
		Capacity: capacity,
		pool:     &pool,
	}
}

// ClearPoolAllocatorCache resets internal cache of pools and makes existing pools
// available for GC. One good use case might be when the application
// changes global buffer size.
func ClearPoolAllocatorCache() {
	cache.Lock()
	defer cache.Unlock()
	cache.pools = map[int]*pool{}
}
