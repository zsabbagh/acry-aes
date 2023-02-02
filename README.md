# AES Implementation

---

> GitHub repo available for implementation
https://github.com/zsabbagh/acry-aes, request access when necessary

Zakaria Sabbagh, zsabbagh@kth.se
> 

# Introduction

I worked on slices of bytes instead of array of words. This made certain operations different. I wrote the algorithm in Go. Sources used are “A Graduate Course in Applied Cryptography” by Dan Boneh and Victor Shoup, and the official [NIST publication from 2001](https://nvlpubs.nist.gov/nistpubs/fips/nist.fips.197.pdf).

The preparation works as follows:

1. Expand input key with Key Expansion algorithm
2. Split input into array of states, each a 16 byte array

Then, with this, do the following for **********************each state**********************

1. Transpose the state
2. Add Round Key
3. For Rounds-1 time, do the following
    1. Sub Bytes according to the [Rijndael S-Box](https://en.wikipedia.org/wiki/Rijndael_S-box)
    2. Shift Rows
    3. Mix Columns
    4. Add Round Key
4. Sub Bytes
5. Shift Rows
6. Add Round Key
7. Transpose the state

# Add Round Key

This is a basic $\oplus$ (xor) operation on the key and state.

```python
for i := 0...16:
	state[i] ^= key[i]
```

# Sub Bytes

This substitutes each byte in accordance with the [Rijndael S-Box](https://en.wikipedia.org/wiki/Rijndael_S-box).

```python
for i := 0...16:
	state[i] = sbox[state[i]]
```

# Shift Rows

This shifts row $i$ with that many steps to the left:

```python
for i := 0...4:
	rotate_word(state[i:i+4], -i)
```

Where rotate simply rotates the slice `i` steps to the left because of the negative value. See repo for more information.

# Mix Columns

This implements and uses the [mix-columns array of Rijndael](https://en.wikipedia.org/wiki/Rijndael_MixColumns). For each column, do the following

```python
a := m1[state[i]] ^ m2[state[4+i]] ^ m3[state[8+i]] ^ m4[state[12+i]]
b := m4[state[i]] ^ m1[state[4+i]] ^ m2[state[8+i]] ^ m3[state[12+i]]
c := m3[state[i]] ^ m4[state[4+i]] ^ m1[state[8+i]] ^ m2[state[12+i]]
d := m2[state[i]] ^ m3[state[4+i]] ^ m4[state[8+i]] ^ m1[state[12+i]]
state[i] = a
state[4+i] = b
state[8+i] = c
state[12+i] = d
```

# Key Expansion

The following is a summary of the algorithm:

1. Create space for 11 (number of rounds is 10+1) sets of 16 byte keys (where each 16 byte key is a set of 4 words).
2. Copy the input key into the first 16 bytes of the space
3. For each new key-set, do the following
    1. Calculate G-Function of the most previous word
    2. XOR the previous key-sets first word with the word from a) and store in the current key-sets first word
    3. For each coming word in the current key, set it equal to the previous key-sets’ word (with the same index) XOR:ed with the most previous word in the key
4. For each key-set, Transpose the set as seen as a matrix

# G-Function

Below `glookup` is the [RCON array found here](https://en.wikipedia.org/wiki/AES_key_schedule)

```python
out := make([]byte, 4)
for i := 0...3:
	out[i] = word[i]
rotate_word(out, -1)
sub_bytes(out, false)
out[0] ^= glookup[i-1]
return out
``
