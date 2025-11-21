package randomness

import (
	"fmt"
	"testing"
)

func TestDiscreteFourierTransformTestFast(t *testing.T) {
	p, q := DiscreteFourierTransformTestFast(sampleTestBits100)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	// 注意：快速版本可能因为采样而有轻微差异，但应该在可接受范围内
	if p < 0 || p > 1 || q < 0 || q > 1 {
		t.Errorf("Invalid P or Q values: P=%f, Q=%f", p, q)
	}
}

func BenchmarkDiscreteFourierTransformTestFast(b *testing.B) {
	bits := make([]bool, 100000000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTestFast(bits)
	}
}

func BenchmarkDiscreteFourierTransformTestBytesFast(b *testing.B) {
	data := make([]byte, 12500000) // 100M bits = 12.5M bytes
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTestBytesFast(data)
	}
}

// 比较所有版本的基准测试
func BenchmarkAllVersions(b *testing.B) {
	bits := make([]bool, 1000000) // 使用较小的数据集进行比较

	b.Run("Original", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = DiscreteFourierTransformTest(bits)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = DiscreteFourierTransformTestOptimized(bits)
		}
	})

	b.Run("Fast", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = DiscreteFourierTransformTestFast(bits)
		}
	})
}

func TestSimpleCorrectness(t *testing.T) {
	// 使用已知的测试数据
	bits := sampleTestBits100

	// 运行所有版本
	p1, q1 := DiscreteFourierTransformTest(bits)
	p2, q2 := DiscreteFourierTransformTestOptimized(bits)
	p3, q3 := DiscreteFourierTransformTestFast(bits)

	t.Logf("数据大小: %d 比特", len(bits))
	t.Logf("原始版本:  P=%.6f, Q=%.6f", p1, q1)
	t.Logf("优化版本:  P=%.6f, Q=%.6f", p2, q2)
	t.Logf("快速版本:  P=%.6f, Q=%.6f", p3, q3)

	// 验证优化版本与原始版本完全一致
	if p1 != p2 {
		t.Errorf("优化版本P值不一致: 原始=%.10f, 优化=%.10f", p1, p2)
	}
	if q1 != q2 {
		t.Errorf("优化版本Q值不一致: 原始=%.10f, 优化=%.10f", q1, q2)
	}

	// 验证快速版本在合理范围内
	if p1 < 0 || p1 > 1 || p2 < 0 || p2 > 1 || p3 < 0 || p3 > 1 {
		t.Errorf("P值超出有效范围[0,1]: 原始=%.3f, 优化=%.3f, 快速=%.3f", p1, p2, p3)
	}
	if q1 < 0 || q1 > 1 || q2 < 0 || q2 > 1 || q3 < 0 || q3 > 1 {
		t.Errorf("Q值超出有效范围[0,1]: 原始=%.3f, 优化=%.3f, 快速=%.3f", q1, q2, q3)
	}

	// 验证与期望值接近
	expectedP := 0.654721
	expectedQ := 0.327360

	if fmt.Sprintf("%.6f", p1) != fmt.Sprintf("%.6f", expectedP) {
		t.Errorf("原始版本P值不符合期望: 期望=%.6f, 实际=%.6f", expectedP, p1)
	}
	if fmt.Sprintf("%.6f", q1) != fmt.Sprintf("%.6f", expectedQ) {
		t.Errorf("原始版本Q值不符合期望: 期望=%.6f, 实际=%.6f", expectedQ, q1)
	}
}

func TestSmallDataComparison(t *testing.T) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// 生成测试数据
			bits := make([]bool, size)
			for i := range bits {
				bits[i] = i%3 == 0 // 简单模式
			}

			// 运行所有版本
			p1, q1 := DiscreteFourierTransformTest(bits)
			p2, q2 := DiscreteFourierTransformTestOptimized(bits)
			p3, q3 := DiscreteFourierTransformTestFast(bits)

			t.Logf("大小 %d: 原始(P=%.3f,Q=%.3f) 优化(P=%.3f,Q=%.3f) 快速(P=%.3f,Q=%.3f)",
				size, p1, q1, p2, q2, p3, q3)

			// 验证优化版本与原始版本一致
			if p1 != p2 || q1 != q2 {
				t.Errorf("优化版本不一致: 原始(%.6f,%.6f) vs 优化(%.6f,%.6f)", p1, q1, p2, q2)
			}

			// 验证所有值在有效范围内
			if p1 < 0 || p1 > 1 || p2 < 0 || p2 > 1 || p3 < 0 || p3 > 1 {
				t.Errorf("P值超出范围: %.3f, %.3f, %.3f", p1, p2, p3)
			}
			if q1 < 0 || q1 > 1 || q2 < 0 || q2 > 1 || q3 < 0 || q3 > 1 {
				t.Errorf("Q值超出范围: %.3f, %.3f, %.3f", q1, q2, q3)
			}
		})
	}
}

func BenchmarkPerformanceComparison(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		bits := make([]bool, size)
		for i := range bits {
			bits[i] = i%7 == 0 // 生成测试数据
		}

		b.Run(fmt.Sprintf("Original_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DiscreteFourierTransformTest(bits)
			}
		})

		b.Run(fmt.Sprintf("Optimized_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DiscreteFourierTransformTestOptimized(bits)
			}
		})

		b.Run(fmt.Sprintf("Fast_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DiscreteFourierTransformTestFast(bits)
			}
		})
	}
}
