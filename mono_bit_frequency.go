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

// MonoBitFrequencyTest 单比特频数检测
func MonoBitFrequencyTest(bits []bool) float64 {
	if len(bits) == 0 {
		fmt.Println("MonoBitFrequencyTest:arg wrong")
		return -1
	}
	n := len(bits)
	S := 0
	var V float64
	var P float64
	for _, bit := range bits {
		if bit {
			S++
		} else {
			S--
		}
	}
	V = math.Abs(float64(S)) / math.Sqrt(float64(n))
	P = math.Erfc(V / math.Sqrt(2))
	return P
}
