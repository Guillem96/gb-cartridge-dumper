package bytes

import (
	"errors"
	"fmt"
	"math"
)

// AddressToBitArray converts  a memory address to a bit array. The resulting bit array contains
// 16 bits, being the bit in the array position 0 the least significant and the one in 15th spot the most
// significant
func AddressToBitArray(addr uint, size int) ([]bool, error) {
	if !isPowerOfTwo(size) || size > 64 {
		errMsg := fmt.Sprintf("size %d not supported", size)
		return nil, errors.New(errMsg)
	}
	mask := math.Pow(float64(size), 2) - 1
	addr &= uint(mask)
	bits := make([]bool, size)
	for i := 0; i < size; i++ {
		bits[i] = TestBit(addr, uint(i))
	}
	return bits, nil
}

// TestBit tests if bit in position pos is active
func TestBit(n uint, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

func isPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
