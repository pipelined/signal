// Package signal provides functionality for manipulate digital signals and its attributes.
package signal

import (
	"math"
	"time"
)

type (
	// Signal is a buffer that contains a digital representation of a
	// physical signal that is a sampled and quantized.
	Signal interface {
		Capacity() int
		Channels() int
		Length() int
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
		Append(Signed) Signed
		Sample(channel, pos int) int64
		setSample(channel, pos int, value int64)
		setLength(int) Signed
	}

	// Unsigned is a digital signal represented with unsigned fixed-point values.
	Unsigned interface {
		Fixed
		Append(Unsigned) Unsigned
		Sample(channel, pos int) uint64
		setSample(channel, pos int, value uint64)
		setLength(int) Unsigned
	}

	// Floating is a digital signal represented with floating-point values.
	Floating interface {
		Signal
		Append(Floating) Floating
		Sample(channel, pos int) float64
		setSample(channel, pos int, value float64)
		setLength(int) Floating
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
	capacity int
	length   int
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

var resolutions [64]uint64

func init() {
	for i := range resolutions {
		resolutions[i] = 1 << i
	}
}

// MaxValue returns the maximum value for a bit depth.
func (b BitDepth) MaxValue(signed bool) uint64 {
	if signed {
		if b == 0 {
			return 1
		}
		b--
	}
	return resolutions[b] - 1
}

// MinValue returns the minimum value for a bit depth.
func (b BitDepth) MinValue() int64 {
	if b == 0 {
		return 1
	}
	b--
	return -int64(resolutions[b])
}

// UnsignedValue limits the unsigned signal value for a given bit depth.
func (b BitDepth) UnsignedValue(val uint64) uint64 {
	var (
		max = b.MaxValue(true)
	)
	switch {
	case val > max:
		return max
	default:
		return val
	}
}

// SignedValue limits the signed signal value for a given bit depth.
func (b BitDepth) SignedValue(val int64) int64 {
	var (
		max = int64(b.MaxValue(true))
		min = b.MinValue()
	)
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

// FloatingAsSigned converts floating-point signal into signed fixed-point.
func FloatingAsSigned(src Floating, dst Signed) Signed {
	channels := min(src.Channels(), dst.Channels())
	if channels == 0 {
		return dst
	}
	// cap length to destination capacity.
	length := min(src.Length(), dst.Capacity())
	if length == 0 {
		return dst
	}
	// determine the multiplier for bit depth conversion
	multiplier := float64(dst.BitDepth().MaxValue(true))

	for channel := 0; channel < channels; channel++ {
		for pos := 0; pos < length; pos++ {
			if sample := src.Sample(channel, pos); sample > 0 {
				dst.setSample(channel, pos, int64(sample*multiplier))
			} else {
				dst.setSample(channel, pos, int64(sample*(multiplier+1)))
			}
		}
	}
	return dst.setLength(length)
}

// SignedAsFloating converts signed fixed-point signal into floating-point.
func SignedAsFloating(src Signed, dst Floating) Floating {
	channels := min(src.Channels(), dst.Channels())
	if channels == 0 {
		return dst
	}
	// cap length to destination capacity.
	length := min(src.Length(), dst.Capacity())
	if length == 0 {
		return dst
	}
	// determine the divider for bit depth conversion.
	divider := float64(src.BitDepth().MaxValue(true))
	for channel := 0; channel < channels; channel++ {
		for pos := 0; pos < length; pos++ {
			if sample := src.Sample(channel, pos); sample > 0 {
				dst.setSample(channel, pos, float64(sample)/divider)
			} else {
				dst.setSample(channel, pos, float64(sample)/(divider+1))
			}
		}
	}
	return dst.setLength(length)
}

func (bd bitDepth) BitDepth() BitDepth {
	return BitDepth(bd)
}

// Length returns actual signal length in signal channel of the buffer.
func (l length) Length() int {
	return int(l)
}

func (c capacity) Capacity() int {
	return int(c)
}

// Channels returns number of channels in the buffer.
func (c channels) Channels() int {
	return int(c)
}

func interPos(channels, channel, pos int) int {
	return channels*pos + channel
}

func interLen(channels, totalLen int) int {
	if totalLen%channels > 0 {
		return totalLen/channels + 1
	}
	return totalLen / channels
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
// If the buffer already contains any data, it will be overwritten.
// The length of enclosing slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length.
func WriteFloat64(s Floating, buf [][]float64) Floating {
	mustSameChannels(s.Channels(), len(buf))
	n := 0
	for channel := 0; channel < s.Channels(); channel++ {
		pos := 0
		for pos < s.Capacity() && pos < len(buf[channel]) {
			s.setSample(channel, pos, buf[channel][pos])
			pos++
		}
		if n < pos {
			n = pos
		}
	}
	return s.setLength(n)
}

// WriteInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteInt(s Signed, buf [][]int) Signed {
	mustSameChannels(s.Channels(), len(buf))
	var n int
	for channel := 0; channel < s.Channels(); channel++ {
		pos := 0
		for pos < s.Capacity() && pos < len(buf[channel]) {
			s.setSample(channel, pos, s.BitDepth().SignedValue(int64(buf[channel][pos])))
			pos++
		}
		if n < pos {
			n = pos
		}
	}
	return s.setLength(n)
}

// WriteInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteInt64(s Signed, buf [][]int64) Signed {
	mustSameChannels(s.Channels(), len(buf))
	var n int
	for channel := 0; channel < s.Channels(); channel++ {
		pos := 0
		for pos < s.Capacity() && pos < len(buf[channel]) {
			s.setSample(channel, pos, s.BitDepth().SignedValue(buf[channel][pos]))
			pos++
		}
		if n < pos {
			n = pos
		}
	}
	return s.setLength(n)
}

// WriteStripedInt writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt(s Signed, buf []int) Signed {
	var n int
	for channel := 0; channel < s.Channels(); channel++ {
		pos := 0
		for pos < s.Capacity() && pos < interLen(s.Channels(), len(buf)) {
			s.setSample(channel, pos, s.BitDepth().SignedValue(int64(buf[interPos(s.Channels(), channel, pos)])))
			pos++
		}
		if n < pos {
			n = pos
		}
	}
	return s.setLength(n)
}

// WriteStripedInt64 writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Length is set to the longest
// nested slice length. Sample values are capped by maximum value of
// the buffer bit depth.
func WriteStripedInt64(s Signed, buf []int64) Signed {
	var n int
	for channel := 0; channel < s.Channels(); channel++ {
		pos := 0
		for pos < s.Capacity() && pos < interLen(s.Channels(), len(buf)) {
			s.setSample(channel, pos, s.BitDepth().SignedValue(buf[interPos(s.Channels(), channel, pos)]))
			pos++
		}
		if n < pos {
			n = pos
		}
	}
	return s.setLength(n)
}

func ReadFloat64(s Floating, buf [][]float64) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Capacity() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = s.Sample(channel, pos)
		}
	}
}

func ReadInt(s Signed, buf [][]int) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Capacity() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = int(s.Sample(channel, pos))
		}
	}
}

func ReadInt64(s Signed, buf [][]int64) {
	mustSameChannels(s.Channels(), len(buf))
	for channel := 0; channel < s.Channels(); channel++ {
		for pos := 0; pos < s.Capacity() && pos < len(buf[channel]); pos++ {
			buf[channel][pos] = s.Sample(channel, pos)
		}
	}
}
