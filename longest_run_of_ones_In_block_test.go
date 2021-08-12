package randomness

import (
	"fmt"
	"testing"
)

func TestLongestRunOfOnesInABlockTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := LongestRunOfOnesInABlockTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
