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

// OverlappingTemplateMatchingTest 重叠子序列检测方法
// bits: 检测序列
// return:
//      p1: P-value1
//      p2: P-value2
func OverlappingTemplateMatchingTest(bits []bool) (p1 float64, p2 float64) {
	n := len(bits)
	if n < 5 {
		fmt.Println("SerialTest:args wrong")
		return -1, -1
	}
	m := 5
	patterns1 := make([]int, 1<<m)
	patterns2 := make([]int, 1<<(m-1))
	patterns3 := make([]int, 1<<(m-2))
	var Phi1, Phi2, Phi3 float64 = 0, 0, 0
	var DPhi2, D2Phi2 float64 = 0, 0

	var mask1 int = (1 << m) - 1
	var mask2 int = (1 << (m - 1)) - 1
	var mask3 int = (1 << (m - 2)) - 1
	var tmp int = 0

	for i := 0; i < m-1; i++ {
		bits = append(bits, bits[i])
	}

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
		patterns1[tmp&mask1]++
		patterns2[tmp&mask2]++
		patterns3[tmp&mask3]++
	}

	for i := 0; i <= mask1; i++ {
		Phi1 += math.Pow(float64(patterns1[i]), 2.0)
	}
	Phi1 *= float64(mask1 + 1)
	Phi1 /= float64(n)
	Phi1 -= float64(n)
	for i := 0; i <= mask2; i++ {
		Phi2 += math.Pow(float64(patterns2[i]), 2.0)
	}
	Phi2 *= float64(mask2 + 1)
	Phi2 /= float64(n)
	Phi2 -= float64(n)
	for i := 0; i <= mask3; i++ {
		Phi3 += math.Pow(float64(patterns3[i]), 2.0)
	}
	Phi3 *= float64(mask3 + 1)
	Phi3 /= float64(n)
	Phi3 -= float64(n)

	DPhi2 = Phi1 - Phi2
	D2Phi2 = Phi1 - 2*Phi2 + Phi3

	_2m := 1 << m
	p1 = igamc(float64(_2m)/4.0, DPhi2/2.0)
	p2 = igamc(float64(_2m)/8.0, D2Phi2/2.0)

	//for i := 0; i < m-1; i++ {
	//	bits = bits[:len(bits)-1]
	//}
	return
}
