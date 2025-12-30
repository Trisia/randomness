// Copyright (c) 2021 Quan guanyu
// randomness is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//          http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package randomness

import (
	"math"
	"math/cmplx"
	"sync"

	"github.com/Trisia/randomness/fft"
)

// FFT缓存表，用于预置常见数据规模的FFT
var (
	fftCache = make(map[int]fft.FFT)
	fftMutex sync.RWMutex
)

// GMT 0005-2021 规范的附录A中的样本长度及检测设置
const (
	// SmallScale 小规模：2*10^4 bit
	SmallScale = 20000
	// MediumScale 中规模：10^6 bit
	MediumScale = 1000000
	// LargeScale 大规模：10^8 bit
	LargeScale = 100000000
)

//// 预置FFT表，初始化常见规模的FFT
//func init() {
//	// 预置常见规模的FFT表
//	scales := []int{SmallScale, MediumScale, LargeScale}
//	for _, scale := range scales {
//		n := ceilPow2(scale)
//		if f, err := fft.New(n); err == nil {
//			fftCache[n] = f
//		}
//	}
//}

// getFFT 获取FFT实例，优先使用缓存的预置表
func getFFT(n int) (fft.FFT, error) {
	fftMutex.RLock()
	if f, exists := fftCache[n]; exists {
		fftMutex.RUnlock()
		return f, nil
	}
	fftMutex.RUnlock()

	// 缓存中没有，创建新的FFT实例
	fftMutex.Lock()
	defer fftMutex.Unlock()

	// 双重检查，防止并发创建
	if f, exists := fftCache[n]; exists {
		return f, nil
	}

	f, err := fft.New(n)
	if err != nil {
		return fft.FFT{}, err
	}
	fftCache[n] = f
	return f, nil
}

// DiscreteFourierTransform 离散傅里叶检测
func DiscreteFourierTransform(data []byte) *TestResult {
	p, q := DiscreteFourierTransformTestBytes(data)
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBytes 离散傅里叶检测
func DiscreteFourierTransformTestBytes(data []byte) (float64, float64) {
	return DiscreteFourierTransformTest(B2bitArr(data))
}

// DiscreteFourierTransformTest 离散傅里叶检测
// 离散傅立叶变换检测使用频谱的方法来检测序列的随机性。对待检序列进行傅立叶变换后可以得
// 到尖峰高度，根据随机性的假设，这个尖峰高度不能超过某个门限值（与序列长度狀有关），否则将其归
// 入不正常的范围；如果不正常的尖峰个数超过了允许值，即可认为待检序列是不随机的。
// 根据GMT 0005-2021规范，常见数据检测规模为10^8、10^6、2*10^4 bit
func DiscreteFourierTransformTest(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// 根据GMT 0005-2021规范的数据规模选择优化策略
	switch {
	case n >= LargeScale:
		return discreteFourierTransformTestOptimized(bits, true)
	case n >= MediumScale:
		return discreteFourierTransformTestOptimized(bits, false)
	case n >= SmallScale:
		return discreteFourierTransformTestOptimized(bits, false)
	default:
		// 小于2*10^4 bit的数据使用标准算法
		return discreteFourierTransformTestSmall(bits)
	}
}

// discreteFourierTransformTest 离散傅里叶检测，非分块处理版本
func discreteFourierTransformTest(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// Step 1, 2
	N := ceilPow2(n)
	rr := make([]complex128, N)
	for i := 0; i < n; i++ {
		if bits[i] {
			rr[i] = complex(1.0, 0)
		} else {
			rr[i] = complex(-1.0, 0)
		}
	}

	// 傅里叶变换
	f, err := fft.New(N)
	if err != nil {
		panic(err)
	}
	f.Transform(rr)

	// Step 4
	T := math.Sqrt(2.995732274 * float64(n))

	// Step 5
	N_0 := 0.95 * float64(n) / 2

	// Step 6
	var N_1 int = 0
	for i := 0; i < n/2-1; i++ {
		if cmplx.Abs(rr[i]) < T {
			N_1++
		}
	}

	// Step 7
	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(2.0*n)/3.8)
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}

// discreteFourierTransformTestSmall 小数据集的优化实现
func discreteFourierTransformTestSmall(bits []bool) (float64, float64) {
	n := len(bits)

	// Step 1, 2
	N := ceilPow2(n)
	rr := make([]complex128, N)

	// 优化数据转换：使用预分配的常量
	ones := complex(1.0, 0)
	minusOnes := complex(-1.0, 0)

	for i := 0; i < n; i++ {
		if bits[i] {
			rr[i] = ones
		} else {
			rr[i] = minusOnes
		}
	}

	// 傅里叶变换
	f, err := fft.New(N)
	if err != nil {
		panic(err)
	}
	f.Transform(rr)

	// Step 4 - 预计算常量
	T := math.Sqrt(2.995732274 * float64(n))

	// Step 5
	N_0 := 0.95 * float64(n) / 2

	// Step 6 - 优化循环，避免重复计算
	var N_1 int = 0
	T_squared := T * T
	limit := n/2 - 1

	for i := 0; i < limit; i++ {
		// 使用平方比较避免开方运算
		real := real(rr[i])
		imag := imag(rr[i])
		if real*real+imag*imag < T_squared {
			N_1++
		}
	}

	// Step 7 - 预计算分母
	denominator := math.Sqrt(0.95 * 0.05 * float64(2.0*n) / 3.8)
	V := (float64(N_1) - N_0) / denominator
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}

// discreteFourierTransformTestOptimized 优化的离散傅里叶检测实现
// 使用预置FFT表加速，支持GMT 0005-2021规范的数据规模
func discreteFourierTransformTestOptimized(bits []bool, isLargeScale bool) (float64, float64) {
	n := len(bits)

	// Step 1, 2 - 计算最接近的2的幂次
	N := ceilPow2(n)
	rr := make([]complex128, N)

	// 优化数据转换：使用预分配常量和批量处理
	ones := complex(1.0, 0)
	minusOnes := complex(-1.0, 0)

	// 根据数据规模选择不同的批量大小
	batchSize := 1024
	if isLargeScale {
		batchSize = 8192 // 大规模数据使用更大的批量
	}

	// 批量转换数据
	for i := 0; i < n; i += batchSize {
		end := i + batchSize
		if end > n {
			end = n
		}

		for j := i; j < end; j++ {
			if bits[j] {
				rr[j] = ones
			} else {
				rr[j] = minusOnes
			}
		}
	}

	// 使用预置FFT表进行傅里叶变换
	f, err := getFFT(N)
	if err != nil {
		panic(err)
	}
	f.Transform(rr)

	// Step 4 - 预计算常量
	T := math.Sqrt(2.995732274 * float64(n))

	// Step 5
	N_0 := 0.95 * float64(n) / 2

	// Step 6 - 优化循环计算
	var N_1 int = 0
	limit := n/2 - 1
	T_squared := T * T // 避免重复计算平方

	for i := 0; i < limit; i++ {
		// 使用平方比较避免开方运算
		real := real(rr[i])
		imag := imag(rr[i])
		if real*real+imag*imag < T_squared {
			N_1++
		}
	}

	// Step 7 - 预计算分母
	denominator := math.Sqrt(0.95 * 0.05 * float64(2.0*n) / 3.8)
	V := (float64(N_1) - N_0) / denominator
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}
