package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how to iterate over the Buffer.
func Example_iterate() {
	// allocate int64 Buffer with 2 channels and capacity of 8 samples per channel
	buf := signal.Alloc[int64](signal.Allocator{
		Channels: 2,
		Capacity: 8,
		Length:   4,
	})

	// write striped data
	signal.WriteStriped([][]int8{{1, 1, 1, 1}, {2, 2, 2, 2}}, buf)

	// iterate over Buffer interleaved data
	for i := 0; i < buf.Len(); i++ {
		fmt.Printf("%d", buf.Sample(i))
	}

	for c := 0; c < buf.Channels(); c++ {
		fmt.Println()
		for i := 0; i < buf.Length(); i++ {
			fmt.Printf("%d", buf.Sample(buf.BufferIndex(c, i)))
		}
	}

	// Output:
	// 12121212
	// 1111
	// 2222
}
