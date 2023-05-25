package randomness

import (
	"fmt"
	"testing"
)

func TestDiscreteFourierTransformTestSample(t *testing.T) {
	p, q := DiscreteFourierTransformTest(sampleTestBits100)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.654721" || fmt.Sprintf("%.6f", q) != "0.327360" {
		t.FailNow()
	}
}

func BenchmarkDiscreteFourierTransformTest(b *testing.B) {
	bits := make([]bool, 100000000)
	b.ReportAllocs()
	b.ResetTimer()
	_, _ = DiscreteFourierTransformTest(bits)
}
