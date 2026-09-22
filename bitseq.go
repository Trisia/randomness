package randomness

import (
	"encoding/binary"
	"io/ioutil"
	"math/bits"
)

// BitSeq 是位序列的紧凑表示：每比特占 1 位，即 []bool 的 1/8 内存。
//
// 位序约定（与 B2bitArr 一致）：
//   - 序列第 i 位 <=> u[i>>6] 的第 (i&63) 位（LSB 起）
//   - 序列第 i 位 <=> data[i/8] 的 (0x80 >> (i%8)) 位（MSB 优先）
//
// 不变量：u 的长度恒为 (n+63)/64，且 u 中所有 i >= n 的位恒为 0。
// BitSeq 是值类型，但内部切片是共享的；若需独立副本请用 Copy。
type BitSeq struct {
	u []uint64
	n int
}

// b2u8 把 bool 转成 0/1。
func b2u8(b bool) uint8 {
	var v uint8
	if b {
		v = 1
	}
	return v
}

// BitSeqFromBytes 把字节序列按 B2bitArr 的位序打包成 BitSeq。
func BitSeqFromBytes(data []byte) BitSeq {
	n := len(data) * 8
	s := BitSeq{u: make([]uint64, (n+63)/64), n: n}
	i := 0
	// 整 8 字节一组：大端加载后整字反转即可让「序列第 i 位」落在第 i 位。
	for ; i+8 <= len(data); i += 8 {
		s.u[i>>3] = bits.Reverse64(binary.BigEndian.Uint64(data[i:]))
	}
	// 不足 8 字节的尾部：高位补齐后同样反转，超出 n 的位自然为 0。
	if rem := len(data) - i; rem > 0 {
		var v uint64
		for j := 0; j < rem; j++ {
			v |= uint64(data[i+j]) << uint(56-8*j)
		}
		s.u[i>>3] = bits.Reverse64(v)
	}
	return s
}

// BitSeqFromBools 把 []bool 打包成 BitSeq。
func BitSeqFromBools(bs []bool) BitSeq {
	n := len(bs)
	s := BitSeq{u: make([]uint64, (n+63)/64), n: n}
	full := n / 64
	for k := 0; k < full; k++ {
		var v uint64
		for j := 0; j < 8; j++ {
			lo := k*64 + j*8
			b := bs[lo : lo+8 : lo+8]
			c := uint64(b2u8(b[0])<<7 | b2u8(b[1])<<6 | b2u8(b[2])<<5 | b2u8(b[3])<<4 |
				b2u8(b[4])<<3 | b2u8(b[5])<<2 | b2u8(b[6])<<1 | b2u8(b[7]))
			// 大端序摆放，配合 Reverse64 得到自然位序
			v |= c << uint(8*(7-j))
		}
		s.u[k] = bits.Reverse64(v)
	}
	if rem := n - full*64; rem > 0 {
		var v uint64
		for j := 0; j*8 < rem; j++ {
			lo := full*64 + j*8
			hi := lo + 8
			if hi > n {
				hi = n
			}
			b := bs[lo:hi:hi]
			var c uint64
			for t := 0; t < len(b); t++ {
				c |= uint64(b2u8(b[t])) << uint(7-t)
			}
			v |= c << uint(8*(7-j))
		}
		s.u[full] = bits.Reverse64(v)
	}
	return s
}

// bit 返回序列第 i 位（0/1）。越界返回 0，便于写出与旧实现等价的环绕逻辑。
func (s BitSeq) bit(i int) uint64 {
	if i < 0 || i >= s.n {
		return 0
	}
	return (s.u[i>>6] >> uint(i&63)) & 1
}

// word 返回第 w 个字；越界返回 0。
func (s BitSeq) word(w int) uint64 {
	if w < 0 || w >= len(s.u) {
		return 0
	}
	return s.u[w]
}

// maskLow 返回低 k 位掩码；k >= 64 时返回全 1。
// 注意不能直接写 uint64(1)<<k - 1：k == 64 时移位会溢出为 0。
func maskLow(k int) uint64 {
	if k >= 64 {
		return ^uint64(0)
	}
	if k <= 0 {
		return 0
	}
	return uint64(1)<<uint(k) - 1
}

// popcountRange 统计 [lo, hi) 区间内 1 的个数。
func (s BitSeq) popcountRange(lo, hi int) int {
	if lo < 0 {
		lo = 0
	}
	if hi > s.n {
		hi = s.n
	}
	c := 0
	for lo < hi {
		w := lo >> 6
		b := lo & 63
		room := 64 - b
		if room > hi-lo {
			room = hi - lo
		}
		c += bits.OnesCount64((s.u[w] >> uint(b)) & maskLow(room))
		lo += room
	}
	return c
}

