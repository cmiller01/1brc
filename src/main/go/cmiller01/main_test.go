package main

import (
	_ "embed"
	"fmt"
	"testing"
)

//go:embed measurements_100000.txt
var measurementsBench []byte

func TestRound(t *testing.T) {
	// weird rounding issue
	res := round((101.79999999999998) / float64(4))
	if res != 25.5 {
		t.Errorf("got %f", res)
		t.Fail()
	}
}

func TestParseTemp(t *testing.T) {
	testTable := []struct {
		input  []byte
		output float64
	}{

		{input: []byte("3.2"), output: 3.2},
		{input: []byte("-20.1"), output: -20.1},
		{input: []byte("-0.3"), output: -0.3},
	}
	for _, tt := range testTable {
		t.Run(fmt.Sprintf("%s", string(tt.input)), func(t *testing.T) {
			res := parseTemp(tt.input)
			if res != tt.output {
				t.Fatalf("expected %f got %f", tt.output, res)
			}
		})
	}
}

func BenchmarkParseTemp(b *testing.B) {
	testCases := [][]byte{
		[]byte("3.2"),
		[]byte("-20.1"),
		[]byte("-0.3"),
		[]byte("99.9"),
		[]byte("-99.9"),
	}
	for b.Loop() {
		for _, tc := range testCases {
			parseTemp(tc)
		}
	}
}

func BenchmarkProcessChunk(b *testing.B) {
	results := make(map[string]*measurements, 10_000)
	for b.Loop() {
		processChunk(measurementsBench, results)
	}
}

func TestProcessChunk(t *testing.T) {
	results := make(map[string]*measurements, 10_000)
	processChunk(measurementsBench, results)
	// spot check one
	measurement := results["Juba"]
	if measurement == nil {
		t.Fatalf("expected measurement not to be nil")
	}
	if measurement.count != 253 {
		t.Fatalf("invalid count")
	}
	if measurement.max != 53.4 {
		t.Fatalf("invalid max, got %f", measurement.max)
	}
	if measurement.sum != 713.5 {
		t.Fatalf("invalid sum, got %f", measurement.sum)
	}
}
