package randomness

import (
	"fmt"
	"github.com/Trisia/randomness/fft"
	"math"
	"math/cmplx"
	"math/rand"
	"testing"
	"time"
)

func TestDiscreteFourierTransformTestSample(t *testing.T) {
	p, q := DiscreteFourierTransformTest(sampleTestBits100)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.654721" || fmt.Sprintf("%.6f", q) != "0.327360" {
		t.FailNow()
	}
}

// TestGMT0005Scales 测试GMT 0005-2021规范中定义的数据规模
func TestGMT0005Scales(t *testing.T) {
	// 测试小规模：2*10^4 bit
	t.Run("SmallScale_20K", func(t *testing.T) {
		bits := make([]bool, SmallScale)
		rand.Seed(time.Now().UnixNano())
		for i := range bits {
			bits[i] = rand.Intn(2) == 1
		}
		p, q := DiscreteFourierTransformTest(bits)
		t.Logf("SmallScale (20K bits): P=%.6f, Q=%.6f", p, q)
		if p < 0 || p > 1 || q < 0 || q > 1 {
			t.Errorf("Invalid P or Q values: P=%f, Q=%f", p, q)
		}
	})

	// 测试中规模：10^6 bit
	t.Run("MediumScale_1M", func(t *testing.T) {
		bits := make([]bool, MediumScale)
		rand.Seed(time.Now().UnixNano())
		for i := range bits {
			bits[i] = rand.Intn(2) == 1
		}
		p, q := DiscreteFourierTransformTest(bits)
		t.Logf("MediumScale (1M bits): P=%.6f, Q=%.6f", p, q)
		if p < 0 || p > 1 || q < 0 || q > 1 {
			t.Errorf("Invalid P or Q values: P=%f, Q=%f", p, q)
		}
	})

	// 测试大规模：10^8 bit (仅验证不崩溃，不实际运行以节省时间)
	t.Run("LargeScale_100M_Validation", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping large scale test in short mode")
		}
		bits := make([]bool, LargeScale)
		rand.Seed(time.Now().UnixNano())
		for i := range bits {
			bits[i] = rand.Intn(2) == 1
		}
		p, q := DiscreteFourierTransformTest(bits)
		t.Logf("LargeScale (100M bits): P=%.6f, Q=%.6f", p, q)
		if p < 0 || p > 1 || q < 0 || q > 1 {
			t.Errorf("Invalid P or Q values: P=%f, Q=%f", p, q)
		}
	})
}

func BenchmarkDiscreteFourierTransformTest(b *testing.B) {
	bits := make([]bool, 100000000)
	b.ReportAllocs()
	b.ResetTimer()
	_, _ = DiscreteFourierTransformTest(bits)
}
func BenchmarkDiscreteFourierTransformOriginal(b *testing.B) {
	bits := make([]bool, 10000000) // 10M bits
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = discreteFourierTransformTest(bits)
	}
}

func BenchmarkDiscreteFourierTransformOptimized(b *testing.B) {
	bits := make([]bool, 10000000) // 10M bits
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTest(bits)
	}
}

func BenchmarkDiscreteFourierTransformLarge(b *testing.B) {
	bits := make([]bool, 100000000) // 100M bits
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DiscreteFourierTransformTest(bits)
	}
}

// 正确性对比测试

// 原始实现（用于对比验证）
func discreteFourierTransformTestOriginal(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// Step 1, 2
	N := ceilPow2(n)
	rr := make([]complex128, N)
	for i := 0; i < n; i++ {
		if bits[i] {
			rr[i] = complex(1.0, 0)
		} else {
			rr[i] = complex(-1.0, 0)
		}
	}

	// 傅里叶变换
	f, err := fft.New(N)
	if err != nil {
		panic(err)
	}
	f.Transform(rr)

	// Step 4
	T := math.Sqrt(2.995732274 * float64(n))

	// Step 5
	N_0 := 0.95 * float64(n) / 2

	// Step 6
	var N_1 int = 0
	for i := 0; i < n/2-1; i++ {
		if cmplx.Abs(rr[i]) < T {
			N_1++
		}
	}

	// Step 7
	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(2.0*n)/3.8)
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}

