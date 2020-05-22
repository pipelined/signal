package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how to use pool to allocate buffers.
func Example_pool() {
	pool := signal.Allocator{Channels: 2, Capacity: 512}.Pool()

	// producer allocates new buffers
	produceFunc := func(allocs int, p *signal.Pool, c chan<- signal.Floating) {
		for i := 0; i < allocs; i++ {
			c <- p.GetFloat64().AppendSample(1.0)
		}
		close(c)
	}
	// consumer processes buffers and puts them back to the pool
	consumeFunc := func(p *signal.Pool, c <-chan signal.Floating, done chan struct{}) {
		for s := range c {
			fmt.Printf("Length: %d Capacity: %d\n", s.Length(), s.Capacity())
		}
		close(done)
	}

	c := make(chan signal.Floating)
	done := make(chan struct{})
	go produceFunc(10, pool, c)
	go consumeFunc(pool, c, done)
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
