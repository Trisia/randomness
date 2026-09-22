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

// coreRuns 是游程总数检测的实现（packed 内核）。
func coreRuns(s BitSeq) (float64, float64) {
	n := s.n
	if n == 0 {
		panic("please provide test bits")
	}
	ones := s.ones()
	trans := 0
	W := len(s.u)
	for w := 0; w < W; w++ {
		m := n - w*64
		if m <= 0 {
			break
		}
		if m > 64 {
			m = 64
		}
		if m >= 2 {
			trans += bits.OnesCount64((s.u[w] ^ (s.u[w] >> 1)) & maskLow(m-1))
		}
		if m == 64 && w+1 < W {
			trans += int((s.u[w] >> 63) ^ (s.u[w+1] & 1))
		}
	}

	Pi := float64(ones) / float64(n)
	// Step 3, 第四、五步的除math.Sqrt(2)，放到这里提前处理，减少math.Sqrt的调用。
	V_obs := float64(1 + trans)
	V := (V_obs - 2.0*float64(n)*Pi*(1.0-Pi)) / (2.0 * math.Sqrt(float64(2*n)) * Pi * (1.0 - Pi))
	// Step 4, 5
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2.0
}

// Runs 游程总数检测
func Runs(data []byte) *TestResult {
	p, q := coreRuns(BitSeqFromBytes(data))
	return &TestResult{Name: "游程总数检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsTestBytes 游程总数检测
func RunsTestBytes(data []byte) (float64, float64) {
	return coreRuns(BitSeqFromBytes(data))
}

// RunsTest 游程总数检测
//
// Deprecated: 请改用 RunsTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 RunsTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func RunsTest(bits []bool) (float64, float64) {
	return coreRuns(BitSeqFromBools(bits))
}

// RunsBitSeq 游程总数检测
func RunsBitSeq(s BitSeq) *TestResult {
	p, q := coreRuns(s)
	return &TestResult{Name: "游程总数检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsTestBitSeq 游程总数检测
func RunsTestBitSeq(s BitSeq) (float64, float64) {
	return coreRuns(s)
}
