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
	"math/cmplx"

	"github.com/Trisia/randomness/ttf"
)

// DiscreteFourierTransform 离散傅里叶检测
func DiscreteFourierTransform(data []byte) *TestResult {
	p, q := DiscreteFourierTransformTestBytes(data)
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBytes 离散傅里叶检测
func DiscreteFourierTransformTestBytes(data []byte) (float64, float64) {
	return DiscreteFourierTransformTest(B2bitArr(data))
}

// DiscreteFourierTransformTest 离散傅里叶检测
func DiscreteFourierTransformTest(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// Step 1, 2
	N := ceilPow2(n)
	rr := make([]complex128, N)
	for i := 0; i < n; i++ {
		if bits[i] {
			rr[i] = complex(1.0, 0)
		} else {
			rr[i] = complex(-1.0, 0)
		}
	}

	// 傅里叶变换
	f, err := ttf.New(N)
	if err != nil {
		panic(err)
	}
	f.Transform(rr)

	// Step 4
	T := math.Sqrt(2.995732274 * float64(n))

	// Step 5
	N_0 := 0.95 * float64(n) / 2

	// Step 6
	var N_1 int = 0
	for i := 0; i < n/2-1; i++ {
		if cmplx.Abs(rr[i]) < T {
			N_1++
		}
	}

	// Step 7
	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(n)/3.8)
	P := math.Erfc(math.Abs(V) / math.Sqrt(2.0))
	Q := math.Erfc(V/math.Sqrt(2.0)) / 2

	return P, Q
}
