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
	"math/bits"
)

var parameters = []struct {
	pi     []float64
	k      int
	m      int
	startV int
}{
	{
		pi:     []float64{0.2148, 0.3672, 0.2305, 0.1875},
		k:      3,
		m:      8,
		startV: 1,
	},
	{
		pi:     []float64{0.1174, 0.2430, 0.2494, 0.1752, 0.1027, 0.1124},
		k:      5,
		m:      128,
		startV: 4,
	},
	{
		pi:     []float64{0.086632, 0.208201, 0.248419, 0.193913, 0.121458, 0.068011, 0.073366},
		k:      6,
		m:      10000,
		startV: 10,
	},
}

func selectParameters(n int) int {
	switch {
	case n >= 750000:
		return 2
	case n >= 6272:
		return 1
	default:
		return 0
	}
}

// longestRunOnesWord 返回一个字内最长连续 1 的长度。
func longestRunOnesWord(w uint64) int {
	c := 0
	for w != 0 {
		w &= w << 1
		c++
	}
	return c
}

// trailingOnes 返回从 bit 0 起连续 1 的个数（即可能从前一字续接过来的游程长度）。
func trailingOnes(w uint64) int { return bits.TrailingZeros64(^w) }

// leadingOnes 返回以 bit lim-1 结尾、向下连续的 1 的个数。
func leadingOnes(w uint64, lim int) int {
	c := 0
	for b := lim - 1; b >= 0; b-- {
		if (w>>uint(b))&1 == 0 {
			break
		}
		c++
	}
	return c
}

// longestRunBlock 返回 words 所表示序列（前 nbits 位）中的最长 1 游程
// （checkOne=false 时为最长 0 游程）。
func longestRunBlock(w []uint64, nbits int, checkOne bool) int {
	best, cur := 0, 0
	for i := 0; i < len(w); i++ {
		lim := nbits - i*64
		if lim <= 0 {
			break
		}
		if lim > 64 {
			lim = 64
		}
		mask := maskLow(lim)
		word := w[i] & mask
		if !checkOne {
			word = ^w[i] & mask
		}

		// 字内最长游程（不含跨字续接）
		if r := longestRunOnesWord(word); r > best {
			best = r
		}
		// 与上一字续接
		if word&1 == 1 && cur > 0 {
			if r := cur + trailingOnes(word); r > best {
				best = r
			}
		}
		// 更新「以本字最高有效位结尾的游程长度」，供下一字续接
		if word == mask {
			cur += lim
		} else if (word>>uint(lim-1))&1 == 1 {
			cur = leadingOnes(word, lim)
		} else {
			cur = 0
		}
		if cur > best {
			best = cur
		}
	}
	return best
}

// coreLongestRunOfOnesInABlock 是块内最大游程检测的实现（packed 内核）。
func coreLongestRunOfOnesInABlock(s BitSeq, checkOne bool) (float64, float64) {
	n := s.n
	if n < 128 {
		panic("please provide valid test bits")
	}
	param := parameters[selectParameters(n)]

	// Step 1
	N := n / param.m

	// Step 2
	blockWords := make([]uint64, (param.m+63)/64)
	v := make([]float64, param.k+1)
	for i := 0; i < N; i++ {
		extractBits(s, i*param.m, param.m, blockWords)
		mlr1 := longestRunBlock(blockWords, param.m, checkOne)
		if mlr1 < param.startV {
			mlr1 = param.startV
		} else if mlr1 > param.startV+param.k {
			mlr1 = param.startV + param.k
		}
		v[mlr1-param.startV]++
	}

	// Step 3
	var V float64 = 0
	NF := float64(N)
	for i := 0; i < param.k+1; i++ {
		V += (v[i] - NF*param.pi[i]) * (v[i] - NF*param.pi[i]) / (NF * param.pi[i])
	}
	// Step 4
	P := igamc(float64(param.k)/2.0, V/2.0)
	return P, P
}

// LongestRunOfOnesInABlock 块内最大游程检测,m=10000, k=6 for bits = 1000_000
func LongestRunOfOnesInABlock(data []byte) *TestResult {
	p, q := coreLongestRunOfOnesInABlock(BitSeqFromBytes(data), true)
	return &TestResult{Name: "块内最大游程检测", P: p, Q: q, Pass: p >= Alpha}
}

// LongestRunOfOnesInABlockTest 块内最大游程检测,m=10000 for bits = 1000_000
//
// Deprecated: 请改用 LongestRunOfOnesInABlockTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LongestRunOfOnesInABlockTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LongestRunOfOnesInABlockTest(bits []bool, checkOne bool) (float64, float64) {
	return coreLongestRunOfOnesInABlock(BitSeqFromBools(bits), checkOne)
}

// LongestRunOfOnesInABlockTestBytes 块内最大游程检测
func LongestRunOfOnesInABlockTestBytes(data []byte, checkOne bool) (float64, float64) {
	return coreLongestRunOfOnesInABlock(BitSeqFromBytes(data), checkOne)
}

// LongestRunOfOnesInABlockProto 块内最大游程检测
// bits: 待检测序列
// m: m长度， m = 10000, k=6 for bits = 1000_000
//
// Deprecated: 请改用 LongestRunOfOnesInABlockTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LongestRunOfOnesInABlockTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LongestRunOfOnesInABlockProto(bits []bool, checkOne bool) (float64, float64) {
	return coreLongestRunOfOnesInABlock(BitSeqFromBools(bits), checkOne)
}

// LongestRunOfOnesInABlockBitSeq 块内最大游程检测（块长按数据规模自动选择）
func LongestRunOfOnesInABlockBitSeq(s BitSeq) *TestResult {
	p, q := coreLongestRunOfOnesInABlock(s, true)
	return &TestResult{Name: "块内最大游程检测", P: p, Q: q, Pass: p >= Alpha}
}

// LongestRunOfOnesInABlockTestBitSeq 块内最大游程检测
// checkOne: true 统计最长 1 游程，false 统计最长 0 游程
func LongestRunOfOnesInABlockTestBitSeq(s BitSeq, checkOne bool) (float64, float64) {
	return coreLongestRunOfOnesInABlock(s, checkOne)
}
