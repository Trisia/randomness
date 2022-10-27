package randomness

import (
	"fmt"
	"testing"
)

func TestFrequencyWithinBlockTest(t *testing.T) {
	//bits := GroupBit()
	bits := ReadGroup("data.bin")
	p, q := FrequencyWithinBlockTest(bits)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
}

func TestFrequencyWithinBlockTestSample(t *testing.T) {
	bits := B2bitArr([]byte{0xc9, 0xf, 0xda, 0xa2, 0x21, 0x68, 0xc2, 0x34, 0xc4, 0xc6, 0x62, 0x8b, 0x80})
	bits = bits[:100]
	p, q := FrequencyWithinBlockTest(bits)
	fmt.Printf("n: %v, P-value: %.6f, Q-value: %.6f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.706438" || fmt.Sprintf("%.6f", q) != "0.706438" {
		t.FailNow()
	}
}
