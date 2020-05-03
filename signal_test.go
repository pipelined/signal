package signal_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"pipelined.dev/signal"
)

var (
	_ signal.Signed = signal.Allocator{}.Int64(signal.MaxBitDepth)
	// _ signal.Unsigned = signal.Allocator{}.Uint64(signal.MaxBitDepth)
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
		assertEqual(t, "capacity", r.Capacity(), ex.capacity)
		assertEqual(t, "length", r.Length(), ex.length)
		assertEqual(t, "slices", result(r), ex.data)
	}
}

func TestSignedAsFloating(t *testing.T) {
	t.Run("8 bits", testOk(
		signal.SignedAsFloating(
			signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 1}.
				Int64(signal.BitDepth8),
				[][]int64{
					{math.MaxInt8},
					{math.MinInt8},
				}),
			signal.Allocator{
				Channels: 2,
				Capacity: 1,
			}.Float64(),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]float64{
				{1},
				{-1},
			},
		},
	))
	t.Run("16 bits", testOk(
		signal.SignedAsFloating(
			signal.WriteStripedInt64(
				signal.Allocator{
					Channels: 2,
					Capacity: 1,
				}.Int64(signal.BitDepth16),
				[][]int64{
					{math.MaxInt16},
					{math.MaxInt16 + 1},
				}),
			signal.Allocator{
				Channels: 2,
				Capacity: 1,
			}.Float64(),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]float64{
				{1},
				{1},
			},
		},
	))
	t.Run("32 bits", testOk(
		signal.SignedAsFloating(
			signal.WriteStripedInt64(
				signal.Allocator{
					Channels: 2,
					Capacity: 1,
				}.Int64(signal.BitDepth32),
				[][]int64{
					{math.MaxInt32},
					{math.MaxInt32 + 1},
				}),
			signal.Allocator{
				Channels: 2,
				Capacity: 1,
			}.Float64(),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]float64{
				{1},
				{1},
			},
		},
	))
	t.Run("24 bits", testOk(
		signal.SignedAsFloating(
			signal.WriteStripedInt64(
				signal.Allocator{
					Channels: 2,
					Capacity: 1,
				}.Int64(signal.BitDepth24),
				[][]int64{
					{1<<23 - 1},
					{(1<<23 - 1) + 1},
				}),
			signal.Allocator{
				Channels: 2,
				Capacity: 1,
			}.Float64(),
		),
		expected{
			length:   1,
			capacity: 1,
			data: [][]float64{
				{1},
				{1},
			},
		},
	))
}

