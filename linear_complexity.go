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
	"math/bits"
	"runtime"
	"sync"
)

// bmComplexity 用位并行 Berlekamp-Massey 求 block 前 M 位的线性复杂度。
func bmComplexity(block, C, B, T, R []uint64, nw, M int) int {
	for i := range C {
		C[i], B[i], T[i], R[i] = 0, 0, 0, 0
	}
	C[0], B[0] = 1, 1
	L, m := 0, -1
	for N := 0; N < M; N++ {
		// 偏差 d = a[N] ^ Σ_{i=1..L} C[i]*a[N-i]
		par := int((block[N>>6] >> uint(N&63)) & 1)
		for w := 0; w < nw; w++ {
			cs := C[w] >> 1
			if w+1 < nw {
				cs |= C[w+1] << 63
			}
			par ^= bits.OnesCount64(cs&R[w]) & 1
		}
		if par&1 == 1 {
			copy(T, C)
			shift := N - m
			sw, sb := shift>>6, uint(shift&63)
			if sb == 0 {
				for w := nw - 1; w >= sw; w-- {
					C[w] ^= B[w-sw]
				}
			} else {
				for w := nw - 1; w >= sw; w-- {
					v := B[w-sw] << sb
					if w-sw-1 >= 0 {
						v |= B[w-sw-1] >> (64 - sb)
					}
					C[w] ^= v
				}
			}
			if L <= N/2 {
				L, m = N+1-L, N
				copy(B, T)
			}
		}
		// 为下一轮准备 R：R_{N+1} = (R_N << 1) | a[N]
		if N+1 < M {
			bit := (block[N>>6] >> uint(N&63)) & 1
			for w := nw - 1; w > 0; w-- {
				R[w] = R[w]<<1 | R[w-1]>>63
			}
			R[0] = R[0]<<1 | bit
		}
	}
	return L
}

// lcClassify 按 T 值落入 7 个区间之一计数。
func lcClassify(v *[7]float64, T float64) {
	switch {
	case T > 2.5:
		v[6]++
	case T > 1.5:
		v[5]++
	case T > 0.5:
		v[4]++
	case T > -0.5:
		v[3]++
	case T > -1.5:
		v[2]++
	case T > -2.5:
		v[1]++
	default:
		v[0]++
	}
}

// lcRange 处理 [lo, hi) 号数据块，把分类计数累加进 v。
func lcRange(s BitSeq, m, nw int, sign, miu float64, lo, hi int, v *[7]float64) {
	C := make([]uint64, nw)
	B := make([]uint64, nw)
	T := make([]uint64, nw)
	R := make([]uint64, nw)
	block := make([]uint64, nw)
	for i := lo; i < hi; i++ {
		extractBits(s, i*m, m, block)
		cplx := bmComplexity(block, C, B, T, R, nw, m)
		lcClassify(v, sign*(float64(cplx)-miu)+2.0/9.0)
	}
}

// coreLinearComplexity 是线型复杂度检测的实现（packed 内核）。
func coreLinearComplexity(s BitSeq, m, workers int) (float64, float64) {
	N := s.n / m
	if N == 0 {
		panic("please provide valid test bits")
	}
	if m <= 0 {
		panic("please provide valid m (> 0)")
	}
	nw := (m + 63) / 64

	// Step 3, miu - 预计算 _1_m
	var _1_m float64 = 1.0
	if m%2 != 0 {
		_1_m = -1.0
	}
	miu := float64(m)/2.0 + (9.0+_1_m)/36.0 - (float64(m)/3.0+2.0/9.0)/math.Pow(2.0, float64(m))

	var v [7]float64
	if workers <= 1 {
		lcRange(s, m, nw, _1_m, miu, 0, N, &v)
	} else {
		if workers > N {
			workers = N
		}
		chunk := (N + workers - 1) / workers
		partials := make([][7]float64, workers)
		var wg sync.WaitGroup
		for wk := 0; wk < workers; wk++ {
			lo, hi := wk*chunk, wk*chunk+chunk
			if hi > N {
				hi = N
			}
			if lo >= hi {
				continue
			}
			wg.Add(1)
			go func(wk, lo, hi int) {
				defer wg.Done()
				lcRange(s, m, nw, _1_m, miu, lo, hi, &partials[wk])
			}(wk, lo, hi)
		}
		wg.Wait()
		for _, p := range partials {
			for i := range v {
				v[i] += p[i]
			}
		}
	}

	// Step 6
	pi := [7]float64{0.010417, 0.03125, 0.12500, 0.5000, 0.25000, 0.06250, 0.020833}
	NF := float64(N)
	var V float64 = 0.0
	for i := 0; i < 7; i++ {
		diff := v[i] - NF*pi[i]
		V += diff * diff / (NF * pi[i])
	}
	// Step 7
	P := igamc(3.0, V/2.0)
	return P, P
}

