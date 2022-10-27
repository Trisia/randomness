package randomness

import (
	"fmt"
	"testing"
)

func TestLongestRunOfOnesInABlockTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := LongestRunOfOnesInABlockTest(bits, true)
	fmt.Printf("n: 1000000, P-value: %f, Q-value: %f\n", p, q)
}

func TestLongestRunOfOnesInABlockTestSample(t *testing.T) {
	bits := B2bitArr([]byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2})
	// check "1"
	p, q := LongestRunOfOnesInABlockTest(bits, true)
	fmt.Printf("Check 1, n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.180598" || fmt.Sprintf("%.6f", q) != "0.180598" {
		t.FailNow()
	}
	// check "0"
	p, q = LongestRunOfOnesInABlockTest(bits, false)
	fmt.Printf("Check 0, n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.839299" || fmt.Sprintf("%.6f", q) != "0.839299" {
		t.FailNow()
	}
}
