// Package signal provides functionality for manipulate digital signals and its attributes.
package signal

import (
	"fmt"
	"math"
	"time"
)

// BitDepth is the number of bits of information in each sample.
type BitDepth uint

const (
	// BitDepth8 is 8 bit depth.
	BitDepth8 = BitDepth(8)
	// BitDepth16 is 16 bit depth.
	BitDepth16 = BitDepth(16)
	// BitDepth24 is 32 bit depth.
	BitDepth24 = BitDepth(24)
	// BitDepth32 is 32 bit depth.
	BitDepth32 = BitDepth(32)
)

func (b BitDepth) String() string {
	return fmt.Sprintf("%d-bit", b)
}

// SampleRate is the number of samples obtained in one second.
type SampleRate uint

// DurationOf returns time duration of samples at this sample rate.
func (rate SampleRate) DurationOf(samples int64) time.Duration {
	return time.Duration(math.Round(float64(time.Second) / float64(rate) * float64(samples)))
}

// SamplesIn returns number of samples for time duration at this sample rate.
func (rate SampleRate) SamplesIn(d time.Duration) int64 {
	return int64(math.Round(float64(rate) / float64(time.Second) * float64(d)))
}

// Float64 is a non-interleaved float64 signal.
type Float64 [][]float64

// Int is a non-interleaved int signal.
type Int [][]int

// InterInt is an interleaved int signal.
type InterInt struct {
	Data        []int
	NumChannels int
	BitDepth
	Unsigned bool
}

// resolution returns a half resolution for a passed bit depth.
// example: bit depth of 8 bits has resolution of (2^8)/2 -1 ie 127.
func resolution(bitDepth BitDepth) int {
	if bitDepth == 0 {
		return 1
	}
	return 1<<(bitDepth-1) - 1
}

// AsFloat64 converts interleaved int signal to float64.
func (ints InterInt) AsFloat64() Float64 {
	if ints.Data == nil || ints.NumChannels == 0 {
		return nil
	}
	floats := make([][]float64, ints.NumChannels)
	bufSize := int(math.Ceil(float64(len(ints.Data)) / float64(ints.NumChannels)))

	// get resolution of bit depth
	res := resolution(ints.BitDepth)
	// determine the divider for bit depth conversion
	divider := float64(res)
	// determine the shift for signed-unsigned conversion
	shift := 0
	if ints.Unsigned {
		shift = res
	}

	for i := range floats {
		floats[i] = make([]float64, bufSize)
		pos := 0
		for j := i; j < len(ints.Data); j = j + ints.NumChannels {
			floats[i][pos] = float64(ints.Data[j]-shift) / divider
			pos++
		}
	}
	return floats
}

// AsInterInt converts float64 signal to interleaved int.
// If unsigned is true, then all values are shifted and result will be in unsigned range.
func (floats Float64) AsInterInt(bitDepth BitDepth, unsigned bool) []int {
	var numChannels int
	if numChannels = len(floats); numChannels == 0 {
		return nil
	}

	// get resolution of bit depth
	res := resolution(bitDepth)
	// determine the multiplier for bit depth conversion
	multiplier := float64(res)
	// determine the shift for signed-unsigned conversion
	shift := 0
	if unsigned {
		shift = res
	}

	ints := make([]int, len(floats[0])*numChannels)

	for j := range floats {
		for i := range floats[j] {
			ints[i*numChannels+j] = int(floats[j][i]*multiplier) + shift
		}
	}
	return ints
}

// Float64Buffer returns an Float64 buffer of specified dimentions.
// Optional value could be provided to mock signal.
func Float64Buffer(numChannels, bufferSize int, value float64) Float64 {
	result := make([][]float64, numChannels)
	for i := range result {
		result[i] = make([]float64, bufferSize)
		if value != 0 {
			for j := range result[i] {
				result[i][j] = value
			}
		}
	}
	return result
}

// NumChannels returns number of channels in this sample slice
func (floats Float64) NumChannels() int {
	return len(floats)
}

// Size returns number of samples in single block in this sample slice
func (floats Float64) Size() int {
	if floats.NumChannels() == 0 {
		return 0
	}
	return len(floats[0])
}

// Append buffers set to existing one one
// new buffer is returned if b is nil
func (floats Float64) Append(source Float64) Float64 {
	if floats == nil {
		floats = make([][]float64, source.NumChannels())
		for i := range floats {
			floats[i] = make([]float64, 0, source.Size())
		}
	}
	for i := range source {
		floats[i] = append(floats[i], source[i]...)
	}
	return floats
}

// Slice creates a new copy of buffer from start position with defined legth
// if buffer doesn't have enough samples - shorten block is returned
//
// if start >= buffer size, nil is returned
// if start + len >= buffer size, len is decreased till the end of slice
// if start < 0, nil is returned
func (floats Float64) Slice(start int, len int) Float64 {
	if floats == nil || start >= floats.Size() || start < 0 {
		return nil
	}
	end := start + len
	result := make([][]float64, floats.NumChannels())
	for i := range floats {
		if end > floats.Size() {
			end = floats.Size()
		}
		result[i] = append(result[i], floats[i][start:end]...)
	}
	return result
}
