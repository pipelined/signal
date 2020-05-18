// +build ignore

// This program generates all signal types.
package main

import (
	"fmt"
	"os"
	"text/template"
	"time"
)

type typeGenerator struct {
	Timestamp   time.Time
	Builtin     string
	Name        string
	MaxBitDepth string // used for fixed-point types only
}

func main() {
	types := map[typeGenerator]templates{
		{
			Builtin:     "int8",
			Name:        "Int8",
			MaxBitDepth: "BitDepth8",
		}: signedTemplates,
		{
			Builtin:     "int16",
			Name:        "Int16",
			MaxBitDepth: "BitDepth16",
		}: signedTemplates,
		{
			Builtin:     "int32",
			Name:        "Int32",
			MaxBitDepth: "BitDepth32",
		}: signedTemplates,
		{
			Builtin:     "int64",
			Name:        "Int64",
			MaxBitDepth: "BitDepth64",
		}: signedTemplates,
		{
			Builtin:     "uint8",
			Name:        "Uint8",
			MaxBitDepth: "BitDepth8",
		}: unsignedTemplates,
		{
			Builtin:     "uint16",
			Name:        "Uint16",
			MaxBitDepth: "BitDepth16",
		}: unsignedTemplates,
		{
			Builtin:     "uint32",
			Name:        "Uint32",
			MaxBitDepth: "BitDepth32",
		}: unsignedTemplates,
		{
			Builtin:     "uint64",
			Name:        "Uint64",
			MaxBitDepth: "BitDepth64",
		}: unsignedTemplates,
		{
			Builtin: "float32",
			Name:    "Float32",
		}: floatingTemplates,
		{
			Builtin: "float64",
			Name:    "Float64",
		}: floatingTemplates,
	}

	for gen, template := range types {
		generate(gen, template)
	}
}

func generate(gen typeGenerator, t templates) {
	gen.Timestamp = time.Now()

	generateFile(fmt.Sprintf("%s.go", gen.Builtin), gen, t.types)
	generateFile(fmt.Sprintf("%s_test.go", gen.Builtin), gen, t.tests)

	// err = t.tests.Execute(f, gen)
	// die(fmt.Sprintf("execute %s tests template for %s type", t.tests.Name(), gen.Name), err)
}

func generateFile(fileName string, gen typeGenerator, t *template.Template) {
	if t == nil {
		return
	}
	f, err := os.Create(fileName)
	die(fmt.Sprintf("create %s file", fileName), err)
	defer f.Close()

	err = t.Execute(f, gen)
	die(fmt.Sprintf("execute %s template for %s type", t.Name(), gen.Name), err)
}

func die(reason string, err error) {
	if err != nil {
		panic(fmt.Sprintf("failed %s: %v", reason, err))
	}
}

type templates struct {
	types *template.Template
	tests *template.Template
}

var (
	floatingTemplates = templates{
		types: template.Must(template.New("floating").Parse(floating)),
		tests: template.Must(template.New("floating tests").Parse(floatingTests)),
	}
	signedTemplates = templates{
		types: template.Must(template.New("signed").Parse(signed)),
		tests: template.Must(template.New("signed tests").Parse(fixedTests)),
	}
	unsignedTemplates = templates{
		types: template.Must(template.New("unsigned").Parse(unsigned)),
		tests: template.Must(template.New("unsigned tests").Parse(fixedTests)),
	}
)

const (
	floating = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package signal

// {{ .Name }} is a sequential {{ .Builtin }} floating-point signal.
type {{ .Name }} struct {
	buffer []{{ .Builtin }}
	channels
}

// {{ .Name }} allocates new sequential {{ .Builtin }} signal buffer.
func (a Allocator) {{ .Name }}() Floating {
	return {{ .Name }}{
		buffer:   make([]{{ .Builtin }}, 0, a.Channels*a.Capacity),
		channels: channels(a.Channels),
	}
}

// Capacity returns capacity of a single channel.
func (s {{ .Name }}) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s {{ .Name }}) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s {{ .Name }}) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s {{ .Name }}) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s {{ .Name }}) AppendSample(value float64) Floating {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, {{ .Builtin }}(value))
	return s
}

// Sample returns signal value for provided channel and position.
func (s {{ .Name }}) Sample(pos int) float64 {
	return float64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s {{ .Name }}) SetSample(pos int, value float64) {
	s.buffer[pos] = {{ .Builtin }}(value)
}

// Slice slices buffer with respect to channels.
func (s {{ .Name }}) Slice(start, end int) Floating {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s {{ .Name }}) Reset() Floating {
	return s.Slice(0, 0)
}

// Append appends data from src buffer to the end of the buffer.
func (s {{ .Name }}) Append(src Floating) Floating {
	mustSameChannels(s.Channels(), src.Channels())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with cap of both sources capacity;
		s.buffer = append(make([]{{ .Builtin }}, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Floating(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// Read{{ .Name }} reads values from the buffer into provided slice.
func Read{{ .Name }}(src Floating, dst []{{ .Builtin }}) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = {{ .Builtin }}(src.Sample(pos))
	}
}

// ReadStriped{{ .Name }} reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStriped{{ .Name }}(src Floating, dst [][]{{ .Builtin }}) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = {{ .Builtin }}(src.Sample(src.ChannelPos(channel, pos)))
		}
	}
}

