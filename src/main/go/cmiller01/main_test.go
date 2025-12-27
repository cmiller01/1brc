package main

import "testing"

func TestRound(t *testing.T) {
	// weird rounding issue
	res := round((101.79999999999998) / float64(4))
	if res != 25.5 {
		t.Errorf("got %f", res)
		t.Fail()
	}
}
