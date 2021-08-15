package randomness

import (
	"fmt"
	"testing"
)

func TestMaurerUniversalTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := MaurerUniversalTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
