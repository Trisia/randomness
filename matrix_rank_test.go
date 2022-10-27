package randomness

import (
	"fmt"
	"testing"
)

func TestMatrixRankTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := MatrixRankTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
}
