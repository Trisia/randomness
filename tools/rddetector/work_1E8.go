package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/Trisia/randomness"
)

// Header_1E8 100_000_000 比特样本检测
const Header_1E8 = "源数据," +
	"[ 1] P 单比特频数检测," +
	"[ 1] Q 单比特频数检测," +
	"[ 2] P 块内频数检测 m=100000," +
	"[ 2] Q 块内频数检测 m=100000," +
	"[ 3] P 扑克检测 m=4," +
	"[ 3] Q 扑克检测 m=4," +
	"[ 3] P 扑克检测 m=8," +
	"[ 3] Q 扑克检测 m=8," +
	"[ 4] P1 重叠子序列检测 m=3," +
	"[ 4] Q1 重叠子序列检测 m=2," +
	"[ 4] P2 重叠子序列检测 m=3," +
	"[ 4] Q2 重叠子序列检测 m=2," +
	"[ 4] P1 重叠子序列检测 m=5," +
	"[ 4] Q1 重叠子序列检测 m=5," +
	"[ 4] P2 重叠子序列检测 m=5," +
	"[ 4] Q2 重叠子序列检测 m=5," +
	"[ 4] P1 重叠子序列检测 m=7," +
	"[ 4] Q1 重叠子序列检测 m=7," +
	"[ 4] P2 重叠子序列检测 m=7," +
	"[ 4] Q2 重叠子序列检测 m=7," +
	"[ 5] P 游程总数检测," +
	"[ 5] Q 游程总数检测," +
	"[ 6] P 游程分布检测," +
	"[ 6] Q 游程分布检测," +
	"[ 7] P 块内最大“1”游程检测 m=10000," +
	"[ 7] Q 块内最大“1”游程检测 m=10000," +
	"[ 7] P 块内最大“0”游程检测 m=10000," +
	"[ 7] Q 块内最大“0”游程检测 m=10000," +
	"[ 8] P 二元推导检测 k=3," +
	"[ 8] Q 二元推导检测 k=3," +
	"[ 8] P 二元推导检测 k=7," +
	"[ 8] Q 二元推导检测 k=7," +
	"[ 8] P 二元推导检测 k=15," +
	"[ 8] Q 二元推导检测 k=15," +
	"[ 9] P 自相关检测 d=1," +
	"[ 9] Q 自相关检测 d=1," +
	"[ 9] P 自相关检测 d=2," +
	"[ 9] Q 自相关检测 d=2," +
	"[ 9] P 自相关检测 d=8," +
	"[ 9] Q 自相关检测 d=8," +
	"[ 9] P 自相关检测 d=16," +
	"[ 9] Q 自相关检测 d=16," +
	"[ 9] P 自相关检测 d=32," +
	"[ 9] Q 自相关检测 d=32," +
	"[10] P 矩阵秩检测," +
	"[10] Q 矩阵秩检测," +
	"[11] P 累加和检测 前向," +
	"[11] Q 累加和检测 前向," +
	"[11] P 累加和检测 后向," +
	"[11] Q 累加和检测 后向," +
	"[12] P 近似熵检测 m=5," +
	"[12] Q 近似熵检测 m=5," +
	"[12] P 近似熵检测 m=7," +
	"[12] Q 近似熵检测 m=7," +
	"[13] P 线性复杂度检测 m=5000," +
	"[13] Q 线性复杂度检测 m=5000," +
	"[14] P Maurer通用统计检测 L=7 Q=1280," +
	"[14] Q Maurer通用统计检测 L=7 Q=1280," +
	"[15] P 离散傅里叶检测," +
	"[15] Q 离散傅里叶检测\n"

