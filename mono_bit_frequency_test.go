package randomness

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func GroupBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
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

// GroupBit 生成一组测试数据 长度为 10^6 比特
func GroupSecBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)

	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)
	for _, b := range buf {
		bits = append(bits, B2bit(b)...)
	}
	return bits
}

func ReadGroup(filename string) []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for _, b := range buf {
		bits = append(bits, B2bit(b)...)
	}
	return bits
}

func TestGroupBit(t *testing.T) {
	bits := GroupBit()
	var tmp []bool
	var buf []byte
	for i := 0; i < len(bits)/8; i++ {
		tmp, bits = bits[:8], bits[8:]
		buf = append(buf, B2Byte(tmp))
	}
	_, _ = rand.Read(buf)
	err := ioutil.WriteFile("data.bin", buf, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMonoBitFrequencyTest(t *testing.T) {
	bits := GroupBit()
	//bits := ReadGroup("data.bin")
	p := MonoBitFrequencyTest(bits)
	fmt.Printf("n: 1000000, P-value: %.6f\n", p)
}
