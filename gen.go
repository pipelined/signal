// +build ignore

// This program generates all signal types.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	Pool        string
	MaxBitDepth string // used for fixed-point types only
}

// properties for interface generation
type InterfaceProps struct {
	Interface  string
	SampleType string
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
			Interface:  "Signed",
			SampleType: "int64",
		}
		unsignedProps = InterfaceProps{
			Interface:  "Unsigned",
			SampleType: "uint64",
		}
		floatingProps = InterfaceProps{
			Interface:  "Floating",
			SampleType: "float64",
		}
	)
	types := map[generator]templateNames{
		{
			InterfaceProps: signedProps,
			Builtin:        "int8",
			Name:           "Int8",
			Pool:           "i8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int16",
			Name:           "Int16",
			Pool:           "i16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int32",
			Name:           "Int32",
			Pool:           "i32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps: signedProps,
			Builtin:        "int64",
			Name:           "Int64",
			Pool:           "i64",
			MaxBitDepth:    "BitDepth64",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint8",
			Name:           "Uint8",
			Pool:           "u8",
			MaxBitDepth:    "BitDepth8",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint16",
			Name:           "Uint16",
			Pool:           "u16",
			MaxBitDepth:    "BitDepth16",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint32",
			Name:           "Uint32",
			Pool:           "u32",
			MaxBitDepth:    "BitDepth32",
		}: fixedTemplates,
		{
			InterfaceProps: unsignedProps,
			Builtin:        "uint64",
			Name:           "Uint64",
			Pool:           "u64",
			MaxBitDepth:    "BitDepth64",
		}: fixedTemplates,
		{
			InterfaceProps: floatingProps,
			Builtin:        "float32",
			Name:           "Float32",
			Pool:           "f32",
		}: floatingTemplates,
		{
			InterfaceProps: floatingProps,
			Builtin:        "float64",
			Name:           "Float64",
			Pool:           "f64",
		}: floatingTemplates,
	}

	for g, t := range types {
		g.Timestamp = time.Now()

		generate(t.signal, g, fmt.Sprintf("%s.go", g.Builtin))
		generate(t.tests, g, fmt.Sprintf("%s_test.go", g.Builtin))
	}
}

func generate(templateName string, gen generator, fileName string) {
	f, err := os.Create(fileName)
	die(fmt.Sprintf("create %s file", fileName), err)
	defer f.Close()

	var b bytes.Buffer
	err = templates.ExecuteTemplate(&b, templateName, gen)
	die(fmt.Sprintf("execute %s template for %s type", templateName, gen.Name), err)
	_, err = io.Copy(f, &b)
	die(fmt.Sprintf("writing file for %s type", templateName, gen.Name), err)
}

func die(reason string, err error) {
	if err != nil {
		panic(fmt.Sprintf("failed %s: %v", reason, err))
	}
}
