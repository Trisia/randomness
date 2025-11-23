package randomness

import (
	"fmt"
	"testing"
	"time"
)

var sampleTestBits100 []bool

func init() {
	sampleTestBits100 = B2bitArr([]byte{0xc9, 0xf, 0xda, 0xa2, 0x21, 0x68, 0xc2, 0x34, 0xc4, 0xc6, 0x62, 0x8b, 0x80})
	sampleTestBits100 = sampleTestBits100[:100]
}

func TestApproximateEntropyTestSample(t *testing.T) {
	p, q := ApproximateEntropyProto(sampleTestBits100, 2)
	fmt.Printf("n: %v, P-value: %f, Q-value: %f\n", len(sampleTestBits100), p, q)
	if fmt.Sprintf("%.6f", p) != "0.235301" || fmt.Sprintf("%.6f", q) != "0.235301" {
		t.FailNow()
	}
}

// 基准测试 - 原始实现
func BenchmarkApproximateEntropyOriginal(b *testing.B) {
	// 生成测试数据
	bits := GroupSecBit() // 1,000,000 bits

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproximateEntropyProto(bits, 5)
	}
}

// 基准测试 - 小数据集
func BenchmarkApproximateEntropySmall(b *testing.B) {
	bits := make([]bool, 10000)
	for i := range bits {
		bits[i] = i%2 == 0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproximateEntropyProto(bits, 5)
	}
}

// 基准测试 - 不同m值
func BenchmarkApproximateEntropyM2(b *testing.B) {
	bits := GroupSecBit()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproximateEntropyProto(bits, 2)
	}
}

func BenchmarkApproximateEntropyM5(b *testing.B) {
	bits := GroupSecBit()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproximateEntropyProto(bits, 5)
	}
}

func BenchmarkApproximateEntropyM10(b *testing.B) {
	bits := GroupSecBit()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ApproximateEntropyProto(bits, 10)
	}
}

// 性能对比测试
func TestApproximateEntropyPerformance(t *testing.T) {
	// 测试不同大小的数据集
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			// 生成测试数据 - 使用更随机的模式
			bits := make([]bool, size)
			for i := range bits {
				// 使用伪随机数生成器
				bits[i] = (i*17+13)%7 < 3
			}

			// 测试优化后的实现
			start := time.Now()
			p, q := ApproximateEntropyProto(bits, 5)
			duration := time.Since(start)

			t.Logf("数据大小: %d bits, 耗时: %v, P值: %f, Q值: %f",
				size, duration, p, q)

			// 验证结果合理性
			if p < 0 || p > 1 || q < 0 || q > 1 {
				t.Errorf("无效的P值或Q值: P=%f, Q=%f", p, q)
			}
		})
	}
}

// 内存使用测试
func TestApproximateEntropyMemory(t *testing.T) {
	bits := GroupSecBit()

	// 测试大m值的内存使用
	for m := 5; m <= 15; m++ {
		t.Run(fmt.Sprintf("M_%d", m), func(t *testing.T) {
			start := time.Now()
			p, q := ApproximateEntropyProto(bits, m)
			duration := time.Since(start)

			t.Logf("m=%d, 耗时: %v, P值: %f, Q值: %f", m, duration, p, q)

			if p < 0 || p > 1 || q < 0 || q > 1 {
				t.Errorf("无效的P值或Q值: P=%f, Q=%f", p, q)
			}
		})
	}
}
