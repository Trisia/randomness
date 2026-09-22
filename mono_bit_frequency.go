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

import (
	"math"
	"math/bits"
)

// coreMonoBitFrequency 是单比特频数检测的 packed 内核。
func coreMonoBitFrequency(s BitSeq) (float64, float64) {
	if s.n == 0 {
		panic("please provide test bits")
	}
	S := 2*s.ones() - s.n
	V := float64(S) / math.Sqrt(float64(2*s.n)) // 除math.Sqrt(2)，放到这里提前处理(n->2*n)，减少math.Sqrt的调用。
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// MonoBitFrequency 单比特频数检测
func MonoBitFrequency(data []byte) *TestResult {
	p, q := MonoBitFrequencyTestBytes(data)
	return &TestResult{Name: "单比特频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// MonoBitFrequencyTestBytes 单比特频数检测，直接对字节处理。
func MonoBitFrequencyTestBytes(data []byte) (float64, float64) {
	if len(data) == 0 {
		panic("please provide test bits")
	}
	n := len(data) * 8
	S := 0
	for _, b := range data {
		S += bits.OnesCount8(b)<<1 - 8 // S += (bits.OnesCount8(b) - (8 - bits.OnesCount8(b)))
	}
	V := float64(S) / math.Sqrt(float64(2*n))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// MonoBitFrequencyTest 单比特频数检测
//
// Deprecated: 请改用 MonoBitFrequencyTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 MonoBitFrequencyTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func MonoBitFrequencyTest(bits []bool) (float64, float64) {
	return coreMonoBitFrequency(BitSeqFromBools(bits))
}

// MonoBitFrequencyBitSeq 单比特频数检测
func MonoBitFrequencyBitSeq(s BitSeq) *TestResult {
	p, q := coreMonoBitFrequency(s)
	return &TestResult{Name: "单比特频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// MonoBitFrequencyTestBitSeq 单比特频数检测
func MonoBitFrequencyTestBitSeq(s BitSeq) (float64, float64) {
	return coreMonoBitFrequency(s)
}
