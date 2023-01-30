package main

func substitute_word(word uint32) uint32 {
	get := func(i uint32) uint32 {
		val := word & (0xFF << (8 * (Nb - i - 1)))
		val = val >> (8 * (Nb - i - 1))
		return uint32(sbox[(val>>4)&0xF][val&0xF])
	}
	return (get(0) << 24) + (get(1) << 16) + (get(2) << 8) + get(3)
}

func sub_bytes(arr *[4]uint32) {
	for row := 0; row < len(*arr); row++ {
		(*arr)[row] = substitute_word((*arr)[row])
	}
}

func rotate_left(word uint32, amount uint32) uint32 {
	get := func(i uint32) uint32 {
		shift := 8 * (Nb - i - 1)
		val := word & (0xFF << shift)
		if amount > i {
			return val >> (shift - 8*(amount-i-1))
		}
		return val << (8 * (amount))
	}
	return get(0) + get(1) + get(2) + get(3)
}

func shift_rows(state [4]uint32) [4]uint32 {
	out := [4]uint32{}
	for r := 0; uint32(r) < Nb; r++ {
		out[r] = rotate_left(state[r], uint32(r))
	}
	return out
}

// mode = 0 for regular, 1 for inverse
func mix_columns(state [4]uint32) [4]uint32 {
	words := [4]uint32{}
	words[0] = mul2[state[0]] ^ mul3[state[1]] ^ state[2] ^ state[3]
	words[1] = state[0] ^ mul2[state[1]] ^ mul3[state[2]] ^ state[3]
	words[2] = state[0] ^ state[1] ^ mul2[state[2]] ^ mul3[state[3]]
	words[3] = mul3[state[0]] ^ state[1] ^ state[2] ^ mul2[state[3]]
	return words
}

func add_round_key() {}

func cipher(in [4]uint32) [4]uint32 {
	// Init state
	state := [4]uint32{}
	for row := 0; uint32(row) < Nb; row++ {
		state[row] = in[row]
	}

	add_round_key()

	for round := 1; uint32(round) < Nr; round++ {
		sub_bytes(&state)
		state = shift_rows(state)
		state = mix_columns(state)
		add_round_key()
	}

	sub_bytes(&state)
	state = shift_rows(state)
	add_round_key()

	// Create out
	out := [4]uint32{}
	for r := 0; uint32(r) < Nb; r++ {
		out[r] = state[r]
	}
	return out
}

func main() {

}
