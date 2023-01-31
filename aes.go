package main

func substitute_word(word uint32) uint32 {
	get := func(i uint32) uint32 {
		val := word & (0xFF << (8 * (Nb - i - 1)))
		val = val >> (8 * (Nb - i - 1))
		return uint32(sbox[(val>>4)&0xF][val&0xF])
	}
	return (get(0) << 24) + (get(1) << 16) + (get(2) << 8) + get(3)
}

func sub_bytes(arr []uint32) {
	for row := 0; row < len(arr); row++ {
		arr[row] = substitute_word(arr[row])
	}
}

func rotate_left(word uint32, amount uint32) uint32 {
	get := func(i uint32) uint32 {
		pos := 8 * (Nb - i - 1)
		val := word & (0xFF << pos)
		if amount > i {
			return val >> (pos - 8*(amount-i-1))
		}
		return val << (8 * (amount))
	}
	return get(0) ^ get(1) ^ get(2) ^ get(3)
}

func rotate_right(word uint32, amount uint32) uint32 {
	get := func(i uint32) uint32 {
		pos := 8 * (Nb - i - 1)
		val := word & (0xFF << pos)
		if amount > (3 - i) {
			return val << (pos - 8*(amount-i-1))
		}
		return val >> (8 * (amount))
	}
	return get(0) ^ get(1) ^ get(2) ^ get(3)
}

func shift_rows(state []uint32) []uint32 {
	out := make([]uint32, 4)
	for r := 0; uint32(r) < Nb; r++ {
		out[r] = rotate_left(state[r], uint32(r))
	}
	return out
}

// mode = 0 for regular, 1 for inverse
func mix_columns(state []uint32) []uint32 {
	words := make([]uint32, 4)
	words[0] = mul2[state[0]] ^ mul3[state[1]] ^ state[2] ^ state[3]
	words[1] = state[0] ^ mul2[state[1]] ^ mul3[state[2]] ^ state[3]
	words[2] = state[0] ^ state[1] ^ mul2[state[2]] ^ mul3[state[3]]
	words[3] = mul3[state[0]] ^ state[1] ^ state[2] ^ mul2[state[3]]
	return words
}

func mix_columns_inverse() {}
func shift_rows_inverse()  {}

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
		sub_bytes(state)
		state = shift_rows(state)
		state = mix_columns(state)
		add_round_key(state, key)
	}

	sub_bytes(state)
	state = shift_rows(state)
	add_round_key(state, key)
}

func decrypt() {}

func main() {
	// state should be initialised with transpose
}