// 生成测试数据
func generateTestDataForDFT(size int) []bool {
	rand.Seed(time.Now().UnixNano())
	bits := make([]bool, size)
	for i := range bits {
		bits[i] = rand.Intn(2) == 1
	}
	return bits
}

// 对比不同规模数据的计算结果
func TestDiscreteFourierTransformCorrectness(t *testing.T) {
	testSizes := []int{
		100,      // 小数据
		1000,     // 中小数据
		10000,    // 中数据
		100000,   // 大数据
		1000000,  // 超大数据
		10000000, // 10M数据
	}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			bits := generateTestDataForDFT(size)

			// 原始算法
			pOrig, qOrig := discreteFourierTransformTestOriginal(bits)

			// 优化算法
			pOpt, qOpt := DiscreteFourierTransformTest(bits)

			// 比较结果（允许浮点数误差）
			const epsilon = 1e-10
			if math.Abs(pOrig-pOpt) > epsilon {
				t.Errorf("P值差异过大: 原始=%.15f, 优化=%.15f, 差异=%.15f", pOrig, pOpt, math.Abs(pOrig-pOpt))
			}
			if math.Abs(qOrig-qOpt) > epsilon {
				t.Errorf("Q值差异过大: 原始=%.15f, 优化=%.15f, 差异=%.15f", qOrig, qOpt, math.Abs(qOrig-qOpt))
			}

			t.Logf("数据规模: %d bits", size)
			t.Logf("原始算法: P=%.15f, Q=%.15f", pOrig, qOrig)
			t.Logf("优化算法: P=%.15f, Q=%.15f", pOpt, qOpt)
			t.Logf("P值差异: %.15f, Q值差异: %.15f", math.Abs(pOrig-pOpt), math.Abs(qOrig-qOpt))
		})
	}
}

// 测试边界情况
func TestDiscreteFourierTransformEdgeCases(t *testing.T) {
	testCases := []struct {
		name string
		bits []bool
	}{
		{
			name: "全0序列",
			bits: make([]bool, 1000),
		},
		{
			name: "全1序列",
			bits: func() []bool {
				b := make([]bool, 1000)
				for i := range b {
					b[i] = true
				}
				return b
			}(),
		},
		{
			name: "交替序列",
			bits: func() []bool {
				b := make([]bool, 1000)
				for i := range b {
					b[i] = i%2 == 0
				}
				return b
			}(),
		},
		{
			name: "最小有效序列",
			bits: []bool{true, false, true, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 原始算法
			pOrig, qOrig := discreteFourierTransformTestOriginal(tc.bits)

			// 优化算法
			pOpt, qOpt := DiscreteFourierTransformTest(tc.bits)

			// 比较结果
			const epsilon = 1e-10
			if math.Abs(pOrig-pOpt) > epsilon {
				t.Errorf("P值差异过大: 原始=%.15f, 优化=%.15f", pOrig, pOpt)
			}
			if math.Abs(qOrig-qOpt) > epsilon {
				t.Errorf("Q值差异过大: 原始=%.15f, 优化=%.15f", qOrig, qOpt)
			}

			t.Logf("%s - P: %.15f, Q: %.15f", tc.name, pOpt, qOpt)
		})
	}
}

// 性能对比测试
func TestDiscreteFourierTransformPerformanceComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能对比测试")
	}

	sizes := []int{1000000, 10000000} // 1M, 10M bits

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Performance_%d", size), func(t *testing.T) {
			bits := generateTestDataForDFT(size)

			// 测试原始算法性能
			start := time.Now()
			pOrig, _ := discreteFourierTransformTestOriginal(bits)
			origDuration := time.Since(start)

			// 测试优化算法性能
			start = time.Now()
			pOpt, _ := DiscreteFourierTransformTest(bits)
			optDuration := time.Since(start)

			// 验证结果一致性
			const epsilon = 1e-10
			if math.Abs(pOrig-pOpt) > epsilon {
				t.Errorf("结果不一致: P差异=%.15f", math.Abs(pOrig-pOpt))
			}

			speedup := float64(origDuration) / float64(optDuration)
			t.Logf("数据规模: %d bits", size)
			t.Logf("原始算法耗时: %v", origDuration)
			t.Logf("优化算法耗时: %v", optDuration)
			t.Logf("性能提升: %.2fx", speedup)
		})
	}
}
