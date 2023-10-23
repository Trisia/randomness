package randomness

import (
	"fmt"
	"testing"
)

func TestMonoBitFrequencyTestSample(t *testing.T) {
	p, q := MonoBitFrequencyTest(sampleTestBits128)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.215925" || fmt.Sprintf("%.6f", q) != "0.892038" {
		t.FailNow()
	}
}

func TestMonoBitFrequencyTestByteSample(t *testing.T) {
	p, q := MonoBitFrequencyTestBytes([]byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2})
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.215925" || fmt.Sprintf("%.6f", q) != "0.892038" {
		t.FailNow()
	}
}
