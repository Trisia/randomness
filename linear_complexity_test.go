package randomness

import (
	"fmt"
	"testing"
)

func TestLinearComplexityTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := LinearComplexityTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
}
