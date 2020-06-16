package signal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"pipelined.dev/signal"
)

var (
	_ signal.Signed = signal.Allocator{}.Int8(signal.MaxBitDepth)
	_ signal.Signed = signal.Allocator{}.Int16(signal.MaxBitDepth)
	_ signal.Signed = signal.Allocator{}.Int32(signal.MaxBitDepth)
	_ signal.Signed = signal.Allocator{}.Int64(signal.MaxBitDepth)

	_ signal.Unsigned = signal.Allocator{}.Uint8(signal.MaxBitDepth)
	_ signal.Unsigned = signal.Allocator{}.Uint16(signal.MaxBitDepth)
	_ signal.Unsigned = signal.Allocator{}.Uint32(signal.MaxBitDepth)
	_ signal.Unsigned = signal.Allocator{}.Uint64(signal.MaxBitDepth)

	_ signal.Floating = signal.Allocator{}.Float32()
	_ signal.Floating = signal.Allocator{}.Float64()
)

type expected struct {
	length   int
	capacity int
	data     interface{}
}

func testOk(r signal.Signal, ex expected) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		if ex.capacity != 0 {
			assertEqual(t, "capacity", r.Capacity(), ex.capacity)
		}
		if ex.length != 0 {
			assertEqual(t, "length", r.Length(), ex.length)
		}
		assertEqual(t, "slices", result(r), ex.data)
	}
}

func ExampleFloatingAsFloating() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}

	// 64-bit floating to 32-bit floating signal
	f64, f32 := alloc.Float64(), alloc.Float32()
	// write float32 values to input
	signal.WriteFloat64([]float64{1, 0, -1}, f64)
	// convert float32 input to float64 output
	signal.FloatingAsFloating(f64, f32)

	result := make([]float32, 3)
	// read result
	signal.ReadFloat32(f32, result)
	fmt.Println(result)
	// Output:
	// [1 0 -1]
}

func ExampleFloatingAsSigned() {
	values := 7
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 64-bit floating to 64-bit signed signal
	f64, i64 := alloc.Float64(), alloc.Int64(signal.BitDepth64)
	// write float64 values to input
	signal.WriteFloat64(
		[]float64{
			math.Nextafter(1, 2),
			1,
			math.Nextafter(1, 0),
			0,
			math.Nextafter(-1, 0),
			-1,
			math.Nextafter(-1, -2),
		},
		f64,
	)
	// convert floating input to signed output
	signal.FloatingAsSigned(f64, i64)

	result := make([]int64, values)
	// read result
	signal.ReadInt64(i64, result)
	fmt.Println(result)
	// Output:
	// [9223372036854775807 9223372036854775807 9223372036854774784 0 -9223372036854774784 -9223372036854775808 -9223372036854775808]
}

func ExampleFloatingAsUnsigned() {
	values := 7
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 64-bit floating to 64-bit unsigned signal
	f64, u64 := alloc.Float64(), alloc.Uint64(signal.BitDepth64)
	// write float64 values to input
	signal.WriteFloat64(
		[]float64{
			math.Nextafter(1, 2),
			1,
			math.Nextafter(1, 0),
			0,
			math.Nextafter(-1, 0),
			-1,
			math.Nextafter(-1, -2),
		},
		f64,
	)
	// convert floating input to unsigned output
	signal.FloatingAsUnsigned(f64, u64)

	result := make([]uint64, values)
	// read result
	signal.ReadUint64(u64, result)
	fmt.Println(result)
	// Output:
	// [18446744073709551615 18446744073709551615 18446744073709550592 9223372036854775808 1024 0 0]
}

func ExampleSignedAsFloating() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}

	// 8-bit signed to 64-bit floating signal
	i8, f64 := alloc.Int8(signal.BitDepth8), alloc.Float64()
	// write int8 values to input
	signal.WriteInt8([]int8{math.MaxInt8, 0, math.MinInt8}, i8)
	// convert signed input to signed output
	signal.SignedAsFloating(i8, f64)

	result := make([]float64, 3)
	// read output to the result
	signal.ReadFloat64(f64, result)
	fmt.Println(result)
	// Output:
	// [1 0 -1]
}

func ExampleSignedAsSigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}
	result := make([]int64, 3)

	// downscale 64-bit signed to 8-bit signed
	{
		i64, i8 := alloc.Int64(signal.BitDepth64), alloc.Int8(signal.BitDepth8)
		// write int64 values to input
		signal.WriteInt64([]int64{math.MaxInt64, 0, math.MinInt64}, i64)
		// convert signed input to signed output
		signal.SignedAsSigned(i64, i8)
		// read output to the result
		signal.ReadInt64(i8, result)
		fmt.Println(result)
	}

	// upscale signed 8-bit to signed 16-bit
	{
		i8, i64 := alloc.Int8(signal.BitDepth8), alloc.Int64(signal.BitDepth16)
		// write int8 values to input
		signal.WriteInt8([]int8{math.MaxInt8, 0, math.MinInt8}, i8)
		// convert signed input to signed output
		signal.SignedAsSigned(i8, i64)
		// read output to the result
		signal.ReadInt64(i64, result)
		fmt.Println(result)
	}
	// Output:
	// [127 0 -128]
	// [32767 0 -32768]
}

func ExampleSignedAsUnsigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}
	result := make([]uint64, 3)

	// downscale 64-bit signed to 8-bit unsigned
	{
		i64, u64 := alloc.Int64(signal.BitDepth64), alloc.Uint64(signal.BitDepth8)
		// write int values to input
		signal.WriteInt64([]int64{math.MaxInt64, 0, math.MinInt64}, i64)
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(i64, u64)
		// read output to the result
		signal.ReadUint64(u64, result)
		fmt.Println(result)
	}

	// upscale 8-bit signed to 16-bit unsigned
	{
		i64, u64 := alloc.Int64(signal.BitDepth8), alloc.Uint64(signal.BitDepth16)
		// write int values to input
		signal.WriteInt64([]int64{math.MaxInt8, 0, math.MinInt8}, i64)
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(i64, u64)
		// read output to the result
		signal.ReadUint64(u64, result)
		fmt.Println(result)
	}
	// Output:
	// [255 128 0]
	// [65535 32768 0]
}

func ExampleUnsignedAsFloating() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}
	result := make([]float64, 3)

	// unsigned 8-bit to 64-bit floating signal
	u64, f64 := alloc.Uint64(signal.BitDepth8), alloc.Float64()
	// write uint values to input
	signal.WriteUint64([]uint64{math.MaxUint8, math.MaxInt8 + 1, 0}, u64)
	// convert unsigned input to floating output
	signal.UnsignedAsFloating(u64, f64)
	// read output to the result
	signal.ReadFloat64(f64, result)
	fmt.Println(result)
	// Output:
	// [1 0 -1]
}

func ExampleUnsignedAsSigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}
	result := make([]int64, 3)

	// downscale 64-bit unsigned to 8-bit signed
	{
		u64, i64 := alloc.Uint64(signal.BitDepth64), alloc.Int64(signal.BitDepth8)
		// write uint values to input
		signal.WriteUint64([]uint64{math.MaxUint64, math.MaxInt64 + 1, 0}, u64)
		// convert unsigned input to signed output
		signal.UnsignedAsSigned(u64, i64)
		// read output to the result
		signal.ReadInt64(i64, result)
		fmt.Println(result)
	}

	// upscale unsigned 8-bit to signed 16-bit
	{
		u64, i64 := alloc.Uint64(signal.BitDepth8), alloc.Int64(signal.BitDepth16)
		// write uint values to input
		signal.WriteUint64([]uint64{math.MaxUint8, math.MaxInt8 + 1, 0}, u64)
		// convert unsigned input to signed output
		signal.UnsignedAsSigned(u64, i64)
		// read output to the result
		signal.ReadInt64(i64, result)
		fmt.Println(result)
	}
	// Output:
	// [127 0 -128]
	// [32767 0 -32768]
}

func ExampleUnsignedAsUnsigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}
	result := make([]uint64, 3)

	// downscale 64-bit unsigned to 8-bit unsigned
	{
		u64, u8 := alloc.Uint64(signal.BitDepth64), alloc.Uint8(signal.BitDepth8)
		// write uint values to input
		signal.WriteUint64([]uint64{math.MaxUint64, math.MaxInt64 + 1, 0}, u64)
		// convert unsigned input to unsigned output
		signal.UnsignedAsUnsigned(u64, u8)
		// read output to the result
		signal.ReadUint64(u8, result)
		fmt.Println(result)
	}

	// upscale 8-bit unsigned to 16-bit unsigned
	{
		u8, u16 := alloc.Uint8(signal.BitDepth8), alloc.Uint16(signal.BitDepth16)
		// write uint values to input
		signal.WriteUint64([]uint64{math.MaxUint16, math.MaxInt8 + 1, 0}, u8)
		// convert unsigned input to unsigned output
		signal.UnsignedAsUnsigned(u8, u16)
		// read output to the result
		signal.ReadUint64(u16, result)
		fmt.Println(result)
	}
	// Output:
	// [255 128 0]
	// [65535 32768 0]
}

func ExampleBitDepth_MaxSignedValue() {
	fmt.Println(signal.BitDepth8.MaxSignedValue())
	// Output:
	// 127
}

func ExampleBitDepth_MinSignedValue() {
	fmt.Println(signal.BitDepth8.MinSignedValue())
	// Output:
	// -128
}

func ExampleBitDepth_MaxUnsignedValue() {
	fmt.Println(signal.BitDepth8.MaxUnsignedValue())
	// Output:
	// 255
}

func ExampleBitDepth_SignedValue() {
	fmt.Println(signal.BitDepth8.SignedValue(math.MaxInt64))
	// Output:
	// 127
}

func ExampleBitDepth_UnsignedValue() {
	fmt.Println(signal.BitDepth8.UnsignedValue(math.MaxUint64))
	// Output:
	// 255
}

