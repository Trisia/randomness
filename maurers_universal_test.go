package randomness

import (
	"fmt"
	"testing"
)

func TestMaurerUniversalTestSample(t *testing.T) {
	bits := ReadGroupInASCIIFormat("data/data.e")
	p, q := MaurerUniversalTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.282568" || fmt.Sprintf("%.6f", q) != "0.282568" {
		t.FailNow()
	}
}
