package randomness

import (
	"fmt"
	"testing"
)

func TestFrequencyWithinBlockTest(t *testing.T) {
	bits := GroupBit()
	p := FrequencyWithinBlockTest(bits)
	fmt.Printf("n: 1000000, P-valye: %f\n", p)
}
