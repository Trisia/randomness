package randomness

import (
	"fmt"
	"testing"
)

func TestOverlappingTemplateMatchingTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p1, p2 := OverlappingTemplateMatchingTest(bits)
	fmt.Printf("n: 1000000, P-value1: %.6f P-value2: %.6f\n", p1, p2)
}
