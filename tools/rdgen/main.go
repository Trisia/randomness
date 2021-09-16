package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	s      int    // 样本数量
	n      int    // 每个样本长度
	output string // 输出目录
)

func init() {
	// 参数解析
	flag.IntVar(&s, "s", 1000, "Sample 样本数量")
	flag.IntVar(&n, "n", 1000_000, "number 每个样本长度(bit)")
	flag.StringVar(&output, "o", "target/data", "output 生成随机数文件存放目录")
	flag.Usage = usage
}

func worker(jobs chan int, source io.Reader, wg *sync.WaitGroup) {
	_ = os.MkdirAll("target/data", os.FileMode(0600))
	buf := make([]byte, n/8)
	for i := range jobs {
		name := fmt.Sprintf("target/data/random%d.bin", i)
		fmt.Println(">> 生成随机数: ", name)
		w, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0600))
		if err != nil {
			_ = w.Close()
			panic(err)
		}
		_, err = source.Read(buf)
		if err != nil {
			_ = w.Close()
			panic(err)
		}
		_, err = w.Write(buf)
		_ = w.Close()
		wg.Done()
		if err != nil {
			panic(err)
		}
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `randomness 随机数生成工具 rdgen 使用说明

rdgen [-s 生成文件数] [-n 每个文件内bit数] [-o 生成位置]

	示例: rdgen -s 1000 -n 1000000 -o ./data

`)
	flag.PrintDefaults()
}

func main() {
	// 解析命令行参数
	flag.Parse()
	var err error
	output, err = filepath.Abs(output)
	if err != nil {
		panic(err)
	}
	_ = os.MkdirAll(output, os.FileMode(0600))
	fmt.Printf(">> 生成文件数 %d 每个样本长度 %d  输出位置: %s\n", s, n, output)

	var wg sync.WaitGroup
	wg.Add(s)
	source := rand.Reader
	jobs := make(chan int)
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(jobs, source, &wg)
	}
	for i := 0; i < s; i++ {
		jobs <- i
	}
	wg.Wait()
	fmt.Printf(">> 随机数测试组生成完成 总计 %d 组 位于: %s\n", s, output)

}
