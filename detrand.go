package randomness

import "math/bits"

// DetRand 是确定性的伪随机数发生器，供库内测试与数据生成工具（tools/rdgen）复用。
//
// 使用 xoshiro256**（SplitMix64 播种）算法，只使用 64 位整数运算，
// 不依赖 unsafe、字节序或 CPU 特性，输出与 Go 版本、操作系统及架构无关。
//
// 位序约定：字节流与 B2bitArr 一致（MSB 优先）。
// DetRand 实现了 io.Reader，可直接作为 detect 包的随机源。
// 零值不可用，请用 NewDetRand 构造。DetRand 不是并发安全的。
type DetRand struct {
	s [4]uint64
	// rem/remN 保存上一次 Read 未能交付的整字后缀，保证 Read 可以接收
	// 任意长度（含非 8 字节对齐）的缓冲区而不丢字节。
	// 若不做缓冲，用 1 字节缓冲区连续读取会得到与整块读取完全不同的流，
	// 而 io.Reader 的调用方（如 detect 包）完全有权这样做。
	rem  uint64 // 左对齐：下一个待交付字节位于 bit 56..63
	remN int    // rem 中待交付的字节数，0..7
}

// SplitMix64 的黄金比例增量与混合常数（与算法原始论文一致）。
const (
	detRandSmInc  uint64 = 0x9E3779B97F4A7C15
	detRandSmMul1 uint64 = 0xBF58476D1CE4E5B9
	detRandSmMul2 uint64 = 0x94D049BB133111EB
)

// NewDetRand 用给定种子构造一个 DetRand。
// 相同种子在任意 Go 版本、任意平台、任意架构上产生完全相同的序列。
func NewDetRand(seed uint64) *DetRand {
	r := &DetRand{}
	r.Reset(seed)
	return r
}

// Reset 重置生成器到给定种子对应的初始状态。
func (r *DetRand) Reset(seed uint64) {
	// 用 SplitMix64 把单个 64 位种子扩展成 xoshiro256** 需要的 4 个字。
	// 直接令 s[0]=seed 其余为 0 会产生质量很差的前若干个输出，故必须做扩展。
	z := seed
	for i := 0; i < 4; i++ {
		z += detRandSmInc
		x := z
		x = (x ^ (x >> 30)) * detRandSmMul1
		x = (x ^ (x >> 27)) * detRandSmMul2
		x = x ^ (x >> 31)
		r.s[i] = x
	}
	r.rem = 0
	r.remN = 0
}

// Uint64 返回序列中的下一个 64 位字。
func (r *DetRand) Uint64() uint64 {
	// xoshiro256** 的标量参考实现。
	result := bits.RotateLeft64(r.s[1]*5, 7) * 9

	t := r.s[1] << 17
	r.s[2] ^= r.s[0]
	r.s[3] ^= r.s[1]
	r.s[1] ^= r.s[2]
	r.s[0] ^= r.s[3]
	r.s[2] ^= t
	r.s[3] = bits.RotateLeft64(r.s[3], 45)

	return result
}

// Read 实现 io.Reader，用伪随机字节填满 p，返回 len(p) 且永不返回错误。
//
// 字节按 MSB 优先从每个 64 位字中取出，与 Bits 的位序一致。
func (r *DetRand) Read(p []byte) (int, error) {
	n := len(p)
	i := 0

	// 先交付上一次残留的字节后缀
	for i < n && r.remN > 0 {
		p[i] = byte(r.rem >> 56)
		r.rem <<= 8
		r.remN--
		i++
	}

	// 整字直出
	for ; i+8 <= n; i += 8 {
		v := r.Uint64()
		p[i+0] = byte(v >> 56)
		p[i+1] = byte(v >> 48)
		p[i+2] = byte(v >> 40)
		p[i+3] = byte(v >> 32)
		p[i+4] = byte(v >> 24)
		p[i+5] = byte(v >> 16)
		p[i+6] = byte(v >> 8)
		p[i+7] = byte(v)
	}

	// 尾部不足 8 字节：交付前 k 字节，其余留到下次
	if i < n {
		v := r.Uint64()
		k := n - i
		for j := 0; j < k; j++ {
			p[i+j] = byte(v >> uint(56-8*j))
		}
		r.rem = v << uint(8*k)
		r.remN = 8 - k
	}
	return n, nil
}

// Bytes 返回恰好 nbits/8 字节的确定性随机数据（nbits 向下取整到 8 的倍数）。
func (r *DetRand) Bytes(nbits int) []byte {
	if nbits < 0 {
		nbits = 0
	}
	b := make([]byte, nbits/8)
	r.Read(b)
	return b
}

// RawBytes 返回恰好 nbytes 个字节。
func (r *DetRand) RawBytes(nbytes int) []byte {
	if nbytes < 0 {
		nbytes = 0
	}
	return r.Bytes(nbytes * 8)
}

// Bits 返回长度恰为 nbits 的布尔序列，位序与 B2bitArr 一致（MSB 优先）。
//
// 它严格建立在 Read 之上，因此与 Bytes 描述同一条字节流：对 8 的倍数 n，
// 恒有 B2bitArr(src.Bytes(n)) == src.Bits(n)。
// n 不是 8 的倍数时，会多消耗一个字节（padding 位含在其中）。
func (r *DetRand) Bits(nbits int) []bool {
	if nbits < 0 {
		nbits = 0
	}
	b := make([]byte, (nbits+7)/8)
	r.Read(b)
	out := make([]bool, nbits)
	for i := 0; i < nbits; i++ {
		out[i] = b[i>>3]&(0x80>>uint(i&7)) != 0
	}
	return out
}
