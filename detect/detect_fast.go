package detect

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/Trisia/randomness"
)

// 工作器
// jobs: 启动参数
// source: 随机源
// n: 读取字节数
// round: 检测方式
// counter: 结果集统计
func worker(jobs chan int, source io.Reader, n int, round func([]byte) []*randomness.TestResult, counter []int32, distributions [][]float64, wait *sync.WaitGroup) {
	buf := make([]byte, n, n*2)
	for i := range jobs {
		_, err := source.Read(buf)
		if err != nil {
			continue
		}
		resArr := round(buf)
		for idx, result := range resArr {
			distributions[idx][i] = result.Q
			if result.Pass {
				atomic.AddInt32(&counter[idx], 1)
			}
		}
		wait.Done()
	}
}

// 根据处理器情况启动worker
// return 控制命令管道, 结束型号器
func bootWorker(source io.Reader, n int, round func([]byte) []*randomness.TestResult, counter []int32, distributions [][]float64) (chan int, *sync.WaitGroup) {
	var wait sync.WaitGroup
	jobs := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(jobs, source, n, round, counter, distributions, &wait)
	}
	return jobs, &wait
}

// FactoryDetectFast 出厂检测，15种检测，每组 10^6比特，分50组
// source: 随机源
func FactoryDetectFast(source io.Reader) (bool, error) {
	s := 50
	t := Threshold(s)
	n := 1000000 / 8
	counters := make([]int32, 15)
	distributions := createDistributions(s, 15)
	jobs, wg := bootWorker(source, n, Round15, counters, distributions)
	wg.Add(s)
	defer close(jobs)
	for i := 0; i < s; i++ {
		jobs <- i
	}
	wg.Wait()
	fmt.Println(counters)
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

// PowerOnDetectFast 上电自检，15种检测，每组 10^6比特，分20组
// source: 随机源
func PowerOnDetectFast(source io.Reader) (bool, error) {
	s := 20
	t := Threshold(s)
	n := 1000000 / 8
	counters := make([]int32, 15)
	distributions := createDistributions(s, 15)
	jobs, wg := bootWorker(source, n, Round15, counters, distributions)
	wg.Add(s)
	defer close(jobs)
	for i := 0; i < s; i++ {
		jobs <- i
	}
	wg.Wait()
	fmt.Println(counters)

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

// PeriodDetectFast 周期性检测，除去离散傅里叶检测、线型复杂度检测、通用统计的12种检测
// 检测 20组，每组 20000比特
// source: 随机源
func PeriodDetectFast(source io.Reader) (bool, error) {
	s := 20
	t := Threshold(s)
	n := 20000 / 8
	counters := make([]int32, 15)
	distributions := createDistributions(s, 15)
	jobs, wg := bootWorker(source, n, Round15, counters, distributions)
	wg.Add(s)
	defer close(jobs)
	for i := 0; i < s; i++ {
		jobs <- i
	}
	wg.Wait()
	fmt.Println(counters)
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
