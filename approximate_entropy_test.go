package randomness

import (
	"fmt"
	"testing"
)

var sampleTestBits100 []bool

func init() {
	sampleTestBits100 = B2bitArr([]byte{0xc9, 0xf, 0xda, 0xa2, 0x21, 0x68, 0xc2, 0x34, 0xc4, 0xc6, 0x62, 0x8b, 0x80})
	sampleTestBits100 = sampleTestBits100[:100]
}

func TestApproximateEntropyTestSample(t *testing.T) {
	p, q := ApproximateEntropyProto(sampleTestBits100, 2)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.235301" || fmt.Sprintf("%.6f", q) != "0.235301" {
		t.FailNow()
	}
}
