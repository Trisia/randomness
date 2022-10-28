package randomness

import (
	"fmt"
	"testing"
)

func TestLinearComplexityTestSample(t *testing.T) {
	bits := ReadGroupInASCIIFormat("data/data.e")
	p, q := LinearComplexityProto(bits, 1000)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.844721" || fmt.Sprintf("%.6f", q) != "0.844721" {
		t.FailNow()
	}
}
