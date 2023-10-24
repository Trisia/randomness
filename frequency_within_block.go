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

// FrequencyWithinBlock 块内频数检测, m = 10000 for bits = 1000_000
func FrequencyWithinBlock(data []byte) *TestResult {
	p, q := FrequencyWithinBlockTest(B2bitArr(data))
	return &TestResult{Name: "块内频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// FrequencyWithinBlockTest 块内频数检测, m = 10000 for bits = 1000_000
func FrequencyWithinBlockTest(bits []bool) (float64, float64) {
	return FrequencyWithinBlockProto(bits, selectM(len(bits)))
}

// FrequencyWithinBlockTestBytes 块内频数检测
func FrequencyWithinBlockTestBytes(data []byte, m int) (float64, float64) {
	return FrequencyWithinBlockProto(B2bitArr(data), m)
}

func selectM(n int) int {
	var m int
	switch {
	case n >= 100000000:
		m = 1000000
	case n >= 1000000:
		m = 10000
	case n >= 10000:
		m = 1000
	case n >= 1000:
		m = 100
	default:
		m = 10
	}
	return m
}

// FrequencyWithinBlockProto 块内频数检测
func FrequencyWithinBlockProto(bits []bool, m int) (float64, float64) {
	n := len(bits)
	N := n / m
	if N == 0 {
		panic("please provide test bits")
	}
	bits = bits[:N*m]

	var Pi float64 = 0
	var V float64 = 0
	var P float64 = 0

	var b bool
	for i := 0; i < N; i++ {
		Pi = 0
		for j := 0; j < m; j++ {
			b, bits = bits[0], bits[1:]
			if b {
				Pi++
			}
		}
		Pi = Pi / float64(m)
		Pi = Pi - 0.5
		V += Pi * Pi
	}
	V *= 2.0 * float64(m)
	P = igamc(float64(N)/2.0, V)
	return P, P
}
