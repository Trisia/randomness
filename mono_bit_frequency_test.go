package randomness

import (
	"fmt"
	"math/rand"
	"testing"
)

// GroupBit 生成一组测试数据 长度为 10^6 比特
func GroupBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)
	for _, b := range buf {
		bits = append(bits, B2bit(b)...)
	}
	return bits
}

func TestMonoBitFrequencyTest(t *testing.T) {
	bits := GroupBit()
	p := MonoBitFrequencyTest(bits)
	fmt.Printf("n: 1000000, P-valye: %0.2f\n", p)
}
