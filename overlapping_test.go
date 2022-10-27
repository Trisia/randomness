package randomness

import (
	"fmt"
	"testing"
)

func TestOverlappingTemplateMatchingTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p1, p2, q1, q2 := OverlappingTemplateMatchingTest(bits)
	fmt.Printf("n: %v, P-value1: %.6f P-value2: %.6f, Q-value1: %.6f Q-value2: %.6f\n", len(bits), p1, p2, q1, q2)
}

func TestOverlappingTemplateMatchingTestSample(t *testing.T) {
	bytes := []byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2}
	p1, p2, q1, q2 := OverlappingTemplateMatchingTestBytes(bytes, 2)
	fmt.Printf("n: %v, P-value1: %.6f P-value2: %.6f, Q-value1: %.6f Q-value2: %.6f\n", len(bytes)*8, p1, p2, q1, q2)
	if fmt.Sprintf("%.6f", p1) != "0.436868" || fmt.Sprintf("%.6f", p2) != "0.723674" {
		t.FailNow()
	}
}
