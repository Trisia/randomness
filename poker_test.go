package randomness

import (
	"fmt"
	"testing"
)

func TestPokerTestSample(t *testing.T) {
	p, q := PokerProto(sampleTestBits128, 4)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.213734" || fmt.Sprintf("%.6f", q) != "0.213734" {
		t.FailNow()
	}
}

func TestPokerTestM8Sample(t *testing.T) {
	p, q := PokerProto(sampleTestBits128, 8)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.221829" || fmt.Sprintf("%.6f", q) != "0.221829" {
		t.FailNow()
	}
}

func TestPokerTestByteSample(t *testing.T) {
	p, q := PokerTestBytes([]byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2}, 4)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.213734" || fmt.Sprintf("%.6f", q) != "0.213734" {
		t.FailNow()
	}
}

func TestPokerTestByteM8Sample(t *testing.T) {
	p, q := PokerTestBytes([]byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2}, 8)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits128), p, q)
	if fmt.Sprintf("%.6f", p) != "0.221829" || fmt.Sprintf("%.6f", q) != "0.221829" {
		t.FailNow()
	}
}
