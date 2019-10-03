package signal_test

import (
	"math"
	"testing"
	"time"

	"github.com/pipelined/signal"
	"github.com/stretchr/testify/assert"
)

func TestInterIntsAsFloat64(t *testing.T) {
	tests := []struct {
		name        string
		ints        []int
		numChannels int
		bitDepth    signal.BitDepth
		unsigned    bool
		expected    [][]float64
	}{
		{
			name:        "Same length",
			ints:        []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
			numChannels: 2,
			expected: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 2},
			},
		},
		{
			name:        "Different length",
			ints:        []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			numChannels: 2,
			expected: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 0},
			},
		},
		{
			name:        "8 bits",
			ints:        []int{math.MaxInt8, math.MaxInt8 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth8,
		},
		{
			name:        "16 bits",
			ints:        []int{math.MaxInt16, math.MaxInt16 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth16,
		},
		{
			name:        "32 bits",
			ints:        []int{math.MaxInt32, math.MaxInt32 * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth32,
		},
		{
			name:        "24 bits",
			ints:        []int{1<<23 - 1, (1<<23 - 1) * 2},
			numChannels: 2,
			expected: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth24,
		},
		{
			name:     "Nil",
			ints:     nil,
			expected: nil,
		},
		{
			name:     "0 channels",
			ints:     []int{1, 2, 3},
			expected: nil,
		},
		{
			name:        "Padding",
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
		{
			name:        "Unsigned",
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

	for i, test := range tests {
		ints := signal.InterInt{
			Data:        test.ints,
			NumChannels: test.numChannels,
			BitDepth:    test.bitDepth,
			Unsigned:    test.unsigned,
		}
		result := ints.AsFloat64()
		assert.Equal(t, len(test.expected), len(result), "Test %v Bit depth %v", i, test.bitDepth)
		for i := range test.expected {
			for j, val := range test.expected[i] {
				assert.Equal(t, val, result[i][j], "Test %v Bit depth %v", i, test.bitDepth)
			}
		}
	}
}

func TestFloat64AsInterInt(t *testing.T) {
	tests := []struct {
		name     string
		floats   [][]float64
		bitDepth signal.BitDepth
		expected []int
		unsigned bool
	}{
		{
			name: "Same length",
			floats: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2, 2, 2},
			},
			expected: []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2},
		},
		{
			name: "Diffirent length",
			floats: [][]float64{
				{1, 1, 1, 1, 1, 1, 1, 1},
				{2, 2, 2, 2, 2, 2},
			},
			expected: []int{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 0, 1, 0},
		},
		{
			name: "8 bits",
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth8,
			expected: []int{1 * math.MaxInt8, 2 * math.MaxInt8},
		},
		{
			name: "16 bits",
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth16,
			expected: []int{1 * math.MaxInt16, 2 * math.MaxInt16},
		},
		{
			name: "32 bits",
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth32,
			expected: []int{1 * math.MaxInt32, 2 * math.MaxInt32},
		},
		{
			name: "24 bits",
			floats: [][]float64{
				{1},
				{2},
			},
			bitDepth: signal.BitDepth24,
			expected: []int{1 * (1<<23 - 1), 2 * (1<<23 - 1)},
		},
		{
			name:     "Nil",
			floats:   nil,
			expected: nil,
		},
		{
			name:     "0 channels",
			floats:   [][]float64{},
			expected: nil,
		},
		{
			name: "Empty channels",
			floats: [][]float64{
				{},
				{},
			},
			expected: []int{},
		},
		{
			name: "5 channels",
			floats: [][]float64{
				{1},
				{2},
				{3},
				{4},
				{5},
			},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name: "Unsigned",
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

	for _, test := range tests {
		floats := signal.Float64(test.floats)
		ints := floats.AsInterInt(test.bitDepth, test.unsigned)
		assert.Equal(t, len(test.expected), len(ints), "Test: %v Bit depth: %v", test.name, test.bitDepth)
		for i := range test.expected {
			assert.Equal(t, test.expected[i], ints[i], "Test: %v Bit depth: %v", test.name, test.bitDepth)
		}
	}
}

func TestSliceFloat64(t *testing.T) {
	var sliceTests = []struct {
		in       signal.Float64
		start    int
		len      int
		expected signal.Float64
	}{
		{
			in:       signal.Float64([][]float64{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    1,
			len:      2,
			expected: signal.Float64([][]float64{{1, 2}, {1, 2}}),
		},
		{
			in:       signal.Float64([][]float64{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    5,
			len:      2,
			expected: signal.Float64([][]float64{{5, 6}, {5, 6}}),
		},
		{
			in:       signal.Float64([][]float64{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    7,
			len:      4,
			expected: signal.Float64([][]float64{{7, 8, 9}, {7, 8, 9}}),
		},
		{
			in:       signal.Float64([][]float64{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    9,
			len:      1,
			expected: signal.Float64([][]float64{{9}}),
		},
		{
			in:       signal.Float64([][]float64{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}}),
			start:    10,
			len:      1,
			expected: nil,
		},
	}

	for _, test := range sliceTests {
		result := test.in.Slice(test.start, test.len)
		assert.Equal(t, test.expected.Size(), result.Size())
		assert.Equal(t, test.expected.NumChannels(), result.NumChannels())
		for i := range test.expected {
			for j := 0; j < len(test.expected[i]); j++ {
				assert.Equal(t, test.expected[i][j], result[i][j])
			}
		}
	}
}

func TestFloat64(t *testing.T) {
	var s signal.Float64
	assert.Equal(t, 0, s.NumChannels())
	assert.Equal(t, 0, s.Size())

	s2 := [][]float64{make([]float64, 1024)}
	s = s.Append(s2)
	assert.Equal(t, 1024, s.Size())
	s2[0] = make([]float64, 1024)
	s = s.Append(s2)
	assert.Equal(t, 2048, s.Size())
	s = s.Append(signal.Float64Buffer(1, 2048))
	assert.Equal(t, 4096, s.Size())
}

func TestDurationOf(t *testing.T) {
	var cases = []struct {
		sampleRate signal.SampleRate
		samples    int
		expected   time.Duration
	}{
		{
			sampleRate: 44100,
			samples:    44100,
			expected:   1 * time.Second,
		},
		{
			sampleRate: 44100,
			samples:    22050,
			expected:   500 * time.Millisecond,
		},
		{
			sampleRate: 44100,
			samples:    50,
			expected:   1133787 * time.Nanosecond,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.sampleRate.DurationOf(c.samples))
	}
}

func TestSamplesIn(t *testing.T) {
	var cases = []struct {
		sampleRate signal.SampleRate
		duration   time.Duration
		expected   int
	}{
		{
			sampleRate: 44100,
			duration:   1 * time.Second,
			expected:   44100,
		},
		{
			sampleRate: 44100,
			duration:   500 * time.Millisecond,
			expected:   22050,
		},
		{
			sampleRate: 44100,
			duration:   1133787 * time.Nanosecond,
			expected:   50,
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.expected, c.sampleRate.SamplesIn(c.duration))
	}
}

func TestMock(t *testing.T) {
	tests := []struct {
		numChannels int
		size        int
		value       float64
		expected    [][]float64
	}{
		{
			numChannels: 1,
			size:        2,
			expected:    [][]float64{{0, 0}},
		},
	}

	for _, test := range tests {
		result := signal.Float64Buffer(test.numChannels, test.size)

		assert.Equal(t, test.numChannels, result.NumChannels())
		assert.Equal(t, test.size, result.Size())
		for i := 0; i < len(result); i++ {
			assert.Equal(t, test.size, len(result[i]))
			for j := 0; j < len(result[i]); j++ {
				assert.Equal(t, test.expected[i][j], result[i][j])
			}
		}
	}
}
