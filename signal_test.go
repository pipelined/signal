package signal_test

import (
	"math"
	"testing"
	"time"

	"pipelined.dev/signal"
)

func TestInterIntsAsFloat64(t *testing.T) {
	tests := map[string]struct {
		ints        []int
		numChannels int
		bitDepth    signal.BitDepth
		unsigned    bool
		expected    [][]float64
	}{
		"Same length": {
			ints:        []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			numChannels: 2,
			expected: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 2},
			},
		},
		"Different length": {
			ints:        []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			numChannels: 2,
			expected: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 0},
			},
		},
		"8 bits": {
			ints:        []int{math.MaxInt8, math.MaxInt8 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth8,
		},
		"16 bits": {
			ints:        []int{math.MaxInt16, math.MaxInt16 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth16,
		},
		"32 bits": {
			ints:        []int{math.MaxInt32, math.MaxInt32 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth32,
		},
		"24 bits": {
			ints:        []int{1<<23 - 1, (1<<23 - 1) * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth24,
		},
		"Nil": {
			ints:     nil,
			expected: nil,
		},
		"0 channels": {
			ints:     []int{1, 2, 3},
			expected: nil,
		},
		"Padding": {
			ints:        []int{1, 2, 3, 4},
			numChannels: 5,
			expected: [][]float64{
				{1},
				{2},
				{3},
				{4},
				{0},
			},
		},
		"Unsigned": {
			ints:        []int{0, math.MaxInt16, math.MaxInt16 * 2},
			numChannels: 3,
			bitDepth:    signal.BitDepth16,
			unsigned:    true,
			expected: [][]float64{
				{-1},
				{0},
				{1},
			},
		},
	}

	for name, test := range tests {
		ints := signal.InterInt{
			Data:        test.ints,
			NumChannels: test.numChannels,
			BitDepth:    test.bitDepth,
			Unsigned:    test.unsigned,
		}
		assertFloat64Buffers(t, name, ints.AsFloat64(), test.expected)
	}
}

func TestInterIntCopyToFloat64(t *testing.T) {
	testPositive := func(ints signal.InterInt, floats, expected signal.Float64) func(*testing.T) {
		return func(t *testing.T) {
			ints.CopyToFloat64(floats)
			assertFloat64Buffers(t, t.Name(), floats, expected)
		}
	}
	testPanic := func(ints signal.InterInt, floats signal.Float64) func(*testing.T) {
		return func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("Didn't panic")
				}
			}()
			ints.CopyToFloat64(floats)
		}
	}

	t.Run("empty ints", testPositive(
		signal.InterInt{NumChannels: 1},
		[][]float64{{0}},
		[][]float64{{0}},
	))
	t.Run("two channels", testPositive(
		signal.InterInt{
			Data:        []int{1, 2, 3, 4},
			NumChannels: 2,
		},
		[][]float64{{0, 0}, {0, 0}},
		[][]float64{{1, 3}, {2, 4}},
	))
	t.Run("two channels padded", testPositive(
		signal.InterInt{
			Data:        []int{1, 2, 3},
			NumChannels: 2,
		},
		[][]float64{{0, 0}, {0, 0}},
		[][]float64{{1, 3}, {2, 0}},
	))

	t.Run("float not enough channels", testPanic(
		signal.InterInt{
			Data:        []int{1, 2, 3, 4},
			NumChannels: 2,
		},
		[][]float64{{0, 0}},
	))
	t.Run("float not enough samples in channel", testPanic(
		signal.InterInt{
			Data:        []int{1, 2, 3, 4},
			NumChannels: 1,
		},
		[][]float64{{0, 0}, {0, 0}},
	))
}

func TestFloat64AsInterInt(t *testing.T) {
	tests := map[string]struct {
		floats   signal.Float64
		bitDepth signal.BitDepth
		expected []int
		unsigned bool
	}{
		"Same length": {
			floats: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 2},
			},
			expected: []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
		},
		"Diffirent length": {
			floats: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2},
			},
			expected: []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 0, 1, 0},
		},
		"8 bits": {
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth8,
			expected: []int{1 * math.MaxInt8, 2 * math.MaxInt8},
		},
		"16 bits": {
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth16,
			expected: []int{1 * math.MaxInt16, 2 * math.MaxInt16},
		},
		"32 bits": {
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth32,
			expected: []int{1 * math.MaxInt32, 2 * math.MaxInt32},
		},
		"24 bits": {
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth24,
			expected: []int{1 * (1<<23 - 1), 2 * (1<<23 - 1)},
		},
		"Nil": {
			floats:   nil,
			expected: nil,
		},
		"0 channels": {
			floats:   [][]float64{},
			expected: nil,
		},
		"Empty channels": {
			floats: [][]float64{
				{},
				{},
			},
			expected: []int{},
		},
		"5 channels": {
			floats: [][]float64{
				{1},
				{2},
				{3},
				{4},
				{5},
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		"Unsigned": {
			floats: [][]float64{
				{-1},
				{0},
				{1},
			},
			bitDepth: signal.BitDepth8,
			unsigned: true,
			expected: []int{0, math.MaxInt8, math.MaxInt8 * 2},
		},
	}

	for name, test := range tests {
		ints := test.floats.AsInterInt(test.bitDepth, test.unsigned)
		assertIntBuffers(t, name, ints.Data, test.expected)
	}
}

