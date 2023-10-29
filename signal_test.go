package signal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"pipelined.dev/signal"
)

type expected struct {
	length   int
	capacity int
	data     interface{}
}

func testOk[T signal.SignalTypes](r *signal.Buffer[T], ex expected) func(t *testing.T) {
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

func ExampleFloatAsFloat() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 64-bit floating to 32-bit floating signal
	f64, f32 := signal.Alloc[float64](alloc), signal.Alloc[float32](alloc)
	// write float32 values to input
	signal.Write(
		[]float64{
			1.5,
			1,
			0,
			-1,
			-1.5,
		},
		f64,
	)
	// convert float32 input to float64 output
	signal.FloatAsFloat(f64, f32)

	result := make([]float32, values)
	// read result
	signal.Read(f32, result)
	fmt.Println(result)
	// Output:
	// [1.5 1 0 -1 -1.5]
}

func ExampleFloatAsSigned() {
	values := 7
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 64-bit floating to 64-bit signed signal
	f64, i64 := signal.Alloc[float64](alloc), signal.Alloc[int64](alloc)
	// write float64 values to input
	signal.Write(
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
	signal.FloatAsSigned(f64, i64)

	result := make([]int64, values)
	// read result
	signal.Read(i64, result)
	fmt.Println(result)
	// Output:
	// [9223372036854775807 9223372036854775807 9223372036854774784 0 -9223372036854774784 -9223372036854775808 -9223372036854775808]
}

func ExampleFloatAsUnsigned() {
	values := 7
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 64-bit floating to 64-bit unsigned signal
	f64, u64 := signal.Alloc[float64](alloc), signal.Alloc[uint64](alloc)
	// write float64 values to input
	signal.Write(
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
	signal.FloatAsUnsigned(f64, u64)

	result := make([]uint64, values)
	// read result
	signal.Read(u64, result)
	fmt.Println(result)
	// Output:
	// [18446744073709551615 18446744073709551615 18446744073709550592 9223372036854775808 1024 0 0]
}

func ExampleSignedAsFloat() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}

	// 8-bit signed to 64-bit floating signal
	i8, f64 := signal.Alloc[int8](alloc), signal.Alloc[float64](alloc)
	// write int8 values to input
	signal.Write(
		[]int8{
			math.MaxInt8,
			math.MaxInt8 - 1,
			0,
			math.MinInt8 + 1,
			math.MinInt8,
		},
		i8,
	)
	// convert signed input to signed output
	signal.SignedAsFloat(i8, f64)

	result := make([]float64, values)
	// read output to the result
	signal.Read(f64, result)
	fmt.Println(result)
	// Output:
	// [1 0.9921259842519685 0 -0.9921875 -1]
}

func ExampleSignedAsSigned() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}
	result := make([]int64, values)

	// downscale 64-bit signed to 8-bit signed
	{
		i64, i8 := signal.Alloc[int64](alloc), signal.Alloc[int8](alloc)
		// write int64 values to input
		signal.Write(
			[]int64{
				math.MaxInt64,
				math.MaxInt64 / 2,
				0,
				math.MinInt64 / 2,
				math.MinInt64,
			},
			i64,
		)
		// convert signed input to signed output
		signal.SignedAsSigned(i64, i8)
		// read output to the result
		signal.Read(i8, result)
		fmt.Println(result)
	}

	// upscale signed 8-bit to signed 16-bit
	{
		i8, i16 := signal.Alloc[int8](alloc), signal.Alloc[int16](alloc)
		// write int8 values to input
		signal.Write([]int8{math.MaxInt8, math.MaxInt8 / 2, 0, math.MinInt8 / 2, math.MinInt8}, i8)
		// convert signed input to signed output
		signal.SignedAsSigned(i8, i16)
		// read output to the result
		signal.Read(i16, result)
		fmt.Println(result)
	}
	// Output:
	// [127 63 0 -64 -128]
	// [32767 16383 0 -16384 -32768]
}

