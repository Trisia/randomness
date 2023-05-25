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
