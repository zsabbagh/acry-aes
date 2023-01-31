package main

import "testing"

func test_value(t *testing.T, got, want uint32) {
	if got != want {
		t.Errorf("Sum was incorrect, got: 0x%x, want: 0x%x.", got, want)
	}
}

func TestSubstituteWord(t *testing.T) {
	total := substitute_word(0x53535353, false)
	want := uint32(0xedededed)
	test_value(t, total, want)
	total = substitute_word(0x4014587f, false)
	want = uint32(0x09fa6ad2)
	test_value(t, total, want)
}

func TestSubBytes(t *testing.T) {
	arr := []uint32{
		0x00005300,
		0x54000000,
		0x59000059,
		0x00000053,
	}
	sub_bytes(arr, false)
	want := uint32(0x6363ed63)
	test_value(t, arr[0], want)
	want = 0x20636363
	test_value(t, arr[1], want)
	want = 0xcb6363cb
	test_value(t, arr[2], want)
	want = 0x636363ed
	test_value(t, arr[3], want)
}

func TestRotateLeft(t *testing.T) {
	got := rotate(0x01020304, 1, false)
	want := uint32(0x02030401)
	test_value(t, got, want)
}

func TestRotateRight(t *testing.T) {
	got := rotate(0x01020304, 1, true)
	want := uint32(0x04010203)
	test_value(t, got, want)
}

func TestRotateInverse(t *testing.T) {
	got := rotate(rotate(0x12345678, 1, true), 1, false)
	test_value(t, got, 0x12345678)
}

func TestTranspose(t *testing.T) {
	state := transpose([]uint32{0x01020304, 0x05060708, 0x09101112, 0x13141516})
	test_value(t, state[0], 0x01050913)
	test_value(t, state[1], 0x02061014)
	test_value(t, state[2], 0x03071115)
	test_value(t, state[3], 0x04081216)
}

func TestShiftRows(t *testing.T) {
	input := []uint32{
		0x8e9f01c6,
		0x4ddc01c6,
		0xa15801c6,
		0xbc9d01c6,
	}
	expected := []uint32{
		0x8e9f01c6,
		0xdc01c64d,
		0x01c6a158,
		0xc6bc9d01,
	}
	input = shift_rows(input, false)
	for i := 0; i < 4; i++ {
		test_value(t, input[0], expected[0])
	}
}
