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

var (
	maximum [65]uint64
	minimum [64]int64
)

func init() {
	for i := range maximum {
		maximum[i] = (1 << i) - 1
	}
	for i := range minimum {
		minimum[i] = (-1) << i
	}
}

// MaxSignedValue returns the maximum signed value for a bit depth.
func (b BitDepth) MaxSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return int64(maximum[b-1])
}

// MaxUnsignedValue returns the maximum unsigned value for a bit depth.
func (b BitDepth) MaxUnsignedValue() uint64 {
	if b == 0 {
		return 0
	}
	return maximum[b]
}

// MinSignedValue returns the minimum signed value for a bit depth.
func (b BitDepth) MinSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return minimum[b-1]
}

// UnsignedValue limits the unsigned signal value for a given bit depth.
func (b BitDepth) UnsignedValue(val uint64) uint64 {
	max := b.MaxUnsignedValue()
	switch {
	case val > max:
		return max
	default:
		return val
	}
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
	default:
		return val
	}
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
	msv := dst.BitDepth().MaxSignedValue()
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample(uint64(int64(sample)*msv + (msv + 1)))
		} else {
			dst = dst.AppendSample(uint64(int64(sample)*(msv+1) + (msv + 1)))
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
	divider := float64(src.BitDepth().MaxSignedValue())
	for pos := 0; pos < length; pos++ {
		if sample := src.Sample(pos); sample > 0 {
			dst = dst.AppendSample(float64(sample) / divider)
		} else {
			dst = dst.AppendSample(float64(sample) / (divider + 1))
		}
	}
	return dst
}

//TODO:
// SignedAsUnsigned
// SignedAsSigned
// UnsingedAsFloating
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

// WriteFloat64 writes values from provided slice into the buffer.
func WriteFloat64(s Floating, buf []float64) Floating {
	length := min(s.Cap()-s.Len(), len(buf))
	for pos := 0; pos < length; pos++ {
		s = s.AppendSample(buf[pos])
	}
	return s
}

// WriteInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteInt(s Signed, buf []int) Signed {
	length := min(s.Cap()-s.Len(), len(buf))
	for pos := 0; pos < length; pos++ {
		s = s.AppendSample(int64(buf[pos]))
	}
	return s
}

// WriteInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func WriteInt64(s Signed, buf []int64) Signed {
	length := min(s.Cap()-s.Len(), len(buf))
	for pos := 0; pos < length; pos++ {
		s = s.AppendSample(buf[pos])
	}
	return s
}

// WriteStripedFloat64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of enclosing slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length.
func WriteStripedFloat64(s Floating, buf [][]float64) Floating {
	mustSameChannels(s.Channels(), len(buf))
	var length int
	for i := range buf {
		if len(buf[i]) > length {
			length = len(buf[i])
		}
	}
	length = min(length, s.Capacity()-s.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < s.Channels(); channel++ {
			if pos < len(buf[channel]) {
				s = s.AppendSample(buf[channel][pos])
			} else {
				s = s.AppendSample(0)
			}
		}
	}
	return s
}

// WriteStripedInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt(s Signed, buf [][]int) Signed {
	mustSameChannels(s.Channels(), len(buf))
	var length int
	for i := range buf {
		if len(buf[i]) > length {
			length = len(buf[i])
		}
	}
	length = min(length, s.Capacity()-s.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < s.Channels(); channel++ {
			if pos < len(buf[channel]) {
				s = s.AppendSample(int64(buf[channel][pos]))
			} else {
				s = s.AppendSample(0)
			}
		}
	}
	return s
}

// WriteStripedInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt64(s Signed, buf [][]int64) Signed {
	mustSameChannels(s.Channels(), len(buf))
	var length int
	for i := range buf {
		if len(buf[i]) > length {
			length = len(buf[i])
		}
	}
	length = min(length, s.Capacity()-s.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < s.Channels(); channel++ {
			if pos < len(buf[channel]) {
				s = s.AppendSample(buf[channel][pos])
			} else {
				s = s.AppendSample(0)
			}
		}
	}
	return s
}

func ReadStripedFloat64(s Floating, buf [][]float64) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < min(s.Length(), len(buf[channel])); pos++ {
			buf[channel][pos] = s.Sample(s.ChannelPos(channel, pos))
		}
	}
}

func ReadStripedInt(s Signed, buf [][]int) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Length() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = int(s.Sample(s.ChannelPos(channel, pos)))
		}
	}
}

func ReadStripedInt64(s Signed, buf [][]int64) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Length() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = s.Sample(s.ChannelPos(channel, pos))
		}
	}
}
