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

// Poker 扑克检测，m=8
func Poker(data []byte) *TestResult {
	p, q := PokerTestBytes(data, 8)
	return &TestResult{Name: "扑克检测", P: p, Q: q, Pass: p >= Alpha}
}

// PokerTest 扑克检测，m=8
func PokerTest(bits []bool) (float64, float64) {
	return PokerProto(bits, 8)
}

// PokerTestBytes 扑克检测
// data: 检测序列
// m: m长度，m=4,8
func PokerTestBytes(data []byte, m int) (float64, float64) {
	if len(data) == 0 {
		panic("please provide valid test bits")
	}
	if m != 4 && m != 8 {
		panic("just support m=4 or m=8")
	}
	// 2^m
	_2m := 1 << m

	patterns := make([]int, _2m)
	N := (len(data) * 8) / m
	var V float64 = 0
	var P float64 = 0

	if m == 8 {
		for i := 0; i < N; i++ {
			patterns[data[i]]++
		}
	} else { // m = 4
		for i := 0; i < len(data); i++ {
			patterns[data[i] >> 4]++
			patterns[data[i]&0x0f]++
		}
	}

	for i := 0; i < _2m; i++ {
		V += float64(patterns[i]) * float64(patterns[i])
	}

	V *= float64(_2m)
	V /= float64(N)
	V -= float64(N)

	P = igamc(float64(_2m-1)/2, V/2)
	return P, P
}

// PokerProto 扑克检测
// bits: 检测序列
// m: m长度，m=4,8
func PokerProto(bits []bool, m int) (float64, float64) {
	n := len(bits)

	if n < 8 {
		panic("please provide valid test bits")
	}
	// 2^m
	_2m := 1 << m

	patterns := make([]int, _2m)
	N := n / m
	var V float64 = 0
	var P float64 = 0

	for i := 0; i < N; i++ {
		patterns[subsequencepattern(bits[i*m:], m)]++
	}

	for i := 0; i < _2m; i++ {
		V += float64(patterns[i]) * float64(patterns[i])
	}

	V *= float64(_2m)
	V /= float64(N)
	V -= float64(N)

	P = igamc(float64(_2m-1)/2, V/2)
	return P, P
}
