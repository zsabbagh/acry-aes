package main

import "testing"

func test_value(t *testing.T, got, want uint32) {
	if got != want {
		t.Errorf("Sum was incorrect, got: 0x%x, want: 0x%x.", got, want)
	}
}

func TestSubstituteWord(t *testing.T) {
	total := substitute_word(0x53535353)
	want := uint32(0xedededed)
	test_value(t, total, want)
	total = substitute_word(0x4014587f)
	want = uint32(0x09fa6ad2)
	test_value(t, total, want)
}

func TestSubBytes(t *testing.T) {
	arr := [4]uint32{
		0x00005300,
		0x54000000,
		0x59000059,
		0x00000053,
	}
	sub_bytes(&arr)
	want := uint32(0x6363ed63)
	test_value(t, arr[0], want)
	want = 0x20636363
	test_value(t, arr[1], want)
	want = 0xcb6363cb
	test_value(t, arr[2], want)
	want = 0x636363ed
	test_value(t, arr[3], want)
}

func TestRotateWord(t *testing.T) {
	got := rotate_left(0x01020304, 1)
	want := uint32(0x02030401)
	test_value(t, got, want)
}
