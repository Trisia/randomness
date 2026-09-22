package detect

import (
	"math"
	"testing"

	"github.com/Trisia/randomness"
)

func sameResult(a, b *randomness.TestResult) bool {
	if a == nil || b == nil {
		return a == b
	}
	eq := func(x, y float64) bool {
		if math.IsNaN(x) || math.IsNaN(y) {
			return math.IsNaN(x) && math.IsNaN(y)
		}
		return x == y
	}
	return a.Name == b.Name && a.Pass == b.Pass &&
		eq(a.P, b.P) && eq(a.Q, b.Q) && eq(a.P2, b.P2) && eq(a.Q2, b.Q2)
}

// 并发版必须与串行版逐项完全一致（含顺序）。
// 线型复杂度内部并行只是把整数计数按块累加，合并顺序不影响结果，
// 因此这里可以用精确相等而不是容差比较。
func TestRoundConcurrentMatchesSerial(t *testing.T) {
	for _, nbits := range []int{20000, 1000000} {
		data := randomness.NewDetRand(uint64(nbits)).RawBytes(nbits / 8)
		cases := []struct {
			name       string
			serial     func([]byte) []*randomness.TestResult
			concurrent func([]byte) []*randomness.TestResult
		}{
			{"Round15", Round15, Round15Concurrent},
			{"Round12", Round12, Round12Concurrent},
		}
		for _, c := range cases {
			want := c.serial(data)
			got := c.concurrent(data)
			if len(want) != len(got) {
				t.Fatalf("%s n=%d: 长度 %d vs %d", c.name, nbits, len(want), len(got))
			}
			for i := range want {
				if !sameResult(want[i], got[i]) {
					t.Fatalf("%s n=%d 第 %d 项不一致:\n serial=%+v\n conc  =%+v",
						c.name, nbits, i, want[i], got[i])
				}
			}
		}
	}
}

func BenchmarkRound15(b *testing.B) {
	data := randomness.NewDetRand(1).RawBytes(1000000 / 8)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Round15(data)
	}
}

func BenchmarkRound15Concurrent(b *testing.B) {
	data := randomness.NewDetRand(1).RawBytes(1000000 / 8)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Round15Concurrent(data)
	}
}
