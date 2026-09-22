package randomness

import (
	"math/bits"
	"os"
	"path/filepath"
	"testing"
)

// 位序契约测试：BitSeq 必须与 B2bitArr 描述同一条比特流。
// 这是整个 packed 化的地基——位序错一位，所有 packed 内核的结果都会随之偏移。
func TestBitSeqMatchesB2bitArr(t *testing.T) {
	// 覆盖各种长度：0、1 字节、8 的倍数/非倍数、64 位边界附近
	lens := []int{0, 1, 2, 3, 7, 8, 9, 15, 16, 17, 63, 64, 65, 127, 128, 129, 255, 256}
	for _, n := range lens {
		data := NewDetRand(uint64(n) + 1).RawBytes(n)
		want := B2bitArr(data)
		got := BitSeqFromBytes(data)

		if len(want) != got.n {
			t.Fatalf("n=%d 字节: B2bitArr=%d 位, BitSeq.n=%d", n, len(want), got.n)
		}
		if len(got.u) != (got.n+63)/64 {
			t.Fatalf("n=%d 字节: u 长度 %d, 期望 %d", n, len(got.u), (got.n+63)/64)
		}
		for i := 0; i < got.n; i++ {
			w := uint64(0)
			if want[i] {
				w = 1
			}
			if g := got.bit(i); g != w {
				t.Fatalf("n=%d 字节: 第 %d 位不一致, B2bitArr=%d BitSeq=%d", n, i, w, g)
			}
		}
		// 同时校验 []bool 入口
		fromBools := BitSeqFromBools(want)
		if fromBools.n != got.n {
			t.Fatalf("n=%d 字节: BitSeqFromBools.n=%d", n, fromBools.n)
		}
		for i := range got.u {
			if fromBools.u[i] != got.u[i] {
				t.Fatalf("n=%d 字节: 字 %d 不一致, fromBytes=%x fromBools=%x",
					n, i, got.u[i], fromBools.u[i])
			}
		}
	}
}

// 不变量：超出 n 的位必须恒为 0。
// popcountRange 的正确性依赖这条（否则尾部补零位会被算进 1 的个数）。
func TestBitSeqPaddingBitsAreZero(t *testing.T) {
	for n := 1; n <= 200; n++ {
		data := NewDetRand(uint64(n)*7 + 3).RawBytes(n)
		s := BitSeqFromBytes(data)
		for w := range s.u {
			start := w * 64
			if start >= s.n {
				if s.u[w] != 0 {
					t.Fatalf("n=%d 字节: 整字 %d 完全越界却非 0 (%x)", n, w, s.u[w])
				}
				continue
			}
			if valid := s.n - start; valid < 64 {
				if hi := s.u[w] >> uint(valid); hi != 0 {
					t.Fatalf("n=%d 字节: 字 %d 的第 %d 位以上应为 0, 实际 %x", n, w, valid, hi)
				}
			}
		}
	}
}

// BitSeqFromBools 必须与逐位打包完全等价（无分支优化的正确性保证）。
func TestBitSeqFromBoolsMatchesNaive(t *testing.T) {
	for _, n := range []int{0, 1, 7, 8, 63, 64, 65, 127, 128, 129, 1000, 4096, 100000} {
		bs := NewDetRand(uint64(n) + 11).Bits(n)
		got := BitSeqFromBools(bs)
		naive := BitSeq{u: make([]uint64, (n+63)/64), n: n}
		for i, b := range bs {
			if b {
				naive.u[i>>6] |= 1 << uint(i&63)
			}
		}
		if len(got.u) != len(naive.u) {
			t.Fatalf("n=%d: u 长度 %d vs %d", n, len(got.u), len(naive.u))
		}
		for i := range naive.u {
			if got.u[i] != naive.u[i] {
				t.Fatalf("n=%d: 字 %d 不一致, got=%x naive=%x", n, i, got.u[i], naive.u[i])
			}
		}
	}
}

func TestMaskLow(t *testing.T) {
	cases := []struct {
		k    int
		want uint64
	}{
		{-1, 0}, {0, 0}, {1, 1}, {2, 3}, {63, 1<<63 - 1}, {64, ^uint64(0)}, {65, ^uint64(0)}, {100, ^uint64(0)},
	}
	for _, c := range cases {
		if got := maskLow(c.k); got != c.want {
			t.Errorf("maskLow(%d)=%x, 期望 %x", c.k, got, c.want)
		}
	}
}

func TestPopcountRange(t *testing.T) {
	const nbits = 20000
	data := NewDetRand(4242).RawBytes(nbits / 8)
	s := BitSeqFromBytes(data)
	ref := B2bitArr(data)

	if got, want := s.ones(), countOnes(ref, 0, nbits); got != want {
		t.Fatalf("ones()=%d, 期望 %d", got, want)
	}
	// 覆盖各种区间：空区间、单点、跨字、非对齐、全区间
	ranges := [][2]int{
		{0, 0}, {0, 1}, {5, 6}, {0, nbits}, {63, 65}, {64, 128}, {1, 63}, {7, 9999},
		{nbits - 1, nbits}, {nbits, nbits + 100}, {-5, 10}, {100, 100},
	}
	for _, r := range ranges {
		want := countOnes(ref, r[0], r[1])
		got := s.popcountRange(r[0], r[1])
		if got != want {
			t.Errorf("popcountRange(%d,%d)=%d, 期望 %d", r[0], r[1], got, want)
		}
	}
}

