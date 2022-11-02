package randomness

import (
	"fmt"
	"testing"
)

func TestRunsDistributionTestSample(t *testing.T) {
	p, q := RunsDistributionTest(sampleTestBits128)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.970152" || fmt.Sprintf("%.6f", q) != "0.970152" {
		t.FailNow()
	}
}
