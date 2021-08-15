package randomness

import (
	"fmt"
	"testing"
)

func TestApproximateEntropyTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := ApproximateEntropyTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
