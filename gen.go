// +build ignore

// This program generates all signal types.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

// template names for signal types and tests
type templateNames struct {
	signal string
	tests  string
}

// data for templates execution
type generator struct {
	InterfaceProps
	Timestamp   time.Time
	Builtin     string
	Name        string
	Allocator   string
	MaxBitDepth string // used for fixed-point types only
	*ChannelTemplates
}

// properties for interface generation
type InterfaceProps struct {
	Interface   string
	SampleType  string
	ChannelType string
}

// data for channel types generation. generated for each widest type -
// float64, int64, uint64.
type ChannelTemplates struct {
	signal string
	tests  string
}

var templates = template.Must(template.ParseGlob("templates/*.tmpl"))

func main() {
	var (
		fixedTemplates = templateNames{
			signal: "fixed",
			tests:  "tests",
		}
		floatingTemplates = templateNames{
			signal: "floating",
			tests:  "tests",
		}
	)
	var (
		signedProps = InterfaceProps{
			Interface:   "Signed",
			SampleType:  "int64",
			ChannelType: "signedChannel",
		}
		unsignedProps = InterfaceProps{
			Interface:   "Unsigned",
			SampleType:  "uint64",
			ChannelType: "unsignedChannel",
		}
		floatingProps = InterfaceProps{
			Interface:   "Floating",
			SampleType:  "float64",
			ChannelType: "floatingChannel",
		}
	)
	channelTemplates := ChannelTemplates{
		signal: "channel",
		tests:  "channel_tests",
	}
	types := map[generator]templateNames{
		{
			InterfaceProps: signedProps,
			Builtin:        "int8",
			Name:           "i8",
			Allocator:      "Int8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int16",
			Name:           "i16",
			Allocator:      "Int16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int32",
			Name:           "i32",
			Allocator:      "Int32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps:   signedProps,
			Builtin:          "int64",
			Name:             "i64",
			Allocator:        "Int64",
			MaxBitDepth:      "BitDepth64",
			ChannelTemplates: &channelTemplates,
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint8",
			Name:           "u8",
			Allocator:      "Uint8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint16",
			Name:           "u16",
			Allocator:      "Uint16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint32",
			Name:           "u32",
			Allocator:      "Uint32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps:   unsignedProps,
			Builtin:          "uint64",
			Name:             "u64",
			Allocator:        "Uint64",
			MaxBitDepth:      "BitDepth64",
			ChannelTemplates: &channelTemplates,
		}: fixedTemplates,
		{
			InterfaceProps: floatingProps,
			Builtin:        "float32",
			Name:           "f32",
			Allocator:      "Float32",
		}: floatingTemplates,
		{
			InterfaceProps:   floatingProps,
			Builtin:          "float64",
			Name:             "f64",
			Allocator:        "Float64",
			ChannelTemplates: &channelTemplates,
		}: floatingTemplates,
	}

	for g, t := range types {
		g.Timestamp = time.Now()

		generate(t.signal, g, fmt.Sprintf("%s.go", g.Builtin))
		generate(t.tests, g, fmt.Sprintf("%s_test.go", g.Builtin))
		if g.ChannelTemplates != nil {
			generate(g.ChannelTemplates.signal, g, fmt.Sprintf("channel_%s.go", strings.ToLower(g.Interface)))
			generate(g.ChannelTemplates.tests, g, fmt.Sprintf("channel_%s_test.go", strings.ToLower(g.Interface)))
		}
	}
}

func generate(templateName string, gen generator, fileName string) {

	var raw bytes.Buffer
	err := templates.ExecuteTemplate(&raw, templateName, gen)
	die(fmt.Sprintf("execute %s template for %s type", templateName, gen.Name), err)
	formatted, err := format.Source(raw.Bytes())
	die(fmt.Sprintf("formatting file for %s type", templateName, gen.Name), err)

	f, err := os.Create(fileName)
	die(fmt.Sprintf("create %s file", fileName), err)
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(formatted))
	die(fmt.Sprintf("writing file for %s type", templateName, gen.Name), err)
}

func die(reason string, err error) {
	if err != nil {
		panic(fmt.Sprintf("failed %s: %v", reason, err))
	}
}
