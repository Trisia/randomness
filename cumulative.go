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
)

// coreCumulative 是累加和检测的实现（packed 内核）。
func coreCumulative(s BitSeq, forward bool) (float64, float64) {
	n := s.n
	if n == 0 {
		panic("please provide test bits")
	}

	S := 0
	Z := 0
	step := func(one bool) {
		if one {
			S++
		} else {
			S--
		}
		if S > Z {
			Z = S
		} else if -S > Z {
			Z = -S
		}
	}

	if forward {
		for w := 0; w < len(s.u); w++ {
			word := s.u[w]
			lim := n - w*64
			if lim > 64 {
				lim = 64
			}
			for b := 0; b < lim; b++ {
				step(word&1 == 1)
				word >>= 1
			}
		}
	} else {
		W := len(s.u)
		for w := W - 1; w >= 0; w-- {
			hi := 64
			if w == W-1 {
				hi = n - w*64
			}
			word := s.u[w]
			for b := hi - 1; b >= 0; b-- {
				step((word>>uint(b))&1 == 1)
			}
		}
	}

	sqrtN := math.Sqrt(float64(n)) // 提前求平方根，避免下面多次求平方根
	P := 1.0
	for i := ((-n / Z) + 1) / 4; i <= ((n/Z)-1)/4; i++ {
		P -= normal_CDF(float64((4*i+1)*Z)/sqrtN) - normal_CDF(float64((4*i-1)*Z)/sqrtN)
	}
	for i := ((-n / Z) - 3) / 4; i <= ((n/Z)-1)/4; i++ {
		P += normal_CDF(float64((4*i+3)*Z)/sqrtN) - normal_CDF(float64((4*i+1)*Z)/sqrtN)
	}
	return P, P
}

// Cumulative 累加和检测
func Cumulative(data []byte) *TestResult {
	p, q := coreCumulative(BitSeqFromBytes(data), true)
	return &TestResult{Name: "累加和检测", P: p, Q: q, Pass: p >= Alpha}
}

// CumulativeTestBytes 累加和检测
// forward: true 前向, false 后向
func CumulativeTestBytes(data []byte, forward bool) (float64, float64) {
	return coreCumulative(BitSeqFromBytes(data), forward)
}

// CumulativeTest 累加和检测
// forward: true 前向, false 后向
//
// Deprecated: 请改用 CumulativeTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 CumulativeTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func CumulativeTest(bits []bool, forward bool) (float64, float64) {
	return coreCumulative(BitSeqFromBools(bits), forward)
}

// CumulativeBitSeq 累加和检测（前向）
func CumulativeBitSeq(s BitSeq) *TestResult {
	p, q := coreCumulative(s, true)
	return &TestResult{Name: "累加和检测", P: p, Q: q, Pass: p >= Alpha}
}

// CumulativeTestBitSeq 累加和检测
// forward: true 前向, false 后向
func CumulativeTestBitSeq(s BitSeq, forward bool) (float64, float64) {
	return coreCumulative(s, forward)
}
