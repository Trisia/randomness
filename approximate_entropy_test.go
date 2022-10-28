package randomness

import (
	"fmt"
	"testing"
)

func TestApproximateEntropyTestSample(t *testing.T) {
	bits := B2bitArr([]byte{0xc9, 0xf, 0xda, 0xa2, 0x21, 0x68, 0xc2, 0x34, 0xc4, 0xc6, 0x62, 0x8b, 0x80})
	bits = bits[:100]
	p, q := ApproximateEntropyProto(bits, 2)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.235301" || fmt.Sprintf("%.6f", q) != "0.235301" {
		t.FailNow()
	}
}
