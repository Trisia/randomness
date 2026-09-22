package randomness

import (
	"testing"
)

// 性能基准门禁：统一在 10^6 bit 规模上测量 15 项检测与完整一轮。
//
// 输入来自 NewDetRand，因此数据可复现；建议在 CI 中记录
// ns/op / B/op / allocs/op 作为回退门禁。
//
//	go test -run XXX -bench . -benchtime 10x .
var benchData = NewDetRand(1).RawBytes(1000000 / 8)

func benchRunner(b *testing.B, f func([]byte) *TestResult) {
	b.Helper()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = f(benchData)
	}
}

func BenchmarkMonoBitFrequency(b *testing.B)     { benchRunner(b, MonoBitFrequency) }
func BenchmarkFrequencyWithinBlock(b *testing.B) { benchRunner(b, FrequencyWithinBlock) }
func BenchmarkPoker(b *testing.B)                { benchRunner(b, Poker) }
func BenchmarkOverlapping(b *testing.B)          { benchRunner(b, OverlappingTemplateMatching) }
func BenchmarkRuns(b *testing.B)                 { benchRunner(b, Runs) }
func BenchmarkRunsDistribution(b *testing.B)     { benchRunner(b, RunsDistribution) }
func BenchmarkLongestRun(b *testing.B)           { benchRunner(b, LongestRunOfOnesInABlock) }
func BenchmarkBinaryDerivative(b *testing.B)     { benchRunner(b, BinaryDerivative) }
func BenchmarkAutocorrelation(b *testing.B)      { benchRunner(b, Autocorrelation) }
func BenchmarkMatrixRank(b *testing.B)           { benchRunner(b, MatrixRank) }
func BenchmarkCumulative(b *testing.B)           { benchRunner(b, Cumulative) }
func BenchmarkApproximateEntropy(b *testing.B)   { benchRunner(b, ApproximateEntropy) }
func BenchmarkLinearComplexity(b *testing.B)     { benchRunner(b, LinearComplexity) }
func BenchmarkMaurerUniversal(b *testing.B)      { benchRunner(b, MaurerUniversal) }
func BenchmarkDiscreteFourierTransform(b *testing.B) {
	benchRunner(b, DiscreteFourierTransform)
}

// BenchmarkTestRound15 代表 rddetector 的真实工作负载。
func BenchmarkTestRound15(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, item := range TestMethodArr {
			_ = item.Runner(benchData)
		}
	}
}

// BenchmarkB2bitArr 保留旧转换路径的基线，用于说明 BitSeq 转换的收益。
func BenchmarkB2bitArr(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = B2bitArr(benchData)
	}
}

// BenchmarkTestRound15Packed 展示「整轮只转换一次」的附加形态收益。
func BenchmarkTestRound15Packed(b *testing.B) {
	s := BitSeqFromBytes(benchData)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MonoBitFrequency(benchData)
		_, _ = coreFrequencyWithinBlock(s, selectM(s.n))
		_ = Poker(benchData)
		_, _, _, _ = coreOverlapping(s, 5)
		_, _ = coreRuns(s)
		_, _ = coreRunsDistribution(s)
		_, _ = coreLongestRunOfOnesInABlock(s, true)
		_, _ = coreBinaryDerivative(s, 7)
		_, _ = coreAutocorrelation(s, 16)
		_, _ = coreMatrixRank(s, 32, 32)
		_, _ = coreCumulative(s, true)
		_, _ = coreApproximateEntropy(s, 5)
		_, _ = coreLinearComplexity(s, 500, 16)
		_, _ = coreMaurerUniversal(s)
		_, _ = coreDiscreteFourierTransform(s)
	}
}

// []bool 入口的成本对比：用于判定「哪些方法值得走 BitSeq」。
// BitSeqFromBools 的转换本身约 0.354ms/10^6 bit，若某个检测在 []bool 上
// 原本就远快于这个数，强制统一就是净损失。
var benchBits = B2bitArr(benchData)

func BenchmarkBoolEntryMonoBit(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MonoBitFrequencyTest(benchBits)
	}
}

func BenchmarkBoolEntryPokerM8(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PokerProto(benchBits, 8)
	}
}

func BenchmarkBoolEntryPokerM4(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PokerProto(benchBits, 4)
	}
}

func BenchmarkBoolEntryRuns(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RunsTest(benchBits)
	}
}

// m 不在 {4,8} 时 PokerTestBytes 会回退到 B2bitArr + PokerProto
func BenchmarkBytesPokerFallbackM3(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PokerTestBytes(benchData, 3)
	}
}
