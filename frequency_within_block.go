package randomness

import (
	"fmt"
)

// FrequencyWithinBlockTest 块内频数监测
func FrequencyWithinBlockTest(bits []bool) float64 {
	n := len(bits)
	if n == 0 {
		fmt.Println("FrequencyTestWithinABlock:args wrong")
		return -1
	}
	m := 100
	N := n / m
	bits = bits[:N*m]

	var Pi float64 = 0
	var V float64 = 0
	var P float64 = 0

	var b bool
	for i := 0; i < N; i++ {
		Pi = 0
		for j := 0; j < m; j++ {
			b, bits = bits[0], bits[1:]
			if b {
				Pi++
			}
		}
		Pi = Pi / float64(m)
		V += (Pi - 0.5) * (Pi - 0.5)
	}
	V *= 4.0 * float64(m)

	P = igamc(float64(N)/2.0, V/2.0)
	return P
}
