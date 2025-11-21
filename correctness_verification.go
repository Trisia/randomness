package randomness

import (
	"fmt"
	"math"
	"testing"
)

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
