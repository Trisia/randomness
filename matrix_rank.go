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

// gf2Rank 在 GF(2) 上求按行位打包矩阵的秩。
//
// mat 有 M 行，每行 qw = ceil(Q/64) 个字；第 r 行第 j 列位于
// mat[r*qw + (j>>6)] 的第 (j&63) 位。
func gf2Rank(mat []uint64, M, Q, qw int) int {
	rank := 0
	for col := 0; col < Q && rank < M; col++ {
		w, b := col>>6, uint(col&63)
		pivot := -1
		for r := rank; r < M; r++ {
			if mat[r*qw+w]>>b&1 == 1 {
				pivot = r
				break
			}
		}
		if pivot < 0 {
			continue
		}
		if pivot != rank {
			for k := 0; k < qw; k++ {
				mat[rank*qw+k], mat[pivot*qw+k] = mat[pivot*qw+k], mat[rank*qw+k]
			}
		}
		base := rank * qw
		for r := 0; r < M; r++ {
			if r == rank {
				continue
			}
			if mat[r*qw+w]>>b&1 == 1 {
				rb := r * qw
				for k := 0; k < qw; k++ {
					mat[rb+k] ^= mat[base+k]
				}
			}
		}
		rank++
	}
	return rank
}

// coreMatrixRank 是矩阵秩检测的实现（packed 内核）。
func coreMatrixRank(s BitSeq, M, Q int) (float64, float64) {
	if M <= 0 || Q <= 0 {
		panic("please provide valid matrix dimension")
	}
	N := s.n / (M * Q)
	if N == 0 {
		panic("please provide valid test bits")
	}
	qw := (Q + 63) / 64
	mat := make([]uint64, M*qw)
	target := M
	if Q < M {
		target = Q
	}

	var Fm, Fm1, Fr int
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			extractBits(s, i*M*Q+j*Q, Q, mat[j*qw:j*qw+qw])
		}
		r := gf2Rank(mat, M, Q, qw)
		switch {
		case r == target:
			Fm++
		case r == target-1:
			Fm1++
		default:
			Fr++
		}
	}

	_N := float64(N)
	// 0.2888/0.5776/0.1336 是 32×32 随机二元矩阵秩为 32/31/其它 的概率。
	V := math.Pow(float64(Fm)-0.2888*_N, 2)/(0.2888*_N) +
		math.Pow(float64(Fm1)-0.5776*_N, 2)/(0.5776*_N) +
		math.Pow(float64(Fr)-0.1336*_N, 2)/(0.1336*_N)
	P := igamc(1, V/2.0)
	return P, P
}

// MatrixRank 矩阵秩检测,M=Q=32
func MatrixRank(data []byte) *TestResult {
	p, q := coreMatrixRank(BitSeqFromBytes(data), 32, 32)
	return &TestResult{Name: "矩阵秩检测", P: p, Q: q, Pass: p >= Alpha}
}

// MatrixRankTest 矩阵秩检测,M=Q=32
//
// Deprecated: 请改用 MatrixRankTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 MatrixRankTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func MatrixRankTest(bits []bool) (float64, float64) {
	return coreMatrixRank(BitSeqFromBools(bits), 32, 32)
}

// MatrixRankTestBytes 矩阵秩检测
func MatrixRankTestBytes(data []byte, M, Q int) (float64, float64) {
	return coreMatrixRank(BitSeqFromBytes(data), M, Q)
}

// MatrixRankProto 矩阵秩检测
// bits: 待检测序列
// M: 矩阵行数
// Q: 矩阵列隶属
//
// Deprecated: 请改用 MatrixRankTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 MatrixRankTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func MatrixRankProto(bits []bool, M, Q int) (float64, float64) {
	return coreMatrixRank(BitSeqFromBools(bits), M, Q)
}

// MatrixRankBitSeq 矩阵秩检测，M=Q=32
func MatrixRankBitSeq(s BitSeq) *TestResult {
	p, q := coreMatrixRank(s, 32, 32)
	return &TestResult{Name: "矩阵秩检测", P: p, Q: q, Pass: p >= Alpha}
}

// MatrixRankTestBitSeq 矩阵秩检测
// M: 矩阵行数
// Q: 矩阵列数
func MatrixRankTestBitSeq(s BitSeq, M, Q int) (float64, float64) {
	return coreMatrixRank(s, M, Q)
}
