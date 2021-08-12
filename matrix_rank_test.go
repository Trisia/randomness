package randomness

import (
	"fmt"
	"testing"
)

func TestMatrixRankTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := MatrixRankTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