func ExampleSignedAsUnsigned() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}
	result := make([]uint64, values)

	// downscale 64-bit signed to 8-bit unsigned
	{
		i64, u8 := signal.Alloc[int64](alloc), signal.Alloc[uint8](alloc)
		// write int values to input
		signal.Write(
			[]int64{
				math.MaxInt64,
				math.MaxInt64 / 2,
				0,
				math.MinInt64 / 2,
				math.MinInt64,
			},
			i64,
		)
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(i64, u8)
		// read output to the result
		signal.Read(u8, result)
		fmt.Println(result)
	}

	// upscale 8-bit signed to 16-bit unsigned
	{
		i8, u16 := signal.Alloc[int8](alloc), signal.Alloc[uint16](alloc)
		// write int values to input
		signal.Write(
			[]int8{
				math.MaxInt8,
				math.MaxInt8 / 2,
				0,
				math.MinInt8 / 2,
				math.MinInt8,
			},
			i8,
		)
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(i8, u16)
		// read output to the result
		signal.Read(u16, result)
		fmt.Println(result)
	}
	// Output:
	// [255 191 128 64 0]
	// [65535 49151 32768 16384 0]
}

func ExampleUnsignedAsFloat() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}
	result := make([]float64, values)

	// unsigned 8-bit to 64-bit floating signal
	u8, f64 := signal.Alloc[uint8](alloc), signal.Alloc[float64](alloc)
	// write uint values to input
	signal.Write(
		[]uint64{
			math.MaxUint8,
			math.MaxUint8 - (math.MaxInt8+1)/2,
			math.MaxInt8 + 1,
			(math.MaxInt8 + 1) / 2,
			0,
		},
		u8,
	)
	// convert unsigned input to floating output
	signal.UnsignedAsFloat(u8, f64)
	// read output to the result
	signal.Read(f64, result)
	fmt.Println(result)
	// Output:
	// [1 0.49606299212598426 0 -0.5039370078740157 -1]
}

func ExampleUnsignedAsSigned() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}
	result := make([]int64, values)

	// downscale 64-bit unsigned to 8-bit signed
	{
		u64, i8 := signal.Alloc[uint64](alloc), signal.Alloc[int8](alloc)
		// write uint values to input
		signal.Write(
			[]uint64{
				math.MaxUint64,
				math.MaxUint64 - (math.MaxInt64+1)/2,
				math.MaxInt64 + 1,
				(math.MaxInt64 + 1) / 2,
				0,
			},
			u64,
		)
		// convert unsigned input to signed output
		signal.UnsignedAsSigned(u64, i8)
		// read output to the result
		signal.Read(i8, result)
		fmt.Println(result)
	}

	// upscale unsigned 8-bit to signed 16-bit
	{
		u8, i16 := signal.Alloc[uint8](alloc), signal.Alloc[int16](alloc)
		// write uint values to input
		signal.Write(
			[]uint64{
				math.MaxUint8,
				math.MaxUint8 - (math.MaxInt8+1)/2,
				math.MaxInt8 + 1,
				(math.MaxInt8 + 1) / 2,
				0,
			},
			u8,
		)
		// convert unsigned input to signed output
		signal.UnsignedAsSigned(u8, i16)
		// read output to the result
		signal.Read(i16, result)
		fmt.Println(result)
	}
	// Output:
	// [127 63 0 -64 -128]
	// [32767 16383 0 -16384 -32768]
}

func ExampleUnsignedAsUnsigned() {
	values := 5
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: values,
		Length:   values,
	}
	result := make([]uint64, values)

	// downscale 64-bit unsigned to 8-bit unsigned
	{
		u64, u8 := signal.Alloc[uint64](alloc), signal.Alloc[uint8](alloc)
		// write uint values to input
		signal.Write(
			[]uint64{
				math.MaxUint64,
				math.MaxUint64 - (math.MaxInt64+1)/2,
				math.MaxInt64 + 1,
				(math.MaxInt64 + 1) / 2,
				0,
			},
			u64,
		)
		// convert unsigned input to unsigned output
		signal.UnsignedAsUnsigned(u64, u8)
		// read output to the result
		signal.Read(u8, result)
		fmt.Println(result)
	}

	// upscale 8-bit unsigned to 16-bit unsigned
	{
		u8, u16 := signal.Alloc[uint8](alloc), signal.Alloc[uint16](alloc)
		// write uint values to input
		signal.Write(
			[]uint64{
				math.MaxUint8,
				math.MaxUint8 - (math.MaxInt8+1)/2,
				math.MaxInt8 + 1,
				(math.MaxInt8 + 1) / 2,
				0,
			},
			u8,
		)
		// convert unsigned input to unsigned output
		signal.UnsignedAsUnsigned(u8, u16)
		// read output to the result
		signal.Read(u16, result)
		fmt.Println(result)
	}
	// Output:
	// [255 191 128 64 0]
	// [65535 49151 32768 16384 0]
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