// ones 返回全部有效位中 1 的个数。
func (s BitSeq) ones() int {
	return s.popcountRange(0, s.n)
}

// extractBits 把从 off 开始的 count 位连续打包写入 dst（dst 长度需 >= (count+63)/64）。
func extractBits(s BitSeq, off, count int, dst []uint64) {
	for i := range dst {
		dst[i] = 0
	}
	out, rem := 0, count
	pos := off
	for rem > 0 {
		take := 64 - (out & 63)
		if v := 64 - (pos & 63); take > v {
			take = v
		}
		if take > rem {
			take = rem
		}
		chunk := (s.u[pos>>6] >> uint(pos&63)) & maskLow(take)
		dst[out>>6] |= chunk << uint(out&63)
		pos += take
		out += take
		rem -= take
	}
}

// toBools 展开成 []bool。
func (s BitSeq) toBools() []bool {
	out := make([]bool, s.n)
	for i := 0; i < s.n; i++ {
		out[i] = s.bit(i) == 1
	}
	return out
}

// NewBitSeq 返回长度为 nbits 的全 0 位序列。
// nbits < 0 时按 0 处理。
func NewBitSeq(nbits int) BitSeq {
	if nbits < 0 {
		nbits = 0
	}
	return BitSeq{u: make([]uint64, (nbits+63)/64), n: nbits}
}

// NewBitSeqFromBytes 把字节序列按 B2bitArr 的位序打包成 BitSeq。
// 它是 BitSeqFromBytes 的构造器风格别名，两者完全等价。
func NewBitSeqFromBytes(data []byte) BitSeq {
	return BitSeqFromBytes(data)
}

// ReadBitSeqFromFile 从文件读取字节数据并按 B2bitArr 位序打包成 BitSeq。
// 文件内容是原始二进制字节（每字节 8 位，MSB 优先），
// 与 BitSeqFromBytes / NewBitSeqFromBytes 使用相同的位序约定。
func ReadBitSeqFromFile(filename string) (BitSeq, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return BitSeq{}, err
	}
	return BitSeqFromBytes(buf), nil
}

// Copy 返回一个拥有独立底层存储的副本。
func (s BitSeq) Copy() BitSeq {
	u := make([]uint64, len(s.u))
	copy(u, s.u)
	return BitSeq{u: u, n: s.n}
}

// Len 返回位序列长度（比特数）。
func (s BitSeq) Len() int { return s.n }

// Bit 返回第 i 位的值；i 越界时返回 false。
func (s BitSeq) Bit(i int) bool { return s.bit(i) == 1 }

// Set 设置第 i 位；i 越界时不做任何事。
//
// 注意 BitSeq 内部切片是共享的，Set 会影响共享同一底层存储的其他 BitSeq。
func (s BitSeq) Set(i int, v bool) {
	if i < 0 || i >= s.n {
		return
	}
	if v {
		s.u[i>>6] |= 1 << uint(i&63)
	} else {
		s.u[i>>6] &^= 1 << uint(i&63)
	}
}

// Ones 返回全部有效位中 1 的个数。
func (s BitSeq) Ones() int { return s.ones() }

// OnesInRange 返回 [lo, hi) 区间内 1 的个数，区间会被裁剪到 [0, Len()]。
func (s BitSeq) OnesInRange(lo, hi int) int { return s.popcountRange(lo, hi) }

// ToBools 展开成 []bool。
func (s BitSeq) ToBools() []bool { return s.toBools() }

// truncate 把序列长度截断为 n（n 不得超过当前长度），并清零超出部分的位。
func (s BitSeq) truncate(n int) BitSeq {
	if n < 0 {
		n = 0
	}
	if n > s.n {
		n = s.n
	}
	s.n = n
	if r := n & 63; r != 0 {
		s.u[n>>6] &= maskLow(r)
	}
	if need := (n + 63) / 64; need < len(s.u) {
		s.u = s.u[:need]
	}
	return s
}

// ToBytes 按 B2bitArr 的位序（每字节 MSB 优先）打包成字节序列。
// 长度不足整字节时，末尾补 0 位。
func (s BitSeq) ToBytes() []byte {
	out := make([]byte, (s.n+7)/8)
	for i := 0; i < s.n; i++ {
		if s.bit(i) == 1 {
			out[i>>3] |= 0x80 >> uint(i&7)
		}
	}
	return out
}
