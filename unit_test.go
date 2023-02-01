package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

func test_value(t *testing.T, got, want uint32) {
	if got != want {
		t.Errorf("Sum was incorrect, got: 0x%x, want: 0x%x.", got, want)
	}
}

func test_array_equals(t *testing.T, got, want []byte) {
	if len(got) != len(want) {
		t.Errorf("Length do not match, got: %x, want: %d.", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("Byte do not match in pos %d, got: 0x%x, want: 0x%x.", i, got[i], want[i])
		}
	}
}

func TestKeyExpansion(t *testing.T) {
	key := "YELLOW SUBMARINE"
	expandedkey := []uint32{0x594f5552, 0x45574249, 0x4c204d4e, 0x4c534145, 0x632c792b,
		0x6a3d7f36, 0x22024f01, 0x4c1f5e1b, 0x6448311a, 0x162b5462, 0x8d8fc0c1, 0xbda2fce7,
		0xca82b3a9, 0x6e451173, 0x19965697, 0x1fbd41a6, 0x4dcf7cd5, 0xe6a3b2c1, 0x3dabfd6a,
		0xcc713096, 0x25ea9643, 0xe447f534, 0xad06fb91, 0xcfbe8e18, 0x1df76122, 0x6522d7e3,
		0x6fd6c, 0xd56be5fd, 0x4cbbdaf8, 0x3517c023, 0x5452afc3, 0x462dc835, 0xea518b73,
		0x1b0cccef, 0xc2903ffc, 0x72ae2d7, 0x2e7ff487, 0xaba76b84, 0xcc5c639f, 0x88a24097,
		0x4738cc4b, 0x70d7bc38, 0x44187be4, 0x9f3d7dea}
	aes := create_aes([]byte(key), Nr, Nb, Nk)
	equals := func(exp uint32, act []byte) bool {
		for i := 0; i < 4; i++ {
			if byte(exp>>(3-i)) != act[i] {
				return false
			}
		}
		return true
	}
	for i := range expandedkey {
		if !equals(expandedkey[i], aes.key[4*i:]) {
			t.Errorf("Expanded key for %s is incorrect\n", key)
			break
		}
	}
}

func TestSubBytes(t *testing.T) {
	arr := []byte{
		0x00, 0x00, 0x53, 0x00,
		0x54, 0x00, 0x00, 0x00,
		0x59, 0x00, 0x00, 0x59,
		0x00, 0x00, 0x00, 0x53,
	}
	sub_bytes(arr, false)
	want := []byte{
		0x63, 0x63, 0xed, 0x63,
		0x20, 0x63, 0x63, 0x63,
		0xcb, 0x63, 0x63, 0xcb,
		0x63, 0x63, 0x63, 0xed,
	}
	test_array_equals(t, arr, want)
}

func TestTranspose(t *testing.T) {
	state := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	transpose(state)
	want := []byte{0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14, 3, 7, 11, 15}
	test_array_equals(t, state, want)
}

func TestShiftRows(t *testing.T) {
	input := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0x4d, 0xdc, 0x01, 0xc6,
		0xa1, 0x58, 0x01, 0xc6,
		0xbc, 0x9d, 0x01, 0xc6,
	}
	expected := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0xdc, 0x01, 0xc6, 0x4d,
		0x01, 0xc6, 0xa1, 0x58,
		0xc6, 0xbc, 0x9d, 0x01,
	}
	shift_rows(input, false)
	test_array_equals(t, input, expected)
}

func TestShiftRowsInverse(t *testing.T) {
	input := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0x4d, 0xdc, 0x01, 0xc6,
		0xa1, 0x58, 0x01, 0xc6,
		0xbc, 0x9d, 0x01, 0xc6,
	}
	expected := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0xdc, 0x01, 0xc6, 0x4d,
		0x01, 0xc6, 0xa1, 0x58,
		0xc6, 0xbc, 0x9d, 0x01,
	}
	shift_rows(expected, true)
	test_array_equals(t, expected, input)
}

func TestWordsToBytes(t *testing.T) {
	input := []uint32{0x01020304, 0x01020304, 0x01020304, 0x010203FF}
	out := [16]byte{}
	words_to_bytes(input, out[:])
	if out[0] != 0x01 {
		t.Errorf("Expected 0x01")
	}
	if out[15] != 0xFF {
		t.Errorf("Expected 0xFF")
	}
}

func TestData(t *testing.T) {
	data, err := os.ReadFile("./data/aes_sample.in")
	if err != nil {
		t.Errorf("Could not open file.")
	}
	key := data[:16]
	expected_key, err := hex.DecodeString("F4C020A0A1F604FD343FAC6A7E6AE0F9")
	if err != nil {
		t.Errorf("Could not parse KEY string")
	}
	test_array_equals(t, key, expected_key)
	aes := create_aes(key, Nr, Nb, Nk)
	fmt.Println(len(data))
	data = data[16:]
	expected, err := hex.DecodeString("F295B9318B994434D93D98A4E449AFD8")
	if err != nil {
		t.Errorf("Could not parse DATA string")
	}
	test_array_equals(t, data, expected)
	for i := 0; i+16 <= len(data); i += 16 {
		aes.encrypt(data[i : i+16])
	}
	expected, err = hex.DecodeString("52E418CBB1BE4949308B381691B109FE")
	if err != nil {
		t.Errorf("Could not parse string")
	}
	test_array_equals(t, data, expected)
}

func TestMixColumns(t *testing.T) {
	input := []byte{
		0xdb, 0xf2, 0x01, 0xc6,
		0x13, 0x0a, 0x01, 0xc6,
		0x53, 0x22, 0x01, 0xc6,
		0x45, 0x5c, 0x01, 0xc6}
	expected := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0x4d, 0xdc, 0x01, 0xc6,
		0xa1, 0x58, 0x01, 0xc6,
		0xbc, 0x9d, 0x01, 0xc6}
	mix_columns(input, false)
	test_array_equals(t, input, expected)
}

func TestMixInverse(t *testing.T) {
	input := []byte{
		0xdb, 0xf2, 0x01, 0xc6,
		0x13, 0x0a, 0x01, 0xc6,
		0x53, 0x22, 0x01, 0xc6,
		0x45, 0x5c, 0x01, 0xc6}
	expected := []byte{
		0x8e, 0x9f, 0x01, 0xc6,
		0x4d, 0xdc, 0x01, 0xc6,
		0xa1, 0x58, 0x01, 0xc6,
		0xbc, 0x9d, 0x01, 0xc6}
	mix_columns(expected, true)
	test_array_equals(t, expected, input)
}
