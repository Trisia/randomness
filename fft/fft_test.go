package fft

import (
	"math"
	"math/cmplx"
	"testing"
)

// 本文件只用预置的输入/输出对照验证 fft 包的公开接口：
//
//	func New(N int) (f FFT, err error)
//	type FFT struct{ N int; E []complex128 }
//	func (FFT) Transform(x []complex128) []complex128
//	func (FFT) Inverse(x []complex128) []complex128
//
// 期望值取自手算结果或闭式解，与实现无关。容差 tol 仅用于吸收浮点舍入：
// 根表由 math.Sincos 生成，cos(π/2) 之类的值并非精确的 0。

const tol = 1e-12

// sqrt1_2 是 √2/2，N=8 的旋转因子用到。
var sqrt1_2 = math.Sqrt2 / 2

// check 逐项比较变换结果与预置期望值。
func check(t *testing.T, name string, got, want []complex128) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: 长度 %d, 期望 %d", name, len(got), len(want))
		return
	}
	for i := range want {
		if cmplx.Abs(got[i]-want[i]) > tol {
			t.Errorf("%s: 第 %d 项 %v, 期望 %v", name, i, got[i], want[i])
		}
	}
}

// TestNew 验证长度向下取整、非法长度报错，以及预置的根表 E。
func TestNew(t *testing.T) {
	t.Run("长度向下取整到 2 的幂", func(t *testing.T) {
		cases := []struct {
			in   int
			want int
		}{
			{2, 2}, {3, 2}, {4, 4}, {5, 4}, {7, 4},
			{8, 8}, {1000, 512}, {1024, 1024}, {1025, 1024},
		}
		for _, c := range cases {
			f, err := New(c.in)
			if err != nil {
				t.Errorf("New(%d) 返回错误: %v", c.in, err)
				continue
			}
			if f.N != c.want {
				t.Errorf("New(%d).N = %d, 期望 %d", c.in, f.N, c.want)
			}
			if len(f.E) != c.want {
				t.Errorf("New(%d): len(E) = %d, 期望 %d", c.in, len(f.E), c.want)
			}
		}
	})

	t.Run("非法长度报错", func(t *testing.T) {
		for _, n := range []int{-4, 0, 1} {
			if _, err := New(n); err == nil {
				t.Errorf("New(%d) 未返回错误", n)
			}
		}
	})

	t.Run("根表 E", func(t *testing.T) {
		cases := []struct {
			n    int
			want []complex128
		}{
			// E[m] = exp(-2πi·m/N)
			{2, []complex128{1, -1}},
			{4, []complex128{1, complex(0, -1), -1, complex(0, 1)}},
			{8, []complex128{
				1, complex(sqrt1_2, -sqrt1_2), complex(0, -1), complex(-sqrt1_2, -sqrt1_2),
				-1, complex(-sqrt1_2, sqrt1_2), complex(0, 1), complex(sqrt1_2, sqrt1_2),
			}},
		}
		for _, c := range cases {
			f, err := New(c.n)
			if err != nil {
				t.Fatalf("New(%d): %v", c.n, err)
			}
			check(t, "E", f.E, c.want)
		}
	})
}

// TestTransform 用有闭式解的预置输入验证正向变换。
func TestTransform(t *testing.T) {
	cases := []struct {
		name string
		n    int
		in   []complex128
		want []complex128
	}{
		{"N=2 直流", 2, []complex128{1, 1}, []complex128{2, 0}},
		{"N=2 斜坡", 2, []complex128{1, 2}, []complex128{3, -1}},
		{"N=2 Nyquist", 2, []complex128{1, -1}, []complex128{0, 2}},

		// 教科书用例：X = [10, -2+2i, -2, -2-2i]
		{"N=4 斜坡", 4,
			[]complex128{1, 2, 3, 4},
			[]complex128{10, complex(-2, 2), -2, complex(-2, -2)}},
		{"N=4 冲激@1", 4,
			[]complex128{0, 1, 0, 0},
			[]complex128{1, complex(0, -1), -1, complex(0, 1)}},
		{"N=4 直流", 4,
			[]complex128{1, 1, 1, 1},
			[]complex128{4, 0, 0, 0}},
		{"N=4 Nyquist", 4,
			[]complex128{1, -1, 1, -1},
			[]complex128{0, 0, 4, 0}},

		// 冲激@1 一次给出 N=8 的全部旋转因子
		{"N=8 冲激@1", 8,
			[]complex128{0, 1, 0, 0, 0, 0, 0, 0},
			[]complex128{
				1, complex(sqrt1_2, -sqrt1_2), complex(0, -1), complex(-sqrt1_2, -sqrt1_2),
				-1, complex(-sqrt1_2, sqrt1_2), complex(0, 1), complex(sqrt1_2, sqrt1_2),
			}},
		{"N=8 直流", 8,
			[]complex128{1, 1, 1, 1, 1, 1, 1, 1},
			[]complex128{8, 0, 0, 0, 0, 0, 0, 0}},
		{"N=8 Nyquist", 8,
			[]complex128{1, -1, 1, -1, 1, -1, 1, -1},
			[]complex128{0, 0, 0, 0, 8, 0, 0, 0}},
	}
	for _, c := range cases {
		f, err := New(c.n)
		if err != nil {
			t.Fatalf("New(%d): %v", c.n, err)
		}
		got := f.Transform(append([]complex128(nil), c.in...))
		check(t, c.name, got, c.want)
	}
}

// TestInverse 用有闭式解的预置输入验证逆向变换。
func TestInverse(t *testing.T) {
	cases := []struct {
		name string
		n    int
		in   []complex128
		want []complex128
	}{
		{"N=2 频域冲激@0", 2,
			[]complex128{1, 0},
			[]complex128{0.5, 0.5}},

		// TestTransform 中 N=4 斜坡用例的逆运算
		{"N=4 还原斜坡", 4,
			[]complex128{10, complex(-2, 2), -2, complex(-2, -2)},
			[]complex128{1, 2, 3, 4}},
		{"N=4 频域直流", 4,
			[]complex128{4, 0, 0, 0},
			[]complex128{1, 1, 1, 1}},
		{"N=4 频域冲激@0", 4,
			[]complex128{1, 0, 0, 0},
			[]complex128{0.25, 0.25, 0.25, 0.25}},
		{"N=4 频域 Nyquist", 4,
			[]complex128{0, 0, 4, 0},
			[]complex128{1, -1, 1, -1}},
	}
	for _, c := range cases {
		f, err := New(c.n)
		if err != nil {
			t.Fatalf("New(%d): %v", c.n, err)
		}
		got := f.Inverse(append([]complex128(nil), c.in...))
		check(t, c.name, got, c.want)
	}
}

// TestPanicsOnWrongLength 验证长度与 FFT 不匹配时按约定 panic。
func TestPanicsOnWrongLength(t *testing.T) {
	f, err := New(4)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		fn   func()
	}{
		{"Transform 过短", func() { f.Transform(make([]complex128, 3)) }},
		{"Transform 过长", func() { f.Transform(make([]complex128, 5)) }},
		{"Inverse 不匹配", func() { f.Inverse(make([]complex128, 8)) }},
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: 期望 panic", c.name)
				}
			}()
			c.fn()
		}()
	}
}
