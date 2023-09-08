package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how read and write data to the buffer.
func Example_readWrite() {
	var output []int

	// allocate int64 buffer with 2 channels and capacity of 8 samples per channel
	buf := signal.AllocInteger[int64](signal.Allocator{
		Channels: 2,
		Capacity: 8,
		Length:   8,
	}, signal.BitDepth64)

	// write striped data
	signal.WriteStriped([][]int8{{1, 1, 1, 1}, {2, 2, 2, 2}}, buf.Slice(0, 4))
	// write interleaved data
	signal.Write([]int16{11, 22, 11, 22, 11, 22, 11, 22}, buf.Slice(4, 8))

	output = make([]int, 16) // reset output
	signal.Read(buf, output) // read data into output
	fmt.Println(output)

	output = make([]int, 16)             // reset output
	signal.Read(buf.Slice(0, 0), output) // reset buffer length to 0 and read data into output
	fmt.Println(output)

	// Output:
	// [1 2 1 2 1 2 1 2 11 22 11 22 11 22 11 22]
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
}
