package randomness

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var initonce sync.Once
var eConstantBits []bool

func getEConstantBits() []bool {
	initonce.Do(func() {
		eConstantBits = ReadGroupInASCIIFormat("data/data_e")
	})
	return eConstantBits
}

func TestLinearComplexityTestSample(t *testing.T) {
	bits := getEConstantBits()
	p, q := LinearComplexityProto(bits, 1000)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(bits), p, q)
	if fmt.Sprintf("%.6f", p) != "0.844721" || fmt.Sprintf("%.6f", q) != "0.844721" {
		t.FailNow()
	}
}

// 生成测试数据
func generateTestData(n int) []bool {
	data := make([]bool, n)
	for i := range data {
		data[i] = i%2 == 0 // 简单的交替模式
	}
	return data
}

// 基准测试 - 并行版本
func BenchmarkLinearComplexityParallel(b *testing.B) {
	data := generateTestData(100000) // 100k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProto(data, m)
	}
}

// 基准测试 - 单线程版本（用于对比）
func BenchmarkLinearComplexitySerial(b *testing.B) {
	data := generateTestData(100000) // 100k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProtoSerial(data, m)
	}
}

// 性能对比测试
func TestLinearComplexityPerformance(t *testing.T) {
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())

	// 测试不同大小的数据
	testSizes := []int{10000, 50000, 100000, 200000}
	m := 1000

	for _, size := range testSizes {
		data := generateTestData(size)

		// 测试并行版本
		start := time.Now()
		pParallel, _ := LinearComplexityProto(data, m)
		parallelTime := time.Since(start)

		// 测试串行版本
		start = time.Now()
		pSerial, _ := LinearComplexityProtoSerial(data, m)
		serialTime := time.Since(start)

		// 验证结果一致性
		if fmt.Sprintf("%.6f", pParallel) != fmt.Sprintf("%.6f", pSerial) {
			t.Errorf("结果不一致! 并行: %f, 串行: %f", pParallel, pSerial)
		}

		speedup := float64(serialTime) / float64(parallelTime)
		fmt.Printf("数据大小: %d bits, 串行: %v, 并行: %v, 加速比: %.2fx\n",
			size, serialTime, parallelTime, speedup)
	}
}

// 基准测试 - 自适应版本（主函数）
func BenchmarkLinearComplexityAdaptive(b *testing.B) {
	data := generateTestData(100000) // 100k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProto(data, m)
	}
}

// 基准测试 - 并行版本（直接调用并行函数）
func BenchmarkLinearComplexityParallelDirect(b *testing.B) {
	data := generateTestData(100000) // 100k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProtoParallel(data, m)
	}
}

// 基准测试 - 单线程版本（直接调用串行函数）
func BenchmarkLinearComplexitySerialDirect(b *testing.B) {
	data := generateTestData(100000) // 100k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProtoSerial(data, m)
	}
}

// 小数据集基准测试 - 应该使用串行
func BenchmarkLinearComplexitySmall(b *testing.B) {
	data := generateTestData(10000) // 10k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProto(data, m)
	}
}

// 大数据集基准测试 - 应该使用并行
func BenchmarkLinearComplexityLarge(b *testing.B) {
	data := generateTestData(200000) // 200k bits
	m := 1000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LinearComplexityProto(data, m)
	}
}

// 性能对比测试
func TestLinearComplexityAdaptivePerformance(t *testing.T) {
	fmt.Printf("CPU核心数: %d\n", runtime.NumCPU())

	// 测试不同大小的数据
	testSizes := []int{10000, 30000, 50000, 100000, 200000}
	m := 1000

	for _, size := range testSizes {
		data := generateTestData(size)
		n := len(data)
		N := n / m

		// 判断应该使用哪种策略
		expectedStrategy := "串行"
		if N >= 50 && n >= 50000 {
			expectedStrategy = "并行"
		}

		// 测试自适应版本
		start := time.Now()
		pAdaptive, _ := LinearComplexityProto(data, m)
		adaptiveTime := time.Since(start)

		// 测试串行版本
		start = time.Now()
		pSerial, _ := LinearComplexityProtoSerial(data, m)
		serialTime := time.Since(start)

		// 测试并行版本
		start = time.Now()
		pParallel, _ := LinearComplexityProtoParallel(data, m)
		parallelTime := time.Since(start)

		// 验证结果一致性
		if fmt.Sprintf("%.6f", pAdaptive) != fmt.Sprintf("%.6f", pSerial) ||
			fmt.Sprintf("%.6f", pAdaptive) != fmt.Sprintf("%.6f", pParallel) {
			t.Errorf("结果不一致! 自适应: %f, 串行: %f, 并行: %f", pAdaptive, pSerial, pParallel)
		}

		// 计算最佳时间
		bestTime := serialTime
		if parallelTime < serialTime {
			bestTime = parallelTime
		}

		efficiency := float64(bestTime) / float64(adaptiveTime)

		fmt.Printf("数据大小: %d bits (块数: %d), 预期策略: %s, 自适应: %v, 串行: %v, 并行: %v, 效率: %.2fx\n",
			size, N, expectedStrategy, adaptiveTime, serialTime, parallelTime, efficiency)
	}
}
