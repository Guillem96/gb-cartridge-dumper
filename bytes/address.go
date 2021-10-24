package bytes

func AddressToBitArray(addr uint16) []bool {
	bits := make([]bool, 16)
	for i := 0; i < 16; i++ {
		bits[i] = TestBit(addr, uint(i))
	}
	return bits
}

func TestBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}
