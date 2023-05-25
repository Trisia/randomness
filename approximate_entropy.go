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
	if n == 0 {
		panic("please provide test bits")
	}
	bits2 := bits
	var pattern []int
	var mask int
	var tmp int = 0
	var Cjm float64
	var phim, phim1 float64 = 0, 0
	var ApEn float64
	var V float64
	var P float64

	// round1
	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}
	pattern = make([]int, 1<<m)
	mask = (1 << m) - 1

	var b bool
	for i := 0; i < m-1; i++ {
		tmp <<= 1
		b, bits = bits[0], bits[1:]
		if b {
			tmp++
		}
	}

	for i := 0; i < n; i++ {
		tmp <<= 1
		b, bits = bits[0], bits[1:]
		if b {
			tmp++
		}
		pattern[tmp&mask]++
	}

	for i := 0; i < (1 << m); i++ {
		Cjm = float64(pattern[i]) / float64(n)
		phim += Cjm * math.Log(Cjm)
	}

	// round2
	bits = bits2
	m++
	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}
	pattern = make([]int, 1<<m)
	mask = (1 << m) - 1

	for i := 0; i < m-1; i++ {
		tmp <<= 1
		b, bits = bits[0], bits[1:]
		if b {
			tmp++
		}
	}
	for i := 0; i < n; i++ {
		tmp <<= 1
		b, bits = bits[0], bits[1:]
		if b {
			tmp++
		}
		pattern[tmp&mask]++
	}

	for i := 0; i < (1 << m); i++ {
		Cjm = float64(pattern[i]) / float64(n)
		phim1 += Cjm * math.Log(Cjm)
	}
	// --
	m--
	ApEn = phim - phim1
	V = 2.0 * float64(n) * (math.Log(2) - ApEn)
	_2m := 1 << m
	P = igamc(float64(_2m)/2.0, V/2.0)
	return P, P
}
