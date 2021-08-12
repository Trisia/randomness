package randomness

import (
	"fmt"
	"math"
)

// CumulativeTest 累加和检测
func CumulativeTest(bits []bool) float64 {
	n := len(bits)

	if n == 0 {
		fmt.Println("CumulativeTest:args wrong")
		return -1
	}
	var S int = 0
	var Z int = 0
	var P float64 = 1.0
	for i := 0; i < n; i++ {
		if bits[i] {
			S++
		} else {
			S--
		}
		Z = max(Z, abs(S))
	}
	_n := float64(n)
	for i := ((-n / Z) + 1) / 4; i <= ((n/Z)-1)/4; i++ {
		P -= normal_CDF(float64((4*i+1)*Z)/math.Sqrt(_n)) - normal_CDF(float64((4*i-1)*Z)/math.Sqrt(_n))
	}
	for i := ((-n / Z) - 3) / 4; i <= ((n/Z)-1)/4; i++ {
		P += normal_CDF(float64((4*i+3)*Z)/math.Sqrt(_n)) - normal_CDF(float64((4*i+1)*Z)/math.Sqrt(_n))
	}
	return P
}
