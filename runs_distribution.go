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

import "fmt"

// RunsDistribution 游程分布检测
func RunsDistribution(data []byte) *TestResult {
	p := RunsDistributionTestBytes(data)
	return &TestResult{Name: "游程分布检测", P: p, Pass: p >= Alpha}
}

// RunsDistributionTestBytes 游程分布检测
func RunsDistributionTestBytes(data []byte) float64 {
	return RunsDistributionTest(B2bitArr(data))
}

// RunsDistributionTest 游程分布检测
func RunsDistributionTest(bits []bool) float64 {
	n := len(bits)
	if n < 100 {
		fmt.Println("RunsDistributionTest:args wrong")
		return -1
	}

	k := 0
	e := make([]float64, 50)
	b := make([]float64, 50)
	g := make([]float64, 50)
	var V float64 = 0
	var P float64 = 0
	var cur bool = bits[0]
	cnt := 0

	for {
		k++
		_2k2 := 1 << int(k+2)
		e[k] = float64(n-k+3) / float64(_2k2)
		if !(e[k] >= 5.0) {
			break
		}
	}
	k--
	bits = append(bits, !(bits[n-1]))

	for i := 0; i <= n; i++ {
		if bits[i] == cur {
			cnt++
		} else {
			if cnt <= k {
				if cur {
					b[cnt]++
				} else {
					g[cnt]++
				}
			}
			cur = bits[i]
			cnt = 1
		}
	}
	//bits.remove(bits.size() - 1);
	for i := 1; i <= k; i++ {
		V += (b[i] - e[i]) * (b[i] - e[i]) / e[i]
		V += (g[i] - e[i]) * (g[i] - e[i]) / e[i]
	}
	P = igamc(float64(k-1), V/2.0)
	return P
}
