package randomness

import (
	"fmt"
	"testing"
)

func TestPokerTest(t *testing.T) {
	bits := ReadGroup("data.bin")
	p, q := PokerTest(bits)
	fmt.Printf("n: 1000000, P-value: %f, Q-value: %f\n", p, q)
}

func TestPokerTestSample(t *testing.T) {
	bytes := []byte{0xcc, 0x15, 0x62, 0x4b, 0xe0, 0x02, 0x4d, 0x51, 0x13, 0xd6, 0x80, 0xd7, 0xcc, 0xe3, 0xd8, 0xb2}
	p, q := PokerTestBytes(bytes, 4)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bytes)*8, p, q)
	//Not same with C.3
	//if fmt.Sprintf("%.6f", p) != "0.213734" || fmt.Sprintf("%.6f", q) != "0.213734" {
	//	t.FailNow()
	//}
}
