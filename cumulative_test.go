package randomness

import (
	"fmt"
	"testing"
)

func TestCumulativeTestSample(t *testing.T) {
	p, q := CumulativeTest(sampleTestBits100, true)
	fmt.Printf("Foward, n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.219194" || fmt.Sprintf("%.6f", q) != "0.219194" {
		t.FailNow()
	}
	p, q = CumulativeTest(sampleTestBits100, false)
	fmt.Printf("Backward, n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.114866" || fmt.Sprintf("%.6f", q) != "0.114866" {
		t.FailNow()
	}
}