func ExampleSampleRate_DurationOf() {
	fmt.Println(signal.SampleRate(44100).DurationOf(88200))
	// Output:
	// 2s
}

func ExampleSampleRate_SamplesIn() {
	fmt.Println(signal.SampleRate(44100).SamplesIn(time.Second * 2))
	// Output:
	// 88200
}

func TestWrite(t *testing.T) {
	allocator := signal.Allocator{
		Capacity: 1,
		Length:   1,
		Channels: 3,
	}
	t.Run("striped int 8-bits overflow", func() func(t *testing.T) {
		buf := allocator.Int64(signal.BitDepth8)
		length := signal.WriteStripedInt(
			[][]int{
				{math.MaxInt32},
				{math.MinInt32},
				{},
			},
			buf)
		return testOk(
			buf,
			expected{
				length:   length,
				capacity: 1,
				data: [][]int64{
					{math.MaxInt8},
					{math.MinInt8},
					{0},
				},
			},
		)
	}())
	t.Run("int 8-bits overflow", func() func(t *testing.T) {
		buf := allocator.Int64(signal.BitDepth8)
		length := signal.WriteInt(
			[]int{
				math.MaxInt32,
				math.MinInt32,
			},
			buf,
		)
		return testOk(
			buf,
			expected{
				length:   length,
				capacity: 1,
				data: [][]int64{
					{math.MaxInt8},
					{math.MinInt8},
					{0},
				},
			},
		)
	}())
	t.Run("striped uint 8-bits overflow", func() func(t *testing.T) {
		buf := allocator.Uint64(signal.BitDepth8)
		length := signal.WriteStripedUint(
			[][]uint{
				{math.MaxUint32},
				{0},
				{},
			},
			buf,
		)
		return testOk(
			buf,
			expected{
				length:   length,
				capacity: 1,
				data: [][]uint64{
					{math.MaxUint8},
					{0},
					{0},
				},
			},
		)
	}())
	t.Run("uint 8-bits overflow", func() func(t *testing.T) {
		buf := allocator.Uint64(signal.BitDepth8)
		length := signal.WriteUint(
			[]uint{
				math.MaxUint32,
				0,
			},
			buf,
		)
		return testOk(
			buf,
			expected{
				length:   length,
				capacity: 1,
				data: [][]uint64{
					{math.MaxUint8},
					{0},
					{0},
				},
			},
		)
	}())
}

func TestAppendPanic(t *testing.T) {
	testPanic := func(appender signal.Signed, data signal.Signed) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			assertPanic(t, func() {
				appender.Append(data)
			})
		}
	}
	t.Run("different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		signal.Allocator{Channels: 1, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
	t.Run("different bit depth", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.BitDepth8),
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}

func result(sig signal.Signal) interface{} {
	switch src := sig.(type) {
	case signal.Int8:
		result := make([][]int8, src.Channels())
		for i := range result {
			result[i] = make([]int8, src.Length())
		}
		signal.ReadStripedInt8(src, result)
		return result
	case signal.Int16:
		result := make([][]int16, src.Channels())
		for i := range result {
			result[i] = make([]int16, src.Length())
		}
		signal.ReadStripedInt16(src, result)
		return result
	case signal.Int32:
		result := make([][]int32, src.Channels())
		for i := range result {
			result[i] = make([]int32, src.Length())
		}
		signal.ReadStripedInt32(src, result)
		return result
	case signal.Int64:
		result := make([][]int64, src.Channels())
		for i := range result {
			result[i] = make([]int64, src.Length())
		}
		signal.ReadStripedInt64(src, result)
		return result
	case signal.Uint8:
		result := make([][]uint8, src.Channels())
		for i := range result {
			result[i] = make([]uint8, src.Length())
		}
		signal.ReadStripedUint8(src, result)
		return result
	case signal.Uint16:
		result := make([][]uint16, src.Channels())
		for i := range result {
			result[i] = make([]uint16, src.Length())
		}
		signal.ReadStripedUint16(src, result)
		return result
	case signal.Uint32:
		result := make([][]uint32, src.Channels())
		for i := range result {
			result[i] = make([]uint32, src.Length())
		}
		signal.ReadStripedUint32(src, result)
		return result
	case signal.Uint64:
		result := make([][]uint64, src.Channels())
		for i := range result {
			result[i] = make([]uint64, src.Length())
		}
		signal.ReadStripedUint64(src, result)
		return result
	case signal.Float32:
		result := make([][]float32, src.Channels())
		for i := range result {
			result[i] = make([]float32, src.Length())
		}
		signal.ReadStripedFloat32(src, result)
		return result
	case signal.Float64:
		result := make([][]float64, src.Channels())
		for i := range result {
			result[i] = make([]float64, src.Length())
		}
		signal.ReadStripedFloat64(src, result)
		return result
	default:
		panic(fmt.Sprintf("unsupported result type: %T", sig))
	}
}
