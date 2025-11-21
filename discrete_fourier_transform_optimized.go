package randomness

import (
	"math"
	"math/cmplx"
	"sync"

	"github.com/Trisia/randomness/fft"
)

// FFTCache 缓存FFT对象以避免重复创建
type FFTCache struct {
	cache map[int]*fft.FFT
	mutex sync.RWMutex
}

var fftCache = &FFTCache{
	cache: make(map[int]*fft.FFT),
}

// getFFT 获取或创建FFT对象
func (c *FFTCache) getFFT(N int) (*fft.FFT, error) {
	c.mutex.RLock()
	if f, exists := c.cache[N]; exists {
		c.mutex.RUnlock()
		return f, nil
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 双重检查
	if f, exists := c.cache[N]; exists {
		return f, nil
	}

	f, err := fft.New(N)
	if err != nil {
		return nil, err
	}
	c.cache[N] = &f
	return &f, nil
}

// DiscreteFourierTransformOptimized 优化版本的离散傅里叶检测
func DiscreteFourierTransformOptimized(data []byte) *TestResult {
	p, q := DiscreteFourierTransformTestBytesOptimized(data)
	return &TestResult{Name: "离散傅里叶检测", P: p, Q: q, Pass: p >= Alpha}
}

// DiscreteFourierTransformTestBytesOptimized 优化版本的离散傅里叶检测
func DiscreteFourierTransformTestBytesOptimized(data []byte) (float64, float64) {
	return DiscreteFourierTransformTestOptimized(B2bitArr(data))
}

// DiscreteFourierTransformTestOptimized 优化版本的离散傅里叶检测
func DiscreteFourierTransformTestOptimized(bits []bool) (float64, float64) {
	n := len(bits)
	if n == 0 {
		panic("please provide test bits")
	}

	// Step 1, 2
	N := ceilPow2(n)

	// 使用sync.Pool复用复数数组
	rr := fftPool.Get(N)
	defer fftPool.Put(rr)

	// 优化的数据转换：批量处理
	convertBitsToComplex(bits, rr)

	// 使用缓存的FFT对象
	f, err := fftCache.getFFT(N)
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
	halfN := n / 2
	for i := 0; i < halfN-1; i++ {
		if cmplx.Abs(rr[i]) < T {
			N_1++
		}
	}

	// Step 7，最后V除math.Sqrt(2)，放到这里提前处理，减少math.Sqrt的调用。
	V := (float64(N_1) - N_0) / math.Sqrt(0.95*0.05*float64(2.0*n)/3.8)
	P := math.Erfc(math.Abs(V))
	Q := math.Erfc(V) / 2

	return P, Q
}

// convertBitsToComplex 优化的比特到复数转换
func convertBitsToComplex(bits []bool, rr []complex128) {
	n := len(bits)

	// 使用循环展开优化
	for i := 0; i < n; i++ {
		if bits[i] {
			rr[i] = 1 + 0i
		} else {
			rr[i] = -1 + 0i
		}
	}

	// 填充剩余部分为0
	for i := n; i < len(rr); i++ {
		rr[i] = 0 + 0i
	}
}

// fftPool 复数数组对象池
var fftPool = &complex128Pool{
	pools: make(map[int]*sync.Pool),
}

type complex128Pool struct {
	pools map[int]*sync.Pool
	mutex sync.RWMutex
}

func (p *complex128Pool) Get(N int) []complex128 {
	p.mutex.RLock()
	if pool, exists := p.pools[N]; exists {
		p.mutex.RUnlock()
		if arr := pool.Get(); arr != nil {
			return arr.([]complex128)
		}
	} else {
		p.mutex.RUnlock()
	}

	// 创建新的池和数组
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// 双重检查
	if pool, exists := p.pools[N]; exists {
		if arr := pool.Get(); arr != nil {
			return arr.([]complex128)
		}
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return make([]complex128, N)
		},
	}
	p.pools[N] = pool
	return make([]complex128, N)
}

func (p *complex128Pool) Put(arr []complex128) {
	N := len(arr)
	p.mutex.RLock()
	if pool, exists := p.pools[N]; exists {
		p.mutex.RUnlock()
		// 清零数组以避免数据泄露
		for i := range arr {
			arr[i] = 0 + 0i
		}
		pool.Put(arr)
		return
	}
	p.mutex.RUnlock()
}
