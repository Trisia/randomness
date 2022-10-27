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

// LinearComplexity 线型复杂度检测,m=500
func LinearComplexity(data []byte) *TestResult {
	p, q := LinearComplexityTestBytes(data, 500)
	return &TestResult{Name: "线型复杂度检测(m=500)", P: p, Q: q, Pass: p >= Alpha}
}

// LinearComplexityTest 线型复杂度检测,m=500
func LinearComplexityTest(bits []bool) (float64, float64) {
	return LinearComplexityProto(bits, 500)
}

// LinearComplexityTestBytes 线型复杂度检测
// data: 待检测序列
// m: m长度
func LinearComplexityTestBytes(data []byte, m int) (float64, float64) {
	return LinearComplexityProto(B2bitArr(data), m)
}

// LinearComplexityProto 线型复杂度检测
// bits: 待检测序列
// m: m长度
func LinearComplexityProto(bits []bool, m int) (float64, float64) {
	n := len(bits)
	N := n / m
	if N == 0 {
		panic("please provide valid test bits")
	}

	var v = []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0}
	var pi = []float64{0.010417, 0.03125, 0.12500, 0.5000, 0.25000, 0.06250, 0.020833}
	var V float64 = 0.0
	var P float64 = 0

	var arr []bool
	var complexity int
	var T float64
	arr = make([]bool, m)

	// Step 3, miu
	var _1_m float64
	if m%2 == 0 {
		_1_m = 1.0
	} else {
		_1_m = -1.0
	}
	miu := float64(m)/2.0 + (9.0+_1_m)/36.0 - (float64(m)/3.0+2.0/9.0)/math.Pow(2.0, float64(m))

	// Step 2, 4, 5
	for i := 0; i < N; i++ {
		for j := 0; j < m; j++ {
			arr[j], bits = bits[0], bits[1:]
		}
		complexity = linearComplexity(arr, m)
		if m%2 == 0 {
			_1_m = 1.0
		} else {
			_1_m = -1.0
		}
		T = _1_m*(float64(complexity)-miu) + 2.0/9.0
		if T <= -2.5 {
			v[0]++
		} else if T <= -1.5 {
			v[1]++
		} else if T <= -0.5 {
			v[2]++
		} else if T <= 0.5 {
			v[3]++
		} else if T <= 1.5 {
			v[4]++
		} else if T <= 2.5 {
			v[5]++
		} else {
			v[6]++
		}
	}

	// Step 6
	for i := 0; i < 7; i++ {
		V += math.Pow(v[i]-float64(N)*pi[i], 2.0) / (float64(N) * pi[i])
	}

	// Step 7
	P = igamc(3.0, V/2.0)

	return P, P
}
