package randomness

import (
	"fmt"
	"testing"
)

func TestBinaryDerivativeTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := BinaryDerivativeTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
