package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

const (
	s = 1000     // 样本数量
	n = 1000_000 // 每个样本长度
)

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

func main() {
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
	fmt.Printf(">> 随机数测试组生成完成 总计 %d 组\n", s)

}
