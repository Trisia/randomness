package randomness

import (
	"fmt"
	"math"
)

// MonoBitFrequencyTest 单比特频数监测
func MonoBitFrequencyTest(bits []bool) float64 {
	if len(bits) == 0 {
		fmt.Println("MonoBitFrequencyTest:arg wrong")
		return -1
	}
	n := len(bits)
	S := 0
	var V float64
	var P float64
	for _, bit := range bits {
		if bit {
			S++
		} else {
			S--
		}
	}
	V = math.Abs(float64(S)) / math.Sqrt(float64(n))
	P = math.Erfc(V / math.Sqrt(2))
	return P
}
