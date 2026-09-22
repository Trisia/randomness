package randomness

import (
	"math/rand"
	"testing"
)

// 差分测试：新的通用 M×Q rank 必须与旧方阵实现在 32×32 上逐矩阵一致。
// 旧实现从 git HEAD 的 utils.go 原样复制而来（仅改名），确保对比的是真身。

func legacyRowEchelon(matrix [][]int, m int) {
	pivotstartrow := 0
	pivotstartcol := 0
	pivotrow := 0
	for i := 0; i < m; i++ {
		found := false
		for k := pivotstartrow; k < m; k++ {
			if matrix[k][pivotstartcol] == 1 {
				found = true
				pivotrow = k
				break
			}
		}
		if found {
			if pivotrow != pivotstartrow {
				for k := 0; k < m; k++ {
					matrix[pivotrow][k] ^= matrix[pivotstartrow][k]
					matrix[pivotstartrow][k] ^= matrix[pivotrow][k]
					matrix[pivotrow][k] ^= matrix[pivotstartrow][k]
				}
			}
			for j := pivotstartrow + 1; j < m; j++ {
				if matrix[j][pivotstartcol] == 1 {
					for k := 0; k < m; k++ {
						matrix[j][k] = matrix[pivotstartrow][k] ^ matrix[j][k]
					}
				}
			}
			pivotstartcol += 1
			pivotstartrow += 1
		} else {
			pivotstartcol += 1
		}
	}
}

func legacyRank(matrix [][]int, m int) int {
	temp := make([][]int, m)
	for i := 0; i < m; i++ {
		temp[i] = make([]int, m)
		for j := 0; j < m; j++ {
			temp[i][j] = matrix[i][j]
		}
	}
	legacyRowEchelon(temp, m)
	r := 0
	for i := 0; i < m; i++ {
		notZero := false
		for j := 0; j < m; j++ {
			if temp[i][j] != 0 {
				notZero = true
			}
		}
		if notZero {
			r++
		}
	}
	return r
}

func TestRankDiffAgainstLegacy(t *testing.T) {
	rng := rand.New(rand.NewSource(20240920))
	const m = 32
	mat := make([][]int, m)
	for i := range mat {
		mat[i] = make([]int, m)
	}

	// 1) 稠密随机矩阵
	dist := map[int]int{}
	for iter := 0; iter < 20000; iter++ {
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				mat[i][j] = rng.Intn(2)
			}
		}
		got := rank(mat, m, m)
		want := legacyRank(mat, m)
		if got != want {
			t.Fatalf("iter=%d rank 不一致: new=%d legacy=%d", iter, got, want)
		}
		dist[got]++
	}
	t.Logf("稠密随机 20000 个矩阵 rank 分布: %v", dist)

	// 2) 稀疏矩阵（更容易出现秩缺陷，覆盖 31/30 等分支）
	dist = map[int]int{}
	for iter := 0; iter < 20000; iter++ {
		p := rng.Intn(30) // 0..29% 的置 1 概率
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				if rng.Intn(100) < p {
					mat[i][j] = 1
				} else {
					mat[i][j] = 0
				}
			}
		}
		got := rank(mat, m, m)
		want := legacyRank(mat, m)
		if got != want {
			t.Fatalf("稀疏 iter=%d (p=%d) rank 不一致: new=%d legacy=%d", iter, p, got, want)
		}
		dist[got]++
	}
	t.Logf("稀疏矩阵 20000 个 rank 分布: %v", dist)

	// 3) 结构化矩阵：全 0、全 1、单位阵、重复行、上三角
	cases := []struct {
		name string
		fill func(i, j int) int
	}{
		{"全0", func(i, j int) int { return 0 }},
		{"全1", func(i, j int) int { return 1 }},
		{"单位阵", func(i, j int) int {
			if i == j {
				return 1
			}
			return 0
		}},
		{"重复行", func(i, j int) int { return j % 2 }},
		{"上三角", func(i, j int) int {
			if j >= i {
				return 1
			}
			return 0
		}},
		{"下三角", func(i, j int) int {
			if j <= i {
				return 1
			}
			return 0
		}},
	}
	for _, c := range cases {
		for i := 0; i < m; i++ {
			for j := 0; j < m; j++ {
				mat[i][j] = c.fill(i, j)
			}
		}
		got := rank(mat, m, m)
		want := legacyRank(mat, m)
		if got != want {
			t.Errorf("%s rank 不一致: new=%d legacy=%d", c.name, got, want)
		}
	}

	// 4) 非方阵：新实现必须能算（旧实现会 panic / 越界），且秩不超过 min(M,Q)
	for _, dim := range [][2]int{{16, 16}, {32, 32}, {8, 32}, {32, 8}, {64, 64}, {1, 1}, {40, 24}} {
		M, Q := dim[0], dim[1]
		mm := make([][]int, M)
		for i := range mm {
			mm[i] = make([]int, Q)
			for j := range mm[i] {
				mm[i][j] = rng.Intn(2)
			}
		}
		r := rank(mm, M, Q)
		if r < 0 || r > M || r > Q {
			t.Errorf("M=%d Q=%d 秩越界: %d", M, Q, r)
		}
	}
}

// MatrixRankTestBytes 对 M>32 不应再 panic，且结果应在 [0,1] 内。
func TestMatrixRankNonSquareNoPanic(t *testing.T) {
	for _, dim := range [][2]int{{16, 16}, {32, 32}, {64, 64}, {8, 32}} {
		M, Q := dim[0], dim[1]
		data := make([]byte, M*Q*10/8)
		rng := rand.New(rand.NewSource(int64(M*100 + Q)))
		rng.Read(data)
		p, q := MatrixRankTestBytes(data, M, Q)
		if p < 0 || p > 1 || q < 0 || q > 1 {
			t.Errorf("M=%d Q=%d P/Q 越界: %v %v", M, Q, p, q)
		}
		t.Logf("M=%d Q=%d -> P=%.6f", M, Q, p)
	}
}
