package randomness

import (
	"bufio"
	rand2 "crypto/rand"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

const (
	MAXLOG float64 = 7.09782712893383996732224e2 // log(MAXNUM)
	biginv float64 = 2.22044604925031308085e-16
	big    float64 = 4.503599627370496e15
	MACHEP float64 = 1.11022302462515654042e-16
)

func igam(a, x float64) float64 {
	var ans, ax, c, r float64

	if (x <= 0) || (a <= 0) {
		return 0.0
	}

	if (x > 1.0) && (x > a) {
		return 1.e0 - igamc(a, x)
	}

	/* Compute  x**a * exp(-x) / gamma(a)  */
	ax = a*math.Log(x) - x - logGamma(a)
	if ax < -MAXLOG {
		return 0.0
	}
	ax = math.Exp(ax)

	/* power series */
	r = a
	c = 1.0
	ans = 1.0

	for {
		r += 1.0
		c *= x / r
		ans += c
		if !(c/ans > MACHEP) {
			break
		}
	}

	return ans * ax / a
}

// logGamma 返回 log|Γ(x)|。
func logGamma(x float64) float64 {
	res, _ := math.Lgamma(x)
	return res
}

func Igamc(a, x float64) float64 {
	return igamc(a, x)
}

func igamc(a, x float64) float64 {
	var ans, ax, c, yc, r, t, y, z float64
	var pk, pkm1, pkm2, qk, qkm1, qkm2 float64

	if (x <= 0) || (a <= 0) {
		return (1.0)
	}

	if (x < 1.0) || (x < a) {
		return (1.e0 - igam(a, x))
	}

	ax = a*math.Log(x) - x - logGamma(a)

	if ax < -MAXLOG {
		return 0.0
	}
	ax = math.Exp(ax)

	/* continued fraction */
	y = 1.0 - a
	z = x + y + 1.0
	c = 0.0
	pkm2 = 1.0
	qkm2 = x
	pkm1 = x + 1.0
	qkm1 = z * x
	ans = pkm1 / qkm1

	for {
		c += 1.0
		y += 1.0
		z += 2.0
		yc = y * c
		pk = pkm1*z - pkm2*yc
		qk = qkm1*z - qkm2*yc
		if qk != 0 {
			r = pk / qk
			t = math.Abs((ans - r) / r)
			ans = r
		} else {
			t = 1.0
		}
		pkm2 = pkm1
		pkm1 = pk
		qkm2 = qkm1
		qkm1 = qk
		if math.Abs(pk) > big {
			pkm2 *= biginv
			pkm1 *= biginv
			qkm2 *= biginv
			qkm1 *= biginv
		}
		if !(t > MACHEP) {
			break
		}
	}
	return ans * ax
}

func normal_CDF(x float64) float64 {
	return (1 + math.Erf(x/math.Sqrt(2))) / 2
}

// rank 计算 GF(2) 上 M 行 Q 列矩阵的秩。
func rank(matrix [][]int, M, Q int) int {
	temp := make([][]int, M)
	for i := 0; i < M; i++ {
		temp[i] = make([]int, Q)
		for j := 0; j < Q; j++ {
			temp[i][j] = matrix[i][j]
		}
	}

	pivot := 0
	for col := 0; col < Q && pivot < M; col++ {
		sel := -1
		for r := pivot; r < M; r++ {
			if temp[r][col] == 1 {
				sel = r
				break
			}
		}
		if sel < 0 {
			continue
		}
		if sel != pivot {
			temp[sel], temp[pivot] = temp[pivot], temp[sel]
		}
		// 消去主元行以下的同列 1；主元行本身之后不再被使用，
		// 因此无需回代——非零行数即秩。
		for r := pivot + 1; r < M; r++ {
			if temp[r][col] == 1 {
				for j := col; j < Q; j++ {
					temp[r][j] ^= temp[pivot][j]
				}
			}
		}
		pivot++
	}
	return pivot
}

func linearComplexity(a []bool, M int) int {
	var N_ int = 0
	var L int = 0
	var m int = -1
	var d int = 0

	// 预分配数组并初始化
	B_ := make([]int, M)
	C := make([]int, M)
	P := make([]int, M)
	T := make([]int, M)

	C[0] = 1
	B_[0] = 1

	for N_ < M {
		// 计算 d = a[N_] + sum(C[i]*a[N_-i]) mod 2
		d = b2i(a[N_])
		for i := 1; i <= L; i++ {
			if C[i] == 1 {
				d ^= b2i(a[N_-i])
			}
		}

		if d == 1 {
			// 复制 C 到 T
			copy(T, C)
			// 清零 P
			for i := range P {
				P[i] = 0
			}

			// 计算 P = B_ shifted by (N_-m)
			shift := N_ - m
			for j := 0; j < M; j++ {
				if B_[j] == 1 && j+shift < M {
					P[j+shift] = 1
				}
			}

			// C = C + P mod 2 (使用 XOR)
			for i := 0; i < M; i++ {
				C[i] ^= P[i]
			}

			if L <= N_/2 {
				L = N_ + 1 - L
				m = N_
				// 复制 T 到 B_
				copy(B_, T)
			}
		}
		N_++
	}
	return L
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ceilPow2(N int) int {
	i := 2
	for {
		if i >= N {
			return i
		}
		i <<= 1
	}
}

// B2bit 字节 转换为 bool数组
func B2bit(b byte) []bool {
	//b&0b10000000 > 0,
	//b&0b01000000 > 0,
	//b&0b00100000 > 0,
	//b&0b00010000 > 0,
	//b&0b00001000 > 0,
	//b&0b00000100 > 0,
	//b&0b00000010 > 0,
	//b&0b00000001 > 0,
	return []bool{
		b&0x80 > 0,
		b&0x40 > 0,
		b&0x20 > 0,
		b&0x10 > 0,
		b&0x08 > 0,
		b&0x04 > 0,
		b&0x02 > 0,
		b&0x01 > 0,
	}
}

// B2Byte bool数组 转换为 字节
func B2Byte(arr []bool) byte {
	var res byte = 0
	var v byte = 0
	for _, b := range arr {
		res <<= 1
		if b {
			v = 1
		} else {
			v = 0
		}
		res += v
	}
	return res
}

// B2bitArr 转换字节数组为比特序列。
func B2bitArr(src []byte) []bool {
	n := len(src) * 8
	res := make([]bool, n)

	for i, b := range src {
		baseIdx := i * 8
		res[baseIdx] = b&0x80 > 0
		res[baseIdx+1] = b&0x40 > 0
		res[baseIdx+2] = b&0x20 > 0
		res[baseIdx+3] = b&0x10 > 0
		res[baseIdx+4] = b&0x08 > 0
		res[baseIdx+5] = b&0x04 > 0
		res[baseIdx+6] = b&0x02 > 0
		res[baseIdx+7] = b&0x01 > 0
	}
	return res
}

// GroupBitSeq 生成一组 nbits 比特的随机检测序列，返回 BitSeq。
func GroupBitSeq(nbits int) (BitSeq, error) {
	if nbits < 0 {
		nbits = 0
	}
	buf := make([]byte, (nbits+7)/8)
	if _, err := rand2.Read(buf); err != nil {
		return BitSeq{}, err
	}
	return BitSeqFromBytes(buf).truncate(nbits), nil
}

// ReadGroupBitSeq 从文件（每字节 8 位，MSB 优先）读取一组二元序列，返回 BitSeq。
func ReadGroupBitSeq(filename string) (BitSeq, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return BitSeq{}, err
	}
	return BitSeqFromBytes(buf), nil
}

// ReadGroupInASCIIFormatBitSeq 从 ASCII 格式文件（每行若干 0/1 字符）读取
// 一组二元序列；非 0/1 字符（空白、换行等）会被跳过。
//
// nbits > 0：只读取前 nbits 位，文件中可用位不足时返回错误。
// nbits <= 0：读取文件中的全部位。
func ReadGroupInASCIIFormatBitSeq(filename string, nbits int) (BitSeq, error) {
	file, err := os.Open(filename)
	if err != nil {
		return BitSeq{}, err
	}
	defer file.Close()

	r := bufio.NewReaderSize(file, 64*1024)
	var words []uint64
	var cur uint64
	var cnt uint
	total := 0
	for nbits <= 0 || total < nbits {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b != '0' && b != '1' {
			continue
		}
		if b == '1' {
			cur |= 1 << cnt
		}
		cnt++
		total++
		if cnt == 64 {
			words = append(words, cur)
			cur, cnt = 0, 0
		}
	}
	if cnt > 0 {
		words = append(words, cur)
	}
	if nbits > 0 && total < nbits {
		return BitSeq{}, fmt.Errorf("%s 仅含 %d 位，不足所需的 %d 位", filename, total, nbits)
	}
	return BitSeq{u: words, n: total}, nil
}

// GroupBit 生成一组 10^6 比特的检测序列。
func GroupBit() []bool {
	s, _ := GroupBitSeq(1000000)
	return s.ToBools()
}

// GroupSecBit 生成一组测试数据 长度为 10^6 比特。
func GroupSecBit() []bool {
	return GroupBit()
}

// ReadGroup 从文件中读取一组二元序列。
func ReadGroup(filename string) []bool {
	s, err := ReadGroupBitSeq(filename)
	if err != nil {
		panic(err)
	}
	return s.ToBools()
}

// ReadGroupInASCIIFormat 从 ASCII 格式文件读取 10^6 bit 的二元序列。
func ReadGroupInASCIIFormat(filename string) []bool {
	const n = 1000000
	s, err := ReadGroupInASCIIFormatBitSeq(filename, n)
	if err != nil {
		panic(err)
	}
	return s.ToBools()
}
