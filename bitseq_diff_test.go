package randomness

import (
	"math"
	"testing"
)

// 本文件用「最直白的逐位实现」作为参照，验证各 packed 内核里位技巧的正确性。
//
// 固定向量只能覆盖有限的输入；而位技巧（整字 popcount、XOR-shift、转移位图、
// 倍增消去）最容易在「尾部不满 64 位」「跨字边界」「非 64 倍数长度」上出错，
// 所以这里用大量随机长度 + 结构化数据做交叉验证。

const diffSeed = uint64(20240920)

// ---------- 参照实现（逐位、最直白） ----------

func bruteRuns(bits []bool) (float64, float64) {
	n := len(bits)
	V_obs := 1
	Pi := 0.0
	for i := 0; i < n-1; i++ {
		if bits[i] != bits[i+1] {
			V_obs++
		}
		if bits[i] {
			Pi++
		}
	}
	if bits[n-1] {
		Pi++
	}
	Pi /= float64(n)
	V := (float64(V_obs) - 2.0*float64(n)*Pi*(1.0-Pi)) / (2.0 * math.Sqrt(float64(2*n)) * Pi * (1.0 - Pi))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2.0
}

// bruteBinaryDerivative 严格照抄优化前的实现（含逐趟收缩的循环上界 n-i-1），
// 用于确认「整字全长 XOR-shift」与「逐位收缩 XOR」在有效区等价。
func bruteBinaryDerivative(bits []bool, k int) (float64, float64) {
	n := len(bits)
	_bits := make([]bool, n)
	copy(_bits, bits)
	for i := 0; i < k; i++ {
		for j := 0; j < n-i-1; j++ {
			_bits[j] = _bits[j] != _bits[j+1]
		}
	}
	S := 0
	for i := 0; i < n-k; i++ {
		if _bits[i] {
			S++
		} else {
			S--
		}
	}
	V := float64(S) / math.Sqrt(2*float64(n-k))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

func bruteAutocorrelation(bits []bool, d int) (float64, float64) {
	n := len(bits)
	Ad := 0
	for i := 0; i < n-d; i++ {
		if bits[i] != bits[i+d] {
			Ad++
		}
	}
	V := 2.0 * (float64(Ad) - (float64(n-d) / 2.0)) / math.Sqrt(2*float64(n-d))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

func bruteLongestRun(bits []bool, checkOne bool) int {
	best, cur := 0, 0
	for _, b := range bits {
		v := b
		if !checkOne {
			v = !b
		}
		if v {
			cur++
			if cur > best {
				best = cur
			}
		} else {
			cur = 0
		}
	}
	return best
}

// bruteRunsDistribution 严格照抄优化前的逐位扫描实现。
func bruteRunsDistribution(bits []bool) (float64, float64) {
	n := len(bits)
	k := 0
	for {
		k++
		_2k2 := 1 << uint(k+2)
		if float64(n-k+3)/float64(_2k2) < 5.0 {
			break
		}
	}
	k--

	b := make([]float64, k)
	g := make([]float64, k)
	cur := bits[0]
	cnt := 0
	for i := 0; i < n; i++ {
		if bits[i] == cur {
			cnt++
		} else {
			if cnt > k {
				cnt = k
			}
			if cur {
				b[cnt-1]++
			} else {
				g[cnt-1]++
			}
			cur = bits[i]
			cnt = 1
		}
	}
	if cnt > k {
		cnt = k
	}
	if cur {
		b[cnt-1]++
	} else {
		g[cnt-1]++
	}

	var T float64 = 0
	for i := 0; i < k; i++ {
		T += b[i] + g[i]
	}
	e := make([]float64, k)
	for i := 0; i < k; i++ {
		if i < k-1 {
			e[i] = T / float64(int(1)<<uint(i+2))
		} else {
			e[i] = T / float64(int(1)<<uint(i+1))
		}
	}
	var V float64 = 0
	for i := 0; i < k; i++ {
		V += (b[i] - e[i]) * (b[i] - e[i]) / e[i]
		V += (g[i] - e[i]) * (g[i] - e[i]) / e[i]
	}
	P := igamc(float64(k-1), V/2.0)
	return P, P
}

// ---------- 测试数据 ----------

// diffCases 返回 (名称, 比特序列) 的集合：随机 + 结构化 + 边界长度。
func diffCases(t *testing.T) []struct {
	name string
	bits []bool
} {
	t.Helper()
	var out []struct {
		name string
		bits []bool
	}
	// 随机长度，重点覆盖非 64 倍数与 64 边界
	lens := []int{128, 129, 191, 192, 193, 255, 256, 257, 511, 512, 999, 1000, 1001, 4095, 4096, 4097, 20000, 10007, 64000}
	for i, n := range lens {
		out = append(out, struct {
			name string
			bits []bool
		}{"random", NewDetRand(diffSeed + uint64(i)).Bits(n)})
	}
	// 结构化数据：全 0、全 1、交替、周期 8、长游程块、25% 偏置
	for _, v := range structuredDiffVectors() {
		out = append(out, struct {
			name string
			bits []bool
		}{v.name, B2bitArr(v.data)})
	}
	return out
}

// structuredDiffVectors 返回对拍用的结构化输入。
//
// 这些序列刻意覆盖 packed 内核最容易走错的分支：全同、周期、长游程与偏置，
// 长度固定为 20000 bit（2500 字节），与随机长度用例互补。
func structuredDiffVectors() []struct {
	name string
	data []byte
} {
	const nbits = 20000
	repeat := func(v byte) []byte {
		b := make([]byte, nbits/8)
		for i := range b {
			b[i] = v
		}
		return b
	}
	// blocks 返回「blockBits 位 0，接着 blockBits 位 1」交替的长游程序列。
	// blockBits 必须是 8 的倍数，以保证字节边界与位边界对齐。
	blocks := func(blockBits int) []byte {
		b := make([]byte, nbits/8)
		nb := blockBits / 8
		for i := range b {
			if (i/nb)%2 == 0 {
				b[i] = 0x00
			} else {
				b[i] = 0xFF
			}
		}
		return b
	}
	// biased25 返回「每比特为 1 的概率恰好为 25%」的确定性序列：
	// 对两个独立随机比特取逻辑与，故 P(1) = 0.5 × 0.5 = 0.25。
	biased25 := func(seed uint64) []byte {
		r := NewDetRand(seed)
		b := make([]byte, nbits/8)
		for i := range b {
			a, c := r.Uint64(), r.Uint64()
			var out byte
			for j := 0; j < 8; j++ {
				if (a>>uint(63-j)&1)&(c>>uint(63-j)&1) == 1 {
					out |= 0x80 >> uint(j)
				}
			}
			b[i] = out
		}
		return b
	}
	return []struct {
		name string
		data []byte
	}{
		{"pattern/zeros", repeat(0x00)},
		{"pattern/ones", repeat(0xFF)},
		{"pattern/alternating-01", repeat(0xAA)},
		{"pattern/period8-0x5A", repeat(0x5A)},
		{"pattern/blocks64", blocks(64)},
		{"biased/ones-25pct", biased25(0xB1A5ED25)},
	}
}

// ---------- 对拍 ----------

func TestPackedRunsAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		wantP, wantQ := bruteRuns(c.bits)
		gotP, gotQ := coreRuns(BitSeqFromBools(c.bits))
		if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
			t.Errorf("%s n=%d: packed=(%v,%v) brute=(%v,%v)", c.name, len(c.bits), gotP, gotQ, wantP, wantQ)
		}
	}
}

func TestPackedBinaryDerivativeAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		for _, k := range []int{1, 2, 3, 7, 16} {
			if len(c.bits) <= k {
				continue
			}
			wantP, wantQ := bruteBinaryDerivative(c.bits, k)
			gotP, gotQ := coreBinaryDerivative(BitSeqFromBools(c.bits), k)
			if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
				t.Errorf("%s n=%d k=%d: packed=(%v,%v) brute=(%v,%v)",
					c.name, len(c.bits), k, gotP, gotQ, wantP, wantQ)
			}
		}
	}
}

func TestPackedAutocorrelationAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		for _, d := range []int{1, 2, 3, 8, 16, 63, 64, 65, 127} {
			if len(c.bits) <= d {
				continue
			}
			wantP, wantQ := bruteAutocorrelation(c.bits, d)
			gotP, gotQ := coreAutocorrelation(BitSeqFromBools(c.bits), d)
			if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
				t.Errorf("%s n=%d d=%d: packed=(%v,%v) brute=(%v,%v)",
					c.name, len(c.bits), d, gotP, gotQ, wantP, wantQ)
			}
		}
	}
}

// longestRunBlock 直接与逐位参照比对最长游程本身（比比对 P 值更精确），
// 并且刻意在任意偏移上切块，模拟真实的块内扫描。
func TestLongestRunBlockAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		n := len(c.bits)
		s := BitSeqFromBools(c.bits)
		for _, blockLen := range []int{1, 7, 8, 63, 64, 65, 128, 129, 500, 1000, 9999} {
			if blockLen > n {
				continue
			}
			// 取三个偏移：0、中间、末尾对齐处
			offs := []int{0, (n - blockLen) / 2, n - blockLen}
			for _, off := range offs {
				seg := c.bits[off : off+blockLen]
				words := make([]uint64, (blockLen+63)/64)
				extractBits(s, off, blockLen, words)
				for _, checkOne := range []bool{true, false} {
					want := bruteLongestRun(seg, checkOne)
					got := longestRunBlock(words, blockLen, checkOne)
					if got != want {
						t.Fatalf("%s n=%d blockLen=%d off=%d checkOne=%v: got=%d want=%d",
							c.name, n, blockLen, off, checkOne, got, want)
					}
				}
			}
		}
	}
}

func TestPackedRunsDistributionAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		if len(c.bits) < 100 {
			continue
		}
		wantP, wantQ := bruteRunsDistribution(c.bits)
		gotP, gotQ := coreRunsDistribution(BitSeqFromBools(c.bits))
		if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
			t.Errorf("%s n=%d: packed=(%v,%v) brute=(%v,%v)", c.name, len(c.bits), gotP, gotQ, wantP, wantQ)
		}
	}
}

// 覆盖 popcountRange 落到整字/非整字边界时的块内频数
func TestPackedFreqBlockAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		n := len(c.bits)
		s := BitSeqFromBools(c.bits)
		for _, m := range []int{1, 8, 10, 64, 65, 100, 1000, 10000} {
			if m > n {
				continue
			}
			N := n / m
			var wantV float64
			for i := 0; i < N; i++ {
				ones := 0
				for _, b := range c.bits[i*m : i*m+m] {
					if b {
						ones++
					}
				}
				Pi := float64(ones)/float64(m) - 0.5
				wantV += Pi * Pi
			}
			wantV *= 2.0 * float64(m)
			wantP := igamc(float64(N)/2.0, wantV)
			gotP, gotQ := coreFrequencyWithinBlock(s, m)
			if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantP) {
				t.Errorf("%s n=%d m=%d: packed=%v brute=%v", c.name, n, m, gotP, wantP)
			}
		}
	}
}

