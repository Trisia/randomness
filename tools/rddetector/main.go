package main

import (
	"fmt"
	"github.com/Trisia/randomness"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type R struct {
	Name string
	P    []float64
}

func worker(jobs <-chan string, out chan<- *R) {
	for filename := range jobs {
		buf, _ := ioutil.ReadFile(filename)
		bits := randomness.B2bitArr(buf)
		buf = nil
		arr := make([]float64, 16)

		p := randomness.MonoBitFrequencyTest(bits)
		arr[0] = p
		p = randomness.FrequencyWithinBlockTest(bits)
		arr[1] = p
		p = randomness.PokerTest(bits)
		arr[2] = p
		p1, p2 := randomness.OverlappingTemplateMatchingTest(bits)
		arr[3] = p1
		arr[4] = p2

		p = randomness.RunsTest(bits)
		arr[5] = p

		p = randomness.RunsDistributionTest(bits)
		arr[6] = p

		p = randomness.LongestRunOfOnesInABlockTest(bits)
		arr[8] = p

		p = randomness.BinaryDerivativeTest(bits)
		arr[9] = p

		p = randomness.AutocorrelationTest(bits)
		arr[10] = p

		p = randomness.MatrixRankTest(bits)
		arr[10] = p

		p = randomness.CumulativeTest(bits)
		arr[11] = p

		p = randomness.ApproximateEntropyTest(bits)
		arr[12] = p

		p = randomness.LinearComplexityTest(bits)
		arr[13] = p

		p = randomness.MaurerUniversalTest(bits)
		arr[14] = p

		p = randomness.DiscreteFourierTransformTest(bits)
		arr[15] = p
		fmt.Printf(">> 文件 %s 测试结束.\n", filename)
		go func() {
			out <- &R{path.Base(filename), arr}
		}()
	}
}

// 结果集写入文件工作器
func resultWriter(in <-chan *R, w io.StringWriter, cnt []int32, wg *sync.WaitGroup) {
	for r := range in {
		_, _ = w.WriteString(r.Name)
		for j := 0; j < len(r.P); j++ {
			if r.P[j] >= 0.01 {
				atomic.AddInt32(&cnt[j], 1)
			}
			_, _ = w.WriteString(fmt.Sprintf(", %0.6f", r.P[j]))
		}
		_, _ = w.WriteString("\n")
		wg.Done()
	}

}

func main() {
	n := runtime.NumCPU()
	out := make(chan *R)
	jobs := make(chan string)
	const inputPath = "target/data"

	reportLoc, _ := filepath.Abs("target/RandomnessTestReport.csv")
	w, _ := os.OpenFile(reportLoc, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0600))
	defer w.Close()
	_, _ = w.WriteString(
		"源数据," +
			"单比特频数检测," +
			"块内频数检测," +
			"扑克检测m=8," +
			"重叠子序列检测m=5 P1,重叠子序列检测m=5 P2," +
			"游程总数检测," +
			"游程分布检测," +
			"块内最大”1“游程检测," +
			"二元推导检测k=7," +
			"自相关检测d=16," +
			"矩阵秩检测," +
			"累加和检测," +
			"近似熵检测m=5," +
			"线性复杂度检测," +
			"Maurer通用统计检测," +
			"离散傅里叶检测\n")
	var wg sync.WaitGroup
	var cnt = make([]int32, 16)
	wg.Add(toBeTestFileNum(inputPath))

	// 启动数据写入消费者
	go resultWriter(out, w, cnt, &wg)
	// 检测工作器
	for i := 0; i < n; i++ {
		go worker(jobs, out)
	}
	// 结果工作器
	go filepath.Walk(inputPath, func(p string, _ fs.FileInfo, _ error) error {
		if !strings.HasSuffix(p, ".bin") {
			return nil
		}
		jobs <- p
		return nil
	})

	wg.Wait()
	_, _ = w.WriteString("总计")
	for i := 0; i < len(cnt); i++ {
		_, _ = w.WriteString(fmt.Sprintf(", %d", cnt[i]))
	}
	_, _ = w.WriteString("\n")
	fmt.Println(">> 检测完成 检测报告:", reportLoc)
}

func toBeTestFileNum(p string) int {
	cnt := 0
	// 结果工作器
	filepath.Walk(p, func(p string, _ fs.FileInfo, _ error) error {
		if !strings.HasSuffix(p, ".bin") {
			return nil
		}
		cnt++
		return nil
	})
	return cnt
}
