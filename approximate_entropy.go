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

// coreApproximateEntropy 是近似熵检测的实现（packed 内核）。
func coreApproximateEntropy(s BitSeq, m int) (float64, float64) {
	n := s.n
	if n == 0 {
		panic("please provide test bits")
	}
	if m < 1 {
		panic("please provide valid m (>= 1)")
	}
	if m >= n {
		panic("block size m must be less than sequence length")
	}

	numOfBlocks := float64(n)
	var ApEn [2]float64

	// Compute phi for blockSize=m and then blockSize=m+1.
	for blockSize := m; blockSize <= m+1; blockSize++ {
		powLen := 1 << uint(blockSize)
		mask := powLen - 1
		pattern := make([]int, powLen)

		// 统计起点 0..n-blockSize 的模式
		cur := 0
		for i := 0; i < n; i++ {
			cur = ((cur << 1) | int(s.bit(i))) & mask
			if i >= blockSize-1 {
				pattern[cur]++
			}
		}
		// 结尾环绕：统计起点 n-blockSize+1..n-1 的模式
		for i := 0; i < blockSize-1; i++ {
			cur = ((cur << 1) | int(s.bit(i))) & mask
			pattern[cur]++
		}

		// Compute the terms of the phi formula
		sum := 0.0
		for i := 0; i < powLen; i++ {
			if pattern[i] > 0 {
				f := float64(pattern[i])
				sum += f * math.Log(f/numOfBlocks)
			}
		}
		ApEn[blockSize-m] = sum / numOfBlocks
	}

	apen := ApEn[0] - ApEn[1]
	V := 2.0 * numOfBlocks * (math.Log(2) - apen)
	pow2m1 := 1 << uint(m-1)
	P := igamc(float64(pow2m1), V/2.0)
	return P, P
}

// ApproximateEntropy 近似熵检测,m=5
func ApproximateEntropy(data []byte) *TestResult {
	p, q := coreApproximateEntropy(BitSeqFromBytes(data), 5)
	return &TestResult{Name: "近似熵检测(m=5)", P: p, Q: q, Pass: p >= Alpha}
}

// ApproximateEntropyTest 近似熵检测,m=5
//
// Deprecated: 请改用 ApproximateEntropyTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 ApproximateEntropyTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func ApproximateEntropyTest(bits []bool) (float64, float64) {
	return coreApproximateEntropy(BitSeqFromBools(bits), 5)
}

// ApproximateEntropyTestBytes 近似熵检测
func ApproximateEntropyTestBytes(data []byte, m int) (float64, float64) {
	return coreApproximateEntropy(BitSeqFromBytes(data), m)
}

// ApproximateEntropyProto 近似熵检测, The purpose of the test is to compare the frequency of
// overlapping blocks of two consecutive/adjacent lengths (m and m+1) against the expected result for a
// random sequence. 这个实现参考自NIST的参考实现。
// Reference:
//
//	https://csrc.nist.gov/CSRC/media/Projects/Random-Bit-Generation/documents/sts-2_1_2.zip
//	https://github.com/arcetri/sts/blob/master/src/tests/approximateEntropy.c
//
// bits: 待检测序列
// m: m长度
//
// Deprecated: 请改用 ApproximateEntropyTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 ApproximateEntropyTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func ApproximateEntropyProto(bits []bool, m int) (float64, float64) {
	return coreApproximateEntropy(BitSeqFromBools(bits), m)
}

// ApproximateEntropyBitSeq 近似熵检测，m=5
func ApproximateEntropyBitSeq(s BitSeq) *TestResult {
	p, q := coreApproximateEntropy(s, 5)
	return &TestResult{Name: "近似熵检测(m=5)", P: p, Q: q, Pass: p >= Alpha}
}

// ApproximateEntropyTestBitSeq 近似熵检测
// m: m长度
func ApproximateEntropyTestBitSeq(s BitSeq, m int) (float64, float64) {
	return coreApproximateEntropy(s, m)
}
