package signal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

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
	}
	result := make([]float64, 3)

	// float64 to float32 signal
	// read output to the result
	signal.ReadFloat64(
		// convert signed input to signed output
		signal.FloatingAsFloating(
			// write int values to input
			signal.WriteFloat64(
				[]float64{1, 0, -1},
				alloc.Float64(),
			),
			alloc.Float32(),
		),
		result,
	)

	fmt.Println(result)
	// Output:
	// [1 0 -1]
}

func ExampleFloatingAsSigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
	}
	result := make([]int64, 3)

	// 64-bit floating to 64-bit signed signal
	// read output to the result
	signal.ReadInt64(
		// convert floating input to signed output
		signal.FloatingAsSigned(
			// write int values to input
			signal.WriteFloat64(
				[]float64{1, 0, -1},
				alloc.Float64(),
			),
			alloc.Int64(signal.BitDepth64),
		),
		result,
	)

	fmt.Println(result)
	// Output:
	// [9223372036854775807 0 -9223372036854775808]
}

func ExampleFloatingAsUnsigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
	}
	result := make([]uint64, 3)

	// 64-bit floating to 64-bit unsigned signal
	// read output to the result
	signal.ReadUint64(
		// convert floating input to unsigned output
		signal.FloatingAsUnsigned(
			// write int values to input
			signal.WriteFloat64(
				[]float64{1, 0, -1},
				alloc.Float64(),
			),
			alloc.Uint64(signal.BitDepth64),
		),
		result,
	)

	fmt.Println(result)
	// Output:
	// [18446744073709551615 9223372036854775808 0]
}

func ExampleSignedAsFloating() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
	}
	result := make([]float64, 3)

	// downscale signed 8-bit to floating signal
	// read output to the result
	signal.ReadFloat64(
		// convert signed input to signed output
		signal.SignedAsFloating(
			// write int values to input
			signal.WriteInt64(
				[]int64{math.MaxInt8, 0, math.MinInt8},
				alloc.Int8(signal.BitDepth8),
			),
			alloc.Float64(),
		),
		result,
	)

	fmt.Println(result)
	// Output:
	// [1 0 -1]
}

func ExampleSignedAsSigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
	}
	result := make([]int64, 3)

	// downscale 64-bit signed to 8-bit signed
	// read output to the result
	signal.ReadInt64(
		// convert signed input to signed output
		signal.SignedAsSigned(
			// write int values to input
			signal.WriteInt64(
				[]int64{math.MaxInt64, 0, math.MinInt64},
				alloc.Int64(signal.BitDepth64),
			),
			alloc.Int64(signal.BitDepth8),
		),
		result,
	)
	fmt.Println(result)

	// upscale signed 8-bit to signed 16-bit
	// read output to the result
	signal.ReadInt64(
		// convert signed input to signed output
		signal.SignedAsSigned(
			// write int values to input
			signal.WriteInt64(
				[]int64{math.MaxInt8, 0, math.MinInt8},
				alloc.Int64(signal.BitDepth8),
			),
			alloc.Int64(signal.BitDepth16),
		),
		result,
	)
	fmt.Println(result)
	// Output:
	// [127 0 -128]
	// [32767 0 -32768]
}

func ExampleSignedAsUnsigned() {
	alloc := signal.Allocator{
		Channels: 1,
		Capacity: 3,
	}
	result := make([]uint64, 3)

	// downscale 64-bit signed to 8-bit unsigned
	// read output to the result
	signal.ReadUint64(
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(
			// write int values to input
			signal.WriteInt64(
				[]int64{math.MaxInt64, 0, math.MinInt64},
				alloc.Int64(signal.BitDepth64),
			),
			alloc.Uint64(signal.BitDepth8),
		),
		result,
	)
	fmt.Println(result)

	// upscale signed 8-bit to unsigned 16-bit
	// read output to the result
	signal.ReadUint64(
		// convert signed input to unsigned output
		signal.SignedAsUnsigned(
			// write int values to input
			signal.WriteInt64(
				[]int64{math.MaxInt8, 0, math.MinInt8},
				alloc.Int64(signal.BitDepth8),
			),
			alloc.Uint64(signal.BitDepth16),
		),
		result,
	)
	fmt.Println(result)
	// Output:
	// [255 128 0]
	// [65535 32768 0]
}

