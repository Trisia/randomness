package randomness

import (
	"fmt"
	"testing"
)

func TestLongestRunOfOnesInABlockTestSample(t *testing.T) {
	// check "1"
	p, q := LongestRunOfOnesInABlockTest(sampleTestBits128, true)
	fmt.Printf("Check 1, n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.180598" || fmt.Sprintf("%.6f", q) != "0.180598" {
		t.FailNow()
	}
	// check "0"
	p, q = LongestRunOfOnesInABlockTest(sampleTestBits128, false)
	fmt.Printf("Check 0, n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.839299" || fmt.Sprintf("%.6f", q) != "0.839299" {
		t.FailNow()
	}
}
