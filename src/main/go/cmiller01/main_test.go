package main

import (
	_ "embed"
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
	if measurement.max != 534 {
		t.Fatalf("invalid max")
	}
	if measurement.sum != 71321 {
		t.Fatalf("invalid sum, got %d", measurement.sum)
	}
}
