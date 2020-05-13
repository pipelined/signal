// Package signal provides functionality for manipulate digital signals and its attributes.
package signal

import (
	"math"
	"time"
)

type (
	// Signal is a buffer that contains a digital representation of a
	// physical signal that is a sampled and quantized.
	// Signal types have semantics of go slices. They can be sliced
	// and appended to each other.
	Signal interface {
		Capacity() int
		Channels() int
		Length() int
		Len() int
		Cap() int
		ChannelPos(int, int) int
	}

	// Fixed is a digital signal represented with fixed-point values.
	Fixed interface {
		Signal
		BitDepth() BitDepth
		MaxBitDepth() BitDepth
	}

	// Signed is a digital signal represented with signed fixed-point values.
	Signed interface {
		Fixed
		Slice(int, int) Signed
		Append(Signed) Signed
		AppendSample(int64) Signed
		Sample(pos int) int64
		SetSample(pos int, value int64)
		Reset() Signed
	}

	// Unsigned is a digital signal represented with unsigned fixed-point values.
	Unsigned interface {
		Fixed
		Slice(int, int) Unsigned
		Append(Unsigned) Unsigned
		AppendSample(uint64) Unsigned
		Sample(pos int) uint64
		SetSample(pos int, value uint64)
		Reset() Unsigned
	}

	// Floating is a digital signal represented with floating-point values.
	Floating interface {
		Signal
		Slice(int, int) Floating
		Append(Floating) Floating
		AppendSample(float64) Floating
		Sample(pos int) float64
		SetSample(pos int, value float64)
		Reset() Floating
	}

	// Allocator provides allocation of various signal buffers.
	Allocator struct {
		Capacity int
		Channels int
	}
)

// types for buffer properties.
type (
	bitDepth BitDepth
	channels int
)

// BitDepth is the number of bits of information in each sample.
type BitDepth uint8

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
	MaxBitDepth BitDepth = BitDepth64
)

// MaxSignedValue returns the maximum signed value for a bit depth.
func (b BitDepth) MaxSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return 1<<(b-1) - 1
}

// MaxUnsignedValue returns the maximum unsigned value for a bit depth.
func (b BitDepth) MaxUnsignedValue() uint64 {
	if b == 0 {
		return 0
	}
	return 1<<b - 1
}

// MinSignedValue returns the minimum signed value for a bit depth.
func (b BitDepth) MinSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return -1 << (b - 1)
}

// UnsignedValue limits the unsigned signal value for a given bit depth.
func (b BitDepth) UnsignedValue(val uint64) uint64 {
	max := b.MaxUnsignedValue()
	switch {
	case val > max:
		return max
	}
	return val
}

// SignedValue limits the signed signal value for a given bit depth.
func (b BitDepth) SignedValue(val int64) int64 {
	max := b.MaxSignedValue()
	min := b.MinSignedValue()
	switch {
	case val < min:
		return min
	case val > max:
		return max
	}
	return val
}

func (b BitDepth) Scale(dst BitDepth) float64 {
	if b == dst {
		return 1
	}
	return float64(b.MaxSignedValue()) / float64(dst.MaxSignedValue())
}

// cap limits bit depth value to max and returns max if it's value is 0.
func (b BitDepth) cap(max BitDepth) bitDepth {
	if b == 0 || b > max {
		return bitDepth(max)
	}
	return bitDepth(b)
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

// FloatingAsFloating converts floating-point signal into floating-point.
func FloatingAsFloating(src Floating, dst Floating) Floating {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	// determine the multiplier for bit depth conversion
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(src.Sample(pos))
	}
	return dst
}

// FloatingAsSigned converts floating-point signal into signed fixed-point.
func FloatingAsSigned(src Floating, dst Signed) Signed {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	// determine the multiplier for bit depth conversion
	msv := dst.BitDepth().MaxSignedValue()
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample(int64(sample) * msv)
		} else {
			dst = dst.AppendSample(int64(sample) * (msv + 1))
		}
	}
	return dst
}

// FloatingAsUnsigned converts floating-point signal into unsigned fixed-point.
func FloatingAsUnsigned(src Floating, dst Unsigned) Unsigned {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	// determine the multiplier for bit depth conversion
	msv := uint64(dst.BitDepth().MaxSignedValue())
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample(uint64(sample)*msv + (msv + 1))
		} else {
			dst = dst.AppendSample(uint64(sample)*(msv+1) + (msv + 1))
		}
	}
	return dst
}

// SignedAsFloating converts signed fixed-point signal into floating-point.
func SignedAsFloating(src Signed, dst Floating) Floating {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	// determine the divider for bit depth conversion.
	msv := float64(src.BitDepth().MaxSignedValue())
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample(float64(sample) / msv)
		} else {
			dst = dst.AppendSample(float64(sample) / (msv + 1))
		}
	}
	return dst
}

// SignedAsSigned converts signed fixed-point signal into signed fixed-point.
func SignedAsSigned(src Signed, dst Signed) Signed {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	scale := src.BitDepth().Scale(dst.BitDepth())
	for pos := 0; pos < length; pos++ {
		if sample := float64(src.Sample(pos)) / scale; sample < math.MaxInt64 {
			dst = dst.AppendSample(int64(sample))
		} else {
			dst = dst.AppendSample(math.MaxInt64)
		}
	}
	return dst
}

// UnsignedAsFloating converts unsigned fixed-point signal into floating-point.
func UnsignedAsFloating(src Unsigned, dst Floating) Floating {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Cap()-dst.Len())
	if length == 0 {
		return dst
	}
	// determine the multiplier for bit depth conversion
	msv := float64(src.BitDepth().MaxSignedValue())
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample((float64(sample) - (msv + 1)) / msv)
		} else {
			dst = dst.AppendSample((float64(sample) - (msv + 1)) / (msv + 1))
		}
	}
	return dst
}

//TODO:
// SignedAsUnsigned
// UnsignedAsSigned
// UnsignedAsUnsigned

func (bd bitDepth) BitDepth() BitDepth {
	return BitDepth(bd)
}

// Channels returns number of channels in the buffer.
func (c channels) Channels() int {
	return int(c)
}

func (c channels) ChannelPos(channel, pos int) int {
	return int(c)*pos + channel
}

func min(v1, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func mustSameChannels(c1, c2 int) {
	if c1 != c2 {
		panic("different number of channels")
	}
}

func mustSameBitDepth(bd1, bd2 BitDepth) {
	if bd1 != bd2 {
		panic("different bit depth")
	}
}

// WriteInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteInt(src []int, dst Signed) Signed {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(int64(src[pos]))
	}
	return dst
}

// WriteStripedInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt(src [][]int, dst Signed) Signed {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(int64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}

func ReadInt(src Signed, dst []int) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = int(src.Sample(pos))
	}
}

func ReadStripedInt(src Signed, dst [][]int) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = int(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
}
