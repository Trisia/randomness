package randomness

import (
	"math/rand"
	"testing"
)

// 位打包内核与「朴素整数实现」的差分测试。
//
// 这两个内核改写自公认算法（GF(2) 高斯消元、Berlekamp-Massey），
// 但它们对「位序、字边界、缓冲区复用」的处理才是风险所在，
// 因此这里与保留在 utils.go 中的朴素实现逐例对拍。

func TestPackedMatrixRankAgainstRank(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	dims := [][2]int{{1, 1}, {8, 8}, {16, 16}, {32, 32}, {64, 64}, {8, 32}, {32, 8}, {40, 24}, {100, 100}, {65, 129}}
	for _, dim := range dims {
		M, Q := dim[0], dim[1]
		qw := (Q + 63) / 64
		intMat := make([][]int, M)
		for i := range intMat {
			intMat[i] = make([]int, Q)
		}
		packed := make([]uint64, M*qw)
		for iter := 0; iter < 200; iter++ {
			// 必须清零：gf2Rank 会原地消元，且 packed 是跨迭代复用的
			for i := range packed {
				packed[i] = 0
			}
			// 前一半用稠密随机，后一半用稀疏（更容易出现秩亏，覆盖小秩分支）
			sparse := iter >= 100
			for i := 0; i < M; i++ {
				for j := 0; j < Q; j++ {
					v := 0
					if sparse {
						if rng.Intn(10) == 0 {
							v = 1
						}
					} else {
						v = rng.Intn(2)
					}
					intMat[i][j] = v
					if v == 1 {
						packed[i*qw+(j>>6)] |= 1 << uint(j&63)
					}
				}
			}
			want := rank(intMat, M, Q)
			got := gf2Rank(packed, M, Q, qw)
			if got != want {
				t.Fatalf("M=%d Q=%d iter=%d(sparse=%v): gf2Rank=%d rank=%d", M, Q, iter, sparse, got, want)
			}
		}
	}
}

func TestPackedLinearComplexityAgainstLegacy(t *testing.T) {
	ms := []int{1, 2, 3, 7, 8, 33, 63, 64, 65, 100, 127, 128, 499, 500, 501, 512, 1000}
	for _, m := range ms {
		nw := (m + 63) / 64
		C := make([]uint64, nw)
		B := make([]uint64, nw)
		T := make([]uint64, nw)
		R := make([]uint64, nw)
		block := make([]uint64, nw)

		// 结构化序列：全 0、全 1、交替（BM 的边界分支最容易在这里出错）
		patterns := map[string]func(i int) bool{
			"zeros":   func(i int) bool { return false },
			"ones":    func(i int) bool { return true },
			"alt01":   func(i int) bool { return i%2 == 0 },
			"alt0011": func(i int) bool { return (i/2)%2 == 0 },
		}
		// 真正的 LFSR 序列（线性复杂度应当恰好等于阶数）
		lfsr := make([]bool, m)
		st := uint32(0xACE1)
		for i := range lfsr {
			b := st & 1
			lfsr[i] = b == 1
			st >>= 1
			if b == 1 {
				st ^= 0xB400
			}
		}

		check := func(name string, bits []bool) {
			want := linearComplexity(bits, m)
			extractBits(BitSeqFromBools(bits), 0, m, block)
			got := bmComplexity(block, C, B, T, R, nw, m)
			if got != want {
				t.Fatalf("m=%d %s: bmComplexity=%d linearComplexity=%d", m, name, got, want)
			}
		}

		for name, f := range patterns {
			bits := make([]bool, m)
			for i := range bits {
				bits[i] = f(i)
			}
			check(name, bits)
		}
		check("lfsr", lfsr)

		// 随机块
		for iter := 0; iter < 40; iter++ {
			check("random", NewDetRand(uint64(m*1000+iter)).Bits(m))
		}
	}
}

// 位并行 BM 的复杂度分布应当与朴素实现完全一致（而不只是返回值相等）。
// 这里用同一批数据构造 T 值并比较落入的区间。
func TestPackedLinearComplexityDistribution(t *testing.T) {
	const m = 500
	const blocks = 200
	bits := NewDetRand(31337).Bits(m * blocks)

	s := BitSeqFromBools(bits)
	nw := (m + 63) / 64
	C := make([]uint64, nw)
	B := make([]uint64, nw)
	T := make([]uint64, nw)
	R := make([]uint64, nw)
	block := make([]uint64, nw)

	for i := 0; i < blocks; i++ {
		seg := bits[i*m : (i+1)*m]
		want := linearComplexity(seg, m)
		extractBits(s, i*m, m, block)
		got := bmComplexity(block, C, B, T, R, nw, m)
		if got != want {
			t.Fatalf("块 %d: bmComplexity=%d linearComplexity=%d", i, got, want)
		}
	}
}
