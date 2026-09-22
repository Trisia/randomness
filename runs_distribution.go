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

// coreRunsDistribution 是游程分布检测的实现（packed 内核）。
func coreRunsDistribution(s BitSeq) (float64, float64) {
	n := s.n
	if n < 100 {
		panic("please provide valid test bits")
	}

	// Step 1, calculate k
	k := 0
	for {
		k++
		_2k2 := 1 << uint(k+2)
		if float64(n-k+3)/float64(_2k2) < 5.0 {
			break
		}
	}
	k--

	// Step 2
	b := make([]float64, k)
	g := make([]float64, k)

	prev := 0            // 当前游程的起始下标
	cur := s.bit(0) == 1 // 当前游程的位值
	// visit 处理位于位置 pos 的转移：即位 pos 与位 pos+1 不同，
	// 因此 [prev, pos] 构成一个完整游程。
	visit := func(pos int) {
		length := pos + 1 - prev
		if length > k {
			length = k
		}
		if cur {
			b[length-1]++
		} else {
			g[length-1]++
		}
		prev = pos + 1
		cur = !cur
	}

	W := len(s.u)
	for w := 0; w < W; w++ {
		m := n - w*64
		if m <= 0 {
			break
		}
		if m > 64 {
			m = 64
		}
		// 字内相邻位对的转移：位 j 与位 j+1 不同 <=> (u ^ (u>>1)) 的第 j 位
		if m >= 2 {
			t := (s.u[w] ^ (s.u[w] >> 1)) & maskLow(m-1)
			for t != 0 {
				j := bits.TrailingZeros64(t)
				t &= t - 1
				visit(w*64 + j)
			}
		}
		// 跨字转移：本字最高有效位 vs 下一字最低位
		if m == 64 && w+1 < W {
			if (s.u[w]>>63)^(s.u[w+1]&1) == 1 {
				visit(w*64 + 63)
			}
		}
	}
	// 最后一个游程（到 n-1 结束）
	length := n - prev
	if length > k {
		length = k
	}
	if cur {
		b[length-1]++
	} else {
		g[length-1]++
	}

	// Step 3
	var T float64 = 0
	for i := 0; i < k; i++ {
		T += b[i] + g[i]
	}

	// Step 4
	e := make([]float64, k)
	for i := 0; i < k; i++ {
		if i < k-1 {
			e[i] = T / float64(int(1)<<uint(i+2))
		} else {
			e[i] = T / float64(int(1)<<uint(i+1))
		}
	}

	// Step 5
	var V float64 = 0
	for i := 0; i < k; i++ {
		V += (b[i] - e[i]) * (b[i] - e[i]) / e[i]
		V += (g[i] - e[i]) * (g[i] - e[i]) / e[i]
	}

	// Step 6
	P := igamc(float64(k-1), V/2.0)
	return P, P
}

// RunsDistribution 游程分布检测
func RunsDistribution(data []byte) *TestResult {
	p, q := coreRunsDistribution(BitSeqFromBytes(data))
	return &TestResult{Name: "游程分布检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsDistributionTestBytes 游程分布检测
func RunsDistributionTestBytes(data []byte) (float64, float64) {
	return coreRunsDistribution(BitSeqFromBytes(data))
}

// RunsDistributionTest 游程分布检测
//
// Deprecated: 请改用 RunsDistributionTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 RunsDistributionTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func RunsDistributionTest(bits []bool) (float64, float64) {
	return coreRunsDistribution(BitSeqFromBools(bits))
}

var _ = math.Sqrt

// RunsDistributionBitSeq 游程分布检测
func RunsDistributionBitSeq(s BitSeq) *TestResult {
	p, q := coreRunsDistribution(s)
	return &TestResult{Name: "游程分布检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsDistributionTestBitSeq 游程分布检测
func RunsDistributionTestBitSeq(s BitSeq) (float64, float64) {
	return coreRunsDistribution(s)
}
