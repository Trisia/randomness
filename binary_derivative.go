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

// coreBinaryDerivative 是二元推导检测的实现（packed 内核）。
func coreBinaryDerivative(s BitSeq, k int) (float64, float64) {
	n := s.n
	if n < 7 {
		panic("please provide valid test bits")
	}
	if k < 1 {
		panic("please provide valid k (>= 1)")
	}
	w := make([]uint64, len(s.u))
	copy(w, s.u)
	for p := 0; p < k; p++ {
		for j := range w {
			sh := w[j] >> 1
			if j+1 < len(w) {
				sh |= w[j+1] << 63
			}
			w[j] ^= sh
		}
	}

	// Step 3：只统计前 n-k 位
	cnt := n - k
	ones := 0
	for j := 0; j*64 < cnt; j++ {
		ones += bits.OnesCount64(w[j] & maskLow(cnt-j*64))
	}
	S := 2*ones - cnt

	// Step 4, 计算 V
	V := float64(S) / math.Sqrt(2*float64(cnt))
	// Step 5, 6
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// BinaryDerivative 二元推导检测， k=7
func BinaryDerivative(data []byte) *TestResult {
	p, q := coreBinaryDerivative(BitSeqFromBytes(data), 7)
	return &TestResult{Name: "二元推导检测(k=7)", P: p, Q: q, Pass: p >= Alpha}
}

// BinaryDerivativeTest 二元推导检测， k=7
//
// Deprecated: 请改用 BinaryDerivativeTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 BinaryDerivativeTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func BinaryDerivativeTest(bits []bool, k int) (float64, float64) {
	return coreBinaryDerivative(BitSeqFromBools(bits), k)
}

// BinaryDerivativeTestBytes 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
func BinaryDerivativeTestBytes(data []byte, k int) (float64, float64) {
	return coreBinaryDerivative(BitSeqFromBytes(data), k)
}

// BinaryDerivativeProto 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
//
// Deprecated: 请改用 BinaryDerivativeTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 BinaryDerivativeTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func BinaryDerivativeProto(bits []bool, k int) (float64, float64) {
	return coreBinaryDerivative(BitSeqFromBools(bits), k)
}

// BinaryDerivativeBitSeq 二元推导检测，k=7
func BinaryDerivativeBitSeq(s BitSeq) *TestResult {
	p, q := coreBinaryDerivative(s, 7)
	return &TestResult{Name: "二元推导检测(k=7)", P: p, Q: q, Pass: p >= Alpha}
}

// BinaryDerivativeTestBitSeq 二元推导检测
// k: 重复次数，k=3,7
func BinaryDerivativeTestBitSeq(s BitSeq, k int) (float64, float64) {
	return coreBinaryDerivative(s, k)
}
