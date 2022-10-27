package randomness

import (
	"fmt"
	"testing"
)

func TestMaurerUniversalTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := MaurerUniversalTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
}
