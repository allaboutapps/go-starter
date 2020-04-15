// +build !go1.12

// The following is copied from math/big/bits/bits.go, licensed under the
// BSD 3-clause license: https://github.com/golang/go/blob/master/LICENSE

package arith

import "math/big"

func Add64(x, y uint64) (sum, carryOut uint64) {
	yc := y + 0 /* carry */
	sum = x + yc
	if sum < x || yc < y {
		carryOut = 1
	}
	return
}

func Sub64(x, y uint64) (diff, borrowOut uint64) {
	yb := y
	diff = x - yb
	if diff > x || yb < y {
		borrowOut = 1
	}
	return
}

func Mul64(x, y uint64) (hi, lo uint64) {
	const mask32 = 1<<32 - 1
	x0 := x & mask32
	x1 := x >> 32
	y0 := y & mask32
	y1 := y >> 32
	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & mask32
	w2 := t >> 32
	w1 += x0 * y1
	hi = x1*y1 + w2 + w1>>32
	lo = x * y
	return
}

func mulWW(x, y big.Word) (z1, z0 big.Word) {
	x0 := x & _M2
	x1 := x >> _W2
	y0 := y & _M2
	y1 := y >> _W2
	w0 := x0 * y0
	t := x1*y0 + w0>>_W2
	w1 := t & _M2
	w2 := t >> _W2
	w1 += x0 * y1
	z1 = x1*y1 + w2 + w1>>_W2
	z0 = x * y
	return
}

func addWW(x, y, c big.Word) (z1, z0 big.Word) {
	yc := y + c
	z0 = x + yc
	if z0 < x || yc < y {
		z1 = 1
	}
	return
}

func subWW(x, y, c big.Word) (z1, z0 big.Word) {
	yc := y + c
	z0 = x - yc
	if z0 > x || yc < y {
		z1 = 1
	}
	return
}
