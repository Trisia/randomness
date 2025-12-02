package main

import (
	"io/ioutil"
	"log"
	"path"

	"github.com/Trisia/randomness"
)

// 数据规模为 100 000 000 个比特的随机数列检测工作器
func worker_1E8(jobs <-chan string, out chan<- *R) {
	for filename := range jobs {
		buf, _ := ioutil.ReadFile(filename)
		bits := randomness.B2bitArr(buf)

		testItems := make([]TestItem, 0, 64)

		log.Printf("[%s] 检测开始...\n", filename)

		// [1] 单比特频数检测
		p, q := randomness.MonoBitFrequencyTestBytes(buf)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "单比特频数检测"})
		log.Printf("[%s] 单比特频数检测 P: %.5f Q: %.5f", filename, p, q)

		// [2] 块内频数检测
		p, q = randomness.FrequencyWithinBlockProto(bits, 100000)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "块内频数检测 m=100000"})
		log.Printf("[%s] 块内频数检测 m=100000 P: %.5f Q: %.5f", filename, p, q)

		// [3] 扑克检测
		p, q = randomness.PokerTestBytes(buf, 4)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "扑克检测 m=4"})
		log.Printf("[%s] 扑克检测 m=4 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.PokerTestBytes(buf, 8)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "扑克检测 m=8"})
		log.Printf("[%s] 扑克检测 m=8 P: %.5f Q: %.5f", filename, p, q)

		// 下文中不再需要比特数组，释放以节约内存。
		buf = nil

		// [4] 重叠子序列检测
		p1, p2, q1, q2 := randomness.OverlappingTemplateMatchingProto(bits, 3)
		testItems = append(testItems, TestItem{PValue: p1, QValue: q1, TestName: "重叠子序列检测 m=3 P1"})
		testItems = append(testItems, TestItem{PValue: p2, QValue: q2, TestName: "重叠子序列检测 m=3 P2"})
		log.Printf("[%s] 重叠子序列检测 m=3 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)
		p1, p2, q1, q2 = randomness.OverlappingTemplateMatchingProto(bits, 5)
		testItems = append(testItems, TestItem{PValue: p1, QValue: q1, TestName: "重叠子序列检测 m=5 P1"})
		testItems = append(testItems, TestItem{PValue: p2, QValue: q2, TestName: "重叠子序列检测 m=5 P2"})
		log.Printf("[%s] 重叠子序列检测 m=5 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)
		p1, p2, q1, q2 = randomness.OverlappingTemplateMatchingProto(bits, 7)
		testItems = append(testItems, TestItem{PValue: p1, QValue: q1, TestName: "重叠子序列检测 m=7 P1"})
		testItems = append(testItems, TestItem{PValue: p2, QValue: q2, TestName: "重叠子序列检测 m=7 P2"})
		log.Printf("[%s] 重叠子序列检测 m=7 P1: %.5f P2: %.5f Q1: %.5f Q2: %.5f", filename, p1, p2, q1, q2)

		// [5] 游程总数检测
		p, q = randomness.RunsTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "游程总数检测"})
		log.Printf("[%s] 游程总数检测 P: %.5f Q: %.5f", filename, p, q)

		// [6] 游程分布检测
		p, q = randomness.RunsDistributionTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "游程分布检测"})
		log.Printf("[%s] 游程分布检测 P: %.5f Q: %.5f", filename, p, q)

		// [7] 块内最大游程检测
		p, q = randomness.LongestRunOfOnesInABlockTest(bits, true)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "块内最大\"1\"游程检测 m=10000"})
		log.Printf("[%s] 块内最大\"1\"游程检测 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.LongestRunOfOnesInABlockTest(bits, false)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "块内最大\"0\"游程检测 m=10000"})
		log.Printf("[%s] 块内最大\"0\"游程检测 P: %.5f Q: %.5f", filename, p, q)

		// [8] 二元推导检测
		p, q = randomness.BinaryDerivativeProto(bits, 3)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "二元推导检测 k=3"})
		log.Printf("[%s] 二元推导检测 k=3 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.BinaryDerivativeProto(bits, 7)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "二元推导检测 k=7"})
		log.Printf("[%s] 二元推导检测 k=7 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.BinaryDerivativeProto(bits, 15)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "二元推导检测 k=15"})
		log.Printf("[%s] 二元推导检测 k=15 P: %.5f Q: %.5f", filename, p, q)

		// [9] 自相关检测
		p, q = randomness.AutocorrelationProto(bits, 1)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "自相关检测 d=1"})
		log.Printf("[%s] 自相关检测 d=1 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 2)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "自相关检测 d=2"})
		log.Printf("[%s] 自相关检测 d=2 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 8)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "自相关检测 d=8"})
		log.Printf("[%s] 自相关检测 d=8 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.AutocorrelationProto(bits, 16)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "自相关检测 d=16"})
		log.Printf("[%s] 自相关检测 d=16 P: %.5f Q: %.5f", filename, p, q)

		// [10] 矩阵秩检测
		p, q = randomness.MatrixRankTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "矩阵秩检测"})
		log.Printf("[%s] 矩阵秩检测 P: %.5f Q: %.5f", filename, p, q)

		// [11] 累加和检测
		p, q = randomness.CumulativeTest(bits, true)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "累加和检测 前向"})
		log.Printf("[%s] 累加和检测 前向 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.CumulativeTest(bits, false)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "累加和检测 后向"})
		log.Printf("[%s] 累加和检测 后向 P: %.5f Q: %.5f", filename, p, q)

		// [12] 近似熵检测
		p, q = randomness.ApproximateEntropyProto(bits, 2)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "近似熵检测 m=2"})
		log.Printf("[%s] 近似熵检测 m=2 P: %.5f Q: %.5f", filename, p, q)
		p, q = randomness.ApproximateEntropyProto(bits, 5)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "近似熵检测 m=5"})
		log.Printf("[%s] 近似熵检测 m=5 P: %.5f Q: %.5f", filename, p, q)

		// [13] 线性复杂度检测
		p, q = randomness.LinearComplexityTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "线性复杂度检测 m=500"})
		log.Printf("[%s] 线性复杂度检测 m=500 P: %.5f Q: %.5f", filename, p, q)

		// [14] Maurer通用统计检测
		p, q = randomness.MaurerUniversalTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "Maurer通用统计检测 L=7 Q=1280"})
		log.Printf("[%s] Maurer通用统计检测 P: %.5f Q: %.5f", filename, p, q)

		// [15] 离散傅里叶检测
		p, q = randomness.DiscreteFourierTransformTest(bits)
		testItems = append(testItems, TestItem{PValue: p, QValue: q, TestName: "离散傅里叶检测"})
		log.Printf("[%s] 离散傅里叶检测 P: %.5f Q: %.5f", filename, p, q)

		out <- &R{Name: path.Base(filename), TestItems: testItems}
	}
}
