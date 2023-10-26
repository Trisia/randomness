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

// OverlappingTemplateMatching 重叠子序列检测方法,m=5
func OverlappingTemplateMatching(data []byte) *TestResult {
	p1, p2, q1, q2 := OverlappingTemplateMatchingTestBytes(data, 5)
	return &TestResult{
		Name: "重叠子序列检测方法",
		P:    p1, P2: p2,
		Q: q1, Q2: q2,
		Pass: math.Min(p1, p2) >= Alpha,
	}
}

// OverlappingTemplateMatchingTest 重叠子序列检测方法,m=5
// bits: 检测序列
// return:
//
//	p1: P-value1
//	p2: P-value2
func OverlappingTemplateMatchingTest(bits []bool) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return OverlappingTemplateMatchingProto(bits, 5)
}

// OverlappingTemplateMatchingTestBytes 重叠子序列检测方法
// data: 检测序列
// m: m长度,m=2,5
// return:
//
//	p1: P-value1
//	p2: P-value2
func OverlappingTemplateMatchingTestBytes(data []byte, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	return OverlappingTemplateMatchingProto(B2bitArr(data), m)
}

// OverlappingTemplateMatchingProto 重叠子序列检测方法
// bits: 检测序列
// m: m长度,m=3,5
// return:
//
//	p1: P-value1
//	p2: P-value2
func OverlappingTemplateMatchingProto(bits []bool, m int) (p1 float64, p2 float64, q1 float64, q2 float64) {
	n := len(bits)
	if n < 5 {
		panic("please provide valid test bits")
	}
	patterns1 := make([]int, 1<<m)
	patterns2 := make([]int, 1<<(m-1))
	patterns3 := make([]int, 1<<(m-2))
	var Phi1, Phi2, Phi3 float64 = 0, 0, 0
	var DPhi2, D2Phi2 float64 = 0, 0

	var mask1 int = (1 << m) - 1
	var mask2 int = (1 << (m - 1)) - 1
	var mask3 int = (1 << (m - 2)) - 1

	// 本来这里需要取bits后面预先插入bits[:m-1]，使得bits[m-1:]的长度依然是n。
	// 现在改成不对bits切片做预处理，而是取位时对索引进行模操作。
	//
	// Step 2
	tmp := subsequencepattern(bits, m-1)

	for i := m - 1; i < n+m-1; i++ {
		tmp <<= 1
		if bits[i%n] { // i % n is used to avoid appending m-1 bits in the end
			tmp++
		}
		patterns1[tmp&mask1]++
		patterns2[tmp&mask2]++
		patterns3[tmp&mask3]++
	}

	// Step 3
	for i := 0; i <= mask1; i++ {
		Phi1 += float64(patterns1[i]) * float64(patterns1[i])
	}
	Phi1 *= float64(mask1 + 1)
	Phi1 /= float64(n)
	Phi1 -= float64(n)
	for i := 0; i <= mask2; i++ {
		Phi2 += float64(patterns2[i]) * float64(patterns2[i])
	}
	Phi2 *= float64(mask2 + 1)
	Phi2 /= float64(n)
	Phi2 -= float64(n)
	for i := 0; i <= mask3; i++ {
		Phi3 += float64(patterns3[i]) * float64(patterns3[i])
	}
	Phi3 *= float64(mask3 + 1)
	Phi3 /= float64(n)
	Phi3 -= float64(n)

	// Step 4
	DPhi2 = Phi1 - Phi2
	D2Phi2 = Phi1 - 2*Phi2 + Phi3

	// Step 5
	p1 = igamc(float64(len(patterns3)), DPhi2/2.0)
	p2 = igamc(float64(len(patterns3))/2.0, D2Phi2/2.0)

	// Step 6
	q1 = p1
	q2 = p2

	return
}
