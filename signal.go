package signal

//go:generate go run gen.go

import (
	"fmt"
	"math"
	"reflect"
	"time"

	"golang.org/x/exp/constraints"
)

type (
	SignalTypes interface {
		constraints.Float | constraints.Integer
	}

	GenSig[T SignalTypes] interface {
		Capacity() int
		Channels() int
		Length() int
		Len() int
		Cap() int
		BufferIndex(channel int, index int) int
		Slice(start int, end int) GenSig[T]
		Channel(channel int) C[T]
		AppendSample(value T)
		Append(GenSig[T])
		Sample(index int) T
		SetSample(index int, value T)
	}

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
		BufferIndex(channel int, index int) int
	}

	// Fixed is a digital signal represented with fixed-point values.
	Fixed interface {
		Signal
		BitDepth() BitDepth
	}

	// Signed is a digital signal represented with signed fixed-point values.
	Signed interface {
		Fixed
		Slice(start int, end int) Signed
		Channel(channel int) Signed
		Append(Signed)
		AppendSample(value int64)
		Sample(index int) int64
		SetSample(index int, value int64)
	}

	// Unsigned is a digital signal represented with unsigned fixed-point values.
	Unsigned interface {
		Fixed
		Slice(start int, end int) Unsigned
		Channel(channel int) Unsigned
		Append(Unsigned)
		AppendSample(value uint64)
		Sample(index int) uint64
		SetSample(index int, value uint64)
	}

	// Floating is a digital signal represented with floating-point values.
	Floating interface {
		Signal
		Slice(start int, end int) Floating
		Channel(channel int) Floating
		Append(Floating)
		AppendSample(value float64)
		Sample(index int) float64
		SetSample(index int, value float64)
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

// MaxSignedValue returns the maximum signed value for the bit depth.
func (b BitDepth) MaxSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return 1<<(b-1) - 1
}

// MaxUnsignedValue returns the maximum unsigned value for the bit depth.
func (b BitDepth) MaxUnsignedValue() uint64 {
	if b == 0 {
		return 0
	}
	return 1<<b - 1
}

// MinSignedValue returns the minimum signed value for the bit depth.
func (b BitDepth) MinSignedValue() int64 {
	if b == 0 {
		return 0
	}
	return -1 << (b - 1)
}

func wrapIntegerSample[T constraints.Integer](bd BitDepth) wrapSampleFunc[T] {
	newT := new(T)
	switch any(newT).(type) {
	case *int, *int8, *int16, *int32, *int64:
		return func(v T) T {
			return T(bd.SignedValue(int64(v)))
		}
	case *uint, *uint8, *uint16, *uint32, *uint64, *uintptr:
		return func(v T) T {
			return T(bd.UnsignedValue(uint64(v)))
		}
	}
	panic(fmt.Sprintf("unsupported type: %T", newT))
}

// UnsignedValue clips the unsigned signal value to the given bit depth
// range.
func (b BitDepth) UnsignedValue(val uint64) uint64 {
	max := b.MaxUnsignedValue()
	if val > max {
		return max
	}
	return val
}

// SignedValue clips the signed signal value to the given bit depth range.
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

// Scale returns scale for bit depth requantization.
func Scale(high, low BitDepth) int64 {
	return int64(1 << (high - low))
}

// Frequency in Hertz is the number of occurrences of a repeating event per
// second. It might represent sample rate or pitch.
type Frequency float64

// Duration returns a total time duration for a number events at this
// frequency.
func (f Frequency) Duration(events int) time.Duration {
	return time.Duration(math.Round(float64(time.Second) / float64(f) * float64(events)))
}

// Events returns a number of events for time duration at this frequency.
func (f Frequency) Events(d time.Duration) int {
	return int(math.Round(float64(f) / float64(time.Second) * float64(d)))
}

