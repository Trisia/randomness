package randomness

import (
	"fmt"
	"testing"
)

//1100110000010101011011000100110011100000000000100100110101010001
//0001001111010110100000001101011111001100111001101101100010110010
var sampleTestBits128 []bool

func init() {
	sampleTestBits128 = B2bitArr([]byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2})
}

func TestAutocorrelationTestSample(t *testing.T) {
	p, q := AutocorrelationTest(sampleTestBits128, 1)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.790080" || fmt.Sprintf("%.6f", q) != "0.395040" {
		t.FailNow()
	}
}
