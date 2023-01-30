package main

import "fmt"

var Nb int = 4  // byte length in 32 bit words
var Nk int = 4  // key length in 32 bit words
var Nr int = 10 // 10 when nk=4, 12 nk=6, 14 nk=8
var sbox = [16][16]byte{{99, 124, 119, 123, 242, 107, 111, 197, 48, 1, 103, 43, 254, 215, 171, 118}, {202, 130, 201, 125, 250, 89, 71, 240, 173, 212, 162, 175, 156, 164, 114, 192}, {183, 253, 147, 38, 54, 63, 247, 204, 52, 165, 229, 241, 113, 216, 49, 21}, {4, 199, 35, 195, 24, 150, 5, 154, 7, 18, 128, 226, 235, 39, 178, 117}, {9, 131, 44, 26, 27, 110, 90, 160, 82, 59, 214, 179, 41, 227, 47, 132}, {83, 209, 0, 237, 32, 252, 177, 91, 106, 203, 190, 57, 74, 76, 88, 207}, {208, 239, 170, 251, 67, 77, 51, 133, 69, 249, 2, 127, 80, 60, 159, 168}, {81, 163, 64, 143, 146, 157, 56, 245, 188, 182, 218, 33, 16, 255, 243, 210}, {205, 12, 19, 236, 95, 151, 68, 23, 196, 167, 126, 61, 100, 93, 25, 115}, {96, 129, 79, 220, 34, 42, 144, 136, 70, 238, 184, 20, 222, 94, 11, 219}, {224, 50, 58, 10, 73, 6, 36, 92, 194, 211, 172, 98, 145, 149, 228, 121}, {231, 200, 55, 109, 141, 213, 78, 169, 108, 86, 244, 234, 101, 122, 174, 8}, {186, 120, 37, 46, 28, 166, 180, 198, 232, 221, 116, 31, 75, 189, 139, 138}, {112, 62, 181, 102, 72, 3, 246, 14, 97, 53, 87, 185, 134, 193, 29, 158}, {225, 248, 152, 17, 105, 217, 142, 148, 155, 30, 135, 233, 206, 85, 40, 223}, {140, 161, 137, 13, 191, 230, 66, 104, 65, 153, 45, 15, 176, 84, 187, 22}}

func substitute(b byte) byte {
	return sbox[(int)((b>>4)&0xF)][(int)(b&0xF)]
}

func mixColumns()  {}
func shiftRows()   {}
func addRoundKey() {}

func subBytes(arr *[4][4]byte) {
	var i, j int
	for row := 0; row < len(*arr); row++ {
		for col := 0; col < len((*arr)[0]); col++ {
			i = (int)((*arr)[row][col]>>4) & 0xFF
			j = (int)((*arr)[row][col] & 0xFF)
			(*arr)[row][col] = sbox[i][j]
		}
	}
}

func cipher(in [16]byte) [16]byte {
	// Init state
	state := [4][4]byte{}
	for row := 0; row < Nb; row++ {
		for col := 0; col < Nb; col++ {
			state[row][col] = in[row+4*col]
		}
	}

	addRoundKey()

	for round := 1; round < Nr; round++ {
		subBytes(&state)
		shiftRows()
		mixColumns()
		addRoundKey()
	}

	subBytes(&state)
	shiftRows()
	addRoundKey()

	// Create out
	out := [16]byte{}
	for r := 0; r < Nb; r++ {
		for c := 0; c < Nb; c++ {
			out[r+4*c] = state[r][c]
		}
	}
	return out
}

func main() {
	fmt.Println(substitute(0))
	fmt.Printf("%x", substitute(0x53))
}