func TestUnsignedAsFloating(t *testing.T) {
	t.Run("8 bits", testOk(
		signal.UnsignedAsFloating(
			signal.WriteStripedUint64(
				[][]uint64{{
					0,
					128,
					255,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth8),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Float64(),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]float64{{
				-1,
				0,
				1,
			}},
		},
	))
	t.Run("64 bits", testOk(
		signal.UnsignedAsFloating(
			signal.WriteStripedUint64(
				[][]uint64{{
					0,
					1 << 63,
					1<<64 - 1,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth64),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Float64(),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]float64{{
				-1,
				0,
				1,
			}},
		},
	))
}

func TestUnsignedAsSigned(t *testing.T) {
	t.Run("64 bits", testOk(
		signal.UnsignedAsSigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint64,
					math.MaxInt64 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth64),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Int64(signal.BitDepth64),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{{
				math.MaxInt64,
				0,
				math.MinInt64,
			}},
		},
	))
	t.Run("64 bits to 8 bits", testOk(
		signal.UnsignedAsSigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint64,
					math.MaxInt64 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth64),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Int64(signal.BitDepth8),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{{
				math.MaxInt8,
				0,
				math.MinInt8,
			}},
		},
	))
	t.Run("8 bits to 64 bits", testOk(
		signal.UnsignedAsSigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint8,
					math.MaxInt8 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth8),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Int64(signal.BitDepth64),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{{
				math.MaxInt64,
				0,
				math.MinInt64,
			}},
		},
	))
	t.Run("8 bits to 16 bits", testOk(
		signal.UnsignedAsSigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint8,
					math.MaxInt8 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth8),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Int64(signal.BitDepth16),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{{
				math.MaxInt16,
				0,
				math.MinInt16,
			}},
		},
	))
}

func TestUnsignedAsUnsigned(t *testing.T) {
	t.Run("64 bits", testOk(
		signal.UnsignedAsUnsigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint64,
					math.MaxInt64 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth64),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Uint64(signal.BitDepth64),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]uint64{{
				math.MaxUint64,
				math.MaxInt64 + 1,
				0,
			}},
		},
	))
	t.Run("64 bits to 8 bits", testOk(
		signal.UnsignedAsUnsigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint64,
					math.MaxInt64 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth64),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Uint64(signal.BitDepth8),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]uint64{{
				math.MaxUint8,
				math.MaxInt8 + 1,
				0,
			}},
		},
	))
	t.Run("8 bits to 64 bits", testOk(
		signal.UnsignedAsUnsigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint8,
					math.MaxInt8 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth8),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Uint64(signal.BitDepth64),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]uint64{{
				math.MaxUint64,
				math.MaxInt64 + 1,
				0,
			}},
		},
	))
	t.Run("8 bits to 16 bits", testOk(
		signal.UnsignedAsUnsigned(
			signal.WriteStripedUint64(
				[][]uint64{{
					math.MaxUint8,
					math.MaxInt8 + 1,
					0,
				}},
				signal.Allocator{
					Channels: 1,
					Capacity: 3,
				}.Uint64(signal.BitDepth8),
			),
			signal.Allocator{
				Channels: 1,
				Capacity: 3,
			}.Uint64(signal.BitDepth16),
		),
		expected{
			length:   3,
			capacity: 3,
			data: [][]uint64{{
				math.MaxUint16,
				math.MaxInt16 + 1,
				0,
			}},
		},
	))
}

func TestWrite(t *testing.T) {
	t.Run("striped int 8-bits overflow", testOk(
		signal.WriteStripedInt(
			[][]int{
				{math.MaxInt32},
				{math.MinInt32},
			},
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int 8-bits overflow", testOk(
		signal.WriteInt(
			[]int{math.MaxInt32, math.MinInt32},
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("striped uint 8-bits overflow", testOk(
		signal.WriteStripedUint(
			[][]uint{
				{math.MaxUint32},
				{0},
			},
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Uint64(signal.BitDepth8),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]uint64{
				{math.MaxUint8},
				{0},
			},
		},
	))
	t.Run("striped uint 8-bits overflow", testOk(
		signal.WriteUint(
			[]uint{math.MaxUint32, 0},
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Uint64(signal.BitDepth8),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]uint64{
				{math.MaxUint8},
				{0},
			},
		},
	))
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
