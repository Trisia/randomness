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

// BinaryDerivative 二元推导检测， k=7
func BinaryDerivative(data []byte) *TestResult {
	p, q := BinaryDerivativeTestBytes(data, 7)
	return &TestResult{Name: "二元推导检测(k=7)", P: p, Q: q, Pass: p >= Alpha}
}

// BinaryDerivativeTest 二元推导检测， k=7
func BinaryDerivativeTest(bits []bool, k int) (float64, float64) {
	return BinaryDerivativeProto(bits, k)
}

// BinaryDerivativeTestBytes 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
func BinaryDerivativeTestBytes(data []byte, k int) (float64, float64) {
	return BinaryDerivativeProto(B2bitArr(data), k)
}

// BinaryDerivativeProto 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
func BinaryDerivativeProto(bits []bool, k int) (float64, float64) {
	n := len(bits)
	if n < 7 {
		panic("please provide valid test bits")
	}

	S := 0
	var V float64 = 0
	_bits := make([]bool, len(bits))
	copy(_bits, bits)

	// Step 1, 2
	for i := 0; i < k; i++ {
		for j := 0; j < n-i-1; j++ {
			_bits[j] = xor(_bits[j], _bits[j+1])
		}
	}

	// Step 3
	for i := 0; i < n-k; i++ {
		if _bits[i] {
			S++
		} else {
			S--
		}
	}
	// Step 4
	V = float64(S) / math.Sqrt(float64(n-k))

	// Step 5
	P := math.Erfc(math.Abs(V) / math.Sqrt(2))

	// Step 6
	Q := math.Erfc(V/math.Sqrt(2)) / 2
	return P, Q
}