func TestFloat64CopyToInterInt(t *testing.T) {
	testPositive := func(floats signal.Float64, ints signal.InterInt, expected []int) func(*testing.T) {
		return func(t *testing.T) {
			floats.CopyToInterInt(ints)
			assertIntBuffers(t, t.Name(), ints.Data, expected)
		}
	}
	testPanic := func(floats signal.Float64, ints signal.InterInt) func(*testing.T) {
		return func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("Didn't panic")
				}
			}()
			floats.CopyToInterInt(ints)
		}
	}
	t.Run("empty floats", testPositive(
		[][]float64{},
		signal.InterInt{
			Data: []int{0},
		},
		[]int{0},
	))
	t.Run("truncate floats", testPositive(
		[][]float64{
			{1, 1}},
		signal.InterInt{
			NumChannels: 1,
			Data:        []int{0},
		},
		[]int{1},
	))
	t.Run("pad floats", testPositive(
		[][]float64{
			{1, 1}},
		signal.InterInt{
			NumChannels: 1,
			Data:        []int{0, 0, 0},
		},
		[]int{1, 1, 0},
	))
	t.Run("two channels", testPositive(
		[][]float64{
			{1, 1},
			{2, 2}},
		signal.InterInt{
			NumChannels: 2,
			Data:        []int{0, 0, 0, 0},
		},
		[]int{1, 2, 1, 2},
	))
	t.Run("ints nil channels match", testPositive(
		[][]float64{{}},
		signal.InterInt{
			NumChannels: 1,
		},
		[]int{},
	))
	t.Run("ints nil channels not match", testPanic(
		[][]float64{},
		signal.InterInt{
			NumChannels: 1,
		},
	))
	t.Run("ints nil floats empty", testPanic(
		[][]float64{{}},
		signal.InterInt{},
	))
	t.Run("ints too short", testPanic(
		[][]float64{{1, 1}, {2, 2}},
		signal.InterInt{
			NumChannels: 2,
			Data:        []int{0, 0, 0},
		},
	))
}

