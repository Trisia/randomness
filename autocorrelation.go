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

// Autocorrelation 自相关检测,d=16
func Autocorrelation(data []byte) *TestResult {
	p, q := AutocorrelationTestBytes(data, 16)
	return &TestResult{Name: "自相关检测(d=16)", P: p, Q: q, Pass: p >= Alpha}
}

// AutocorrelationTest 自相关检测,d=16
func AutocorrelationTest(bits []bool, d int) (float64, float64) {
	return AutocorrelationProto(bits, d)
}

// AutocorrelationTestBytes 自相关检测
// data: 待检测序列
// d: d=1,2,8,16
func AutocorrelationTestBytes(data []byte, d int) (float64, float64) {
	return AutocorrelationProto(B2bitArr(data), d)
}

// AutocorrelationProto 自相关检测
// bits: 待检测序列
// d: d=1,2,8,16
func AutocorrelationProto(bits []bool, d int) (float64, float64) {
	n := len(bits)
	if n < 16 {
		panic("please provide valid test bits")
	}

	Ad := 0
	var V float64 = 0

	for i := 0; i < n-d; i++ {
		if xor(bits[i], bits[i+d]) {
			Ad++
		}
	}

	V = 2.0 * (float64(Ad) - (float64(n-d) / 2.0)) / math.Sqrt(2*float64(n-d)) // 提前对V除以2的平方根，避免求P Q时再求解
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2
	return P, Q
}
