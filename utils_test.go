package randomness

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestB2bit(t *testing.T) {

	tests := []struct {
		name string
		arg  byte
		want []bool
	}{
		{"单字节 1", 0x01, []bool{false, false, false, false, false, false, false, true}},
		{"单字节 2", 0x02, []bool{false, false, false, false, false, false, true, false}},
		{"单字节 4", 0x04, []bool{false, false, false, false, false, true, false, false}},
		{"单字节 8", 0x08, []bool{false, false, false, false, true, false, false, false}},
		{"单字节 16", 0x0C, []bool{false, false, false, false, true, true, false, false}},
		{"单字节 255", 0xFF, []bool{true, true, true, true, true, true, true, true}},
		{"单字节 45", 0x2D, []bool{false, false, true, false, true, true, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := B2bit(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("B2bit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCeilPow2(t *testing.T) {
	tests := []struct {
		name string
		N    int
		want int
	}{
		{"1", 1, 2},
		{"2", 2, 2},
		{"3", 3, 4},
		{"4", 4, 4},
		{"5", 5, 8},
		{"6", 6, 8},
		{"7", 7, 8},
		{"8", 8, 8},
		{"9", 9, 16},
		{"10", 10, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ceilPow2(tt.N); got != tt.want {
				t.Errorf("ceilPow2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestB2Byte(t *testing.T) {

	tests := []struct {
		name string
		args []bool
		want byte
	}{
		{"单字节 1  ", []bool{false, false, false, false, false, false, false, true}, 0x01},
		{"单字节 2  ", []bool{false, false, false, false, false, false, true, false}, 0x02},
		{"单字节 4  ", []bool{false, false, false, false, false, true, false, false}, 0x04},
		{"单字节 8  ", []bool{false, false, false, false, true, false, false, false}, 0x08},
		{"单字节 16 ", []bool{false, false, false, false, true, true, false, false}, 0x0C},
		{"单字节 255", []bool{true, true, true, true, true, true, true, true}, 0xFF},
		{"单字节 45 ", []bool{false, false, true, false, true, true, false, true}, 0x2D},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := B2Byte(tt.args); got != tt.want {
				t.Errorf("B2Byte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupBit(t *testing.T) {
	bits := GroupBit()
	var tmp []bool
	var buf []byte
	for {
		if len(bits) < 8 {
			break
		}
		tmp, bits = bits[:8], bits[8:]
		buf = append(buf, B2Byte(tmp))
	}
	_, _ = rand.Read(buf)
	err := ioutil.WriteFile("data/data.bin", buf, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