func bruteCumulative(bits []bool, forward bool) (float64, float64) {
	n := len(bits)
	S, Z := 0, 0
	for i := 0; i < n; i++ {
		idx := i
		if !forward {
			idx = n - 1 - i
		}
		if bits[idx] {
			S++
		} else {
			S--
		}
		if S > Z {
			Z = S
		} else if -S > Z {
			Z = -S
		}
	}
	sqrtN := math.Sqrt(float64(n))
	P := 1.0
	for i := ((-n / Z) + 1) / 4; i <= ((n/Z)-1)/4; i++ {
		P -= normal_CDF(float64((4*i+1)*Z)/sqrtN) - normal_CDF(float64((4*i-1)*Z)/sqrtN)
	}
	for i := ((-n / Z) - 3) / 4; i <= ((n/Z)-1)/4; i++ {
		P += normal_CDF(float64((4*i+3)*Z)/sqrtN) - normal_CDF(float64((4*i+1)*Z)/sqrtN)
	}
	return P, P
}

// 累加和的位序最容易搞错：前向要按 bit 0 -> bit n-1，后向要反过来。
// 用逐位参照在随机与结构化数据上双向校验。
// brutePoker 用最直白的方式按 m 位窗口建直方图。
func brutePoker(bits []bool, m int) (float64, float64) {
	n := len(bits)
	_2m := 1 << uint(m)
	patterns := make([]int, _2m)
	N := n / m
	for i := 0; i < N; i++ {
		v := 0
		for j := 0; j < m; j++ {
			v <<= 1
			if bits[i*m+j] {
				v++
			}
		}
		patterns[v]++
	}
	var V float64
	for i := 0; i < _2m; i++ {
		V += float64(patterns[i]) * float64(patterns[i])
	}
	V *= float64(_2m)
	V /= float64(N)
	V -= float64(N)
	P := igamc(float64(_2m-1)/2, V/2)
	return P, P
}

// corePoker 必须与逐窗口直方图完全一致，覆盖 m 的各个取值
// （尤其 m 不是 8 的约数、以及 m 跨越 64 位字边界的情形）。
func TestPackedPokerAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		// m 覆盖到 maxPokerM 以下；m=9/13/16/17 在非对齐起点上必然跨越 64 位字边界，
		// 这正是 corePoker 的位提取循环要覆盖的分支。
		for _, m := range []int{1, 2, 3, 4, 5, 7, 8, 9, 13, 16, 17, 20} {
			if m > len(c.bits) {
				continue
			}
			wantP, wantQ := brutePoker(c.bits, m)
			gotP, gotQ := corePoker(BitSeqFromBools(c.bits), m)
			if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
				t.Errorf("%s n=%d m=%d: packed=(%v,%v) brute=(%v,%v)",
					c.name, len(c.bits), m, gotP, gotQ, wantP, wantQ)
			}
		}
	}
}

func bruteMonoBit(bits []bool) (float64, float64) {
	S := 0
	for _, b := range bits {
		if b {
			S++
		} else {
			S--
		}
	}
	V := float64(S) / math.Sqrt(2*float64(len(bits)))
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

func TestPackedMonoBitAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		wantP, wantQ := bruteMonoBit(c.bits)
		gotP, gotQ := coreMonoBitFrequency(BitSeqFromBools(c.bits))
		if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
			t.Errorf("%s n=%d: packed=(%v,%v) brute=(%v,%v)", c.name, len(c.bits), gotP, gotQ, wantP, wantQ)
		}
	}
}

func TestPackedCumulativeAgainstBrute(t *testing.T) {
	for _, c := range diffCases(t) {
		n := len(c.bits)
		for _, forward := range []bool{true, false} {
			wantP, wantQ := bruteCumulative(c.bits, forward)
			gotP, gotQ := coreCumulative(BitSeqFromBools(c.bits), forward)
			if !floatsClose(gotP, wantP) || !floatsClose(gotQ, wantQ) {
				t.Errorf("%s n=%d forward=%v: packed=(%v,%v) brute=(%v,%v)",
					c.name, n, forward, gotP, gotQ, wantP, wantQ)
			}
		}
	}
}

func floatsClose(a, b float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.IsNaN(a) && math.IsNaN(b)
	}
	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return a == b
	}
	return math.Abs(a-b) <= 1e-12
}
