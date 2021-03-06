// Copyright (c) 2021 Quan guanyu
// randomness is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package randomness

import "fmt"

// Poker 扑克检测，m=8
func Poker(data []byte) *TestResult {
	p := PokerTestBytes(data, 8)
	return &TestResult{Name: "扑克检测", P: p, Pass: p >= Alpha}
}

// PokerTest 扑克检测，m=8
func PokerTest(bits []bool) float64 {
	return PokerProto(bits, 8)
}

// PokerTestBytes 扑克检测
// data: 检测序列
// m: m长度，m=4,8
func PokerTestBytes(data []byte, m int) float64 {
	return PokerProto(B2bitArr(data), m)
}

// PokerProto 扑克检测
// bits: 检测序列
// m: m长度，m=4,8
func PokerProto(bits []bool, m int) float64 {
	n := len(bits)

	if n < 8 {
		fmt.Println("PokerTest:args wrong")
		return -1
	}
	// 2^m
	_2m := 1 << m

	patterns := make([]int, _2m)
	N := n / m
	var V float64 = 0
	var P float64 = 0
	tmp := 0

	var b bool
	for i := 0; i < N; i++ {
		tmp = 0
		for j := 0; j < m; j++ {
			tmp <<= 1
			b, bits = bits[0], bits[1:]
			if b {
				tmp++
			}
		}
		patterns[tmp]++
	}

	for i := 0; i < _2m; i++ {
		V += float64(patterns[i]) * float64(patterns[i])
	}

	V *= float64(_2m)
	V /= float64(N)
	V -= float64(N)
	P = igamc(float64((_2m-1)>>1), V/2)
	return P
}
