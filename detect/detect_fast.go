package detect

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Trisia/randomness"
)

// job 是一个待检测样本：index 为样本序号，seq 为该样本的位序列。
type job struct {
	index int
	seq   randomness.BitSeq
}

// runFast 是 FactoryDetectFast / PowerOnDetectFast / PeriodDetectFast 的共同实现。
//
// s: 样本组数；methods: 检测项数；n: 每组字节数；round: 检测轮次（接受 BitSeq）。
func runFast(source io.Reader, s, methods, n int, round func(randomness.BitSeq) []*randomness.TestResult) (bool, error) {
	t := Threshold(s)
	counters := make([]int32, methods)
	distributions := createDistributions(s, methods)

	jobs := make(chan job)
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				resArr := round(j.seq)
				for idx, result := range resArr {
					distributions[idx][j.index] = result.Q
					if result.Pass {
						atomic.AddInt32(&counters[idx], 1)
					}
				}
			}
		}()
	}

	// 主协程顺序读取并分发。一旦出错就停止投递并关闭 channel。
	var rerr error
	for i := 0; i < s; i++ {
		buf := make([]byte, n)
		if _, err := io.ReadFull(source, buf); err != nil {
			rerr = err
			break
		}
		jobs <- job{index: i, seq: randomness.BitSeqFromBytes(buf)}
	}
	close(jobs)
	wg.Wait()

	if rerr != nil {
		return false, fmt.Errorf("读取随机源失败: %v", rerr)
	}

	for i, itemCnt := range counters {
		if int(itemCnt) < t {
			return false, fmt.Errorf("%s %d/%d", randomness.TestMethodArr[i].Name, itemCnt, s)
		}
	}
	for i := range distributions {
		Pt := ThresholdQ(distributions[i])
		if Pt < randomness.AlphaT {
			return false, fmt.Errorf("%s %f", randomness.TestMethodArr[i].Name, Pt)
		}
	}
	return true, nil
}

// FactoryDetectFast 出厂检测，15种检测，每组 10^6比特，分50组
// source: 随机源
func FactoryDetectFast(source io.Reader) (bool, error) {
	return runFast(source, 50, 15, 1000000/8, Round15BitSeq)
}

// PowerOnDetectFast 上电自检，15种检测，每组 10^6比特，分20组
// source: 随机源
func PowerOnDetectFast(source io.Reader) (bool, error) {
	return runFast(source, 20, 15, 1000000/8, Round15BitSeq)
}

// PeriodDetectFast 周期性检测，除去离散傅里叶检测、线型复杂度检测、通用统计的12种检测
// 检测 20组，每组 20000比特
// source: 随机源
func PeriodDetectFast(source io.Reader) (bool, error) {
	return runFast(source, 20, 12, 20000/8, Round12BitSeq)
}
