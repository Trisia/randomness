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
	"math/cmplx"
	"randomness/ttf"
)

// DiscreteFourierTransformTest 离散傅里叶检测
func DiscreteFourierTransformTest(bits []bool) float64 {
	n := len(bits)
	if n == 0 {
		fmt.Println("DiscreteFourierTransformTest:args wrong")
		return -1
	}

	r := make([]float64, n)
	T := math.Sqrt(2.995732274 * float64(n))
	N_0 := 0.95 * float64(n) / 2
	var N_1 int = 0

	for i := 0; i < n; i++ {
		if bits[i] {
			r[i] = 1.0
		} else {
			r[i] = -1.0
		}
	}
	r = pow2DoubleArr(r)
	rr := make([]complex128, len(r))
	for i := range r {
		rr[i] = complex(r[i], 0)
	}

	// 傅里叶变换
	f, err := ttf.New(len(r))
	if err != nil {
		panic(err)
	}
	result := f.Transform(rr)
	//FastFourierTransformer fft = new FastFourierTransformer(DftNormalization.STANDARD);
	//Complex[] result = fft.transform(r, TransformType.FORWARD);
	for i := 0; i < n/2-1; i++ {
		if cmplx.Abs(result[i]) < T {
			N_1++
		}
	}

	if math.Abs(r[0]) < T {
		N_1++
	}

	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(n)/2.0)
	P := math.Erfc(math.Abs(V) / math.Sqrt(2.0))

	return P
}
