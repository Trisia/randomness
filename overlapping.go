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

// coreOverlapping 是重叠子序列检测的实现（packed 内核）。
func coreOverlapping(s BitSeq, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	n := s.n
	if n < 5 {
		panic("please provide valid test bits")
	}
	if m < 2 || m > 30 {
		panic("please provide valid m (2..30)")
	}
	size1 := 1 << uint(m)
	mask1 := size1 - 1
	patterns1 := make([]int, size1)

	// Step 2：本来这里需要预先在结尾插入 bits[:m-1] 使长度仍为 n；
	// 现在改成主体 + 尾循环两段，避免逐位取模。
	tmp := 0
	for j := 0; j < m-1; j++ {
		tmp = (tmp << 1) | int(s.bit(j))
	}
	for i := m - 1; i < n; i++ {
		tmp = ((tmp << 1) | int(s.bit(i))) & mask1
		patterns1[tmp]++
	}
	for j := 0; j < m-1; j++ {
		tmp = ((tmp << 1) | int(s.bit(j))) & mask1
		patterns1[tmp]++
	}

	// patterns2/patterns3 由 patterns1 的边缘和导出
	size2 := size1 >> 1
	patterns2 := make([]int, size2)
	for x := 0; x < size2; x++ {
		patterns2[x] = patterns1[x] + patterns1[x+size2]
	}
	size3 := size1 >> 2
	patterns3 := make([]int, size3)
	for x := 0; x < size3; x++ {
		patterns3[x] = patterns2[x] + patterns2[x+size3]
	}

	// Step 3
	fn := float64(n)
	phi := func(pat []int, size int) float64 {
		acc := 0.0
		for i := 0; i < size; i++ {
			acc += float64(pat[i]) * float64(pat[i])
		}
		return acc*float64(size)/fn - fn
	}
	phi1 := phi(patterns1, size1)
	phi2 := phi(patterns2, size2)
	phi3 := phi(patterns3, size3)

	// Step 4
	DPhi2 := phi1 - phi2
	D2Phi2 := phi1 - 2*phi2 + phi3

	// Step 5, 6
	p1 = igamc(float64(size3), DPhi2/2.0)
	p2 = igamc(float64(size3)/2.0, D2Phi2/2.0)
	q1 = p1
	q2 = p2
	return
}

// OverlappingTemplateMatching 重叠子序列检测方法,m=5
func OverlappingTemplateMatching(data []byte) *TestResult {
	p1, p2, q1, q2 := coreOverlapping(BitSeqFromBytes(data), 5)
	return &TestResult{
		Name: "重叠子序列检测方法",
		P:    p1, P2: p2,
		Q: q1, Q2: q2,
		Pass: math.Min(p1, p2) >= Alpha,
	}
}

// OverlappingTemplateMatchingTest 重叠子序列检测方法,m=5
// bits: 检测序列
// return:
//
//	p1: P-value1
//	p2: P-value2
//
// Deprecated: 请改用 OverlappingTemplateMatchingTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 OverlappingTemplateMatchingTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func OverlappingTemplateMatchingTest(bits []bool) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return coreOverlapping(BitSeqFromBools(bits), 5)
}

// OverlappingTemplateMatchingTestBytes 重叠子序列检测方法
// data: 检测序列
// m: m长度,m=2,5
// return:
//
//	p1: P-value1
//	p2: P-value2
func OverlappingTemplateMatchingTestBytes(data []byte, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return coreOverlapping(BitSeqFromBytes(data), m)
}

// OverlappingTemplateMatchingProto 重叠子序列检测方法
// bits: 检测序列
// m: m长度,m=3,5
// return:
//
//	p1: P-value1
//	p2: P-value2
//
// Deprecated: 请改用 OverlappingTemplateMatchingTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 OverlappingTemplateMatchingTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func OverlappingTemplateMatchingProto(bits []bool, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return coreOverlapping(BitSeqFromBools(bits), m)
}

// OverlappingTemplateMatchingBitSeq 重叠子序列检测方法,m=5
func OverlappingTemplateMatchingBitSeq(s BitSeq) *TestResult {
	p1, p2, q1, q2 := coreOverlapping(s, 5)
	return &TestResult{
		Name: "重叠子序列检测方法",
		P:    p1, P2: p2,
		Q: q1, Q2: q2,
		Pass: math.Min(p1, p2) >= Alpha,
	}
}

// OverlappingTemplateMatchingTestBitSeq 重叠子序列检测方法
// m: m长度,m=2,3,5,7
// return:
//
//	p1: P-value1
//	p2: P-value2
func OverlappingTemplateMatchingTestBitSeq(s BitSeq, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return coreOverlapping(s, m)
}
