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

// coreAutocorrelation 是自相关检测的实现（packed 内核）。
func coreAutocorrelation(s BitSeq, d int) (float64, float64) {
	n := s.n
	if n < 16 {
		panic("please provide valid test bits")
	}
	cnt := n - d
	total := 0
	// d >= n 时循环不执行，随后按同样的公式计算。
	if cnt > 0 {
		sh, sb := d>>6, uint(d&63)
		for w := 0; w*64 < cnt; w++ {
			var shifted uint64
			if sb == 0 {
				shifted = s.word(w + sh)
			} else {
				shifted = s.word(w+sh)>>sb | s.word(w+sh+1)<<(64-sb)
			}
			total += bits.OnesCount64((s.u[w] ^ shifted) & maskLow(cnt-w*64))
		}
	}

	// V 除以 2 的平方根。
	V := 2.0 * (float64(total) - (float64(cnt) / 2.0)) / math.Sqrt(2*float64(cnt))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// Autocorrelation 自相关检测,d=16
func Autocorrelation(data []byte) *TestResult {
	p, q := coreAutocorrelation(BitSeqFromBytes(data), 16)
	return &TestResult{Name: "自相关检测(d=16)", P: p, Q: q, Pass: p >= Alpha}
}

// AutocorrelationTest 自相关检测,d=16
//
// Deprecated: 请改用 AutocorrelationTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 AutocorrelationTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func AutocorrelationTest(bits []bool, d int) (float64, float64) {
	return coreAutocorrelation(BitSeqFromBools(bits), d)
}

// AutocorrelationTestBytes 自相关检测
// data: 待检测序列
// d: d=1,2,8,16
func AutocorrelationTestBytes(data []byte, d int) (float64, float64) {
	return coreAutocorrelation(BitSeqFromBytes(data), d)
}

// AutocorrelationProto 自相关检测
// bits: 待检测序列
// d: d=1,2,8,16
//
// Deprecated: 请改用 AutocorrelationTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 AutocorrelationTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func AutocorrelationProto(bits []bool, d int) (float64, float64) {
	return coreAutocorrelation(BitSeqFromBools(bits), d)
}

// AutocorrelationBitSeq 自相关检测，d=16
func AutocorrelationBitSeq(s BitSeq) *TestResult {
	p, q := coreAutocorrelation(s, 16)
	return &TestResult{Name: "自相关检测(d=16)", P: p, Q: q, Pass: p >= Alpha}
}

// AutocorrelationTestBitSeq 自相关检测
// d: d=1,2,8,16
func AutocorrelationTestBitSeq(s BitSeq, d int) (float64, float64) {
	return coreAutocorrelation(s, d)
}
