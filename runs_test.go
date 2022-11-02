package randomness

import (
	"fmt"
	"testing"
)

func TestRunsTestSample(t *testing.T) {
	p, q := RunsTest(sampleTestBits128)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.620729" || fmt.Sprintf("%.6f", q) != "0.310364" {
		t.FailNow()
	}
}
