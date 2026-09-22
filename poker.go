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

import "math/bits"

// maxPokerM 是扑克检测支持的 m 上界。
const maxPokerM = 20

// corePoker 是扑克检测的 packed 内核，处理任意 m。
func corePoker(s BitSeq, m int) (float64, float64) {
	n := s.n
	if n < 8 {
		panic("please provide valid test bits")
	}
	// m 在 [1, maxPokerM] 范围内。
	if m < 1 || m > maxPokerM {
		panic("please provide valid m (1..20)")
	}
	_2m := 1 << uint(m)
	patternCounts := make([]int, _2m)
	N := n / m

	for i := 0; i < N; i++ {
		// 取 [i*m, i*m+m) 这 m 位并按 MSB 优先组成模式值。
		pos := i * m
		w := pos >> 6
		b := uint(pos & 63)
		var raw uint64
		if b == 0 {
			raw = s.word(w)
		} else {
			raw = s.word(w)>>b | s.word(w+1)<<(64-b)
		}
		patternCounts[int(bits.Reverse64(raw&maskLow(m))>>uint(64-m))]++
	}

	var V float64 = 0
	for i := 0; i < _2m; i++ {
		V += float64(patternCounts[i]) * float64(patternCounts[i])
	}
	V *= float64(_2m)
	V /= float64(N)
	V -= float64(N)

	P := igamc(float64(_2m-1)/2, V/2)
	return P, P
}

// Poker 扑克检测，m=8
func Poker(data []byte) *TestResult {
	p, q := PokerTestBytes(data, 8)
	return &TestResult{Name: "扑克检测", P: p, Q: q, Pass: p >= Alpha}
}

// PokerTest 扑克检测，m=8
//
// Deprecated: 请改用 PokerTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 PokerTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func PokerTest(bits []bool) (float64, float64) {
	return corePoker(BitSeqFromBools(bits), 8)
}

// PokerTestBytes 扑克检测
// data: 检测序列
// m: m长度，m=4,8
func PokerTestBytes(data []byte, m int) (float64, float64) {
	if len(data) == 0 {
		panic("please provide valid test bits")
	}
	// m=4/8 走字节直读快速路径；其余 m 走 packed 内核。
	if m != 4 && m != 8 {
		if m < 1 || m > maxPokerM {
			panic("please provide valid m (1..20)")
		}
		return corePoker(BitSeqFromBytes(data), m)
	}
	// 2^m
	_2m := 1 << uint(m)

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
			patterns[data[i]>>4]++
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
//
// Deprecated: 请改用 PokerTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 PokerTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func PokerProto(bits []bool, m int) (float64, float64) {
	return corePoker(BitSeqFromBools(bits), m)
}

// PokerBitSeq 扑克检测，m=8
func PokerBitSeq(s BitSeq) *TestResult {
	p, q := corePoker(s, 8)
	return &TestResult{Name: "扑克检测", P: p, Q: q, Pass: p >= Alpha}
}

// PokerTestBitSeq 扑克检测
// m: m长度，m=4,8
func PokerTestBitSeq(s BitSeq, m int) (float64, float64) {
	return corePoker(s, m)
}
