package randomness

import (
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

const (
	MAXLOG float64 = 7.09782712893383996732224e2 // log(MAXNUM)
	biginv float64 = 2.22044604925031308085e-16
	big    float64 = 4.503599627370496e15
	MACHEP float64 = 1.11022302462515654042e-16
)

func subsequencepattern(bits []bool, m int) int {
	tmp := 0
	var b bool
	for j := 0; j < m; j++ {
		tmp <<= 1
		b, bits = bits[0], bits[1:]
		if b {
			tmp++
		}
	}
	return tmp
}

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

func logGamma(x float64) float64 {
	res, sign := math.Lgamma(x)
	return res * float64(sign)
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

func rank(matrix [][]int, m int) int {
	temp := make([][]int, m)
	for i := 0; i < m; i++ {
		temp[i] = make([]int, m)
		for j := 0; j < m; j++ {
			temp[i][j] = matrix[i][j]
		}
	}

	rowEchelon(temp, m)
	rank := 0
	for i := 0; i < m; i++ {
		notZero := false
		for j := 0; j < m; j++ {
			if temp[i][j] != 0 {
				notZero = true
			}
		}
		if notZero {
			rank++
		}
	}
	return rank
}

func rowEchelon(matrix [][]int, m int) {
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

func linearComplexity(a []bool, M int) int {
	var N_ int = 0
	var L int = 0
	var m int = -1
	var d int = 0
	var B_, C, P, T []int

	B_ = make([]int, M)
	C = make([]int, M)
	P = make([]int, M)
	T = make([]int, M)

	for i := 0; i < M; i++ {
		B_[i] = 0
		C[i] = 0
		T[i] = 0
		P[i] = 0
	}
	C[0] = 1
	B_[0] = 1
	for {
		if !(N_ < M) {
			break
		}
		d = b2i(a[N_])
		for i := 1; i <= L; i++ {
			d += C[i] * b2i(a[N_-i])
		}
		d = d % 2
		if d == 1 {
			for i := 0; i < M; i++ {
				T[i] = C[i]
				P[i] = 0
			}
			for j := 0; j < M; j++ {
				if B_[j] == 1 {
					P[j+N_-m] = 1
				}
			}
			for i := 0; i < M; i++ {
				C[i] = (C[i] + P[i]) % 2
			}
			if L <= N_/2 {
				L = N_ + 1 - L
				m = N_
				for i := 0; i < M; i++ {
					B_[i] = T[i]
				}
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

func pow2DoubleArr(data []float64) []float64 {
	// 创建新数组
	var newData []float64

	dataLength := len(data)

	sumNum := 2
	for sumNum < dataLength {
		sumNum = sumNum * 2
	}
	addLength := sumNum - dataLength

	if addLength != 0 {
		newData = make([]float64, sumNum)
		copy(newData, data)
		for i := dataLength; i < sumNum; i++ {
			newData[i] = 0
		}
	} else {
		newData = data
	}

	return newData
}

// B2bit 字节 转换为 bool数组
func B2bit(b byte) []bool {
	return []bool{
		b&0b10000000 > 0,
		b&0b01000000 > 0,
		b&0b00100000 > 0,
		b&0b00010000 > 0,
		b&0b00001000 > 0,
		b&0b00000100 > 0,
		b&0b00000010 > 0,
		b&0b00000001 > 0,
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

func xor(x, y bool) bool {
	return x != y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x > 0 {
		return x
	}
	return -x
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// B2bitArr 转换字节数组为比特序列
func B2bitArr(src []byte) []bool {
	res := make([]bool, 0, len(src)*8)
	for _, b := range src {
		res = append(res, B2bit(b)...)
	}
	return res
}

// GroupBit 生成一组 10^6 比特的监测序列
func GroupBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		if rand.Int()%2 == 1 {
			bits = append(bits, true)
		} else {
			bits = append(bits, false)
		}
	}
	return bits
}

// GroupSecBit 生成一组测试数据 长度为 10^6 比特
func GroupSecBit() []bool {
	n := 1000_000
	bits := make([]bool, 0, n)

	buf := make([]byte, n/8)
	_, _ = rand.Read(buf)
	for _, b := range buf {
		bits = append(bits, B2bit(b)...)
	}
	return bits
}

// ReadGroup 从文件中读取一组二元序列
func ReadGroup(filename string) []bool {
	n := 1000_000
	bits := make([]bool, 0, n)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for _, b := range buf {
		bits = append(bits, B2bit(b)...)
	}
	return bits
}
