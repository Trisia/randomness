package randomness

import (
	"fmt"
	"testing"
)

func TestMatrixRankTestSample(t *testing.T) {
	bits := getEConstantBits()
	p, q := MatrixRankTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.307543" || fmt.Sprintf("%.6f", q) != "0.307543" {
		t.FailNow()
	}
}
