package randomness

import (
	"fmt"
	"testing"
)

func TestCumulativeTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := CumulativeTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
