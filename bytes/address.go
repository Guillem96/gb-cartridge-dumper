package bytes

// AddressToBitArray converts  a memory address to a bit array. The resulting bit array contains
// 16 bits, being the bit in the array position 0 the least significant and the one in 15th spot the most
// significant
func AddressToBitArray(addr uint16) []bool {
	bits := make([]bool, 16)
	for i := 0; i < 16; i++ {
		bits[i] = TestBit(addr, uint(i))
	}
	return bits
}

// TestBit tests if bit in position pos is active
func TestBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}
