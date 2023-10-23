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

// MonoBitFrequency 单比特频数检测
func MonoBitFrequency(data []byte) *TestResult {
	p, q := MonoBitFrequencyTestBytes(data)
	return &TestResult{Name: "单比特频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// MonoBitFrequencyTestBytes 单比特频数检测
func MonoBitFrequencyTestBytes(data []byte) (float64, float64) {
	if len(data) == 0 {
		panic("please provide test bits")
	}
	n := len(data) * 8
	S := 0
	var V, P, Q float64

	for _, b := range data {
		S += bits.OnesCount8(b)<<1 - 8
	}
	V = float64(S) / math.Sqrt(float64(n))
	P = math.Erfc(math.Abs(V) / math.Sqrt(2))
	Q = math.Erfc(V/math.Sqrt(2)) / 2
	return P, Q
}

// MonoBitFrequencyTest 单比特频数检测
func MonoBitFrequencyTest(bits []bool) (float64, float64) {
	if len(bits) == 0 {
		panic("please provide test bits")
	}
	n := len(bits)
	S := 0
	var V, P, Q float64

	for _, bit := range bits {
		if bit {
			S++
		} else {
			S--
		}
	}
	V = float64(S) / math.Sqrt(float64(n))
	P = math.Erfc(math.Abs(V) / math.Sqrt(2))
	Q = math.Erfc(V/math.Sqrt(2)) / 2
	return P, Q
}