func TestSliceFloat64(t *testing.T) {
	var sliceTests = map[string]struct {
		in       signal.Float64
		start    int
		len      int
		expected signal.Float64
	}{
		"slice start": {
			in: signal.Float64([][]float64{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start: 0,
			len:   2,
			expected: signal.Float64([][]float64{
				{0, 1},
				{0, 1}}),
		},
		"slice middle": {
			in: signal.Float64([][]float64{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start: 5,
			len:   2,
			expected: signal.Float64([][]float64{
				{5, 6},
				{5, 6}}),
		},
		"slice end padded": {
			in: signal.Float64([][]float64{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start: 7,
			len:   4,
			expected: signal.Float64([][]float64{
				{7, 8, 9},
				{7, 8, 9}}),
		},
		"slice last": {
			in: signal.Float64([][]float64{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start: 9,
			len:   1,
			expected: signal.Float64([][]float64{
				{9}}),
		},
		"slice after": {
			in: signal.Float64([][]float64{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    10,
			len:      1,
			expected: nil,
		},
	}

	for name, test := range sliceTests {
		assertFloat64Buffers(t, name, test.in.Slice(test.start, test.len), test.expected)
	}
}

func TestFloat64Append(t *testing.T) {
	tests := map[string]struct {
		buffer   signal.Float64
		addition signal.Float64
		expected signal.Float64
	}{
		"slice of slices": {
			addition: [][]float64{
				make([]float64, 1024),
			},
			expected: [][]float64{
				make([]float64, 1024),
			},
		},
		"single chan buffer": {
			buffer: [][]float64{
				make([]float64, 0),
			},
			addition: signal.Float64Buffer(1, 1024),
			expected: signal.Float64Buffer(1, 1024),
		},
		"multiple chan buffer": {
			buffer: [][]float64{
				make([]float64, 0),
			},
			addition: signal.Float64Buffer(2, 1024),
			expected: signal.Float64Buffer(1, 1024),
		},
	}

	for name, test := range tests {
		assertFloat64Buffers(t, name, test.buffer.Append(test.addition), test.expected)
	}
}

func TestDurationOf(t *testing.T) {
	var cases = map[string]struct {
		sampleRate signal.SampleRate
		samples    int
		expected   time.Duration
	}{
		"second: ": {
			sampleRate: 44100,
			samples:    44100,
			expected:   1 * time.Second,
		},
		"millis": {
			sampleRate: 44100,
			samples:    22050,
			expected:   500 * time.Millisecond,
		},
		"nanos": {
			sampleRate: 44100,
			samples:    50,
			expected:   1133787 * time.Nanosecond,
		},
	}
	for name, test := range cases {
		result := test.sampleRate.DurationOf(test.samples)
		if test.expected != result {
			t.Fatalf("%v invalid duration: %v expected: %v", name, result, test.expected)
		}
	}
}

func TestSamplesIn(t *testing.T) {
	var cases = map[string]struct {
		sampleRate signal.SampleRate
		duration   time.Duration
		expected   int
	}{
		"second": {
			sampleRate: 44100,
			duration:   1 * time.Second,
			expected:   44100,
		},
		"millis": {
			sampleRate: 44100,
			duration:   500 * time.Millisecond,
			expected:   22050,
		},
		"nanos": {
			sampleRate: 44100,
			duration:   1133787 * time.Nanosecond,
			expected:   50,
		},
	}
	for name, test := range cases {
		result := test.sampleRate.SamplesIn(test.duration)
		if test.expected != result {
			t.Fatalf("%v invalid samples: %v expected: %v", name, result, test.expected)
		}
	}
}

func TestFloat64Buffer(t *testing.T) {
	tests := map[string]struct {
		numChannels int
		size        int
		value       float64
		expected    signal.Float64
	}{
		"one channel": {
			numChannels: 1,
			size:        2,
			expected: [][]float64{
				{0, 0},
			},
		},
		"two channels": {
			numChannels: 2,
			size:        2,
			expected: [][]float64{
				{0, 0},
				{0, 0},
			},
		},
		"zero channels": {
			numChannels: 0,
			size:        2,
			expected:    [][]float64{},
		},
	}

	for name, test := range tests {
		assertFloat64Buffers(t, name, signal.Float64Buffer(test.numChannels, test.size), test.expected)
	}
}

func TestFloat64Sum(t *testing.T) {
	tests := map[string]struct {
		buffer   signal.Float64
		addition signal.Float64
		expected signal.Float64
	}{
		"add nil": {
			buffer: [][]float64{
				{1, 1}},
			expected: [][]float64{
				{1, 1}},
		},
		"add to nil": {
			buffer: nil,
			addition: [][]float64{
				{1, 1}},
			expected: nil,
		},
		"add same": {
			buffer: [][]float64{
				{1},
				{1}},
			addition: [][]float64{
				{2, 2}},
			expected: [][]float64{
				{3},
				{1}},
		},
		"add smaller": {
			buffer: [][]float64{
				{1, 1},
				{1, 1}},
			addition: [][]float64{
				{2}},
			expected: [][]float64{
				{3, 1},
				{1, 1}},
		},
		"add to smaller": {
			buffer: [][]float64{
				{2}},
			addition: [][]float64{
				{1, 1},
				{1, 1}},
			expected: [][]float64{
				{3}},
		},
	}

	for name, test := range tests {
		assertFloat64Buffers(t, name, test.buffer.Sum(test.addition), test.expected)
	}
}

func assertIntBuffers(t *testing.T, name string, result, expected []int) {
	if len(expected) != len(result) {
		t.Fatalf("%v: invalid buffer size: %v expected: %v", name, len(expected), len(result))
	}
	for i := range expected {
		if expected[i] != result[i] {
			t.Fatalf("%v: invalid value: %v expected: %v", name, result[i], expected[i])
		}
	}
}

func assertFloat64Buffers(t *testing.T, name string, result, expected signal.Float64) {
	if expected.NumChannels() != result.NumChannels() {
		t.Fatalf("%v: invalid num channels: %v expeced %v", name, result.NumChannels(), expected.NumChannels())
	}
	if expected.Size() != result.Size() {
		t.Fatalf("%v: invalid buffer size: %v expeced %v", name, result.Size(), expected.Size())
	}

	for i := range expected {
		if len(expected[i]) != len(result[i]) {
			t.Fatalf("%v: invalid channel %d size: %v expeced %v", name, i, len(result[i]), len(expected[i]))
		}
		for j := range expected[i] {
			if expected[i][j] != result[i][j] {
				t.Fatalf("%v: invalid value: %v expeced %v", name, result[i][j], expected[i][j])
			}
		}
	}
}
