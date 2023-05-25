package randomness

import (
	"fmt"
	"testing"
)

func TestBinaryDerivativeTestSample(t *testing.T) {
	p, q := BinaryDerivativeTest(sampleTestBits128, 3)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.039669" || fmt.Sprintf("%.6f", q) != "0.980166" {
		t.FailNow()
	}
}
