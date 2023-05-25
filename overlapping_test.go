package randomness

import (
	"fmt"
	"testing"
)

func TestOverlappingTemplateMatchingTestSample(t *testing.T) {
	p1, p2, q1, q2 := OverlappingTemplateMatchingProto(sampleTestBits128, 2)
	fmt.Printf("n: %v, P-value1: %.6f P-value2: %.6f, Q-value1: %.6f Q-value2: %.6f\n", len(sampleTestBits128), p1, p2, q1, q2)
	if fmt.Sprintf("%.6f", p1) != "0.436868" || fmt.Sprintf("%.6f", p2) != "0.723674" {
		t.FailNow()
	}
}
