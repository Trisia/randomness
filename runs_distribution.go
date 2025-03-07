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

// RunsDistribution 游程分布检测
func RunsDistribution(data []byte) *TestResult {
	p, q := RunsDistributionTestBytes(data)
	return &TestResult{Name: "游程分布检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsDistributionTestBytes 游程分布检测
func RunsDistributionTestBytes(data []byte) (float64, float64) {
	return RunsDistributionTest(B2bitArr(data))
}

// RunsDistributionTest 游程分布检测
func RunsDistributionTest(bits []bool) (float64, float64) {
	n := len(bits)
	if n < 100 {
		panic("please provide valid test bits")
	}

	// Step 1, calculate k
	k := 0
	for {
		k++
		_2k2 := 1 << uint(k+2)
		if float64(n-k+3)/float64(_2k2) < 5.0 {
			break
		}
	}
	k--

	// Step 2
	e := make([]float64, k)
	b := make([]float64, k)
	g := make([]float64, k)
	var V float64 = 0
	var cur bool = bits[0]
	cnt := 0

	for i := 0; i < n; i++ {
		if bits[i] == cur {
			cnt++
		} else {
			if cnt > k {
				cnt = k
			}
			if cur {
				b[cnt-1]++
			} else {
				g[cnt-1]++
			}
			cur = bits[i]
			cnt = 1
		}
	}
	// 特殊处理结尾
	if cnt > k {
		cnt = k
	}
	if cur {
		b[cnt-1]++
	} else {
		g[cnt-1]++
	}

	// Step 3
	var T float64 = 0
	for i := 0; i < k; i++ {
		T += b[i] + g[i]
	}

	// Step 4
	for i := 0; i < k; i++ {
		if i < k-1 {
			e[i] = T / float64(int(1)<<uint(i+2))
		} else {
			e[i] = T / float64(int(1)<<uint(i+1))
		}
	}

	// Step 5
	for i := 0; i < k; i++ {
		V += (b[i] - e[i]) * (b[i] - e[i]) / e[i]
		V += (g[i] - e[i]) * (g[i] - e[i]) / e[i]
	}

	// Step 6
	P := igamc(float64(k-1), V/2.0)

	return P, P
}
