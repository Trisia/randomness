package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/Trisia/randomness"
)

var (
	s      int    // 样本数量
	n      int    // 每个样本长度(bit)
	output string // 输出目录
	seed   uint64 // 确定性种子；0 表示使用 crypto/rand
	quiet  bool   // 降低日志输出
)

func init() {
	// 参数解析
	flag.IntVar(&s, "s", 1000, "Sample 样本数量")
	flag.IntVar(&n, "n", 1000000, "number 每个样本长度(bit)，必须是 8 的倍数")
	flag.StringVar(&output, "o", "target/data", "output 生成随机数文件存放目录")
	flag.Uint64Var(&seed, "seed", 0, "确定性种子；非 0 时使用 randomness.NewDetRand，同一 -seed/-s/-n 生成完全相同的文件内容（0 表示使用 crypto/rand）")
	flag.BoolVar(&quiet, "q", false, "安静模式，只输出汇总信息")
	flag.Usage = usage
}

func usage() {
	fmt.Fprint(os.Stderr, `randomness 随机数生成工具 rdgen 使用说明

rdgen [-s 生成文件数] [-n 每个文件内bit数] [-o 生成位置] [-seed 种子] [-q]

	示例: rdgen -s 1000 -n 1000000 -o ./data

	默认使用 crypto/rand，每次运行结果不同。
	指定 -seed <非0> 后改用库内的确定性发生器 randomness.NewDetRand（xoshiro256**）：
	第 i 个文件由 seed+i 生成，因此同一组参数永远产出相同的数据，
	便于把检测结果固化为可复现的回归基线。

`)
	flag.PrintDefaults()
}

// generateOne 生成第 index 个样本文件。
func generateOne(source io.Reader, path string, nbits int, index int) error {
	buf := make([]byte, nbits/8)
	if seed != 0 {
		randomness.NewDetRand(seed + uint64(index)).Read(buf)
	} else if _, err := io.ReadFull(source, buf); err != nil {
		return fmt.Errorf("读取随机源失败: %v", err)
	}
	return ioutil.WriteFile(path, buf, os.FileMode(0600))
}

func main() {
	flag.Parse()

	if s < 0 {
		fmt.Fprintf(os.Stderr, "	-s 不能为负\n\n")
		os.Exit(2)
	}
	if n < 8 || n%8 != 0 {
		fmt.Fprintf(os.Stderr, "	-n 必须是 >= 8 且为 8 的倍数（当前 %d），否则无法整字节写入\n\n", n)
		os.Exit(2)
	}

	outDir, err := filepath.Abs(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "	解析输出目录失败: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(outDir, os.FileMode(0755)); err != nil {
		fmt.Fprintf(os.Stderr, "	创建输出目录失败: %v\n", err)
		os.Exit(1)
	}

	src := "crypto/rand"
	if seed != 0 {
		src = fmt.Sprintf("内置确定性发生器 seed=%d", seed)
	}
	fmt.Printf(">> 生成文件数 %d 每个样本长度 %d bit  输出位置: %s  随机源: %s\n", s, n, outDir, src)

	var source io.Reader = rand.Reader
	jobs := make(chan int)
	var wg sync.WaitGroup

	// errOnce 记录第一个错误；出错后不再继续投递新任务。
	var mu sync.Mutex
	var firstErr error
	recordErr := func(err error) {
		mu.Lock()
		if firstErr == nil {
			firstErr = err
		}
		mu.Unlock()
	}

	workers := runtime.NumCPU()
	if workers > s {
		workers = s
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				name := filepath.Join(outDir, fmt.Sprintf("random%d.bin", idx))
				if err := generateOne(source, name, n, idx); err != nil {
					recordErr(fmt.Errorf("%s: %v", name, err))
					return
				}
				if !quiet {
					fmt.Println(">> 生成随机数: ", name)
				}
			}
		}()
	}

	for i := 0; i < s; i++ {
		mu.Lock()
		failed := firstErr != nil
		mu.Unlock()
		if failed {
			break
		}
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	if firstErr != nil {
		fmt.Fprintf(os.Stderr, "生成失败: %v\n", firstErr)
		os.Exit(1)
	}
	fmt.Printf(">> 随机数测试组生成完成 总计 %d 组 位于: %s\n", s, outDir)
}
