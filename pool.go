package signal

import "sync"

var cache = struct {
	sync.Mutex
	pools map[int]*Pool
}{
	pools: map[int]*Pool{},
}

// Pool allows to decrease a number of allocations at runtime. Internally
// it relies on sync.Pool to manage objects in memory. It provides a pool
// per signal buffer type.
type Pool struct {
	allocSize int
	channels  int
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

// GetPool returns pool for provided buffer dimensions. Pools are cached
// internally, so multiple calls for with same dimentions will return the
// same pool instance.
func GetPool(channels, capacity int) *Pool {
	size := channels * capacity
	cache.Lock()
	defer cache.Unlock()

	if p, ok := cache.pools[size]; ok {
		return p
	}

	p := Pool{
		channels:  channels,
		allocSize: size,
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
	cache.pools[size] = &p
	return &p
}

// ResetPoolCache resets internal cache of pools and makes existing pools
// available for GC. One good use case might be when the application
// changes global buffer size.
func ResetPoolCache() {
	cache.Lock()
	defer cache.Unlock()
	cache.pools = map[int]*Pool{}
}
