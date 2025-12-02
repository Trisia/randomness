package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// TestItem 检测项目结果
type TestItem struct {
	PValue   float64 `json:"P值"`
	QValue   float64 `json:"Q值"`
	TestName string  `json:"检测项目"`
}

// R 检测结果结构体
type R struct {
	Name      string     `json:"文件名"`
	TestItems []TestItem `json:"检测项目结果"`
}

// AnalysisResult 分析报告结果
type AnalysisResult struct {
	TestName    string  `json:"检测项目"`
	PassCount   int     `json:"通过数"`
	TotalCount  int     `json:"检测数"`
	PassRate    float64 `json:"通过率"`
	Requirement float64 `json:"满足随机性要求"`
	IsPassed    bool    `json:"是否通过"`
}

// 结果集写入文件工作器
func resultWriter(in <-chan *R, collector *ReportCollector, wg *sync.WaitGroup) {
	for r := range in {
		collector.AddResult(r)
		wg.Done()
	}
}

// Version 软件版本号
const Version = "1.5.2"

var (
	inputPath     string  // 参数文件输入路径
	reportPath    string  // 生成的监测报告位置
	NumWorkers    int     // 工作线程数
	VersionFlag   bool    // 版本号
	analysisPath  string  // 分析报告路径
	outputFormat  string  // 输出格式 (csv/json)
	passThreshold float64 // 通过判定阈值
)

func init() {
	flag.BoolVar(&VersionFlag, "v", false, "检测工具版本")
	flag.StringVar(&inputPath, "i", "", "待检测随机数文件位置")
	flag.StringVar(&reportPath, "o", "RandomnessTestReport.csv", "生成的检测报告位置")
	flag.StringVar(&analysisPath, "a", "", "生成的分析报告位置（可选）")
	flag.StringVar(&outputFormat, "f", "csv", "输出格式 (csv/json/xml)")
	flag.Float64Var(&passThreshold, "t", 0.981, "通过判定阈值（默认98.1%）")
	flag.IntVar(&NumWorkers, "n", runtime.NumCPU(), "工作线程数 (在大数据检测时通过该参数控制并行数量防止内存不足问题)")
	flag.Usage = usage

	log.SetPrefix("[rddetector] ")
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `randomness 随机性检测 rddetector v%s 使用说明

rddetector -i 待检测数据目录 [-o 生成报告位置] [-a 分析报告位置] [-f 输出格式] [-t 通过阈值]

	示例: rddetector -i /data/target/ -o RandomnessTestReport.csv -a AnalysisReport.csv -f csv -t 0.981
	示例: rddetector -i /data/target/ -o RandomnessTestReport.json -a AnalysisReport.json -f json

	数据规模将由程序自动推断，支持单文件规模 [20 000 bit, 1 000 000 bit, 100 000 000 bit]

`, Version)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if VersionFlag {
		_, _ = fmt.Fprintf(os.Stderr, "rddetector v%s\n", Version)
		return
	}

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
	switch sbit {
	case 2e4:
		worker = worker_2E4
	case 1e6:
		worker = worker_1E6
	case 1e8:
		worker = worker_1E8
	default:
		_, _ = fmt.Fprintf(os.Stderr, "无法识别待检测数据规模 %d 程序退出, 支持单文件规模 [20 000, 1 000 000, 100 000 000]\n\n", sbit)
		return
	}

	start := time.Now()

	// 创建统一数据收集器
	collector := NewReportCollector(outputFormat, reportPath, analysisPath, passThreshold)

	// 启动数据写入消费者
	go resultWriter(out, collector, &wg)

	// 检测工作器
	for i := 0; i < n; i++ {
		go worker(jobs, out)
	}
	// 结果工作器
	go filepath.Walk(inputPath, func(p string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(p, ".bin") || strings.HasSuffix(p, ".dat") {
			jobs <- p
		}
		return nil
	})

	wg.Wait()
	close(out)

	// 生成所有报告
	err := collector.GenerateReports()
	if err != nil {
		log.Fatalf("生成报告失败: %v\n", err)
	}
	log.Printf("检测完成 耗时 %s\n", time.Since(start))
	if reportPath != "" {
		log.Printf("检测报告: %s\n", reportPath)
	}
	if analysisPath != "" {
		log.Printf("分析报告: %s\n", analysisPath)
	}
}

// 获取待检测文件数量和规模
func toBeTestFileNum(p string) (samples int, bits int64) {
	// 结果工作器
	_ = filepath.Walk(p, func(p string, fInfo os.FileInfo, _ error) error {
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
