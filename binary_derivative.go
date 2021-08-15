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

// BinaryDerivativeTest 二元推导检测
func BinaryDerivativeTest(bits []bool) float64 {
	n := len(bits)
	if n < 7 {
		fmt.Println("BinaryDerivativeTest:args wrong")
		return -1
	}
	k := 7

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