func TestWrite(t *testing.T) {
	testFloatingOk := func(s signal.Floating, ex expected) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			assertEqual(t, "slices", result(s), ex.data)
			assertEqual(t, "length", s.Length(), ex.length)
		}
	}

	t.Run("int64 int64 8 bits", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
			[][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{127},
				{-128},
			},
		},
	))
	t.Run("int64 striped int full buffer", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 striped int signle channel", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int{
				{},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]int64{
				{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 striped int short buffer", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int{
				{1, 2, 3},
				{11, 12, 0},
			}),
		expected{
			length:   3,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 0},
			},
		},
	))
	t.Run("int64 striped int long buffer", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 3,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 striped int 8-bits overflow", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
			[][]int{
				{math.MaxInt32},
				{math.MinInt32},
			}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 striped int 16-bits overflow", testOk(
		signal.WriteStripedInt(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth16),
			[][]int{
				{math.MaxInt64},
				{math.MinInt64},
			}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("int64 striped int64 full buffer", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 striped int64 short buffer", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2, 3},
				{11, 12, 13},
			}),
		expected{
			length:   3,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 striped int64 long buffer", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 3,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 striped int64 8-bits overflow", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
			[][]int64{
				{math.MaxInt32},
				{math.MinInt32},
			}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 striped int64 16-bits overflow", testOk(
		signal.WriteStripedInt64(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth16),
			[][]int64{
				{math.MaxInt64},
				{math.MinInt64},
			}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("int64 int full buffer", testOk(
		signal.WriteInt(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 int short buffer", testOk(
		signal.WriteInt(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int{1, 11, 2, 12, 3, 13}),
		expected{
			length:   3,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 int long buffer", testOk(
		signal.WriteInt(
			signal.Allocator{
				Capacity: 3,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20}),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 int 8-bits overflow", testOk(
		signal.WriteInt(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
			[]int{math.MaxInt32, math.MinInt32}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 int 16-bits overflow", testOk(
		signal.WriteInt(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth16),
			[]int{math.MaxInt64, math.MinInt64}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("int64 int64 full buffer", testOk(
		signal.WriteInt64(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("int64 int64 short buffer", testOk(
		signal.WriteInt64(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int64{1, 11, 2, 12, 3, 13}),
		expected{
			length:   3,
			capacity: 10,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 int64 long buffer", testOk(
		signal.WriteInt64(
			signal.Allocator{
				Capacity: 3,
				Channels: 2,
			}.Int64(signal.MaxBitDepth),
			[]int64{1, 11, 2, 12, 3, 13, 4, 14, 5, 15, 6, 16, 7, 17, 8, 18, 9, 19, 10, 20}),
		expected{
			length:   3,
			capacity: 3,
			data: [][]int64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
	t.Run("int64 8-bits overflow", testOk(
		signal.WriteInt64(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth8),
			[]int64{math.MaxInt32, math.MinInt32}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt8},
				{math.MinInt8},
			},
		},
	))
	t.Run("int64 16-bits overflow", testOk(
		signal.WriteInt64(
			signal.Allocator{
				Capacity: 1,
				Channels: 2,
			}.Int64(signal.BitDepth16),
			[]int64{math.MaxInt64, math.MinInt64}),
		expected{
			length:   1,
			capacity: 1,
			data: [][]int64{
				{math.MaxInt16},
				{math.MinInt16},
			},
		},
	))
	t.Run("float64 striped full buffer", testFloatingOk(
		signal.WriteStripedFloat64(
			signal.Allocator{
				Capacity: 10,
				Channels: 2,
			}.Float64(),
			[][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   10,
			capacity: 10,
			data: [][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		},
	))
	t.Run("float64 multiple buffers", testFloatingOk(
		signal.WriteFloat64(
			signal.WriteFloat64(
				signal.Allocator{
					Capacity: 10,
					Channels: 2,
				}.Float64(),
				[]float64{1, 11, 2, 12, 3, 13},
			),
			[]float64{101, 111, 102, 112, 103, 113, 104, 114, 105, 115, 106, 116, 107, 117, 108, 118, 109, 119, 110, 120},
		),
		expected{
			length:   10,
			capacity: 10,
			data: [][]float64{
				{1, 2, 3, 101, 102, 103, 104, 105, 106, 107},
				{11, 12, 13, 111, 112, 113, 114, 115, 116, 117},
			},
		},
	))
	t.Run("float64 striped multiple buffers", testFloatingOk(
		signal.WriteStripedFloat64(
			signal.WriteStripedFloat64(
				signal.Allocator{
					Capacity: 10,
					Channels: 2,
				}.Float64(),
				[][]float64{
					{1, 2, 3},
					{11, 12, 13},
				}),
			[][]float64{
				{101, 102, 103, 104, 105, 106, 107, 108, 109, 110},
				{111, 112, 113, 114, 115, 116, 117, 118, 119, 120},
			},
		),
		expected{
			length:   10,
			capacity: 10,
			data: [][]float64{
				{1, 2, 3, 101, 102, 103, 104, 105, 106, 107},
				{11, 12, 13, 111, 112, 113, 114, 115, 116, 117},
			},
		},
	))
	t.Run("float64 striped long buffer", testFloatingOk(
		signal.WriteStripedFloat64(
			signal.Allocator{
				Capacity: 3,
				Channels: 2,
			}.Float64(),
			[][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			}),
		expected{
			length:   3,
			capacity: 3,
			data: [][]float64{
				{1, 2, 3},
				{11, 12, 13},
			},
		},
	))
}

func TestAppend(t *testing.T) {
	testPanic := func(appender signal.Signal, data signal.Signal) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			switch a := appender.(type) {
			case signal.Signed:
				d := data.(signal.Int64)
				assertPanic(t, func() {
					a.Append(d)
				})
			case signal.Floating:
				d := data.(signal.Float64)
				assertPanic(t, func() {
					a.Append(d)
				})
			default:
				panic(fmt.Sprintf("unsupported append panic type %T", appender))
			}
		}
	}

	t.Run("int64 single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth).
			Append(
				signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
					[][]int64{
						{1, 2},
						{11, 12},
					},
				),
			),
		expected{
			capacity: 2,
			length:   2,
			data: [][]int64{
				{1, 2},
				{11, 12},
			},
		},
	))
	t.Run("int64 multiple slices", testOk(
		signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2},
				{1, 2},
			},
		).Append(
			signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
				[][]int64{
					{3, 4},
					{3, 4},
				},
			),
		),
		expected{
			capacity: 4,
			length:   4,
			data: [][]int64{
				{1, 2, 3, 4},
				{1, 2, 3, 4},
			},
		},
	))
	t.Run("int64 different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
		signal.Allocator{Channels: 1, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
	t.Run("int64 different bit depth", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.BitDepth8),
		signal.Allocator{Channels: 2, Capacity: 2}.Int64(signal.MaxBitDepth),
	))
	t.Run("float64 single slice", testOk(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64().
			Append(
				signal.WriteStripedFloat64(
					signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
					[][]float64{
						{1, 2},
						{11, 12},
					},
				),
			),
		expected{
			length:   2,
			capacity: 2,
			data: [][]float64{
				{1, 2},
				{11, 12},
			},
		},
	))
	t.Run("float64 multiple slices", testOk(
		signal.WriteStripedFloat64(
			signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
			[][]float64{
				{1, 2},
				{11, 12},
			},
		).Append(
			signal.WriteStripedFloat64(
				signal.Allocator{Channels: 2, Capacity: 3}.Float64(),
				[][]float64{
					{3, 4},
					{13, 14},
				},
			),
		),
		expected{
			length:   4,
			capacity: 5,
			data: [][]float64{
				{1, 2, 3, 4},
				{11, 12, 13, 14},
			},
		},
	))
	t.Run("float64 different channels", testPanic(
		signal.Allocator{Channels: 2, Capacity: 2}.Float64(),
		signal.Allocator{Channels: 1, Capacity: 2}.Float64(),
	))
}

func TestSlice(t *testing.T) {
	t.Run("int64", testOk(
		signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 10}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		).Slice(3, 6),
		expected{
			length:   3,
			capacity: 7,
			data: [][]int64{
				{4, 5, 6},
				{14, 15, 16},
			},
		},
	))
	t.Run("int64", testOk(
		signal.WriteStripedInt64(signal.Allocator{Channels: 2, Capacity: 10}.Int64(signal.MaxBitDepth),
			[][]int64{
				{1, 2, 3, 4, 0},
				{11, 12, 13, 14, 0},
			},
		).Slice(3, 6),
		expected{
			length:   3,
			capacity: 7,
			data: [][]int64{
				{4, 0, 0},
				{14, 0, 0},
			},
		},
	))
	t.Run("float64", testOk(
		signal.WriteStripedFloat64(signal.Allocator{Channels: 2, Capacity: 10}.Float64(),
			[][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			},
		).Slice(3, 6),
		expected{
			length:   3,
			capacity: 7,
			data: [][]float64{
				{4, 5, 6},
				{14, 15, 16},
			},
		},
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
	switch s := sig.(type) {
	case signal.Signed:
		result := make([][]int64, s.Channels())
		for i := range result {
			result[i] = make([]int64, s.Length())
		}
		signal.ReadStripedInt64(s, result)
		return result
	case signal.Unsigned:
		return nil
	case signal.Floating:
		result := make([][]float64, s.Channels())
		for i := range result {
			result[i] = make([]float64, s.Length())
		}
		signal.ReadStripedFloat64(s, result)
		return result
	default:
		panic(fmt.Sprintf("unsupported result type: %T", sig))
	}
}
