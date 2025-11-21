package randomness

import (
	"fmt"
	"testing"
)

// BenchmarkComparison 比较所有版本的基准测试
func BenchmarkComparison(b *testing.B) {
	// 测试不同大小的数据集
	sizes := []int{10000, 100000, 1000000, 10000000}

	for _, size := range sizes {
		bits := make([]bool, size)

		b.Run(fmt.Sprintf("Original_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = DiscreteFourierTransformTest(bits)
			}
		})

		b.Run(fmt.Sprintf("Optimized_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = DiscreteFourierTransformTestOptimized(bits)
			}
		})

		b.Run(fmt.Sprintf("Fast_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = DiscreteFourierTransformTestFast(bits)
			}
		})
	}
}
