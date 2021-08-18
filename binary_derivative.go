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
	"fmt"
	"math"
)

// BinaryDerivative 二元推导检测， k=7
func BinaryDerivative(data []byte) *TestResult {
	p := BinaryDerivativeTestBytes(data, 7)
	return &TestResult{Name: "二元推导检测", P: p, Pass: p >= Alpha}
}

// BinaryDerivativeTest 二元推导检测， k=7
func BinaryDerivativeTest(bits []bool) float64 {
	return BinaryDerivativeProto(bits, 7)
}

// BinaryDerivativeTestBytes 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
func BinaryDerivativeTestBytes(data []byte, k int) float64 {
	return BinaryDerivativeProto(B2bitArr(data), k)
}

// BinaryDerivativeProto 二元推导检测
// bits: 待检测序列
// k: 重复次数，k=3,7
func BinaryDerivativeProto(bits []bool, k int) float64 {
	n := len(bits)
	if n < 7 {
		fmt.Println("BinaryDerivativeTest:args wrong")
		return -1
	}

	S := 0
	var V float64 = 0
	var P float64 = 0
	_bits := make([]bool, len(bits))
	copy(_bits, bits)
	for i := 0; i < k; i++ {
		for j := 0; j < n-i-1; j++ {
			_bits[j] = xor(_bits[j], _bits[j+1])
		}
	}

	for i := 0; i < n-k; i++ {
		if _bits[i] {
			S++
		} else {
			S--
		}
	}
	V = math.Abs(float64(S)) / math.Sqrt(float64(n-k))
	P = math.Erfc(V / math.Sqrt(2))
	return P
}
