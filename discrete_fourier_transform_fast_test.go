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
			pOptimized, qOptimized := DiscreteFourierTransformTestOptimized(bits)
			pFast, qFast := DiscreteFourierTransformTestFast(bits)

			// 验证结果
			t.Logf("数据大小: %d 比特", size)
			t.Logf("原始版本:    P=%.6f, Q=%.6f", pOriginal, qOriginal)
			t.Logf("优化版本:    P=%.6f, Q=%.6f", pOptimized, qOptimized)
			t.Logf("快速版本:    P=%.6f, Q=%.6f", pFast, qFast)

			// 验证优化版本与原始版本的一致性
			if !floatEqual(pOriginal, pOptimized, 1e-10) {
				t.Errorf("优化版本P值不一致: 原始=%.10f, 优化=%.10f", pOriginal, pOptimized)
			}
			if !floatEqual(qOriginal, qOptimized, 1e-10) {
				t.Errorf("优化版本Q值不一致: 原始=%.10f, 优化=%.10f", qOriginal, qOptimized)
			}

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
			if pOriginal < 0 || pOriginal > 1 || pOptimized < 0 || pOptimized > 1 || pFast < 0 || pFast > 1 {
				t.Errorf("P值超出有效范围[0,1]: 原始=%.3f, 优化=%.3f, 快速=%.3f",
					pOriginal, pOptimized, pFast)
			}
			if qOriginal < 0 || qOriginal > 1 || qOptimized < 0 || qOptimized > 1 || qFast < 0 || qFast > 1 {
				t.Errorf("Q值超出有效范围[0,1]: 原始=%.3f, 优化=%.3f, 快速=%.3f",
					qOriginal, qOptimized, qFast)
			}
		})
	}
}

// TestCorrectnessWithKnownData 使用已知数据验证正确性
func TestCorrectnessWithKnownData(t *testing.T) {
	// 使用已知的测试数据
	testCases := []struct {
		name      string
		bits      []bool
		expectedP float64
		expectedQ float64
		tolerance float64
	}{
		{
			name:      "全零序列",
			bits:      make([]bool, 100),
			expectedP: 0.0, // 全零序列应该完全失败
			expectedQ: 0.0,
			tolerance: 1e-10,
		},
		{
			name:      "全一序列",
			bits:      createAllOnes(100),
			expectedP: 0.0, // 全一序列也应该完全失败
			expectedQ: 0.0,
			tolerance: 1e-10,
		},
		{
			name:      "交替序列",
			bits:      createAlternating(100),
			expectedP: 1.0, // 完美的交替序列应该通过
			expectedQ: 0.5,
			tolerance: 0.1, // 允许一定误差
		},
		{
			name:      "样本数据100",
			bits:      sampleTestBits100,
			expectedP: 0.654721,
			expectedQ: 0.327360,
			tolerance: 1e-6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试原始版本
			pOrig, qOrig := DiscreteFourierTransformTest(tc.bits)

			// 测试优化版本
			pOpt, qOpt := DiscreteFourierTransformTestOptimized(tc.bits)

			// 测试快速版本
			pFast, qFast := DiscreteFourierTransformTestFast(tc.bits)

			t.Logf("测试用例: %s", tc.name)
			t.Logf("期望值:     P=%.6f, Q=%.6f", tc.expectedP, tc.expectedQ)
			t.Logf("原始版本:   P=%.6f, Q=%.6f", pOrig, qOrig)
			t.Logf("优化版本:   P=%.6f, Q=%.6f", pOpt, qOpt)
			t.Logf("快速版本:   P=%.6f, Q=%.6f", pFast, qFast)

			// 验证原始版本
			if !floatEqual(pOrig, tc.expectedP, tc.tolerance) {
				t.Errorf("原始版本P值不符合期望: 期望=%.6f, 实际=%.6f", tc.expectedP, pOrig)
			}
			if !floatEqual(qOrig, tc.expectedQ, tc.tolerance) {
				t.Errorf("原始版本Q值不符合期望: 期望=%.6f, 实际=%.6f", tc.expectedQ, qOrig)
			}

			// 验证优化版本与原始版本一致
			if !floatEqual(pOrig, pOpt, 1e-10) {
				t.Errorf("优化版本P值与原始版本不一致: 原始=%.10f, 优化=%.10f", pOrig, pOpt)
			}
			if !floatEqual(qOrig, qOpt, 1e-10) {
				t.Errorf("优化版本Q值与原始版本不一致: 原始=%.10f, 优化=%.10f", qOrig, qOpt)
			}
		})
	}
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		bits        []bool
		shouldPanic bool
	}{
		{
			name:        "空数据",
			bits:        []bool{},
			shouldPanic: true,
		},
		{
			name:        "单个比特",
			bits:        []bool{true},
			shouldPanic: false,
		},
		{
			name:        "两个比特",
			bits:        []bool{true, false},
			shouldPanic: false,
		},
		{
			name:        "最小有效长度",
			bits:        make([]bool, 8),
			shouldPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("期望panic但没有panic")
					}
				}()
				DiscreteFourierTransformTest(tc.bits)
			} else {
				// 所有版本都应该能正常处理
				p1, q1 := DiscreteFourierTransformTest(tc.bits)
				p2, q2 := DiscreteFourierTransformTestOptimized(tc.bits)
				p3, q3 := DiscreteFourierTransformTestFast(tc.bits)

				// 验证结果在有效范围内
				if p1 < 0 || p1 > 1 || p2 < 0 || p2 > 1 || p3 < 0 || p3 > 1 {
					t.Errorf("P值超出有效范围: %.3f, %.3f, %.3f", p1, p2, p3)
				}
				if q1 < 0 || q1 > 1 || q2 < 0 || q2 > 1 || q3 < 0 || q3 > 1 {
					t.Errorf("Q值超出有效范围: %.3f, %.3f, %.3f", q1, q2, q3)
				}

				t.Logf("边界测试 %s: P=%.3f,%.3f,%.3f Q=%.3f,%.3f,%.3f",
					tc.name, p1, p2, p3, q1, q2, q3)
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
