package randomness

import (
	"fmt"
	"testing"
)

func TestDiscreteFourierTransformTestOptimized(t *testing.T) {
	p, q := DiscreteFourierTransformTestOptimized(sampleTestBits100)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.654721" || fmt.Sprintf("%.6f", q) != "0.327360" {
		t.FailNow()
	}
}

func BenchmarkDiscreteFourierTransformTestOriginal(b *testing.B) {
	bits := make([]bool, 100000000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTest(bits)
	}
}

func BenchmarkDiscreteFourierTransformTestOptimized(b *testing.B) {
	bits := make([]bool, 100000000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTestOptimized(bits)
	}
}

func BenchmarkDiscreteFourierTransformTestBytesOriginal(b *testing.B) {
	data := make([]byte, 12500000) // 100M bits = 12.5M bytes
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTestBytes(data)
	}
}

func BenchmarkDiscreteFourierTransformTestBytesOptimized(b *testing.B) {
	data := make([]byte, 12500000) // 100M bits = 12.5M bytes
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTestBytesOptimized(data)
	}
}
