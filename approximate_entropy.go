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

// ApproximateEntropy 近似熵检测,m=5
func ApproximateEntropy(data []byte) *TestResult {
	p, q := ApproximateEntropyTestBytes(data, 5)
	return &TestResult{Name: "近似熵检测(m=5)", P: p, Q: q, Pass: p >= Alpha}
}

// ApproximateEntropyTest 近似熵检测,m=5
func ApproximateEntropyTest(bits []bool) (float64, float64) {
	return ApproximateEntropyProto(bits, 5)
}

// ApproximateEntropyTestBytes 近似熵检测
func ApproximateEntropyTestBytes(data []byte, m int) (float64, float64) {
	return ApproximateEntropyProto(B2bitArr(data), m)
}

// ApproximateEntropyProto 近似熵检测, The purpose of the test is to compare the frequency of
// overlapping blocks of two consecutive/adjacent lengths (m and m+1) against the expected result for a
// random sequence. 这个实现参考自NIST的参考实现。
// Reference:
//   https://csrc.nist.gov/CSRC/media/Projects/Random-Bit-Generation/documents/sts-2_1_2.zip
//   https://github.com/arcetri/sts/blob/master/src/tests/approximateEntropy.c
//
// bits: 待检测序列
// m: m长度
func ApproximateEntropyProto(bits []bool, m int) (float64, float64) {
	n := len(bits)
	numOfBlocks := float64(n)
	if n == 0 {
		panic("please provide test bits")
	}
	var pattern []int
	var ApEn [2]float64
	var V float64
	var P float64
	r := 0

	// Compute phi for blockSize=m and then blockSize=m+1.
	// 初始实现中，按照《GM/T 0005-2021 随机性检测规范》，第一步要构造新的位序列：添加最开始的blockSize-1位数据到结尾，
	// 目前的实现中，这一步被省去了。
	for blockSize := m; blockSize <= m+1; blockSize++ {
		// Compute how many counters are needed, i.e. how many different possible m-bit sub-sequences can possibly exist.
		powLen := 1 << uint(blockSize)

		pattern = make([]int, powLen)
		// Compute the frequency of all the overlapping sub-sequences
		// 这里的算法也可以采用重叠子序列检测方法实现的方式。
		for i := 0; i < n; i++ {
			k := 1
			for j := 0; j < blockSize; j++ {
				k <<= 1
				if bits[(i+j)%n] { // (i+j) % n is used to avoid appending blockSize-1 bits in the end
					k++
				}
			}
			pattern[k-powLen]++
		}
		// Compute the the terms of the phi formula
		sum := float64(0.0)
		for i := 0; i < powLen; i++ {
			if pattern[i] > 0 {
				sum += float64(pattern[i]) * math.Log(float64(pattern[i])/numOfBlocks)
			}
		}
		sum /= numOfBlocks
		ApEn[r] = sum
		r++
	}
	apen := ApEn[0] - ApEn[1]
	V = 2.0 * numOfBlocks * (math.Log(2) - apen)
	_2mMinus1 := 1 << uint(m-1)
	P = igamc(float64(_2mMinus1), V/2.0)
	return P, P
}
