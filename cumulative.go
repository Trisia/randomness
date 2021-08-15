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

// CumulativeTest 累加和检测
func CumulativeTest(bits []bool) float64 {
	n := len(bits)

	if n == 0 {
		fmt.Println("CumulativeTest:args wrong")
		return -1
	}
	var S int = 0
	var Z int = 0
	var P float64 = 1.0
	for i := 0; i < n; i++ {
		if bits[i] {
			S++
		} else {
			S--
		}
		Z = max(Z, abs(S))
	}
	_n := float64(n)
	for i := ((-n / Z) + 1) / 4; i <= ((n/Z)-1)/4; i++ {
		P -= normal_CDF(float64((4*i+1)*Z)/math.Sqrt(_n)) - normal_CDF(float64((4*i-1)*Z)/math.Sqrt(_n))
	}
	for i := ((-n / Z) - 3) / 4; i <= ((n/Z)-1)/4; i++ {
		P += normal_CDF(float64((4*i+3)*Z)/math.Sqrt(_n)) - normal_CDF(float64((4*i+1)*Z)/math.Sqrt(_n))
	}
	return P
}