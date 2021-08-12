package randomness

import (
	"fmt"
	"testing"
)

func TestAutocorrelationTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p := AutocorrelationTest(bits)
	fmt.Printf("n: 1000000, P-value: %f\n", p)
}
