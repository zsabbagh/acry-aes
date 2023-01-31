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

func add_round_key(state, key []uint32) {
	for i := 0; uint32(i) < Nb; i++ {
		state[i] ^= key[i]
	}
}

func transpose(in []uint32) []uint32 {
	state := make([]uint32, 4)
	get := func(row, col int) uint32 {
		val := in[row] & (0xFF000000 >> (8 * col))
		val <<= (8 * col)
		return val >> (8 * row)
	}
	for row := 0; uint32(row) < Nb; row++ {
		state[row] = get(0, row) ^ get(1, row) ^ get(2, row) ^ get(3, row)
	}
	return state
}

func encrypt(state, key []uint32, rounds int) {

	add_round_key(state, key)

	for round := 1; uint32(round) < Nr; round++ {
		sub_bytes(state, false)
		state = shift_rows(state, false)
		state = mix_columns(state, false)
		add_round_key(state, key)
	}

	sub_bytes(state, false)
	state = shift_rows(state, false)
	add_round_key(state, key)
}

func decrypt(state, key []uint32, rounds int) {

	add_round_key(state, key)

	for round := 1; uint32(round) < Nr; round++ {
		state = shift_rows(state, true)
		sub_bytes(state, true)
		add_round_key(state, key)
		state = mix_columns(state, true)
	}

	state = shift_rows(state, true)
	sub_bytes(state, true)
	add_round_key(state, key)
}

func main() {
	// state should be initialised with transpose
}
