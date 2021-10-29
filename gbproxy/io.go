package gbproxy

import "math"

func writeToPins(value uint, pins []GameBoyPin) {
	for i, p := range pins {
		ab := uint(math.Pow(2, float64(i)))
		if (value & ab) >= uint(1) {
			p.High()
		} else {
			p.Low()
		}
	}
}
