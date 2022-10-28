package randomness

import (
	"fmt"
	"testing"
)

func TestPokerTestSample(t *testing.T) {
	bytes := []byte{0xcc, 0x15, 0x6c, 0x4c, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe6, 0xd8, 0xb2}
	p, q := PokerTestBytes(bytes, 4)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bytes)*8, p, q)
	if fmt.Sprintf("%.6f", p) != "0.213734" || fmt.Sprintf("%.6f", q) != "0.213734" {
		t.FailNow()
	}
}
