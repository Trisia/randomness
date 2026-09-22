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

// coreDiscreteFourierTransform 是离散傅里叶检测的实现（packed 内核）。
func coreDiscreteFourierTransform(s BitSeq) (float64, float64) {
	n := s.n
	if n == 0 {
		panic("please provide test bits")
	}

	N := ceilPow2(n)
	// 统计区间是 [0, n/2-1)
	limit := n/2 - 1
	N_1 := 0

	// limit > 0 时执行变换；n <= 3 时区间为空，无需变换。
	if limit > 0 {
		H := N >> 1
		z := make([]complex128, H)
		for j := 0; j < H; j++ {
			var re, im float64
			if i := 2 * j; i < n {
				if s.bit(i) == 1 {
					re = 1
				} else {
					re = -1
				}
			}
			if i := 2*j + 1; i < n {
				if s.bit(i) == 1 {
					im = 1
				} else {
					im = -1
				}
			}
			z[j] = complex(re, im)
		}

		f, err := getFFT(H)
		if err != nil {
			panic(err)
		}
		f.Transform(z)

		// Step 4, 5
		W0 := complex(math.Cos(-2*math.Pi/float64(N)), math.Sin(-2*math.Pi/float64(N)))
		T2 := 2.995732274 * float64(n) // T = sqrt(2.995732274*n)，比较平方避免开方
		for k := 0; k < limit; k++ {
			zk := z[k]
			zr, zi := real(zk), imag(zk)
			mr, mi := real(z[(H-k)%H]), -imag(z[(H-k)%H])
			Pr, Pi := zr+mr, zi+mi
			Mr, Mi := zr-mr, zi-mi
			W := f.E[k>>1]
			if k&1 == 1 {
				W *= W0
			}
			Qr := real(W)*Mr - imag(W)*Mi
			Qi := real(W)*Mi + imag(W)*Mr
			re := Pr + Qi
			im := Pi - Qr
			if (re*re+im*im)/4 < T2 {
				N_1++
			}
		}
	}

	// Step 5, 7
	N_0 := 0.95 * float64(n) / 2
	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(2.0*n)/3.8)
	return math.Erfc(math.Abs(V)), math.Erfc(V) / 2
}

// DiscreteFourierTransform 离散傅里叶检测
func DiscreteFourierTransform(data []byte) *TestResult {
	p, q := coreDiscreteFourierTransform(BitSeqFromBytes(data))
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBytes 离散傅里叶检测
func DiscreteFourierTransformTestBytes(data []byte) (float64, float64) {
	return coreDiscreteFourierTransform(BitSeqFromBytes(data))
}

// DiscreteFourierTransformTest 离散傅里叶检测
// 离散傅立叶变换检测使用频谱的方法来检测序列的随机性。对待检序列进行傅立叶变换后可以得
// 到尖峰高度，根据随机性的假设，这个尖峰高度不能超过某个门限值（与序列长度狀有关），否则将其归
// 入不正常的范围；如果不正常的尖峰个数超过了允许值，即可认为待检序列是不随机的。
// 根据GMT 0005-2021规范，常见数据检测规模为10^8、10^6、2*10^4 bit
//
// Deprecated: 请改用 DiscreteFourierTransformTestBitSeq——本函数接受 []bool（1 字节/位），并在内部再打包成 BitSeq，
// 同一份数据被转换两次。推荐写法是转换一次后复用：
// s := randomness.BitSeqFromBytes(buf)，之后调用 DiscreteFourierTransformTestBitSeq 系列。
// 若手上已经是 []bool，可用 randomness.BitSeqFromBools 转换一次后同样复用。
func DiscreteFourierTransformTest(bits []bool) (float64, float64) {
	return coreDiscreteFourierTransform(BitSeqFromBools(bits))
}

// DiscreteFourierTransformBitSeq 离散傅里叶检测
func DiscreteFourierTransformBitSeq(s BitSeq) *TestResult {
	p, q := coreDiscreteFourierTransform(s)
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBitSeq 离散傅里叶检测
func DiscreteFourierTransformTestBitSeq(s BitSeq) (float64, float64) {
	return coreDiscreteFourierTransform(s)
}
