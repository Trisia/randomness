package randomness

import (
	"fmt"
	"testing"
)

func TestLinearComplexityTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := LinearComplexityTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
