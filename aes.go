package main

func substitute_word(word uint32, inverse bool) uint32 {
	var box *[256]byte
	if inverse {
		box = &sbox_inverse
	} else {
		box = &sbox
	}
	get := func(i uint32) uint32 {
		shift_amount := 8 * (Nb - i - 1)
		val := word & (0xFF << shift_amount)
		val = val >> shift_amount
		return uint32((*box)[val])
	}
	return (get(0) << 24) + (get(1) << 16) + (get(2) << 8) + get(3)
}

type AES struct {
	key    []uint32
	rounds uint32
}

func create_aes(key []byte, rounds uint32) AES {
	new_key := make([]uint32, 4)
	get := func(row, col uint32) uint32 {
		return uint32(key[row+4*col]) << (8 * (3 - col))
	}
	for i := uint32(0); i < 4; i++ {
		new_key[i] = get(i, 0) ^ get(i, 1) ^ get(i, 2) ^ get(i, 3)
	}
	aes := AES{new_key, rounds}
	return aes
}

func sub_bytes(arr []uint32, inverse bool) {
	for row := 0; row < len(arr); row++ {
		arr[row] = substitute_word(arr[row], inverse)
	}
}

// If inverse = true, then right, else direction is left
func rotate(word uint32, amount uint32, inverse bool) uint32 {
	get := func(i uint32) uint32 {
		pos := 8 * (Nb - i - 1)
		val := word & (0xFF << pos)
		if inverse {
			if amount > (3 - i) {
				return val << (pos - 8*(amount-i-1))
			}
			return val >> (8 * (amount))
		}
		if amount > i {
			return val >> (pos - 8*(amount-i-1))
		}
		return val << (8 * (amount))
	}

	return get(0) ^ get(1) ^ get(2) ^ get(3)
}

// Inverse is for decryption
func shift_rows(state []uint32, inverse bool) []uint32 {
	out := make([]uint32, 4)
	for r := 0; uint32(r) < Nb; r++ {
		out[r] = rotate(state[r], uint32(r), inverse)
	}
	return out
}

// mode = 0 for regular, 1 for inverse
func mix_columns(state []uint32, inverse bool) []uint32 {
	words := make([]uint32, 4)
	to_bytes := func(val uint32) uint8, uint8, uint8, uint8 {
		a := uint8(val & 0xFF000000 >> 24)
		b := uint8(val & 0x00FF0000 >> 16)
		c := uint8(val & 0x0000FF00 >> 8)
		d := uint8(val & 0x000000FF)
		return a,b,c,d
	}
	if inverse {
		words[0] = mul14[state[0]] ^ mul11[state[1]] ^ mul13[state[2]] ^ mul9[state[3]]
		words[1] = mul9[state[0]] ^ mul14[state[1]] ^ mul11[state[2]] ^ mul13[state[3]]
		words[2] = mul13[state[0]] ^ mul9[state[1]] ^ mul14[state[2]] ^ mul11[state[3]]
		words[3] = mul11[state[0]] ^ mul13[state[1]] ^ mul9[state[2]] ^ mul14[state[3]]
	} else {
		words[0] = mul2[state[0]] ^ mul3[state[1]] ^ state[2] ^ state[3]
		words[1] = state[0] ^ mul2[state[1]] ^ mul3[state[2]] ^ state[3]
		words[2] = state[0] ^ state[1] ^ mul2[state[2]] ^ mul3[state[3]]
		words[3] = mul3[state[0]] ^ state[1] ^ state[2] ^ mul2[state[3]]
	}
	return words
}

func (aes *AES) add_round_key(state []uint32) {
	for i := 0; uint32(i) < Nb; i++ {
		state[i] ^= aes.key[i]
	}
}

// In place transpose of a 4-length slice of uint32,
// seen as a square matrix of bytes
// where each uint32 is a row of 4 bytes
func transpose(in []uint32) {
	column_to_row := func(col uint32) uint32 {
		var shift_amount uint32 = 8 * col
		a := (in[0] & (0xFF000000 >> shift_amount)) << shift_amount
		b := (in[1] & (0xFF000000 >> shift_amount)) << shift_amount
		b >>= 8
		c := (in[2] & (0xFF000000 >> shift_amount)) << shift_amount
		c >>= 16
		d := (in[3] & (0xFF000000 >> shift_amount)) << shift_amount
		d >>= 24
		return a ^ b ^ c ^ d
	}
	a := column_to_row(0)
	b := column_to_row(1)
	c := column_to_row(2)
	d := column_to_row(3)
	in[0] = a
	in[1] = b
	in[2] = c
	in[3] = d
}

func init_state(in []byte) []uint32 {
	out := make([]uint32, 4)
	get := func(row, move uint32) uint32 {
		return uint32(in[row+4*move]) << (8 * (3 - move))
	}
	for i := uint32(0); i < 4; i++ {
		out[i] = get(i, 0) ^ get(i, 1) ^ get(i, 2) ^ get(i, 3)
	}
	return out
}

func words_to_bytes(state []uint32, out []byte) {
	var shift uint32
	for i := 0; i < len(out); i++ {
		shift = uint32(8 * (i % 4))
		out[i] = byte((state[i/4] & (0xFF000000 >> shift)) >> (24 - shift))
	}
}

func (aes *AES) encrypt(in []byte) {
	state := init_state(in)
	aes.add_round_key(state)

	for round := uint32(1); round < aes.rounds; round++ {
		sub_bytes(state, false)
		state = shift_rows(state, false)
		state = mix_columns(state, false)
		aes.add_round_key(state)
	}

	sub_bytes(state, false)
	state = shift_rows(state, false)
	aes.add_round_key(state)

	words_to_bytes(state, in)
}

func (aes *AES) decrypt(in []byte) {

	state := init_state(in)
	aes.add_round_key(state)

	for round := uint32(1); round < aes.rounds; round++ {
		state = shift_rows(state, true)
		sub_bytes(state, true)
		aes.add_round_key(state)
		state = mix_columns(state, true)
	}

	state = shift_rows(state, true)
	sub_bytes(state, true)
	aes.add_round_key(state)

	words_to_bytes(state, in)
}

func main() {
	// state should be initialised with transpose
}