// Write{{ .Name }} writes values from provided slice into the buffer.
func Write{{ .Name }}(src []{{ .Builtin }}, dst Floating) Floating {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(float64(src[pos]))
	}
	return dst
}

// WriteStriped{{ .Name }} appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended.
func WriteStriped{{ .Name }}(src [][]{{ .Builtin }}, dst Floating) Floating {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(float64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}`

	signed = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package signal

// {{ .Name }} is {{ .Builtin }} signed fixed signal.
type {{ .Name }} struct {
	buffer []{{ .Builtin }}
	channels
	bitDepth
}

// {{ .Name }} allocates new sequential {{ .Builtin }} signal buffer.
func (a Allocator) {{ .Name }}(bd BitDepth) Signed {
	return {{ .Name }}{
		buffer:   make([]{{ .Builtin }}, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: defaultBitDepth(bd, {{ .MaxBitDepth }}),
	}
}

// Capacity returns capacity of a single channel.
func (s {{ .Name }}) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s {{ .Name }}) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s {{ .Name }}) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s {{ .Name }}) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s {{ .Name }}) AppendSample(value int64) Signed {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, {{ .Builtin }}(s.BitDepth().SignedValue(value)))
	return s
}

// Sample returns signal value for provided channel and position.
func (s {{ .Name }}) Sample(pos int) int64 {
	return int64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s {{ .Name }}) SetSample(pos int, value int64) {
	s.buffer[pos] = {{ .Builtin }}(s.BitDepth().SignedValue(value))
}

// Slice slices buffer with respect to channels.
func (s {{ .Name }}) Slice(start, end int) Signed {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s {{ .Name }}) Reset() Signed {
	return s.Slice(0, 0)
}

// Append appends [0:Length] data from src to current buffer and returns new
// Signed buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic. If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s {{ .Name }}) Append(src Signed) Signed {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with sources cap
		s.buffer = append(make([]{{ .Builtin }}, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Signed(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// Read{{ .Name }} reads values from the buffer into provided slice.
func Read{{ .Name }}(src Signed, dst []{{ .Builtin }}) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = {{ .Builtin }}({{ .MaxBitDepth }}.SignedValue(src.Sample(pos)))
	}
}

// ReadStriped{{ .Name }} reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStriped{{ .Name }}(src Signed, dst [][]{{ .Builtin }}) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = {{ .Builtin }}({{ .MaxBitDepth }}.SignedValue(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
}

// Write{{ .Name }} writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func Write{{ .Name }}(src []{{ .Builtin }}, dst Signed) Signed {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(int64(src[pos]))
	}
	return dst
}

// WriteStriped{{ .Name }} appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended. Sample values are capped by maximum value
// of the buffer bit depth.
func WriteStriped{{ .Name }}(src [][]{{ .Builtin }}, dst Signed) Signed {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(int64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}`

	unsigned = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package signal

// {{ .Name }} is {{ .Builtin }} signed fixed signal.
type {{ .Name }} struct {
	buffer []{{ .Builtin }}
	channels
	bitDepth
}

// {{ .Name }} allocates new sequential {{ .Builtin }} signal buffer.
func (a Allocator) {{ .Name }}(bd BitDepth) Unsigned {
	return {{ .Name }}{
		buffer:   make([]{{ .Builtin }}, 0, a.Capacity*a.Channels),
		channels: channels(a.Channels),
		bitDepth: defaultBitDepth(bd, BitDepth64),
	}
}

// Capacity returns capacity of a single channel.
func (s {{ .Name }}) Capacity() int {
	return cap(s.buffer) / int(s.channels)
}

// Length returns length of a single channel.
func (s {{ .Name }}) Length() int {
	return len(s.buffer) / int(s.channels)
}

// Cap returns capacity of whole buffer.
func (s {{ .Name }}) Cap() int {
	return cap(s.buffer)
}

// Len returns length of whole buffer.
func (s {{ .Name }}) Len() int {
	return len(s.buffer)
}

// AppendSample appends sample at the end of the buffer.
// Sample is not appended if buffer capacity is reached.
func (s {{ .Name }}) AppendSample(value uint64) Unsigned {
	if len(s.buffer) == cap(s.buffer) {
		return s
	}
	s.buffer = append(s.buffer, {{ .Builtin }}(s.BitDepth().UnsignedValue(value)))
	return s
}

// Sample returns signal value for provided channel and position.
func (s {{ .Name }}) Sample(pos int) uint64 {
	return uint64(s.buffer[pos])
}

// SetSample sets sample value for provided position.
func (s {{ .Name }}) SetSample(pos int, value uint64) {
	s.buffer[pos] = {{ .Builtin }}(s.BitDepth().UnsignedValue(value))
}

// Slice slices buffer with respect to channels.
func (s {{ .Name }}) Slice(start, end int) Unsigned {
	start = s.ChannelPos(0, start)
	end = s.ChannelPos(0, end)
	s.buffer = s.buffer[start:end]
	return s
}

// Reset sets length of the buffer to zero.
func (s {{ .Name }}) Reset() Unsigned {
	return s.Slice(0, 0)
}

// Append appends data from src to current buffer and returns new
// Unsigned buffer. Both buffers must have same number of channels and bit depth,
// otherwise function will panic.  If current buffer doesn't have enough capacity,
// new buffer will be allocated with capacity of both sources.
func (s {{ .Name }}) Append(src Unsigned) Unsigned {
	mustSameChannels(s.Channels(), src.Channels())
	mustSameBitDepth(s.BitDepth(), src.BitDepth())
	if s.Cap() < s.Len()+src.Len() {
		// allocate and append buffer with sources cap
		s.buffer = append(make([]{{ .Builtin }}, 0, s.Cap()+src.Cap()), s.buffer...)
	}
	result := Unsigned(s)
	for pos := 0; pos < src.Len(); pos++ {
		result = result.AppendSample(src.Sample(pos))
	}
	return result
}

// Read{{ .Name }} reads values from the buffer into provided slice.
func Read{{ .Name }}(src Unsigned, dst []{{ .Builtin }}) {
	length := min(src.Len(), len(dst))
	for pos := 0; pos < length; pos++ {
		dst[pos] = {{ .Builtin }}({{ .MaxBitDepth }}.UnsignedValue(src.Sample(pos)))
	}
}

// ReadStriped{{ .Name }} reads values from the buffer into provided slice.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, no values for
// that channel will be appended.
func ReadStriped{{ .Name }}(src Unsigned, dst [][]{{ .Builtin }}) {
	mustSameChannels(src.Channels(), len(dst))
	for channel := 0; channel < src.Channels(); channel++ {
		for pos := 0; pos < src.Length() && pos < len(dst[channel]); pos++ {
			dst[channel][pos] = {{ .Builtin }}({{ .MaxBitDepth }}.UnsignedValue(src.Sample(src.ChannelPos(channel, pos))))
		}
	}
}

// Write{{ .Name }} writes values from provided slice into the buffer.
// If the buffer already contains any data, it will be overwritten.
// Sample values are capped by maximum value of the buffer bit depth.
func Write{{ .Name }}(src []{{ .Builtin }}, dst Unsigned) Unsigned {
	length := min(dst.Cap()-dst.Len(), len(src))
	for pos := 0; pos < length; pos++ {
		dst = dst.AppendSample(uint64(src[pos]))
	}
	return dst
}

// WriteStriped{{ .Name }} appends values from provided slice into the buffer.
// The length of provided slice must be equal to the number of channels,
// otherwise function will panic. Nested slices can be nil, zero values for
// that channel will be appended. Sample values are capped by maximum value
// of the buffer bit depth.
func WriteStriped{{ .Name }}(src [][]{{ .Builtin }}, dst Unsigned) Unsigned {
	mustSameChannels(dst.Channels(), len(src))
	var length int
	for i := range src {
		if len(src[i]) > length {
			length = len(src[i])
		}
	}
	length = min(length, dst.Capacity()-dst.Length())
	for pos := 0; pos < length; pos++ {
		for channel := 0; channel < dst.Channels(); channel++ {
			if pos < len(src[channel]) {
				dst = dst.AppendSample(uint64(src[channel][pos]))
			} else {
				dst = dst.AppendSample(0)
			}
		}
	}
	return dst
}`

	fixedTests = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func Test{{ .Name }}(t *testing.T) {
	t.Run("{{ .Builtin }}", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.{{ .Name }}(signal.{{ .MaxBitDepth }}).
			Append(signal.WriteStriped{{ .Name }}(
				[][]{{ .Builtin }}{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.{{ .Name }}(signal.{{ .MaxBitDepth }})),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]{{ .Builtin }}{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}`
	floatingTests = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
package signal_test

import (
	"testing"

	"pipelined.dev/signal"
)

func Test{{ .Name }}(t *testing.T) {
	t.Run("{{ .Builtin }}", testOk(
		signal.Allocator{
			Channels: 3,
			Capacity: 2,
		}.{{ .Name }}().
			Append(signal.WriteStriped{{ .Name }}(
				[][]{{ .Builtin }}{
					{},
					{1, 2, 3},
					{11, 12, 13, 14},
				},
				signal.Allocator{
					Channels: 3,
					Capacity: 3,
				}.{{ .Name }}()),
			).
			Slice(1, 3),
		expected{
			length:   2,
			capacity: 4,
			data: [][]{{ .Builtin }}{
				{0, 0},
				{2, 3},
				{12, 13},
			},
		},
	))
}`
)
