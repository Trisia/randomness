package main

import (
	"github.com/Trisia/randomness"
	"log"
	"os"
	"path"
)

// Header_2E4 20_000 比特样本检测
const Header_2E4 = "源数据," +
	"[ 1] P 单比特频数检测," +
	"[ 1] Q 单比特频数检测," +
	"[ 2] P 块内频数检测 m=1000," +
	"[ 2] Q 块内频数检测 m=1000," +
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
	"[ 5] P 游程总数检测," +
	"[ 5] Q 游程总数检测," +
	"[ 6] P 游程分布检测," +
	"[ 6] Q 游程分布检测," +
	"[ 7] P 块内最大“1”游程检测 m=128," +
	"[ 7] Q 块内最大“1”游程检测 m=128," +
	"[ 7] P 块内最大“0”游程检测 m=128," +
	"[ 7] Q 块内最大“0”游程检测 m=128," +
	"[ 8] P 二元推导检测 k=3," +
	"[ 8] Q 二元推导检测 k=3," +
	"[ 8] P 二元推导检测 k=7," +
	"[ 8] Q 二元推导检测 k=7," +
	"[ 9] P 自相关检测 d=2," +
	"[ 9] Q 自相关检测 d=2," +
	"[ 9] P 自相关检测 d=8," +
	"[ 9] Q 自相关检测 d=8," +
	"[ 9] P 自相关检测 d=16," +
	"[ 9] Q 自相关检测 d=16," +
	"[10] P 累加和检测 前向," +
	"[10] Q 累加和检测 前向," +
	"[10] P 累加和检测 后向," +
	"[10] Q 累加和检测 后向," +
	"[11] P 近似熵检测 m=2," +
	"[11] Q 近似熵检测 m=2," +
	"[11] P 近似熵检测 m=5," +
	"[11] Q 近似熵检测 m=5," +
	"[12] P 离散傅里叶检测," +
	"[12] Q 离散傅里叶检测\n"

// 数据规模为 20 000 个比特的随机数列检测工作器
func worker_2E4(jobs <-chan string, out chan<- *R) {
	for filename := range jobs {
		buf, _ := os.ReadFile(filename)
		bits := randomness.B2bitArr(buf)
		PArr := make([]float64, 0, 64)
		QArr := make([]float64, 0, 64)

		log.Printf("[%s] 检测开始...\n", filename)

		// [1] 单比特频数检测
		//p, q := randomness.MonoBitFrequencyTest(bits)
		p, q := randomness.MonoBitFrequencyTestBytes(buf)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 单比特频数检测 P: %.5f Q: %.5f", filename, p, q)

		// [2] 块内频数检测
		p, q = randomness.FrequencyWithinBlockProto(bits, 1000)
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

		// [7] 块内最大游程检测
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
		p, _ = randomness.BinaryDerivativeProto(bits, 7)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 二元推导检测 m=7 P: %.5f Q: %.5f", filename, p, q)

		// [9] 自相关检测
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

		// [10] 累加和检测
		p, q = randomness.CumulativeTest(bits, true)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 累加和检测 前向 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.CumulativeTest(bits, false)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 累加和检测 后向 P: %.5f Q: %.5f", filename, p, q)

		// [11] 近似熵检测
		p, q = randomness.ApproximateEntropyProto(bits, 2)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 近似熵检测 m=2 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.ApproximateEntropyProto(bits, 5)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 近似熵检测 m=5 P: %.5f Q: %.5f", filename, p, q)

		// [12] 离散傅里叶变换检测
		p, q = randomness.DiscreteFourierTransformTest(bits)
		PArr = append(PArr, p)
		QArr = append(QArr, q)
		log.Printf("[%s] 离散傅里叶变换检测 P: %.5f Q: %.5f", filename, p, q)

		log.Printf("[%s] 检测结束\n", filename)

		go func(file string) {
			out <- &R{path.Base(file), PArr, QArr}
		}(filename)
	}
}
