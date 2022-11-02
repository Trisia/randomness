package randomness

import (
	"fmt"
	"testing"
)

func TestFrequencyWithinBlockTestSample(t *testing.T) {
	p, q := FrequencyWithinBlockTest(sampleTestBits100)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.706438" || fmt.Sprintf("%.6f", q) != "0.706438" {
		t.FailNow()
	}
}