// FloatAsFloat writes floating-point samples to the floating-point
// destination buffer. Both buffers must have the same number of channels,
// otherwise function will panic. Returns a number of samples written per
// channel.
func FloatAsFloat[S, D constraints.Float](src *Float[S], dst *Float[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	// determine the multiplier for bit depth conversion
	for i := 0; i < length; i++ {
		dst.SetSample(i, D(src.Sample(i)))
	}
	return min(src.Length(), dst.Length())
}

// FloatAsSigned converts floating-point samples into signed fixed-point
// and appends them to the destination buffer. The floating sample range
// [-1,1] is mapped to signed [-2^(bitDepth-1), 2^(bitDepth-1)-1]. Floating
// values beyond the range will be clipped. Buffers must have the same
// number of channels, otherwise function will panic. Returns a number of
// samples written per channel.
func FloatAsSigned[S constraints.Float, D constraints.Signed](src *Float[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	// determine the multiplier for bit depth conversion
	msv := D(dst.BitDepth().MaxSignedValue())
	for i := 0; i < length; i++ {
		var sample D
		if f := float64(src.Sample(i)); f > 0 {
			// detect overflow
			if D(f) == 0 {
				sample = D(f * float64(msv))
			} else {
				sample = msv
			}
		} else {
			// no overflow here
			sample = D(f * (float64(msv) + 1))
		}
		dst.SetSample(i, sample)
	}
	return min(src.Length(), dst.Length())
}

// FloatAsUnsigned converts floating-point samples into unsigned
// fixed-point and appends them to the destination buffer. The floating
// sample range [-1,1] is mapped to unsigned [0, 2^bitDepth-1]. Floating
// values beyond the range will be clipped. Buffers must have the same
// number of channels, otherwise function will panic.
func FloatAsUnsigned[S constraints.Float, D constraints.Unsigned](src *Float[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	// determine the multiplier for bit depth conversion
	msv := D(dst.BitDepth().MaxSignedValue())
	offset := msv + 1
	for i := 0; i < length; i++ {
		var sample D
		if f := float64(src.Sample(i)); f > 0 {
			// detect overflow
			if int64(f) == 0 {
				sample = D(f*float64(msv)) + offset
			} else {
				sample = msv + offset
			}
		} else {
			// no overflow here
			sample = D(f*(float64(msv)+1)) + offset
		}
		dst.SetSample(i, sample)
	}
	return min(src.Length(), dst.Length())
}

// SignedAsFloat converts signed fixed-point samples into floating-point
// and appends them to the destination buffer. The signed sample range
// [-2^(bitDepth-1), 2^(bitDepth-1)-1] is mapped to floating [-1,1].
// Buffers must have the same number of channels, otherwise function will
// panic.
func SignedAsFloat[S constraints.Signed, D constraints.Float](src *Integer[S], dst *Float[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	// determine the divider for bit depth conversion.
	msv := D(src.BitDepth().MaxSignedValue())
	for i := 0; i < length; i++ {
		if sample := src.Sample(i); sample > 0 {
			dst.SetSample(i, D(sample)/msv)
		} else {
			dst.SetSample(i, D(sample)/(msv+1))
		}
	}
	return min(src.Length(), dst.Length())
}

// SignedAsSigned appends signed fixed-point samples to the signed
// fixed-point destination buffer. The samples are quantized to the
// destination bit depth. Buffers must have the same number of channels,
// otherwise function will panic.
func SignedAsSigned[S, D constraints.Signed](src *Integer[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}

	// downscale
	if src.BitDepth() >= dst.BitDepth() {
		scale := S(Scale(src.BitDepth(), dst.BitDepth()))
		for i := 0; i < length; i++ {
			dst.SetSample(i, D(src.Sample(i)/scale))
		}
		return min(src.Length(), dst.Length())
	}

	// upscale
	scale := D(Scale(dst.BitDepth(), src.BitDepth()))
	for i := 0; i < length; i++ {
		if sample := src.Sample(i); sample > 0 {
			dst.SetSample(i, ((D(src.Sample(i))+1)*scale)-1)
		} else {
			dst.SetSample(i, D(src.Sample(i))*scale)
		}
	}
	return min(src.Length(), dst.Length())
}

// SignedAsUnsigned converts signed fixed-point samples into unsigned
// fixed-point and appends them to the destination buffer. The samples are
// quantized to the destination bit depth. The signed sample range
// [-2^(bitDepth-1), 2^(bitDepth-1)-1] is mapped to unsigned [0,
// 2^bitDepth-1]. Buffers must have the same number of channels, otherwise
// function will panic.
func SignedAsUnsigned[S constraints.Signed, D constraints.Unsigned](src *Integer[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}

	msv := D(dst.BitDepth().MaxSignedValue())
	// downscale
	if src.BitDepth() >= dst.BitDepth() {
		scale := S(Scale(src.BitDepth(), dst.BitDepth()))
		for i := 0; i < length; i++ {
			dst.SetSample(i, D(src.Sample(i)/scale)+msv+1)
		}
		return min(src.Length(), dst.Length())
	}

	// upscale
	scale := S(Scale(dst.BitDepth(), src.BitDepth()))
	for i := 0; i < length; i++ {
		if sample := src.Sample(i); sample > 0 {
			dst.SetSample(i, D((src.Sample(i)+1)*scale)+msv)
		} else {
			dst.SetSample(i, D(src.Sample(i)*scale)+msv+1)
		}
	}
	return min(src.Length(), dst.Length())
}

// UnsignedAsFloat converts unsigned fixed-point samples into
// floating-point and appends them to the destination buffer. The unsigned
// sample range [0, 2^bitDepth-1] is mapped to floating [-1,1]. Buffers
// must have the same number of channels, otherwise function will panic.
func UnsignedAsFloat[S constraints.Unsigned, D constraints.Float](src *Integer[S], dst *Float[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	// determine the multiplier for bit depth conversion
	msv := D(src.BitDepth().MaxSignedValue())
	for i := 0; i < length; i++ {
		if sample := src.Sample(i); sample > 0 {
			dst.SetSample(i, (D(sample)-(msv+1))/msv)
		} else {
			dst.SetSample(i, (D(sample)-(msv+1))/(msv+1))
		}
	}
	return min(src.Length(), dst.Length())
}

// UnsignedAsSigned converts unsigned fixed-point samples into signed
// fixed-point and appends them to the destination buffer. The samples are
// quantized to the destination bit depth. The unsigned sample range [0,
// 2^bitDepth-1] is mapped to signed [-2^(bitDepth-1), 2^(bitDepth-1)-1].
// Buffers must have the same number of channels, otherwise function will
// panic.
func UnsignedAsSigned[S constraints.Unsigned, D constraints.Signed](src *Integer[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}
	msv := S(src.BitDepth().MaxSignedValue())
	// downscale
	if src.BitDepth() >= dst.BitDepth() {
		scale := D(Scale(src.BitDepth(), dst.BitDepth()))
		for i := 0; i < length; i++ {
			dst.SetSample(i, D(src.Sample(i)-(msv+1))/scale)
		}
		return min(src.Length(), dst.Length())
	}

	// upscale
	scale := Scale(dst.BitDepth(), src.BitDepth())
	for i := 0; i < length; i++ {
		if sample := int64(src.Sample(i) - (msv + 1)); sample > 0 {
			dst.SetSample(i, D((sample+1)*scale-1))
		} else {
			dst.SetSample(i, D(sample*scale))
		}
	}
	return min(src.Length(), dst.Length())
}

// UnsignedAsUnsigned appends unsigned fixed-point samples to the unsigned
// fixed-point destination buffer. The samples are quantized to the
// destination bit depth. Buffers must have the same number of channels,
// otherwise function will panic.
func UnsignedAsUnsigned[S, D constraints.Unsigned](src *Integer[S], dst *Integer[D]) int {
	mustSameChannels(src.Channels(), dst.Channels())
	// cap length to destination capacity.
	length := min(src.Len(), dst.Len())
	if length == 0 {
		return 0
	}

	// downscale
	if src.BitDepth() >= dst.BitDepth() {
		scale := S(Scale(src.BitDepth(), dst.BitDepth()))
		for i := 0; i < length; i++ {
			dst.SetSample(i, D(src.Sample(i)/scale))
		}
		return min(src.Length(), dst.Length())
	}

	// upscale
	scale := D(Scale(dst.BitDepth(), src.BitDepth()))
	msv := S(src.BitDepth().MaxSignedValue())
	for i := 0; i < length; i++ {
		var sample S
		if sample = src.Sample(i); sample > msv+1 {
			dst.SetSample(i, D(sample+1)*scale-1)
		} else {
			dst.SetSample(i, D(sample)*scale)
		}
	}
	return min(src.Length(), dst.Length())
}

// BitDepth returns bit depth of the buffer.
func (bd bitDepth) BitDepth() BitDepth {
	return BitDepth(bd)
}

// Channels returns number of channels in the buffer.
func (c channels) Channels() int {
	return int(c)
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

func mustSameCapacity(c1, c2 int) {
	if c1 != c2 {
		panic("different buffer capacity")
	}
}

// ChannelLength calculates a channel length for provided buffer length and
// number of channels.
func ChannelLength(sliceLen, channels int) int {
	return int(math.Ceil(float64(sliceLen) / float64(channels)))
}

// BufferIndex calculates sample index in the buffer based on number of
// channels in the buffer, channel of the sample and sample index in the
// channel.
func (c channels) BufferIndex(channel, idx int) int {
	return int(c)*idx + channel
}

// WriteInt writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteInt(src []int, dst Signed) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, int64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedInt writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedInt(src [][]int, dst Signed) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for c := 0; c < dst.Channels(); c++ {
		for i := 0; i < written; i++ {
			if i < len(src[c]) {
				dst.SetSample(dst.BufferIndex(c, i), int64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}

// WriteUint writes values from provided slice into the buffer.
// Returns a number of samples written per channel.
func WriteUint(src []uint, dst Unsigned) int {
	length := min(dst.Len(), len(src))
	for i := 0; i < length; i++ {
		dst.SetSample(i, uint64(src[i]))
	}
	return ChannelLength(length, dst.Channels())
}

// WriteStripedUint writes values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be written. Returns a number of samples written for
// the longest channel.
func WriteStripedUint(src [][]uint, dst Unsigned) (written int) {
	mustSameChannels(dst.Channels(), len(src))
	// determine the length of longest nested slice
	for i := range src {
		if len(src[i]) > written {
			written = len(src[i])
		}
	}
	// limit a number of writes to the length of the buffer
	written = min(written, dst.Length())
	for c := 0; c < dst.Channels(); c++ {
		for i := 0; i < written; i++ {
			if i < len(src[c]) {
				dst.SetSample(dst.BufferIndex(c, i), uint64(src[c][i]))
			} else {
				dst.SetSample(dst.BufferIndex(c, i), 0)
			}
		}
	}
	return
}

// ReadInt reads values from the buffer into provided slice.
// Returns number of samples read per channel.
func ReadInt(src Signed, dst []int) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = int(src.Sample(i))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedInt reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedInt(src Signed, dst [][]int) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = int(src.Sample(src.BufferIndex(c, i)))
		}
	}
	return
}

// ReadUint reads values from the buffer into provided slice.
func ReadUint(src Unsigned, dst []uint) int {
	length := min(src.Len(), len(dst))
	for i := 0; i < length; i++ {
		dst[i] = uint(src.Sample(i))
	}
	return ChannelLength(length, src.Channels())
}

// ReadStripedUint reads values from the buffer into provided slice. The
// length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be read. Returns a number of samples read for the
// longest channel.
func ReadStripedUint(src Unsigned, dst [][]uint) (read int) {
	mustSameChannels(src.Channels(), len(dst))
	for c := 0; c < src.Channels(); c++ {
		length := min(len(dst[c]), src.Length())
		if length > read {
			read = length
		}
		for i := 0; i < length; i++ {
			dst[c][i] = uint(src.Sample(src.BufferIndex(c, i)))
		}
	}
	return
}

// alignCapacity ensures that buffer capacity is aligned with number of
// channels.
func alignCapacity(s interface{}, channels, c int) {
	reflect.ValueOf(s).Elem().SetCap(c - c%channels)
}