func countOnes(bits []bool, lo, hi int) int {
	if lo < 0 {
		lo = 0
	}
	if hi > len(bits) {
		hi = len(bits)
	}
	c := 0
	for i := lo; i < hi; i++ {
		if bits[i] {
			c++
		}
	}
	return c
}

// extractBits 必须保证 dst 中第 k 位就是序列第 off+k 位（连续、无空洞）。
func TestExtractBits(t *testing.T) {
	const nbits = 5000
	data := NewDetRand(99).RawBytes(nbits / 8)
	s := BitSeqFromBytes(data)
	ref := B2bitArr(data)

	for _, off := range []int{0, 1, 7, 31, 32, 63, 64, 65, 100, 499, 500, 1000} {
		for _, count := range []int{1, 7, 8, 31, 32, 33, 63, 64, 65, 499, 500, 501, 1000} {
			if off+count > nbits {
				continue
			}
			nw := (count + 63) / 64
			dst := make([]uint64, nw)
			extractBits(s, off, count, dst)
			for k := 0; k < count; k++ {
				want := uint64(0)
				if ref[off+k] {
					want = 1
				}
				got := (dst[k>>6] >> uint(k&63)) & 1
				if got != want {
					t.Fatalf("off=%d count=%d k=%d: got=%d want=%d", off, count, k, got, want)
				}
			}
			// 超出 count 的位必须为 0
			for k := count; k < nw*64; k++ {
				if (dst[k>>6]>>uint(k&63))&1 != 0 {
					t.Fatalf("off=%d count=%d: 超出部分的第 %d 位非 0", off, count, k)
				}
			}
		}
	}
}

func TestBitSeqToBoolsRoundTrip(t *testing.T) {
	for _, n := range []int{0, 1, 63, 64, 65, 1000} {
		bs := NewDetRand(uint64(n) + 5).Bits(n)
		got := BitSeqFromBools(bs).toBools()
		if len(got) != len(bs) {
			t.Fatalf("n=%d: 长度 %d vs %d", n, len(got), len(bs))
		}
		for i := range bs {
			if got[i] != bs[i] {
				t.Fatalf("n=%d: 第 %d 位往返不一致", n, i)
			}
		}
	}
}

// 转换代价的可观测基线，同时防止日后有人把无分支实现改回逐位分支版。
func BenchmarkBitSeqFromBytes(b *testing.B) {
	data := NewDetRand(1).RawBytes(125000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BitSeqFromBytes(data)
	}
}

func BenchmarkBitSeqFromBools(b *testing.B) {
	bs := B2bitArr(NewDetRand(1).RawBytes(125000))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BitSeqFromBools(bs)
	}
}

func BenchmarkB2bitArrBaseline(b *testing.B) {
	data := NewDetRand(1).RawBytes(125000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = B2bitArr(data)
	}
}

// NewBitSeqFromBytes 必须与 BitSeqFromBytes 完全等价。
func TestNewBitSeqFromBytes(t *testing.T) {
	data := NewDetRand(42).RawBytes(1000)
	want := BitSeqFromBytes(data)
	got := NewBitSeqFromBytes(data)
	if got.n != want.n {
		t.Fatalf("n 不一致: %d vs %d", got.n, want.n)
	}
	if len(got.u) != len(want.u) {
		t.Fatalf("u 长度不一致: %d vs %d", len(got.u), len(want.u))
	}
	for i := range got.u {
		if got.u[i] != want.u[i] {
			t.Fatalf("字 %d 不一致: %x vs %x", i, got.u[i], want.u[i])
		}
	}
}

// ReadBitSeqFromFile 必须与 BitSeqFromBytes 读取相同数据等价。
func TestReadBitSeqFromFile(t *testing.T) {
	data := NewDetRand(99).RawBytes(500)
	want := BitSeqFromBytes(data)

	tmpFile := filepath.Join(t.TempDir(), "test.bin")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ReadBitSeqFromFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if got.n != want.n {
		t.Fatalf("n 不一致: %d vs %d", got.n, want.n)
	}
	if len(got.u) != len(want.u) {
		t.Fatalf("u 长度不一致: %d vs %d", len(got.u), len(want.u))
	}
	for i := range got.u {
		if got.u[i] != want.u[i] {
			t.Fatalf("字 %d 不一致: %x vs %x", i, got.u[i], want.u[i])
		}
	}
}

// ReadBitSeqFromFile 对空文件也能正常处理。
func TestReadBitSeqFromFileEmpty(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "empty.bin")
	if err := os.WriteFile(tmpFile, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadBitSeqFromFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if got.n != 0 {
		t.Fatalf("空文件应产生 0 位, 实际 n=%d", got.n)
	}
	if len(got.u) != 0 {
		t.Fatalf("空文件的 u 长度应为 0, 实际 %d", len(got.u))
	}
}

// ReadBitSeqFromFile 对不存在的文件应返回错误。
func TestReadBitSeqFromFileNotExist(t *testing.T) {
	_, err := ReadBitSeqFromFile("/nonexistent/path/bits.bin")
	if err == nil {
		t.Fatal("不存在的文件应返回错误")
	}
}

var _ = bits.OnesCount64
