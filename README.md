# Signal

[![GoDoc](https://godoc.org/pipelined.dev/signal?status.svg)](https://godoc.org/pipelined.dev/signal)
[![Go Report Card](https://goreportcard.com/badge/pipelined.dev/signal)](https://goreportcard.com/report/pipelined.dev/signal)
[![Build Status](https://travis-ci.org/pipelined/signal.svg?branch=master)](https://travis-ci.org/pipelined/signal)
[![codecov](https://codecov.io/gh/pipelined/signal/branch/master/graph/badge.svg)](https://codecov.io/gh/pipelined/signal)

Manipulate digital signal with ease. Check godoc for more information.

// TODO: consider proper datatypes for signal buffers - uint/int/float64
// TODO: generic interfaces for signal buffers:

```go
type Fixed interface {
    Buffer
    BitDepth() BitDepth
    MaxBitDepth() BitDepth
    Sample(channel, pos int) int64
    SetSample(channel, pos int, value int64)
}

type Floating interface {
    Buffer
    Sample(channel, pos int) float64
    SetSample(channel, pos int, value float64)
}

type Buffer interface {
    Length() int
    Size() int
    NumChannels() int
}

// in fixed pckg
func AsFloating(src signal.Fixed, dst signal.Floating) {}

// in floating pckg
func AsFixed(src signal.Floating, dst signal.Fixed) {}

type Fixed struct {
    buffer [][]int64
    length int
}

type FixedInterleaved struct {
    BitDepth
    buffer   []int64
    channels int
    length   int
}
```