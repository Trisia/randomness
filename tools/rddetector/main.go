package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Trisia/randomness"
)

type R struct {
	Name string
	P    []float64
}

func worker(jobs <-chan string, out chan<- *R) {
	for filename := range jobs {
		buf, _ := os.ReadFile(filename)
		bits := randomness.B2bitArr(buf)
		buf = nil
		arr := make([]float64, 0, 25)

		p, _ := randomness.MonoBitFrequencyTest(bits)
		arr = append(arr, p)
		p, _ = randomness.FrequencyWithinBlockTest(bits)
		arr = append(arr, p)
		p, _ = randomness.PokerProto(bits, 4)
		arr = append(arr, p)
		p, _ = randomness.PokerProto(bits, 8)
		arr = append(arr, p)

		p1, p2, _, _ := randomness.OverlappingTemplateMatchingProto(bits, 3)
		arr = append(arr, p1, p2)
		p1, p2, _, _ = randomness.OverlappingTemplateMatchingProto(bits, 5)
		arr = append(arr, p1, p2)

		p, _ = randomness.RunsTest(bits)
		arr = append(arr, p)
		p, _ = randomness.RunsDistributionTest(bits)
		arr = append(arr, p)
		p, _ = randomness.LongestRunOfOnesInABlockTest(bits, true)
		arr = append(arr, p)

		p, _ = randomness.BinaryDerivativeProto(bits, 3)
		arr = append(arr, p)
		p, _ = randomness.BinaryDerivativeProto(bits, 7)
		arr = append(arr, p)

		p, _ = randomness.AutocorrelationProto(bits, 1)
		arr = append(arr, p)
		p, _ = randomness.AutocorrelationProto(bits, 2)
		arr = append(arr, p)
		p, _ = randomness.AutocorrelationProto(bits, 8)
		arr = append(arr, p)
		p, _ = randomness.AutocorrelationProto(bits, 16)
		arr = append(arr, p)

		p, _ = randomness.MatrixRankTest(bits)
		arr = append(arr, p)
		p, _ = randomness.CumulativeTest(bits, true)
		arr = append(arr, p)
		p, _ = randomness.ApproximateEntropyProto(bits, 2)
		arr = append(arr, p)
		p, _ = randomness.ApproximateEntropyProto(bits, 5)
		arr = append(arr, p)
		p, _ = randomness.LinearComplexityProto(bits, 500)
		arr = append(arr, p)
		p, _ = randomness.LinearComplexityProto(bits, 1000)
		arr = append(arr, p)
		p, _ = randomness.MaurerUniversalTest(bits)
		arr = append(arr, p)
		p, _ = randomness.DiscreteFourierTransformTest(bits)
		arr = append(arr, p)

		fmt.Printf(">> 检测结束 文件 %s\n", filename)
		go func(file string) {
			out <- &R{path.Base(file), arr}
		}(filename)
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

var (
	inputPath  string // 参数文件输入路径
	reportPath string // 生成的监测报告位置
)

func init() {
	flag.StringVar(&inputPath, "i", "", "待检测随机数文件位置")
	flag.StringVar(&reportPath, "o", "RandomnessTestReport.csv", "待检测随机数文件位置")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, `randomness 随机性检测 rddetector 使用说明

rddetector -i 待检测数据目录 [-o 生成报告位置]

	示例: rddetector -i /data/target/ -o RandomnessTestReport.csv

`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if inputPath == "" {
		fmt.Fprintf(os.Stderr, "	-i 参数缺失\n\n")
		flag.Usage()
		return
	}
	_ = os.MkdirAll(filepath.Dir(reportPath), os.FileMode(0600))

	n := runtime.NumCPU()
	out := make(chan *R)
	jobs := make(chan string)

	w, err := os.OpenFile(reportPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		fmt.Fprint(os.Stderr, "无法打开写入文件 "+reportPath)
		return
	}
	defer w.Close()
	_, _ = w.WriteString(
		"源数据," +
			"单比特频数检测," +
			"块内频数检测 m=10000," +
			"扑克检测 m=4," +
			"扑克检测 m=8," +
			"重叠子序列检测 m=3 P1,重叠子序列检测 m=2 P2," +
			"重叠子序列检测 m=5 P1,重叠子序列检测 m=5 P2," +
			"游程总数检测," +
			"游程分布检测," +
			"块内最大游程检测 m=10000," +
			"二元推导检测 k=3," +
			"二元推导检测 k=7," +
			"自相关检测 d=1," +
			"自相关检测 d=2," +
			"自相关检测 d=8," +
			"自相关检测 d=16," +
			"矩阵秩检测," +
			"累加和检测," +
			"近似熵检测 m=2," +
			"近似熵检测 m=5," +
			"线性复杂度检测 m=500," +
			"线性复杂度检测 m=1000," +
			"Maurer通用统计检测 L=7 Q=1280," +
			"离散傅里叶检测\n")
	var wg sync.WaitGroup
	var cnt = make([]int32, 25)
	s := toBeTestFileNum(inputPath)
	fmt.Printf(">> 开始执行随机性检测，待检测样本数 s = %d\n", s)
	wg.Add(s)

	// 启动数据写入消费者
	go resultWriter(out, w, cnt, &wg)
	// 检测工作器
	for i := 0; i < n; i++ {
		go worker(jobs, out)
	}
	// 结果工作器
	go filepath.Walk(inputPath, func(p string, _ fs.FileInfo, _ error) error {
		if strings.HasSuffix(p, ".bin") || strings.HasSuffix(p, ".dat") {
			jobs <- p
		}
		return nil
	})

	wg.Wait()
	_, _ = w.WriteString("总计")
	for i := 0; i < len(cnt); i++ {
		_, _ = w.WriteString(fmt.Sprintf(", %d", cnt[i]))
	}
	_, _ = w.WriteString("\n")
	fmt.Println(">> 检测完成 检测报告:", reportPath)
}

func toBeTestFileNum(p string) int {
	cnt := 0
	// 结果工作器
	filepath.Walk(p, func(p string, _ fs.FileInfo, _ error) error {
		if strings.HasSuffix(p, ".bin") || strings.HasSuffix(p, ".dat") {
			cnt++
		}

		return nil
	})
	return cnt
}
