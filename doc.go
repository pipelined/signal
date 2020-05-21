/*
Package signal provides primitives for digital signal manipulations.
Floating-point, fixed signed and unsigned signal types are supported. This
package focuses on allocation optimisations as it's one of the most
important aspects for DSP applications.

Basics

Digital signal is a representation of a physical signal that is a sampled
and quantized. It is discrete in time and amplitude. When analog signal is
converted to digital, it goes through two steps: discretization and
quantization. Discretization means that the signal is divided into equal
intervals of time, and each interval is represented by a single measurement
of amplitude. Quantization means each amplitude measurement is approximated
by a value from a finite set.

The finite set of quantization values determines the digital signal
representation. It differs from type to type:

    floating:   [-1, 1]
    signed      [-2^(bitDepth-1), 2^(bitDepth-1)-1]
    unsigned    [0, 2^bitDepth-1]

Floating-point signal can exceed this range without loosing signal data,
but fixed-point signals will be clipped and meaningful signal values will
be lost.

Signal buffer types

In order to allocate any of signal buffers, an allocator should be used.
Allocator defines what number of channels and capacity per channel
allocated buffers will have:

    alloc := signal.Allocator{Channels: 2, Capacity: 512}

This package offers types that represent floating-point and both
signed/unsigned fixed-point signal buffers. They implement Floating, Signed
and Unsigned interfaces respectively. Internally, signal buffers use a
slice of built-in type to hold the data.

Fixed-point buffers require a bit depth to be provided at allocation time.
It allows to use the same type to hold values of various bit depths:

    alloc := signal.Allocator{Channels: 2, Capacity: 512}
    _ = alloc.Int64(signal.BitDepth32)      // int-64 buffer with 32-bits per sample
    _ = alloc.Uint32(signal.BitDepth24)     // uint-32 buffer with 24-bits per sample
    _ = alloc.Float64()                     // float-64 buffer

Signal buffers have semantics of golang-slices - they can be sliced or
apended one to another. All operations respect number of channels within
buffer, so slicing and appending always happens for all channels.

Write/Read values

There are multiple ways to write/read values from the signal buffers.

WriteT/ReadT functions allows to write or read the data in the format of
single slice, where samples for different channels are interleaved:

    []T{R, L, R, L}

WriteStripedT/ReadStripedT functions, on the other hand, can use slice of
slices, where each nested slice represents a single channel of signal:

    [][]T{{R, R}, {L, L}}

It's possible to append samples to the buffers using AppendSample fucntion.
However, in order to have more control over allocations, this function
won't let the buffer grow beyond it's initial capacity. To achieve this,
another buffer needs to be explicitly allocated and appended to the buffer.

The one can also iterate over signal buffers. Please, refer to examples for
more details.
*/
package signal
