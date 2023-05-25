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

var parameters = []struct {
	pi     []float64
	k      int
	m      int
	startV int
}{
	{
		pi:     []float64{0.2148, 0.3672, 0.2305, 0.1875},
		k:      3,
		m:      8,
		startV: 1,
	},
	{
		pi:     []float64{0.1174, 0.2430, 0.2494, 0.1752, 0.1027, 0.1124},
		k:      5,
		m:      128,
		startV: 4,
	},
	{
		pi:     []float64{0.086632, 0.208201, 0.248419, 0.193913, 0.121458, 0.068011, 0.073366},
		k:      6,
		m:      10000,
		startV: 10,
	},
}

func selectParameters(n int) int {
	switch {
	case n >= 750000:
		return 2
	case n >= 6272:
		return 1
	default:
		return 0
	}
}

// LongestRunOfOnesInABlock 块内最大游程检测,m=10000, k=6 for bits = 1000_000
func LongestRunOfOnesInABlock(data []byte) *TestResult {
	p, q := LongestRunOfOnesInABlockTestBytes(data, true)
	return &TestResult{Name: "块内最大游程检测", P: p, Q: q, Pass: p >= Alpha}
}

// LongestRunOfOnesInABlockTest 块内最大游程检测,m=10000 for bits = 1000_000
func LongestRunOfOnesInABlockTest(bits []bool, checkOne bool) (float64, float64) {
	return LongestRunOfOnesInABlockProto(bits, checkOne)
}

// LongestRunOfOnesInABlockTestBytes 块内最大游程检测
func LongestRunOfOnesInABlockTestBytes(data []byte, checkOne bool) (float64, float64) {
	return LongestRunOfOnesInABlockProto(B2bitArr(data), checkOne)
}

// LongestRunOfOnesInABlockProto 块内最大游程检测
// bits: 待检测序列
// m: m长度， m = 10000, k=6 for bits = 1000_000
func LongestRunOfOnesInABlockProto(bits []bool, checkOne bool) (float64, float64) {
	n := len(bits)

	if n < 128 {
		panic("please provide valid test bits")
	}

	param := parameters[selectParameters(n)]

	// Step 1
	N := n / param.m

	// Step 2
	v := make([]float64, param.k+1)
	var lr1, mlr1 int
	var b bool
	for i := 0; i < N; i++ {
		lr1 = 0
		mlr1 = 0

		for j := 0; j < param.m; j++ {
			b, bits = bits[0], bits[1:]
			if checkOne {
				if b {
					lr1++
					mlr1 = max(mlr1, lr1)
				} else {
					lr1 = 0
				}
			} else {
				if b {
					lr1 = 0
				} else {
					lr1++
					mlr1 = max(mlr1, lr1)
				}
			}
		}
		if mlr1 < param.startV {
			mlr1 = param.startV
		} else if mlr1 > param.startV+param.k {
			mlr1 = param.startV + param.k
		}
		v[mlr1-param.startV]++
	}

	// Step 3
	var V float64 = 0
	for i := 0; i < param.k+1; i++ {
		V += (v[i] - float64(N)*param.pi[i]) * (v[i] - float64(N)*param.pi[i]) / (float64(N) * param.pi[i])
	}

	// Step 4
	P := igamc(float64(param.k)/2.0, V/2.0)
	return P, P
}
