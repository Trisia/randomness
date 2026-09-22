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

// coreFrequencyWithinBlock 是块内频数检测的实现（packed 内核）。
func coreFrequencyWithinBlock(s BitSeq, m int) (float64, float64) {
	if m <= 0 {
		panic("please provide valid m (> 0)")
	}
	N := s.n / m
	if N == 0 {
		panic("please provide test bits")
	}
	var V float64
	for i := 0; i < N; i++ {
		ones := s.popcountRange(i*m, i*m+m)
		Pi := float64(ones)/float64(m) - 0.5
		V += Pi * Pi
	}
	V *= 2.0 * float64(m) // 这一步本来V要乘以4，现在改为2，免得后一步再除以2。
	P := igamc(float64(N)/2.0, V)
	return P, P
}

// FrequencyWithinBlock 块内频数检测, m = 10000 for bits = 1000_000
func FrequencyWithinBlock(data []byte) *TestResult {
	p, q := coreFrequencyWithinBlock(BitSeqFromBytes(data), selectM(len(data)*8))
	return &TestResult{Name: "块内频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// FrequencyWithinBlockTest 块内频数检测, m = 10000 for bits = 1000_000
//
// Deprecated: 请改用 FrequencyWithinBlockTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 FrequencyWithinBlockTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func FrequencyWithinBlockTest(bits []bool) (float64, float64) {
	return coreFrequencyWithinBlock(BitSeqFromBools(bits), selectM(len(bits)))
}

// FrequencyWithinBlockTestBytes 块内频数检测
func FrequencyWithinBlockTestBytes(data []byte, m int) (float64, float64) {
	return coreFrequencyWithinBlock(BitSeqFromBytes(data), m)
}

// FrequencyWithinBlockProto 块内频数检测
//
// Deprecated: 请改用 FrequencyWithinBlockTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 FrequencyWithinBlockTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func FrequencyWithinBlockProto(bits []bool, m int) (float64, float64) {
	return coreFrequencyWithinBlock(BitSeqFromBools(bits), m)
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

// FrequencyWithinBlockBitSeq 块内频数检测（m 按数据规模自动选择）
func FrequencyWithinBlockBitSeq(s BitSeq) *TestResult {
	p, q := coreFrequencyWithinBlock(s, selectM(s.n))
	return &TestResult{Name: "块内频数检测", P: p, Q: q, Pass: p >= Alpha}
}

// FrequencyWithinBlockTestBitSeq 块内频数检测
func FrequencyWithinBlockTestBitSeq(s BitSeq, m int) (float64, float64) {
	return coreFrequencyWithinBlock(s, m)
}
