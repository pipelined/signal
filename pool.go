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
func GetPoolAllocator(channels, length, capacity int) *PoolAllocator {
	size := channels * capacity
	cache.Lock()
	defer cache.Unlock()

	if p, ok := cache.pools[size]; ok {
		return &PoolAllocator{
			Channels: channels,
			Length:   length,
			Capacity: capacity,
			pool:     p,
		}
	}

	p := pool{
		i8:  int8pool(size),
		i16: int16pool(size),
		i32: int32pool(size),
		i64: int64pool(size),
		u8:  uint8pool(size),
		u16: uint16pool(size),
		u32: uint32pool(size),
		u64: uint64pool(size),
		f32: float32pool(size),
		f64: float64pool(size),
	}

	cache.pools[size] = &p
	return &PoolAllocator{
		Channels: channels,
		Length:   length,
		Capacity: capacity,
		pool:     &p,
	}
}

// ClearPoolAllocatorCache resets internal cache of pools and makes
// existing pools available for GC. One good use case might be when the
// application changes global buffer size.
func ClearPoolAllocatorCache() {
	cache.Lock()
	defer cache.Unlock()
	cache.pools = map[int]*pool{}
}
