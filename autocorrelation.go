package randomness

import (
	"fmt"
	"math"
)

// AutocorrelationTest 自相关检测
func AutocorrelationTest(bits []bool) float64 {
	n := len(bits)
	if n < 16 {
		fmt.Println("AutocorrelationTest:args wrong")
		return -1
	}
	d := 16

	Ad := 0
	var V float64 = 0
	var P float64 = 0

	for i := 0; i < n-d; i++ {
		if xor(bits[i], bits[i+d]) {
			Ad++
		}
	}

	V = 2.0 * (float64(Ad) - (float64(n-d) / 2.0)) / math.Sqrt(float64(n-d))
	P = math.Erfc(math.Abs(V) / math.Sqrt(2))
	return P
}
