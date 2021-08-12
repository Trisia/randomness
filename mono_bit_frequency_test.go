package randomness

import (
	"fmt"
	"testing"
)

func TestMonoBitFrequencyTest(t *testing.T) {
	bits := GroupBit()
	//bits := ReadGroup("data.bin")
	p := MonoBitFrequencyTest(bits)
	fmt.Printf("n: 1000000, P-value: %.6f\n", p)
}
