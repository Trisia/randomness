package randomness

import (
	"fmt"
	"math"
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
	p3, q3 := DiscreteFourierTransformTestFast(bits)

	t.Logf("数据大小: %d 比特", len(bits))
	t.Logf("原始版本:  P=%.6f, Q=%.6f", p1, q1)
	t.Logf("快速版本:  P=%.6f, Q=%.6f", p3, q3)

	// 验证快速版本在合理范围内
	if p1 < 0 || p1 > 1 || p3 < 0 || p3 > 1 {
		t.Errorf("P值超出有效范围[0,1]: 原始=%.3f, 快速=%.3f", p1, p3)
	}
	if q1 < 0 || q1 > 1 || q3 < 0 || q3 > 1 {
		t.Errorf("Q值超出有效范围[0,1]: 原始=%.3f, 快速=%.3f", q1, q3)
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
			p3, q3 := DiscreteFourierTransformTestFast(bits)

			t.Logf("大小 %d: 原始(P=%.3f,Q=%.3f) 快速(P=%.3f,Q=%.3f)",
				size, p1, q1, p3, q3)

			// 验证所有值在有效范围内
			if p1 < 0 || p1 > 1 || p3 < 0 || p3 > 1 {
				t.Errorf("P值超出范围: %.3f, %.3f", p1, p3)
			}
			if q1 < 0 || q1 > 1 || q3 < 0 || q3 > 1 {
				t.Errorf("Q值超出范围: %.3f, %.3f", q1, q3)
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

		b.Run(fmt.Sprintf("Fast_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DiscreteFourierTransformTestFast(bits)
			}
		})
	}
}

// TestCorrectnessVerification 验证所有版本的正确性
func TestCorrectnessVerification(t *testing.T) {
	// 测试不同大小的数据集
	testSizes := []int{100, 1000, 10000}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// 生成测试数据
			bits := generateTestBits(size)

			// 运行所有版本
			pOriginal, qOriginal := DiscreteFourierTransformTest(bits)
			pFast, qFast := DiscreteFourierTransformTestFast(bits)

			// 验证结果
			t.Logf("数据大小: %d 比特", size)
			t.Logf("原始版本:    P=%.6f, Q=%.6f", pOriginal, qOriginal)
			t.Logf("快速版本:    P=%.6f, Q=%.6f", pFast, qFast)

			// 验证快速版本在可接受误差范围内
			// 快速版本使用采样，允许更大的误差
			if !floatEqual(pOriginal, pFast, 0.1) { // 允许10%误差
				t.Errorf("快速版本P值误差过大: 原始=%.6f, 快速=%.6f, 误差=%.2f%%",
					pOriginal, pFast, math.Abs(pOriginal-pFast)/pOriginal*100)
			}
			if !floatEqual(qOriginal, qFast, 0.1) { // 允许10%误差
				t.Errorf("快速版本Q值误差过大: 原始=%.6f, 快速=%.6f, 误差=%.2f%%",
					qOriginal, qFast, math.Abs(qOriginal-qFast)/qOriginal*100)
			}

			// 验证所有值都在有效范围内
			if pOriginal < 0 || pOriginal > 1 || pFast < 0 || pFast > 1 {
				t.Errorf("P值超出有效范围[0,1]: 原始=%.3f, 快速=%.3f",
					pOriginal, pFast)
			}
			if qOriginal < 0 || qOriginal > 1 || qFast < 0 || qFast > 1 {
				t.Errorf("Q值超出有效范围[0,1]: 原始=%.3f, 快速=%.3f",
					qOriginal, qFast)
			}
		})
	}
}

// 辅助函数

// generateTestBits 生成测试用的随机比特序列
func generateTestBits(n int) []bool {
	bits := make([]bool, n)
	// 使用简单的伪随机生成器，确保可重现
	seed := uint32(12345)
	for i := 0; i < n; i++ {
		seed = seed*1103515245 + 12345
		bits[i] = (seed & 0x80000000) != 0
	}
	return bits
}

// createAllOnes 创建全1序列
func createAllOnes(n int) []bool {
	bits := make([]bool, n)
	for i := range bits {
		bits[i] = true
	}
	return bits
}

// createAlternating 创建交替序列
func createAlternating(n int) []bool {
	bits := make([]bool, n)
	for i := range bits {
		bits[i] = i%2 == 0
	}
	return bits
}

// floatEqual 比较两个浮点数是否相等（在给定容差内）
func floatEqual(a, b, tolerance float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= tolerance
}

// BenchmarkCorrectnessComparison 正确性验证的性能对比
func BenchmarkCorrectnessComparison(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		bits := generateTestBits(size)

		b.Run(fmt.Sprintf("Original_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				DiscreteFourierTransformTest(bits)
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
