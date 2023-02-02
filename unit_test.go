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
	key, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	expected_key, err := hex.DecodeString("000102030405060708090a0b0c0d0e0fd6aa74fdd2af72fadaa678f1d6ab76feb692cf0b643dbdf1be9bc5006830b3feb6ff744ed2c2c9bf6c590cbf0469bf4147f7f7bc95353e03f96c32bcfd058dfd3caaa3e8a99f9deb50f3af57adf622aa5e390f7df7a69296a7553dc10aa31f6b14f9701ae35fe28c440adf4d4ea9c02647438735a41c65b9e016baf4aebf7ad2549932d1f08557681093ed9cbe2c974e13111d7fe3944a17f307a78b4d2b30c5")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	transpose_sets(expected_key)
	aes := create_aes([]byte(key), Nr, Nb, Nk)
	test_array_equals(t, aes.key, expected_key)
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

func TestCase(t *testing.T) {
	key, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	expected_key, err := hex.DecodeString("000102030405060708090a0b0c0d0e0fd6aa74fdd2af72fadaa678f1d6ab76feb692cf0b643dbdf1be9bc5006830b3feb6ff744ed2c2c9bf6c590cbf0469bf4147f7f7bc95353e03f96c32bcfd058dfd3caaa3e8a99f9deb50f3af57adf622aa5e390f7df7a69296a7553dc10aa31f6b14f9701ae35fe28c440adf4d4ea9c02647438735a41c65b9e016baf4aebf7ad2549932d1f08557681093ed9cbe2c974e13111d7fe3944a17f307a78b4d2b30c5")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	transpose_sets(expected_key)
	aes := create_aes([]byte(key), Nr, Nb, Nk)
	test_array_equals(t, aes.key, expected_key)
	in, err := hex.DecodeString("00112233445566778899aabbccddeeff")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	expected, err := hex.DecodeString("69c4e0d86a7b0430d8cdb78070b4c55a")
	if err != nil {
		t.Fatalf("Could not decode string")
	}
	out := aes.Encrypt(in)
	test_array_equals(t, out, expected)
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
	encrypted_data := aes.Encrypt(data)
	expected, err = hex.DecodeString("52E418CBB1BE4949308B381691B109FE")
	if err != nil {
		t.Errorf("Could not parse string")
	}
	test_array_equals(t, encrypted_data, expected)
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
