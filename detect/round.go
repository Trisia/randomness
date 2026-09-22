package detect

import (
	"sync"

	"github.com/Trisia/randomness"
)

// Round15 15种方法测试轮
// data: 待检测数据，推荐长度： 10^6 bit =>  125,000 byte
func Round15(data []byte) []*randomness.TestResult {
	results := make([]*randomness.TestResult, 15)
	for i, method := range randomness.TestMethodArr {
		results[i] = method.Runner(data)
	}
	return results
}

// Round12 12种方法测试轮（除去：离散傅里叶检测、线型复杂度检测、通用统计）
// data: 待检测数据，推荐长度： 20000 bit =>  2,500
func Round12(data []byte) []*randomness.TestResult {
	results := make([]*randomness.TestResult, 12)
	arr := randomness.TestMethodArr[:12]
	for i, method := range arr {
		results[i] = method.Runner(data)
	}
	return results
}

// roundConcurrent 并发执行给定检测项，返回顺序与串行版完全一致。
func roundConcurrent(data []byte, methods []randomness.TestItem) []*randomness.TestResult {
	results := make([]*randomness.TestResult, len(methods))
	var wg sync.WaitGroup
	for i := range methods {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = methods[i].Runner(data)
		}(i)
	}
	wg.Wait()
	return results
}

// Round15Concurrent 与 Round15 返回完全相同的结果，但把 15 项检测并发执行。
// 适用场景：一次只处理一个（或少量）大样本。
// 不适用场景：已在样本级并行的调用方（FactoryDetectFast / PowerOnDetectFast / PeriodDetectFast）。
func Round15Concurrent(data []byte) []*randomness.TestResult {
	return roundConcurrent(data, randomness.TestMethodArr)
}

// Round12Concurrent 与 Round12 返回完全相同的结果，但把 12 项检测并发执行。
// 适用/不适用场景同 Round15Concurrent。
func Round12Concurrent(data []byte) []*randomness.TestResult {
	return roundConcurrent(data, randomness.TestMethodArr[:12])
}

// Round15BitSeq 与 Round15 返回完全相同的结果，但直接接受 BitSeq。
func Round15BitSeq(s randomness.BitSeq) []*randomness.TestResult {
	return roundSeq(s, randomness.TestMethodSeqArr)
}

// Round12BitSeq 与 Round12 返回完全相同的结果，但直接接受 BitSeq。
func Round12BitSeq(s randomness.BitSeq) []*randomness.TestResult {
	return roundSeq(s, randomness.TestMethodSeqArr[:12])
}

func roundSeq(s randomness.BitSeq, methods []randomness.TestSeqItem) []*randomness.TestResult {
	results := make([]*randomness.TestResult, len(methods))
	for i := range methods {
		results[i] = methods[i].Runner(s)
	}
	return results
}

// Round15BitSeqConcurrent 与 Round15Concurrent 等价，但直接接受 BitSeq。
// 适用/不适用场景同 Round15Concurrent。
func Round15BitSeqConcurrent(s randomness.BitSeq) []*randomness.TestResult {
	return roundConcurrentSeq(s, randomness.TestMethodSeqArr)
}

// Round12BitSeqConcurrent 与 Round12Concurrent 等价，但直接接受 BitSeq。
func Round12BitSeqConcurrent(s randomness.BitSeq) []*randomness.TestResult {
	return roundConcurrentSeq(s, randomness.TestMethodSeqArr[:12])
}

func roundConcurrentSeq(s randomness.BitSeq, methods []randomness.TestSeqItem) []*randomness.TestResult {
	results := make([]*randomness.TestResult, len(methods))
	var wg sync.WaitGroup
	for i := range methods {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i] = methods[i].Runner(s)
		}(i)
	}
	wg.Wait()
	return results
}
