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

// ApproximateEntropyProto 近似熵检测
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

	for blockSize := m; blockSize <= m+1; blockSize++ {
		powLen := 1<<blockSize
		pattern = make([]int, powLen)
		for i := 0; i < n; i++ {
			k := 1
			for j := 0; j < blockSize; j++ {
				k <<= 1
				if bits[(i+j)%n] {
					k++
				}
			}
			pattern[k-powLen]++
		}
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
	_2mMinus1 := 1 << (m - 1)
	P = igamc(float64(_2mMinus1), V/2.0)
	return P, P
}
