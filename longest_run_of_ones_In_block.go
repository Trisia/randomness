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
)

var pi = []float64{0.0882, 0.2092, 0.2483, 0.1933, 0.1208, 0.0675, 0.0727}

// LongestRunOfOnesInABlockTest 块内最大“1”游程检测
func LongestRunOfOnesInABlockTest(bits []bool) float64 {
	n := len(bits)

	if n < 10000 {
		fmt.Println("testForTheLongestRunOfOnesInABlock:args wrong")
		return -1
	}
	m := 10000
	N := n / m
	v := make([]float64, 7)
	var V float64 = 0
	var P float64 = 9

	var lr1, mlr1 int
	var b bool
	for i := 0; i < N; i++ {
		lr1 = 0
		mlr1 = 0
		for j := 0; j < m; j++ {
			b, bits = bits[0], bits[1:]
			if b {
				lr1++
				mlr1 = max(mlr1, lr1)
			} else {
				lr1 = 0
			}
		}
		if mlr1 <= 10 {
			v[0]++
		}
		if mlr1 >= 16 {
			v[6]++
		}
		if 10 < mlr1 && mlr1 < 16 {
			v[mlr1-10]++
		}
	}

	for i := 0; i < 7; i++ {
		V += (v[i] - float64(N)*pi[i]) * (v[i] - float64(N)*pi[i]) / (float64(N) * pi[i])
	}
	P = igamc(3, V/2.0)
	return P
}
