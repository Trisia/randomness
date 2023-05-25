package randomness

import (
	"fmt"
	"testing"
)

func TestMonoBitFrequencyTestSample(t *testing.T) {
	p, q := MonoBitFrequencyTest(sampleTestBits128)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.215925" || fmt.Sprintf("%.6f", q) != "0.892038" {
		t.FailNow()
	}
}
