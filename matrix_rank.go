package randomness

import (
	"fmt"
	"math"
)

// MatrixRankTest 矩阵秩检测
func MatrixRankTest(bits []bool) float64 {
	n := len(bits)
	if n == 0 {
		fmt.Println("BinaryMatrixRankTest:args wrong")
		return -1
	}
	M := 32
	Q := 32

	N := n / (M * Q)
	//int n_disc = n % (M * Q);
	var Fm, Fm1, Fr = 0, 0, 0
	var matrix = make([][]int, 32)
	for i := 0; i < 32; i++ {
		matrix[i] = make([]int, 32)
	}
	var V, P float64
	var r int
	var b bool

	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			for k := 0; k < Q; k++ {
				b, bits = bits[0], bits[1:]
				if b {
					matrix[j][k] = 1
				} else {
					matrix[j][k] = 0
				}
			}
		}
		r = rank(matrix, M)

		if r == min(M, Q) {
			Fm++
		} else if r == (min(M, Q) - 1) {
			Fm1++
		} else {
			Fr++
		}
	}
	_N := float64(N)
	V = math.Pow(float64(Fm)-0.2888*_N, 2.0)/(0.2888*_N) +
		math.Pow(float64(Fm1)-0.5776*_N, 2.0)/(0.5776*_N) +
		math.Pow(float64(Fr)-0.1336*_N, 2.0)/(0.1336*_N)

	P = igamc(1, V/2.0)

	return P
}
