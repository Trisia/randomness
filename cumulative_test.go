package randomness

import (
	"fmt"
	"testing"
)

func TestCumulativeTestSample(t *testing.T) {
	bits := B2bitArr([]byte{0xc9, 0xf, 0xda, 0xa2, 0x21, 0x68, 0xc2, 0x34, 0xc4, 0xc6, 0x62, 0x8b, 0x80})
	bits = bits[:100]
	p, q := CumulativeTest(bits, true)
	fmt.Printf("Foward, n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.219194" || fmt.Sprintf("%.6f", q) != "0.219194" {
		t.FailNow()
	}
	p, q = CumulativeTest(bits, false)
	fmt.Printf("Backward, n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.114866" || fmt.Sprintf("%.6f", q) != "0.114866" {
		t.FailNow()
	}
}
