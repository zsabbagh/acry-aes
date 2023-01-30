package main

import "testing"

func TestSum(t *testing.T) {
	total := substitute(0x53)
	if total != sbox[5][3] {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, sbox[5][3])
	}
}
