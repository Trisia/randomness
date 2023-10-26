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
	"math"
)

// Runs 游程总数检测
func Runs(data []byte) *TestResult {
	p, q := RunsTestBytes(data)
	return &TestResult{Name: "游程总数检测", P: p, Q: q, Pass: p >= Alpha}
}

// RunsTestBytes 游程总数检测
func RunsTestBytes(data []byte) (float64, float64) {
	return RunsTest(B2bitArr(data))
}

// RunsTest 游程总数检测
func RunsTest(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	var Pi float64 = 0
	var V_obs int = 1
	var P, Q float64 = 0, 0

	// Step 1, 2
	for i := 0; i < n-1; i++ {
		if bits[i] != bits[i+1] {
			V_obs++
		}
		if bits[i] {
			Pi++
		}
	}
	if bits[n-1] {
		Pi++
	}
	Pi /= float64(n)

	// Step 3, 第四、五步的除math.Sqrt(2)，放到这里提前处理，减少math.Sqrt的调用。
	V := (float64(V_obs) - 2.0*float64(n)*Pi*(1.0-Pi)) / (2.0 * math.Sqrt(float64(2*n)) * Pi * (1.0 - Pi))

	// Step 4
	P = math.Erfc(math.Abs(V))

	// Step 5
	Q = math.Erfc(V) / 2.0
	return P, Q
}
