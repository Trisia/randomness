package randomness

import (
	"fmt"
	"testing"
)

func TestPokerTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := PokerTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
