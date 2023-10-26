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
// 这里直接对字节直接处理，避免字节切片到位切片的转换，同时提高效率。
func MonoBitFrequencyTestBytes(data []byte) (float64, float64) {
	if len(data) == 0 {
		panic("please provide test bits")
	}
	n := len(data) * 8
	S := 0
	var V, P, Q float64

	for _, b := range data {
		S += bits.OnesCount8(b)<<1 - 8  // S += (bits.OnesCount8(b) - (8 - bits.OnesCount8(b)))
	}
	V = float64(S) / math.Sqrt(float64(2*n)) // 除math.Sqrt(2)，放到这里提前处理(n->2*n)，减少math.Sqrt的调用。
	P = math.Erfc(math.Abs(V))
	Q = math.Erfc(V) / 2
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
	V = float64(S) / math.Sqrt(float64(2*n)) // 除math.Sqrt(2)，放到这里提前处理(n->2*n)，减少math.Sqrt的调用。
	P = math.Erfc(math.Abs(V))
	Q = math.Erfc(V) / 2
	return P, Q
}
