package main

type AES struct {
	key_init []byte // the initialisation key
	key      []byte // the expanded key
	nrounds  int
	nbytes   int
	nkey     int
	index    int
}

func copy(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = src[i]
	}
	return dst
}

func create_aes(key []byte, nr, nb, nk int) AES {
	aes := AES{key, []byte{}, nr, nb, nk, 0}
	aes.key_expansion()
	return aes
}

// Variable sized state
func sub_bytes(state []byte, inverse bool) {
	for i := 0; i < len(state); i++ {
		if inverse {
			state[i] = sbox_inverse[state[i]]
		} else {
			state[i] = sbox[state[i]]
		}
	}
}

// rotates word (coming 4 bytes) with amount, neg is left, pos is right
func rotate_word(state []byte, amount int) {
	a, b, c, d := state[0], state[1], state[2], state[3]
	state[(4+amount)%4] = a
	state[(4+amount+1)%4] = b
	state[(4+amount+2)%4] = c
	state[(4+amount+3)%4] = d
}

// State must be [16]byte
func shift_rows(state []byte, inverse bool) {
	var shift int
	for i := 0; i < 4; i++ {
		if inverse {
			shift = i
		} else {
			shift = -i
		}
		rotate_word(state[4*i:4*(i+1)], shift)
	}
}

// mode = 0 for regular, 1 for inverse
func mix_columns(state []byte, inverse bool) {
	var m1, m2, m3, m4 []byte
	if inverse {
		m1 = mul14[:]
		m2 = mul11[:]
		m3 = mul13[:]
		m4 = mul9[:]
	} else {
		m1 = mul2[:]
		m2 = mul3[:]
		m3 = id[:]
		m4 = id[:]
	}
	manipulate := func(i int) {
		a := m1[state[i]] ^ m2[state[4+i]] ^ m3[state[8+i]] ^ m4[state[12+i]]
		b := m4[state[i]] ^ m1[state[4+i]] ^ m2[state[8+i]] ^ m3[state[12+i]]
		c := m3[state[i]] ^ m4[state[4+i]] ^ m1[state[8+i]] ^ m2[state[12+i]]
		d := m2[state[i]] ^ m3[state[4+i]] ^ m4[state[8+i]] ^ m1[state[12+i]]
		state[i] = a
		state[4+i] = b
		state[8+i] = c
		state[12+i] = d
	}
	for i := 0; i < 4; i++ {
		manipulate(i)
	}
}

func (aes *AES) add_round_key(state []byte, inverse bool) {
	for i := 0; i < len(state); i++ {
		state[i] ^= aes.key[4*aes.index+i]
	}
	if inverse {
		aes.index -= aes.nkey
	} else {
		aes.index += aes.nkey
	}
}

// state must be 16 in length
func transpose(state []byte) {
	var a, b int
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			a = 4*j + i
			b = 4*i + j
			state[a], state[b] = state[b], state[a]
		}
	}
}

func (aes *AES) key_expansion() {

	aes.key = make([]byte, 4*aes.nbytes*(aes.nrounds+1))
	// the key occupies the first nwords slots of the expanded key
	var i int
	copy_word := func(src, dst []byte) {
		for i := 0; i < 4; i++ {
			dst[i] = src[i]
		}
	}
	for i = 0; i < aes.nkey; i++ {
		copy_word(aes.key_init[4*i:], aes.key[4*i:])
	}
	for i = aes.nkey; i < aes.nbytes*(aes.nrounds+1); i++ {
		copy_word(aes.key[4*(i-1):], aes.key[4*i:])
		if (i/aes.nbytes)%aes.nkey == 0 {
			rotate_word(aes.key[4*i:], -1)
			sub_bytes(aes.key[4*i:4*(i+1)], false)
			aes.key[4*i] ^= rcon[i/aes.nkey]
		} else if aes.nkey > 6 && i%4 == 0 {
			sub_bytes(aes.key[4*i:4*(i+1)], false)
		}
		for j := 0; j < aes.nbytes; j++ {
			aes.key[4*i+j] ^= aes.key[4*(i-aes.nkey)+j]
		}
	}

	for i := 0; i+16 <= len(aes.key); i += 4 {
		transpose(aes.key[i : i+16])
	}
}

func words_to_bytes(state []uint32, out []byte) {
	var shift uint32
	for i := 0; i < len(out); i++ {
		shift = uint32(8 * (i % 4))
		out[i] = byte((state[i/4] & (0xFF000000 >> shift)) >> (24 - shift))
	}
}

func (aes *AES) encrypt(state []byte) {
	aes.index = 0
	transpose(state)
	aes.add_round_key(state, false)

	for round := 1; round < aes.nrounds; round++ {
		sub_bytes(state, false)
		shift_rows(state, false)
		mix_columns(state, false)
		aes.add_round_key(state, false)
	}

	sub_bytes(state, false)
	shift_rows(state, false)
	aes.add_round_key(state, false)
}

func (aes *AES) decrypt(state []byte) {
	aes.index = len(aes.key) - aes.nkey
	transpose(state)
	aes.add_round_key(state, true)

	for round := 1; round < aes.nrounds; round++ {
		shift_rows(state, true)
		sub_bytes(state, true)
		aes.add_round_key(state, true)
		mix_columns(state, true)
	}

	shift_rows(state, true)
	sub_bytes(state, true)
	aes.add_round_key(state, true)
}

func main() {
	// state should be initialised with transpose
}
