package randomness

import (
	"math"
	"math/cmplx"
)

// DiscreteFourierTransformFast 使用Go标准库实现的快速版本
func DiscreteFourierTransformFast(data []byte) *TestResult {
	p, q := DiscreteFourierTransformTestBytesFast(data)
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBytesFast 快速版本的离散傅里叶检测
func DiscreteFourierTransformTestBytesFast(data []byte) (float64, float64) {
	return DiscreteFourierTransformTestFast(B2bitArr(data))
}

// DiscreteFourierTransformTestFast 快速版本的离散傅里叶检测
// 使用简化的DFT实现，避免外部依赖
func DiscreteFourierTransformTestFast(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// 对于大数据集，使用采样方法来提高性能
	const maxSampleSize = 1048576 // 1M samples max
	var sampleBits []bool
	var sampleN int

	if n <= maxSampleSize {
		sampleBits = bits
		sampleN = n
	} else {
		// 均匀采样
		sampleN = maxSampleSize
		sampleBits = make([]bool, sampleN)
		step := float64(n) / float64(sampleN)
		for i := 0; i < sampleN; i++ {
			idx := int(float64(i) * step)
			sampleBits[i] = bits[idx]
		}
	}

	// 使用优化的DFT实现
	N := nextPowerOf2(sampleN)
	X := make([]complex128, N)

	// 转换比特到复数
	for i := 0; i < sampleN; i++ {
		if sampleBits[i] {
			X[i] = 1 + 0i
		} else {
			X[i] = -1 + 0i
		}
	}

	// 填充零
	for i := sampleN; i < N; i++ {
		X[i] = 0 + 0i
	}

	// 快速傅里叶变换
	fastFFT(X)

	// 计算统计量
	T := math.Sqrt(2.995732274 * float64(sampleN))
	N_0 := 0.95 * float64(sampleN) / 2

	var N_1 int
	halfN := sampleN / 2
	for i := 0; i < halfN-1; i++ {
		if cmplx.Abs(X[i]) < T {
			N_1++
		}
	}

	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(2.0*sampleN)/3.8)
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}

// fastFFT 快速傅里叶变换实现（原地计算）
func fastFFT(x []complex128) {
	n := len(x)
	if n <= 1 {
		return
	}

	// 位反转置换
	for i, j := 0, 0; i < n; i++ {
		if i < j {
			x[i], x[j] = x[j], x[i]
		}
		m := n >> 1
		for ; j&m != 0; m >>= 1 {
			j &= ^m
		}
		j |= m
	}

	// Cooley-Tukey FFT
	for length := 2; length <= n; length <<= 1 {
		half := length >> 1
		angle := -2 * math.Pi / float64(length)
		wlen := complex(math.Cos(angle), math.Sin(angle))

		for i := 0; i < n; i += length {
			w := 1 + 0i
			for j := 0; j < half; j++ {
				u := x[i+j]
				v := x[i+j+half] * w
				x[i+j] = u + v
				x[i+j+half] = u - v
				w *= wlen
			}
		}
	}
}

// nextPowerOf2 返回大于等于n的最小2的幂
func nextPowerOf2(n int) int {
	if n <= 1 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
