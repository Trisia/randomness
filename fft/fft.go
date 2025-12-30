// Package fft provides a fast discrete Fourier transformation algorithm.
//
// Implemented is the 1-dimensional DFT of complex input data
// for with input lengths which are powers of 2.
//
// The algorithm is non-recursive and works in-place overwriting
// the input array.
//
// Before doing the transform on acutal data, allocate
// an FFT object with t := fft.New(N) where N is the length of the
// input array.
// Then multiple calls to t.Transform(x) can be done with
// different input vectors having the same length.

package fft

import (
	"fmt"
	"math"
)

/*
LICENSE from https://github.com/ktye/fft
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
*/

// ALGORITHM
// Example of the alogrithm with N=8 (P=3)
// Butterfly diagram:
//
//      1st stage p=0               2nd state p=1              3rd stage p=2
// IN +------------------+     +--------------------+     +-----------------------+
//                     overwrite                  overwrite                       output
// x0 -\/- x0 + E20 * x1 -> x0 -\  /- x0 + E40 * x2 -> x0 -\     /- x0 + E80 * x4 -> x0
// x1 -/\- x0 + E21 * x1 -> x1 -\\//- x1 + E41 * x3 -> x1 -\\   //- x1 + E81 * x5 -> x1
//                               /\                         \\ //
// x2 -\/- x0 + E22 * x1 -> x2 -//\\- x0 + E42 * x2 -> x2 -\ \\/ /- x2 + E82 * x6 -> x2
// x3 -/\- x0 + E23 * x1 -> x3 -/  \- x1 + E43 * x3 -> x3 -\\/\\//- x3 + E83 * x7 -> x3
//                                                          \\/\\
// x4 -\/- x0 + E24 * x1 -> x4 -\  /- x4 + E44 * x6 -> x4 -//\\/\\- x0 + E84 * x4 -> x4
// x5 -/\- x0 + E25 * x1 -> x5 -\\//- x5 + E45 * x7 -> x5 -/ /\\ \- x1 + E85 * x5 -> x5
//                               /\                         // \\
// x6 -\/- x0 + E26 * x1 -> x6 -//\\- x4 + E46 * x6 -> x6 -//   \\- x2 + E86 * x6 -> x6
// x7 -/\- x0 + E27 * x1 -> x7 -/  \- x5 + E47 * x7 -> x7 -/     \- x3 + E87 * x7 -> x7
//
// Enk are the N complex roots of 1 which were precomputed in E[0]..E[N-1].
// The stride s is N/n, and the index in E is k*s mod N,
// so   E21 of the first stage  is E[1*8/2 mod 8] = E[4]. These are +/- 1 alternating.
// and  E45 of the second stage is E[5*8/4 mod 8] = E[2]. These are 1,-i,-1,i and again.
// E8k  are all the roots (with stride=1) in increasing order: E[k].
//
// Before starting with the first stage, the input array must be
// permutated by the bit-inverted order.

type FFT struct {
	N    int          // Fft length, power of 2.
	p    int          // Base-2 exponent: 2^p = N.
	E    []complex128 // Precomputed roots table, length N.
	perm []int        // Index permutation vector for the input array.
}

func New(N int) (f FFT, err error) {
	var p int
	N, p, err = lastPow2(N)
	if err != nil {
		return f, err
	}
	f = FFT{
		N:    N,
		p:    p,
		E:    roots(N),
		perm: permutationIndex(p),
	}
	return f, nil
}

// Transform Forward transform.
// The forward transform overwrites the input array.
func (f FFT) Transform(x []complex128) []complex128 {
	if len(x) != f.N {
		panic("Input dimension mismatches: FFT is not initialized, or called with wrong input.")
	}

	inputPermutation(x, f.perm)

	butterfly := func(k, o, l, s int) {
		i := k + o
		j := i + l
		x[i], x[j] = x[i]+f.E[k*s]*x[j], x[i]+f.E[s*(k+l)]*x[j]
	}

	n := 1
	s := f.N
	for p := 1; p <= f.p; p++ {
		s >>= 1
		for b := 0; b < s; b++ {
			o := 2 * b * n
			for k := 0; k < n; k++ {
				butterfly(k, o, n, s)
			}
		}
		n <<= 1
	}
	return x
}

// Inverse is the backwards transform.
func (f FFT) Inverse(x []complex128) []complex128 {
	if len(x) != f.N {
		panic("FFT is not initialized, or called with wrong input. Input dimension mismatches.")
	}
	// Reverse the input vector
	for i := 1; i < f.N/2; i++ {
		j := f.N - i
		x[i], x[j] = x[j], x[i]
	}

	// Do the transform.
	f.Transform(x)

	// Scale the output by 1/N
	invN := 1.0 / float64(f.N)
	for i := range x {
		x[i] *= complex(invN, 0)
	}
	return x
}

// permutationIndex builds the bit-inverted index vector,
// which is needed to permutate the input data.
func permutationIndex(P int) []int {
	N := 1 << uint(P)
	index := make([]int, N)
	index[0] = 0 // Initial sequence for N=1
	n := 1
	// For every next power of two, the
	// sequence is multiplied by 2 inplace.
	// Then the result is also appended to the
	// end and increased by one.
	for p := 0; p < P; p++ {
		for i := 0; i < n; i++ {
			index[i] <<= 1
			index[i+n] = index[i] + 1
		}
		n <<= 1
	}
	return index
}

// inputPermutation permutes the input vector in the order
// needed for the transformation.
func inputPermutation(x []complex128, p []int) {
	for i := range p {
		if k := p[i]; i < k {
			x[i], x[k] = x[k], x[i]
		}
	}
}

// roots computes the complex-roots-of 1 table of length N.
func roots(N int) []complex128 {
	E := make([]complex128, N)
	for n := 0; n < N; n++ {
		phi := -2.0 * math.Pi * float64(n) / float64(N)
		s, c := math.Sincos(phi)
		E[n] = complex(c, s)
	}
	return E
}

// lastPow2 return the last power of 2 smaller or equal
// to the given N, and it's base-2 logarithm.
func lastPow2(N int) (n, p int, err error) {
	// On 32-bit systems, complex128 arrays are limited in size due to address space constraints.
	// complex128 is 16 bytes, so limit N to 2^25 on 32-bit (~512MB) to avoid allocation panic.
	maxdim := 1 << 27
	if ^uint(0)>>63 == 0 { // 32-bit system
		maxdim = 1 << 25
	}
	if N < 2 {
		return n, p, fmt.Errorf("fft input length must be >= 2")
	} else if N > maxdim {
		return n, p, fmt.Errorf("fft input length must be < %d. It is: %d", maxdim, N)
	}
	i := 2
	for p = 1; ; p++ {
		j := i << 1
		if j > N {
			return i, p, nil
		}
		i = j
	}
}
