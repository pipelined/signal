package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how to use pool to allocate buffers.
func Example_pool() {
	pool := signal.GetFloatPool[float64](signal.Allocator{2, 0, 512})

	// producer allocates new buffers
	produceFunc := func(allocs int, p *signal.PFloatAllocator[float64], c chan<- *signal.Float[float64]) {
		for i := 0; i < allocs; i++ {
			buf := p.Get()
			buf.AppendSample(1.0)
			c <- buf
		}
		close(c)
	}
	// consumer processes buffers and puts them back to the pool
	consumeFunc := func(p *signal.PFloatAllocator[float64], c <-chan *signal.Float[float64], done chan struct{}) {
		for s := range c {
			fmt.Printf("Length: %d Capacity: %d\n", s.Length(), s.Capacity())
		}
		close(done)
	}

	c := make(chan *signal.Float[float64])
	done := make(chan struct{})
	go produceFunc(10, &pool, c)
	go consumeFunc(&pool, c, done)
	<-done
	// Output:
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
	// Length: 1 Capacity: 512
}
