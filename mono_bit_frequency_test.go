package randomness

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// GroupBit 生成一组测试数据 长度为 10^6 比特
func GroupBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
	//buf := make([]byte, n/8)
	//_, _ = rand.Read(buf)
	//for _, b := range buf {
	//	bits = append(bits, B2bit(b)...)
	//}
	//return bits
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		if rand.Int()%2 == 1 {
			bits = append(bits, true)
		} else {
			bits = append(bits, false)
		}
	}
	return bits
}

func TestMonoBitFrequencyTest(t *testing.T) {
	bits := GroupBit()
	p := MonoBitFrequencyTest(bits)
	fmt.Printf("n: 1000000, P-valye: %0.2f\n", p)
}
