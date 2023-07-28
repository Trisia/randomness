package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type R struct {
	Name string
	P    []float64 // 样板通过率	  should >= 0.01
	Q    []float64 // 样本分布均匀性 should >= 0.0001
}

// 结果集写入文件工作器
func resultWriter(in <-chan *R, w io.StringWriter, wg *sync.WaitGroup) {
	for r := range in {
		_, _ = w.WriteString(r.Name)
		for j := 0; j < len(r.P); j++ {
			_, _ = w.WriteString(fmt.Sprintf(", %0.6f, %0.6f", r.P[j], r.Q[j]))
		}
		_, _ = w.WriteString("\n")
		wg.Done()
	}

}

var (
	inputPath  string // 参数文件输入路径
	reportPath string // 生成的监测报告位置
	NumWorkers int    // 工作线程数
)

func init() {
	flag.StringVar(&inputPath, "i", "", "待检测随机数文件位置")
	flag.StringVar(&reportPath, "o", "RandomnessTestReport.csv", "待检测随机数文件位置")
	flag.IntVar(&NumWorkers, "n", runtime.NumCPU(), "工作线程数 (在大数据检测时通过该参数控制并行数量防止内存不足问题)")
	flag.Usage = usage

	log.SetPrefix("[rddetector] ")
}
func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `randomness 随机性检测 rddetector 使用说明

rddetector -i 待检测数据目录 [-o 生成报告位置]

	示例: rddetector -i /data/target/ -o RandomnessTestReport.csv

	数据规模将由程序自动推断，支持单文件规模 [20 000 bit, 1 000 000 bit, 100 000 000 bit]

`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if inputPath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "	-i 参数缺失\n\n")
		flag.Usage()
		return
	}

	n := NumWorkers
	out := make(chan *R)
	jobs := make(chan string)

	var wg sync.WaitGroup
	s, sbit := toBeTestFileNum(inputPath)
	log.Printf("启动 随机性检测，待检测样本总数 s = %d 样本数据规模 bits = %d\n", s, sbit)
	wg.Add(s)

	var worker func(jobs <-chan string, out chan<- *R) = nil
	hdr := ""
	switch sbit {
	case 2e4:
		hdr = Header_2E4
		worker = worker_2E4
	case 1e6:
		hdr = Header_1E6
		worker = worker_1E6
	case 1e8:
		hdr = Header_1E8
		worker = worker_1E8
	default:
		_, _ = fmt.Fprintf(os.Stderr, "无法识别待检测数据规模 %d 程序退出, 支持单文件规模 [20 000, 1 000 000, 100 000 000]\n\n", sbit)
		return
	}

	_ = os.MkdirAll(filepath.Dir(reportPath), os.FileMode(0600))
	w, err := os.OpenFile(reportPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.FileMode(0600))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "无法打开写入文件 "+reportPath)
		return
	}
	defer w.Close()

	_, _ = w.WriteString(hdr)
	start := time.Now()

	// 启动数据写入消费者
	go resultWriter(out, w, &wg)
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

	log.Printf("检测完成 耗时 %s 检测报告: %s\n", time.Since(start), reportPath)
}

// 获取待检测文件数量和规模
func toBeTestFileNum(p string) (samples int, bits int64) {
	// 结果工作器
	_ = filepath.Walk(p, func(p string, fInfo fs.FileInfo, _ error) error {
		if fInfo == nil || fInfo.IsDir() {
			return nil
		}

		if strings.HasSuffix(p, ".bin") || strings.HasSuffix(p, ".dat") {
			samples++
			if fInfo.Size()*8 > bits {
				bits = fInfo.Size() * 8
			}
		}
		return nil
	})
	return samples, bits
}
