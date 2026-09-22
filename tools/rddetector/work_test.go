package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/Trisia/randomness"
)

// 这两个 worker 此前没有任何测试。它们现在通过「每样本只转换一次 BitSeq」
// 的新接口驱动全部检测项，这里做一次端到端回归：
// 用确定性数据跑完整流程，校验报告的行列结构。
func TestWorkersProduceReport(t *testing.T) {
	cases := []struct {
		name     string
		worker   func(<-chan string, chan<- *R)
		bytes    int
		wantCols int
	}{
		{"2E4", worker_2E4, 20000 / 8, 45},
		{"1E6", worker_1E6, 1000000 / 8, 53},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			const files = 2
			paths := make([]string, 0, files)
			for i := 0; i < files; i++ {
				p := filepath.Join(dir, "random"+strconv.Itoa(i)+".bin")
				// 每个文件用独立种子，避免两份完全相同的数据
				data := randomness.NewDetRand(uint64(i) + 100).RawBytes(tc.bytes)
				if err := os.WriteFile(p, data, os.FileMode(0600)); err != nil {
					t.Fatal(err)
				}
				paths = append(paths, p)
			}

			report := filepath.Join(dir, "report.csv")
			collector := NewReportCollector("csv", report, "", 0.981)
			out := make(chan *R)
			jobs := make(chan string)

			var wg sync.WaitGroup
			wg.Add(files)
			go resultWriter(out, collector, &wg)
			go tc.worker(jobs, out)
			for _, p := range paths {
				jobs <- p
			}
			close(jobs)
			wg.Wait()
			close(out)

			if err := collector.GenerateReports(); err != nil {
				t.Fatalf("生成报告失败: %v", err)
			}

			f, err := os.Open(report)
			if err != nil {
				t.Fatalf("报告未生成: %v", err)
			}
			defer f.Close()

			var lines []string
			sc := bufio.NewScanner(f)
			sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
			for sc.Scan() {
				lines = append(lines, sc.Text())
			}
			if len(lines) != files+1 {
				t.Fatalf("报告应有 %d 行（表头+%d 条数据），实际 %d 行", files+1, files, len(lines))
			}
			header := strings.Split(lines[0], ",")
			if len(header) != tc.wantCols {
				t.Errorf("表头列数 %d，期望 %d", len(header), tc.wantCols)
			}
			for i, line := range lines[1:] {
				cols := strings.Split(line, ",")
				if len(cols) != tc.wantCols {
					t.Errorf("第 %d 行数据列数 %d，期望 %d", i+1, len(cols), tc.wantCols)
				}
				// 每个数值列都必须是合法浮点且在 [0,1]
				for j, c := range cols[1:] {
					v, err := strconv.ParseFloat(c, 64)
					if err != nil {
						t.Fatalf("第 %d 行第 %d 列 %q 不是合法浮点: %v", i+1, j+1, c, err)
					}
					if v < 0 || v > 1 {
						t.Errorf("第 %d 行第 %d 列 P/Q 越界: %v", i+1, j+1, v)
					}
				}
			}
		})
	}
}

// 同一种子、同一参数下 worker 结果必须完全一致（确定性数据 + 确定性算法）。
func TestWorkerDeterministic(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "a.bin")
	data := randomness.NewDetRand(20240920).RawBytes(20000 / 8)
	if err := os.WriteFile(p, data, os.FileMode(0600)); err != nil {
		t.Fatal(err)
	}

	run := func() []TestItem {
		collector := NewReportCollector("json", "", "", 0.981)
		out := make(chan *R)
		jobs := make(chan string)
		var wg sync.WaitGroup
		wg.Add(1)
		go resultWriter(out, collector, &wg)
		go worker_2E4(jobs, out)
		jobs <- p
		close(jobs)
		wg.Wait()
		close(out)
		res := collector.GetResults()
		if len(res) != 1 {
			t.Fatalf("期望 1 条结果，实际 %d", len(res))
		}
		return res[0].TestItems
	}

	a, b := run(), run()
	if len(a) != len(b) {
		t.Fatalf("两次运行的检测项数不同: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("第 %d 项结果不确定:\n %+v\n %+v", i, a[i], b[i])
		}
	}
}

// toBeTestFileNum 只统计 .bin/.dat 文件，并取最大文件长度作为样本规模。
func TestToBeTestFileNum(t *testing.T) {
	dir := t.TempDir()
	must := func(name string, n int) {
		if err := os.WriteFile(filepath.Join(dir, name), make([]byte, n), os.FileMode(0600)); err != nil {
			t.Fatal(err)
		}
	}
	must("a.bin", 2500)   // 20000 bit
	must("b.dat", 125000) // 1000000 bit
	must("c.txt", 10)     // 应被忽略
	if err := os.MkdirAll(filepath.Join(dir, "sub"), os.FileMode(0755)); err != nil {
		t.Fatal(err)
	}
	must2 := func(name string, n int) {
		if err := os.WriteFile(filepath.Join(dir, "sub", name), make([]byte, n), os.FileMode(0600)); err != nil {
			t.Fatal(err)
		}
	}
	must2("d.bin", 12500000) // 100000000 bit

	s, bits := toBeTestFileNum(dir)
	if s != 3 {
		t.Errorf("样本数 %d，期望 3", s)
	}
	if bits != 12500000*8 {
		t.Errorf("样本规模 %d bit，期望 %d", bits, 12500000*8)
	}
}