// 数据规模为 100 000 000 个比特的随机数列检测工作器
func worker_1E8(jobs <-chan string, out chan<- *R) {
	for filename := range jobs {
		buf, _ := ioutil.ReadFile(filename)
		bits := randomness.B2bitArr(buf)
		PArr := make([]float64, 0, 64)
		QArr := make([]float64, 0, 64)

		log.Printf("[%s] 检测开始...\n", filename)

		// [1] 单比特频数检测
		//p, q := randomness.MonoBitFrequencyTest(bits)
		p, q := randomness.MonoBitFrequencyTestBytes(buf)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 单比特频数检测 P: %.5f Q: %.5f", filename, p, p)

		// [2] 块内频数检测
		p, q = randomness.FrequencyWithinBlockProto(bits, 100000)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 块内频数检测 m=100_000 P: %.5f Q: %.5f", filename, p, q)

		// [3] 扑克检测
		//p, q = randomness.PokerProto(bits, 4)
		p, q = randomness.PokerTestBytes(buf, 4)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 扑克检测 m=4 P: %.5f Q: %.5f", filename, p, q)
		//p, q = randomness.PokerProto(bits, 8)
		p, q = randomness.PokerTestBytes(buf, 8)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 扑克检测 m=8 P: %.5f Q: %.5f", filename, p, q)

		// 下文中不再需要比特数组，释放以节约内存。
		buf = nil

		// [4] 重叠子序列检测
		p1, p2, q1, q2 := randomness.OverlappingTemplateMatchingProto(bits, 3)
		PArr = append(PArr, p1, p2)
		QArr = append(QArr, q1, q2)
		log.Printf("[%s] 重叠子序列检测 m=3 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)
		p1, p2, q1, q2 = randomness.OverlappingTemplateMatchingProto(bits, 5)
		PArr = append(PArr, p1, p2)
		QArr = append(QArr, q1, q2)
		log.Printf("[%s] 重叠子序列检测 m=5 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)
		p1, p2, q1, q2 = randomness.OverlappingTemplateMatchingProto(bits, 7)
		PArr = append(PArr, p1, p2)
		QArr = append(QArr, q1, q1)
		log.Printf("[%s] 重叠子序列检测 m=7 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)

		// [5] 游程总数检测
		p, q = randomness.RunsTest(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 游程总数检测 P: %.5f Q: %.5f", filename, p, q)

		// [6] 游程分布检测
		p, q = randomness.RunsDistributionTest(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 游程分布检测 P: %.5f Q: %.5f", filename, p, q)

		// [7] 块内最大游程检测 m=10000
		p, q = randomness.LongestRunOfOnesInABlockTest(bits, true)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 块内最大“1”游程检测 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.LongestRunOfOnesInABlockTest(bits, false)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 块内最大“0”游程检测 P: %.5f Q: %.5f", filename, p, q)

		// [8] 二元推导检测
		p, q = randomness.BinaryDerivativeProto(bits, 3)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 二元推导检测 m=3 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.BinaryDerivativeProto(bits, 7)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 二元推导检测 m=7 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.BinaryDerivativeProto(bits, 15)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 二元推导检测 m=15 P: %.5f Q: %.5f", filename, p, q)

		// [9] 自相关检测
		p, q = randomness.AutocorrelationProto(bits, 1)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 自相关检测 m=1 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 2)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 自相关检测 m=2 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 8)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 自相关检测 m=8 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 16)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 自相关检测 m=16 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 32)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 自相关检测 m=32 P: %.5f Q: %.5f", filename, p, q)

		// [10] 矩阵秩检测
		p, q = randomness.MatrixRankTest(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 矩阵秩检测 P: %.5f Q: %.5f", filename, p, q)

		// [11] 累加和检测
		p, q = randomness.CumulativeTest(bits, true)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 累加和检测 前向 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.CumulativeTest(bits, false)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 累加和检测 后向 P: %.5f Q: %.5f", filename, p, q)

		// [12] 近似熵检测
		p, q = randomness.ApproximateEntropyProto(bits, 5)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 近似熵检测 m=5 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.ApproximateEntropyProto(bits, 7)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 近似熵检测 m=7 P: %.5f Q: %.5f", filename, p, q)

		// [13] 线性复杂度检测
		p, q = randomness.LinearComplexityProto(bits, 5000)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 线性复杂度检测 m=5000 P: %.5f Q: %.5f", filename, p, q)

		// [14] 通用统计检测
		p, q = randomness.MaurerUniversalTest(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] Maurer通用统计检测 P: %.5f Q: %.5f", filename, p, q)

		// [15] 离散傅里叶变换检测
		p, q = randomness.DiscreteFourierTransformTestFast(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 离散傅里叶变换检测 P: %.5f Q: %.5f", filename, p, q)

		log.Printf("[%s] 检测结束\n", filename)

		go func(file string) {
			out <- &R{path.Base(file), PArr, QArr}
		}(filename)
	}
}
