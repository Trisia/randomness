package randomness

import (
	"fmt"
	"math"
)

// RunsTest 游程总数检测
func RunsTest(bits []bool) float64 {
	n := len(bits)
	if n == 0 {
		fmt.Println("RunsTest:args wrong")
		return -1
	}

	var Pi float64 = 0
	var V_obs int = 1
	var P float64 = 0

	for i := 0; i < n-1; i++ {
		if bits[i] != bits[i+1] {
			V_obs++
		}
		if bits[i] {
			Pi++
		}
	}
	if bits[n-1] {
		Pi++
	}
	Pi /= float64(n)
	P = math.Erfc(math.Abs(float64(V_obs)-2.0*float64(n)*Pi*(1.0-Pi)) / (2.0 * math.Sqrt(2.0*float64(n)) * Pi * (1.0 - Pi)))
	return P
}
