# Signal

[![GoDoc](https://godoc.org/pipelined.dev/signal?status.svg)](https://godoc.org/pipelined.dev/signal)
[![Go Report Card](https://goreportcard.com/badge/pipelined.dev/signal)](https://goreportcard.com/report/pipelined.dev/signal)
[![Build Status](https://travis-ci.org/pipelined/signal.svg?branch=master)](https://travis-ci.org/pipelined/signal)
[![codecov](https://codecov.io/gh/pipelined/signal/branch/master/graph/badge.svg)](https://codecov.io/gh/pipelined/signal)

This package provides functionality to manipulate digital signals and its
attributes.

It contains structures for various signal types and allows
conversions from one to another:

* Fixed-point signed
* Fixed-point unsigned
* Floating-point

Signal types have semantics of golang slices - they can be appended or
sliced with respect to channels layout.

On top of that, this package was desinged to simplify control on
allocations. Check godoc for examples.
