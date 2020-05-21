package signal_test

import (
	"fmt"

	"pipelined.dev/signal"
)

// This example demonstrates how read and write data to the buffer.
func Example() {
	var output []int

	// allocate int64 buffer with 2 channels and capacity of 8 samples per channel
	buf := signal.Allocator{Channels: 2, Capacity: 8}.Int64(signal.BitDepth64)

	// write interleaved data
	buf = signal.WriteStripedInt8([][]int8{{1, 1, 1, 1}, {2, 2, 2, 2}}, buf)
	// write striped data
	buf = signal.WriteInt16([]int16{11, 22, 11, 22, 11, 22, 11, 22}, buf)

	output = make([]int, 16) // reset output

	signal.ReadInt(buf, output) // read data into output
	fmt.Println(output)

	output = make([]int, 16) // reset output
	buf = buf.Reset()        // reset buffer length to 0

	signal.ReadInt(buf, output) // read data into output
	fmt.Println(output)

	// Output:
	// [1 2 1 2 1 2 1 2 11 22 11 22 11 22 11 22]
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
}