func ExampleFrequency_Duration() {
	fmt.Println(signal.Frequency(44100).Duration(88200))
	// Output:
	// 2s
}

func ExampleFrequency_Events() {
	fmt.Println(signal.Frequency(44100).Events(time.Second * 2))
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
		buf := signal.Alloc[int64](allocator)
		length := signal.WriteStriped(
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
					{math.MaxInt32},
					{math.MinInt32},
					{0},
				},
			},
		)
	}())
	t.Run("int 8-bits overflow", func() func(t *testing.T) {
		buf := signal.Alloc[int64](allocator)
		length := signal.Write(
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
					{math.MaxInt32},
					{math.MinInt32},
					{0},
				},
			},
		)
	}())
	t.Run("striped uint 8-bits overflow", func() func(t *testing.T) {
		buf := signal.Alloc[uint64](allocator)
		length := signal.WriteStriped(
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
					{math.MaxUint32},
					{0},
					{0},
				},
			},
		)
	}())
	t.Run("uint 8-bits overflow", func() func(t *testing.T) {
		buf := signal.Alloc[uint64](allocator)
		length := signal.Write(
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
					{math.MaxUint32},
					{0},
					{0},
				},
			},
		)
	}())
}

func TestAppendPanic(t *testing.T) {
	testPanic := func(appender *signal.Buffer[int64], data *signal.Buffer[int64]) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			assertPanic(t, func() {
				appender.Append(data)
			})
		}
	}
	t.Run("different channels", testPanic(
		signal.Alloc[int64](signal.Allocator{
			Capacity: 1,
			Channels: 2,
		}),
		signal.Alloc[int64](signal.Allocator{
			Capacity: 1,
			Channels: 1,
		}),
	))
}

func TestConversions(t *testing.T) {
	t.Skip()
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
		Length:   3,
	}

	// floating buf
	floating := signal.Alloc[float64](alloc)
	floats := [][]float64{
		{-1, 0, 1},
	}
	signal.WriteStriped(floats, floating)

	// signed buf
	signed := signal.Alloc[int64](alloc)
	ints := [][]int64{
		{math.MinInt64, 0, math.MaxInt64},
	}
	signal.WriteStriped(ints, signed)

	// unsigned buf
	unsigned := signal.Alloc[uint64](alloc)
	uints := [][]uint64{
		{0, math.MaxInt64 + 1, math.MaxUint64},
	}
	signal.WriteStriped(uints, unsigned)

	t.Run("floating", func() func(*testing.T) {
		output := signal.Alloc[float64](alloc)
		return func(t *testing.T) {
			signal.FloatAsFloat(floating, output)
			assertEqual(t, "floating ", result(output), floats)
			signal.SignedAsFloat(signed, output)
			assertEqual(t, "signed", result(output), floats)
			signal.UnsignedAsFloat(unsigned, output)
			assertEqual(t, "unsigned", result(output), floats)
		}
	}())
	t.Run("signed", func() func(*testing.T) {
		output := signal.Alloc[int64](alloc)
		return func(t *testing.T) {
			signal.FloatAsSigned(floating, output)
			assertEqual(t, "floating ", result(output), ints)
			signal.SignedAsSigned(signed, output)
			assertEqual(t, "signed", result(output), ints)
			signal.UnsignedAsSigned(unsigned, output)
			assertEqual(t, "unsigned", result(output), ints)
		}
	}())
	t.Run("unsigned", func() func(*testing.T) {
		output := signal.Alloc[uint64](alloc)
		return func(t *testing.T) {
			signal.FloatAsUnsigned(floating, output)
			assertEqual(t, "floating ", result(output), uints)
			signal.SignedAsUnsigned(signed, output)
			assertEqual(t, "signed", result(output), uints)
			signal.UnsignedAsUnsigned(unsigned, output)
			assertEqual(t, "unsigned", result(output), uints)
		}
	}())
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

func result[T signal.SignalTypes](sig *signal.Buffer[T]) [][]T {
	result := make([][]T, sig.Channels())
	for i := range result {
		result[i] = make([]T, sig.Length())
	}
	signal.ReadStriped[T](sig, result)
	return result
}
