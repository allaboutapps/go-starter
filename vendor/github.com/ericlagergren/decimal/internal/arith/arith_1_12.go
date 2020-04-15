// +build go1.12

package arith

import (
	"math/big"
	"math/bits"
)

func Add64(x, y uint64) (sum, carryOut uint64) {
	return bits.Add64(x, y, 0)
}

func Sub64(x, y uint64) (diff, borrowOut uint64) {
	return bits.Sub64(x, y, 0)
}

func Mul64(x, y uint64) (hi, lo uint64) {
	return bits.Mul64(x, y)
}

func mulWW(x, y big.Word) (z1, z0 big.Word) {
	zz1, zz0 := bits.Mul(uint(x), uint(y))
	return big.Word(zz1), big.Word(zz0)
}

func addWW(x, y, c big.Word) (z1, z0 big.Word) {
	zz1, zz0 := bits.Add(uint(x), uint(y), uint(c))
	return big.Word(zz0), big.Word(zz1)
}

func subWW(x, y, c big.Word) (z1, z0 big.Word) {
	zz1, zz0 := bits.Sub(uint(x), uint(y), uint(c))
	return big.Word(zz0), big.Word(zz1)
}
