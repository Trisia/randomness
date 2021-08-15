package randomness

import (
	"fmt"
	"testing"
)

func TestDiscreteFourierTransformTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := DiscreteFourierTransformTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