// lcDispatch 按数据规模选择串行或并行：块数少于 50 或总长度小于 50000 bit 时用串行。
func lcDispatch(s BitSeq, m int) (float64, float64) {
	if s.n/m < 50 || s.n < 50000 {
		return coreLinearComplexity(s, m, 1)
	}
	return coreLinearComplexity(s, m, runtime.NumCPU())
}

// LinearComplexity 线型复杂度检测,m=500
func LinearComplexity(data []byte) *TestResult {
	p, q := lcDispatch(BitSeqFromBytes(data), 500)
	return &TestResult{Name: "线型复杂度检测(m=500)", P: p, Q: q, Pass: p >= Alpha}
}

// LinearComplexityTest 线型复杂度检测,m=500
//
// Deprecated: 请改用 LinearComplexityTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LinearComplexityTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LinearComplexityTest(bits []bool) (float64, float64) {
	return lcDispatch(BitSeqFromBools(bits), 500)
}

// LinearComplexityTestBytes 线型复杂度检测
// data: 待检测序列
// m: m长度
func LinearComplexityTestBytes(data []byte, m int) (float64, float64) {
	return lcDispatch(BitSeqFromBytes(data), m)
}

// LinearComplexityProto 线型复杂度检测
// bits: 待检测序列
// m: m长度
//
// Deprecated: 请改用 LinearComplexityTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LinearComplexityTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LinearComplexityProto(bits []bool, m int) (float64, float64) {
	return lcDispatch(BitSeqFromBools(bits), m)
}

// LinearComplexityProtoSerial 串行版本的线性复杂度检测
//
// Deprecated: 请改用 LinearComplexityTestBitSeqSerial——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LinearComplexityTestBitSeqSerial 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LinearComplexityProtoSerial(bits []bool, m int) (float64, float64) {
	return coreLinearComplexity(BitSeqFromBools(bits), m, 1)
}

// LinearComplexityProtoParallel 并行版本的线性复杂度检测
//
// Deprecated: 请改用 LinearComplexityTestBitSeqParallel——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 LinearComplexityTestBitSeqParallel 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func LinearComplexityProtoParallel(bits []bool, m int) (float64, float64) {
	return coreLinearComplexity(BitSeqFromBools(bits), m, runtime.NumCPU())
}

// LinearComplexityBitSeq 线型复杂度检测，m=500
func LinearComplexityBitSeq(s BitSeq) *TestResult {
	p, q := lcDispatch(s, 500)
	return &TestResult{Name: "线型复杂度检测(m=500)", P: p, Q: q, Pass: p >= Alpha}
}

// LinearComplexityTestBitSeq 线型复杂度检测
// m: m长度
func LinearComplexityTestBitSeq(s BitSeq, m int) (float64, float64) {
	return lcDispatch(s, m)
}

// LinearComplexityTestBitSeqSerial 线型复杂度检测（强制串行）。
func LinearComplexityTestBitSeqSerial(s BitSeq, m int) (float64, float64) {
	return coreLinearComplexity(s, m, 1)
}

// LinearComplexityTestBitSeqParallel 线型复杂度检测（强制并行）。
func LinearComplexityTestBitSeqParallel(s BitSeq, m int) (float64, float64) {
	return coreLinearComplexity(s, m, runtime.NumCPU())
}
