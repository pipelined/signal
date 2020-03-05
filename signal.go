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
	// BitDepth4 is 4 bit depth.
	BitDepth4 BitDepth = 1 << (iota + 2)
	// BitDepth8 is 8 bit depth.
	BitDepth8
	// BitDepth16 is 16 bit depth.
	BitDepth16
	// BitDepth32 is 32 bit depth.
	BitDepth32
	// BitDepth64 is 64 bit depth.
	BitDepth64
	// BitDepth24 is 24 bit depth.
	BitDepth24 BitDepth = 24
	// MaxBitDepth is a maximum supported bit depth.
	MaxBitDepth BitDepth = 64
)

var resolutions [64]uint64

func init() {
	for i := 0; i < len(resolutions); i++ {
		resolutions[i] = 1 << i
	}
}

func (b BitDepth) String() string {
	return fmt.Sprintf("%d bits", b)
}

// SampleRate is the number of samples obtained in one second.
type SampleRate uint

// DurationOf returns time duration of samples at this sample rate.
func (rate SampleRate) DurationOf(samples int) time.Duration {
	return time.Duration(math.Round(float64(time.Second) / float64(rate) * float64(samples)))
}

// SamplesIn returns number of samples for time duration at this sample rate.
func (rate SampleRate) SamplesIn(d time.Duration) int {
	return int(math.Round(float64(rate) / float64(time.Second) * float64(d)))
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

// SignedResolution returns the signed bit resolution for a bit depth.
func (b BitDepth) SignedResolution() uint64 {
	if b == 0 {
		return 1
	}
	b--
	return resolutions[b]
}

// Size of non-interleaved data.
func (ints InterInt) Size() int {
	return int(math.Ceil(float64(len(ints.Data)) / float64(ints.NumChannels)))
}

// AsFloat64 allocates new Float64 buffer of the same
// size and copies signal values there.
func (ints InterInt) AsFloat64() Float64 {
	if ints.Data == nil || ints.NumChannels == 0 {
		return nil
	}
	floats := make([][]float64, ints.NumChannels)

	for i := range floats {
		floats[i] = make([]float64, ints.Size())
	}
	ints.CopyToFloat64(floats)
	return floats
}

// CopyToFloat64 buffer the values of InterInt buffer.
// If number of channels is not equal, function will panic.
func (ints InterInt) CopyToFloat64(floats Float64) {
	if ints.NumChannels != floats.NumChannels() {
		panic(fmt.Errorf("unexpected number of channels in destination buffer: expected %v got %v", ints.NumChannels, floats.NumChannels()))
	}
	// get resolution of bit depth.
	res := ints.BitDepth.SignedResolution()
	// determine the divider for bit depth conversion.
	divider := float64(res)
	// determine the shift for signed-unsigned conversion.
	var shift uint64
	if ints.Unsigned {
		shift = res - 1
	}

	for i := range floats {
		for pos, j := i, 0; pos < len(ints.Data) && j < len(floats[i]); pos, j = pos+ints.NumChannels, j+1 {
			floats[i][j] = float64(ints.Data[pos]-int(shift)) / divider
		}
	}
}

// AsInterInt allocates new interleaved int buffer of
// the same size and copies signal values there.
// If unsigned is true, then all values are shifted
// and result will be in unsigned range.
func (floats Float64) AsInterInt(bitDepth BitDepth, unsigned bool) InterInt {
	numChannels := floats.NumChannels()
	if numChannels == 0 {
		return InterInt{}
	}

	ints := InterInt{
		Data:        make([]int, len(floats[0])*numChannels),
		NumChannels: numChannels,
		BitDepth:    bitDepth,
		Unsigned:    unsigned,
	}

	floats.CopyToInterInt(ints)
	return ints
}

// CopyToInterInt buffer the values of Float64 buffer.
// If number of channels is not equal, function will panic.
func (floats Float64) CopyToInterInt(ints InterInt) {
	if floats.NumChannels() != ints.NumChannels {
		panic(fmt.Errorf("unexpected number of channels in destination buffer: expected %v got %v", floats.NumChannels(), ints.NumChannels))
	}
	// get resolution of bit depth
	res := ints.BitDepth.SignedResolution()
	// determine the multiplier for bit depth conversion
	multiplier := float64(res)
	// determine the shift for signed-unsigned conversion
	var shift uint64
	if ints.Unsigned {
		shift = res - 1
	}

	size := ints.Size()
	for j := range floats {
		for i := 0; i < len(floats[j]) && i < size; i++ {
			ints.Data[i*ints.NumChannels+j] = int(floats[j][i]*multiplier) + int(shift)
		}
	}
}

// Float64Buffer returns an Float64 buffer of specified dimentions.
func Float64Buffer(numChannels, bufferSize int) Float64 {
	result := make([][]float64, numChannels)
	for i := range result {
		result[i] = make([]float64, bufferSize)
	}
	return result
}

// NumChannels returns number of channels in this sample slice.
func (floats Float64) NumChannels() int {
	return len(floats)
}

// Size returns number of samples in single block in this sample slice.
func (floats Float64) Size() int {
	if floats.NumChannels() == 0 {
		return 0
	}
	return len(floats[0])
}

// Append buffers to existing one.
// New buffer is returned if b is nil.
func (floats Float64) Append(source Float64) Float64 {
	if floats == nil {
		floats = make([][]float64, source.NumChannels())
		for i := range floats {
			floats[i] = make([]float64, 0, source.Size())
		}
	}
	for i := range floats {
		floats[i] = append(floats[i], source[i]...)
	}
	return floats
}

// Slice creates a new buffer that refers to floats data from start
// position with defined legth. Shorten block is returned if buffer
// doesn't have enough samples. If start is less than 0 or more than
// buffer size, nil is returned. If len goes beyond the buffer size,
// it's truncated up to length of the buffer.
func (floats Float64) Slice(start, len int) Float64 {
	if floats == nil || start >= floats.Size() || start < 0 {
		return nil
	}
	end := start + len
	result := make([][]float64, floats.NumChannels())
	for i := range floats {
		if end > floats.Size() {
			end = floats.Size()
		}
		result[i] = floats[i][start:end]
	}
	return result
}

// Sum adds values from one buffer to another.
// The lesser dimensions are used.
func (floats Float64) Sum(b Float64) Float64 {
	if floats == nil {
		return nil
	}

	for i := 0; i < len(floats) && i < len(b); i++ {
		for j := 0; j < len(floats[i]) && j < len(b[i]); j++ {
			floats[i][j] += b[i][j]
		}
	}
	return floats
}
