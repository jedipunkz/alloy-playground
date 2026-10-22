package main

import "testing"

func TestParseLoadavg(t *testing.T) {
	got, err := parseLoadavg("0.52 0.58 0.59 1/389 12345\n")
	if err != nil {
		t.Fatal(err)
	}
	if want := [3]float64{0.52, 0.58, 0.59}; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if _, err := parseLoadavg("0.52"); err == nil {
		t.Fatal("expected error for short input")
	}
}
