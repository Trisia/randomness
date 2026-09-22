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

func mutFactorC(L, K int) float64 {
	var v float64
	v = 0.7
	v -= 0.8 / float64(L)
	v += (4.0 + 32.0/float64(L)) * (math.Pow(float64(K), -3.0/float64(L)) / 15.0)
	return v
}

// coreMaurerUniversal 是 Maurer 通用统计检测的实现（packed 内核）。
func coreMaurerUniversal(s BitSeq) (float64, float64) {
	n := s.n
	if n == 0 {
		panic("please provide test bits")
	}
	L := 7
	Q := 1280
	T := make([]int, 1<<uint(L))
	mask := (1 << uint(L)) - 1

	K := n/L - Q
	var sum float64 = 0.0
	expected_value := []float64{0, 0, 0, 0, 0, 0, 5.2177052, 6.1962507, 7.1836656,
		8.1764248, 9.1723243, 10.170032, 11.168765,
		12.168070, 13.167693, 14.167488, 15.167379}
	variance := []float64{0, 0, 0, 0, 0, 0, 2.954, 3.125, 3.238, 3.311, 3.356, 3.384,
		3.401, 3.410, 3.416, 3.419, 3.421}

	pos := 0
	next := func() int {
		v := 0
		for j := 0; j < L; j++ {
			v = v<<1 | int(s.bit(pos))
			pos++
		}
		return v
	}

	for i := 1; i <= Q; i++ {
		T[next()&mask] = i
	}
	for i := Q + 1; i <= Q+K; i++ {
		w := next() & mask
		sum += math.Log2(float64(i) - float64(T[w]))
		T[w] = i
	}

	sigma := math.Sqrt(variance[L]/float64(K)) * mutFactorC(L, K)
	// 避免求p q时V再除以math.Sqrt(2.0)
	V := (sum/float64(K) - expected_value[L]) / (sigma * math.Sqrt(2.0))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// MaurerUniversal Maurer通用统计检测方法
func MaurerUniversal(data []byte) *TestResult {
	p, q := coreMaurerUniversal(BitSeqFromBytes(data))
	return &TestResult{Name: "Maurer通用统计检测方法", P: p, Q: q, Pass: p >= Alpha}
}

// MaurerUniversalTestBytes Maurer通用统计检测方法
func MaurerUniversalTestBytes(data []byte) (float64, float64) {
	return coreMaurerUniversal(BitSeqFromBytes(data))
}

// MaurerUniversalTest Maurer通用统计检测方法
//
// Deprecated: 请改用 MaurerUniversalTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 MaurerUniversalTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func MaurerUniversalTest(bits []bool) (float64, float64) {
	return coreMaurerUniversal(BitSeqFromBools(bits))
}

// MaurerUniversalBitSeq Maurer通用统计检测方法
func MaurerUniversalBitSeq(s BitSeq) *TestResult {
	p, q := coreMaurerUniversal(s)
	return &TestResult{Name: "Maurer通用统计检测方法", P: p, Q: q, Pass: p >= Alpha}
}

// MaurerUniversalTestBitSeq Maurer通用统计检测方法
func MaurerUniversalTestBitSeq(s BitSeq) (float64, float64) {
	return coreMaurerUniversal(s)
}
