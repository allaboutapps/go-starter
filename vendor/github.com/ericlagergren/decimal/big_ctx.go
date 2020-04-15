package decimal

import (
	"math"
	"math/big"

	"github.com/ericlagergren/decimal/internal/arith"
	cst "github.com/ericlagergren/decimal/internal/c"
)

// Add sets z to x + y and returns z.
func (c Context) Add(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	if x.IsFinite() && y.IsFinite() {
		z.form = c.add(z, x, x.form, y, y.form)
		return c.round(z)
	}

	// NaN + NaN
	// NaN + y
	// x + NaN
	if z.checkNaNs(x, y, addition) {
		return z
	}

	if x.form&inf != 0 {
		if y.form&inf != 0 && x.form^y.form == signbit {
			// +Inf + -Inf
			// -Inf + +Inf
			return z.setNaN(InvalidOperation, qnan, addinfinf)
		}
		// ±Inf + y
		// +Inf + +Inf
		// -Inf + -Inf
		return z.Set(x)
	}
	// x + ±Inf
	return z.Set(y)
}

// add sets z to x + y, where x and y are both finite.
//
// The (*Big).form fields are ignored and must be provided as separate arguments
// in order to facilitate Context.Sub.
func (c Context) add(z, x *Big, xform form, y *Big, yform form) form {
	// addCompact, addMixed, and addBig all require X be the "shifted" number,
	// which means X must have the greater exponent.
	hi, lo := x, y
	hisign, losign := xform, yform
	if hi.exp < lo.exp {
		hi, lo = lo, hi
		hisign, losign = losign, hisign
	}

	if sign, ok := c.tryTinyAdd(z, hi, hisign, lo, losign); ok {
		return sign
	}

	var sign form
	if hi.isCompact() {
		if lo.isCompact() {
			sign = c.addCompact(z, hi.compact, hisign, lo.compact, losign, uint64(hi.exp-lo.exp))
		} else {
			sign = c.addMixed(z, &lo.unscaled, losign, lo.exp, hi.compact, hisign, hi.exp)
		}
	} else if lo.isCompact() {
		sign = c.addMixed(z, &hi.unscaled, hisign, hi.exp, lo.compact, losign, lo.exp)
	} else {
		sign = c.addBig(z, &hi.unscaled, hisign, &lo.unscaled, losign, uint64(hi.exp-lo.exp))
	}
	z.exp = lo.exp
	return sign
}

// tryTinyAdd attempts to set z to X + Y, but only if the addition requires such
// a large shift that the result the addition would be the same if Y is replaced
// with a smaller value.
//
// For example, given
//
//    X = 5 * 10^0      // 5
//    Y = 3 * 10^-99999 // 3e-99999
//
// X would have to be shifted (multiplied) by
//
//    shift = 10 ^ (0 - (-99999)) =
//            10 ^ (0 + 99999)    =
//            10^99999
//
// which is a *large* number.
//
// If the desired precision for the addition is 16, the end result will be
// rounded down to
//
//    5.0000000000000000
//
// making the shift entirely useless.
//
// Instead, Y can be replaced with a smaller number that rounds down to the same
// result and avoids large shifts.
//
// tryTinyAdd reports whether the "tiny" addition was performed.
func (c Context) tryTinyAdd(z *Big, X *Big, Xsign form, Y *Big, Ysign form) (form, bool) {
	if X.isZero() {
		return 0, false
	}

	exp := X.exp - 1
	if xp, zp := X.Precision(), precision(c); xp <= zp {
		exp += xp - zp - 1
	}

	if Y.adjusted() >= exp {
		return 0, false
	}

	var tiny uint64
	if Y.compact != 0 {
		tiny = 1
	}

	var sign form
	if X.isCompact() {
		sign = c.addCompact(z, X.compact, Xsign, tiny, Ysign, uint64(X.exp-exp))
	} else {
		sign = c.addMixed(z, &X.unscaled, Xsign, X.exp, tiny, Ysign, exp)
	}
	z.exp = exp
	return sign, true
}

// addCompact sets z to X + Y where
//
//    X = X0 * 10^shift
//
// and returns the resulting signbit.
func (c Context) addCompact(z *Big, X0 uint64, Xsign form, Y uint64, Ysign form, shift uint64) form {
	// Test whether X0 * 10^shift fits inside a uint64. If not, fall back to
	// big.Ints.
	X, ok := arith.MulPow10(X0, shift)
	if !ok {
		X0 := z.unscaled.SetUint64(X0)
		X := arith.MulBigPow10(X0, X0, shift)
		// Because hi was promoted to a big.Int, it by definition is larger than
		// lo.
		//
		// Therefore, the resulting signbit is the same as hi's signbit.
		//
		// Furthermore, we do not need to check if the result of the operation
		// is zero.
		if Xsign == Ysign {
			z.precision = arith.BigLength(arith.Add(&z.unscaled, X, Y))
			z.compact = cst.Inflated
		} else {
			arith.Sub(&z.unscaled, X, Y)
			z.norm()
		}
		return Xsign
	}

	// If the signs are the same, then X + Y = ℤ≠0.
	if Ysign == Xsign {
		if sum, c := arith.Add64(X, Y); c == 0 {
			z.compact = sum
			if sum == cst.Inflated {
				z.unscaled.SetUint64(cst.Inflated)
			}
			z.precision = arith.Length(z.compact)
		} else {
			arith.Set(&z.unscaled, c, sum)
			z.precision = 20
			z.compact = cst.Inflated
		}
		return Xsign
	}

	sign := Xsign
	// X + (-Y) == X - Y == -(Y - X)
	// (-X) + Y == Y - X == -(X - Y)
	diff, b := arith.Sub64(X, Y)
	if b != 0 {
		sign ^= signbit
		diff = Y - X
	}

	if diff != 0 {
		z.compact = diff
		z.precision = arith.Length(z.compact)
		return sign
	}

	// On a zero result:
	//
	//    Otherwise, the sign of a zero result is 0 unless either both operands
	//    were negative or the signs of the operands were different and the
	//    rounding is round-floor.
	//
	// - http://speleotrove.com/decimal/daops.html#refaddsub
	if c.RoundingMode == ToNegativeInf {
		sign = Xsign ^ Ysign // either 0 or 1
	}
	sign |= Xsign & Ysign

	z.compact = 0
	z.precision = 1
	return sign
}

// addMixed sets z to X + Y where
//
//    X = X * 10^shift
//
// and returns the resulting signbit.
func (c Context) addMixed(z *Big, X *big.Int, Xform form, xs int, Y uint64, Yform form, ys int) form {
	if xs < ys {
		shift := uint64(ys - xs)
		Y0, ok := arith.MulPow10(Y, shift)
		if !ok {
			yb := alias(&z.unscaled, X).SetUint64(Y)
			yb = arith.MulBigPow10(yb, yb, shift)
			return c.addBig(z, X, Xform, yb, Yform, 0)
		}
		Y = Y0
	} else if xs > ys {
		X = arith.MulBigPow10(&z.unscaled, X, uint64(xs-ys))
	}

	if Xform == Yform {
		arith.Add(&z.unscaled, X, Y)
		z.precision = arith.BigLength(&z.unscaled)
		z.compact = cst.Inflated
	} else {
		// X > Y
		arith.Sub(&z.unscaled, X, Y)
		z.norm()
	}
	return Xform
}

// addBig sets z to X + Y where
//
//    X = X0 * 10^shift
//
// and returns the resulting signbit.
func (c Context) addBig(z *Big, X *big.Int, Xsign form, Y *big.Int, Ysign form, shift uint64) form {
	// Guard the call so we don't allocate (via alias) if we don't need to.
	if shift != 0 {
		X = arith.MulBigPow10(alias(&z.unscaled, Y), X, shift)
	}

	if Xsign == Ysign {
		z.unscaled.Add(X, Y)
		z.compact = cst.Inflated
		z.precision = arith.BigLength(&z.unscaled)
		return Xsign
	}

	sign := Xsign
	// X + (-Y) == X - Y == -(Y - X)
	// (-X) + Y == Y - X == -(X - Y)
	if X.Cmp(Y) >= 0 {
		z.unscaled.Sub(X, Y)
	} else {
		sign ^= signbit
		z.unscaled.Sub(Y, X)
	}

	if z.unscaled.Sign() == 0 {
		z.compact = 0
		z.precision = 1
		if c.RoundingMode == ToNegativeInf {
			sign = Xsign ^ Ysign // either 0 or 1
		}
		return sign | Xsign&Ysign
	}

	z.norm()
	return sign
}

// FMA sets z to (x * y) + u without any intermediate rounding.
func (c Context) FMA(z, x, y, u *Big) *Big {
	if z.invalidContext(c) {
		return z
	}
	// Create a temporary receiver if z == u so we handle the z.FMA(x, y, z)
	// without clobbering z partway through.
	z0 := z
	if z == u {
		z0 = WithContext(c)
	}
	c.mul(z0, x, y)
	if z0.Context.Conditions&InvalidOperation != 0 {
		return z.setShared(z0)
	}
	c.Add(z0, z0, u)
	return z.setShared(z0)
}

// Mul sets z to x * y and returns z.
func (c Context) Mul(z, x, y *Big) *Big {
	if z.invalidContext(c) {
		return z
	}
	return c.round(c.mul(z, x, y))
}

// mul is the implementation of Mul.
func (c Context) mul(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}

	sign := x.form&signbit ^ y.form&signbit

	if x.IsFinite() && y.IsFinite() {
		z.form = finite | sign
		z.exp = x.exp + y.exp

		// Multiplication is simple, so inline it.
		if x.isCompact() {
			if y.isCompact() {
				hi, lo := arith.Mul64(x.compact, y.compact)
				if hi == 0 {
					z.compact = lo
					if lo == cst.Inflated {
						z.unscaled.SetUint64(cst.Inflated)
					}
					z.precision = arith.Length(lo)
					return z
				}
				arith.Set(&z.unscaled, hi, lo)
			} else { // y.isInflated
				arith.Mul(&z.unscaled, &y.unscaled, x.compact)
			}
		} else if y.isCompact() { // x.isInflated
			arith.Mul(&z.unscaled, &x.unscaled, y.compact)
		} else {
			z.unscaled.Mul(&x.unscaled, &y.unscaled)
		}
		return z.norm()
	}

	// NaN * NaN
	// NaN * y
	// x * NaN
	if z.checkNaNs(x, y, multiplication) {
		return z
	}

	if (x.IsInf(0) && !y.isZero()) ||
		(y.IsInf(0) && !x.isZero()) ||
		(y.IsInf(0) && x.IsInf(0)) {
		// ±Inf * y
		// x * ±Inf
		// ±Inf * ±Inf
		return z.SetInf(sign != 0)
	}

	// 0 * ±Inf
	// ±Inf * 0
	return z.setNaN(InvalidOperation, qnan, mul0inf)
}

// Quantize sets z to the number equal in value and sign to z with the scale, n.
//
// In order to perform truncation, set the Context's RoundingMode to ToZero.
func (c Context) Quantize(z *Big, n int) *Big {
	if debug {
		z.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	n = -n
	if z.isSpecial() {
		if z.form&inf != 0 {
			return z.setNaN(InvalidOperation, qnan, quantinf)
		}
		z.checkNaNs(z, z, quantization)
		return z
	}

	if n > c.maxScale() || n < c.etiny() {
		return z.setNaN(InvalidOperation, qnan, quantminmax)
	}

	if z.isZero() {
		z.exp = n
		return z
	}

	shift := z.exp - n
	if z.Precision()+shift > precision(c) {
		return z.setNaN(InvalidOperation, qnan, quantprec)
	}

	z.exp = n
	if shift == 0 {
		return z
	}

	if shift < 0 {
		z.Context.Conditions |= Rounded
	}

	m := c.RoundingMode
	neg := z.form & signbit
	if z.isCompact() {
		if shift > 0 {
			if zc, ok := arith.MulPow10(z.compact, uint64(shift)); ok {
				return z.setTriple(zc, neg, n)
			}
			// shift < 0
		} else if yc, ok := arith.Pow10(uint64(-shift)); ok {
			z.quo(m, z.compact, neg, yc, 0)
			return z
		}
		z.unscaled.SetUint64(z.compact)
		z.compact = cst.Inflated
	}

	if shift > 0 {
		arith.MulBigPow10(&z.unscaled, &z.unscaled, uint64(shift))
		z.precision = arith.BigLength(&z.unscaled)
	} else {
		var r big.Int
		z.quoBig(m, &z.unscaled, neg, arith.BigPow10(uint64(-shift)), 0, &r)
	}
	return z
}

// Quo sets z to x / y and returns z.
func (c Context) Quo(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	sign := (x.form & signbit) ^ (y.form & signbit)
	if x.isSpecial() || y.isSpecial() {
		// NaN / NaN
		// NaN / y
		// x / NaN
		if z.checkNaNs(x, y, division) {
			return z
		}

		if x.form&inf != 0 {
			if y.form&inf != 0 {
				// ±Inf / ±Inf
				return z.setNaN(InvalidOperation, qnan, quoinfinf)
			}
			// ±Inf / y
			return z.SetInf(sign != 0)
		}
		// x / ±Inf
		z.Context.Conditions |= Clamped
		return z.setZero(sign, c.etiny())
	}

	if y.isZero() {
		if x.isZero() {
			// 0 / 0
			return z.setNaN(InvalidOperation|DivisionUndefined, qnan, quo00)
		}
		// x / 0
		z.Context.Conditions |= DivisionByZero
		return z.SetInf(sign != 0)
	}
	if x.isZero() {
		// 0 / y
		return c.fix(z.setZero(sign, x.exp-y.exp))
	}

	var (
		ideal = x.exp - y.exp // preferred exponent.
		m     = c.RoundingMode
		yp    = y.Precision() // stored since we might decrement it.
		zp    = precision(c)  // stored because of overhead.
	)
	if zp == UnlimitedPrecision {
		m = unnecessary
		zp = x.Precision() + int(math.Ceil(10*float64(yp)/3))
	}

	if x.isCompact() && y.isCompact() {
		if cmpNorm(x.compact, x.Precision(), y.compact, yp) {
			yp--
		}

		shift := zp + yp - x.Precision()
		z.exp = (x.exp - y.exp) - shift
		expadj := ideal - z.exp
		if shift > 0 {
			if sx, ok := arith.MulPow10(x.compact, uint64(shift)); ok {
				if z.quo(m, sx, x.form, y.compact, y.form) && expadj > 0 {
					c.simpleReduce(z)
				}
				return z
			}
			xb := z.unscaled.SetUint64(x.compact)
			xb = arith.MulBigPow10(xb, xb, uint64(shift))
			yb := new(big.Int).SetUint64(y.compact)
			if z.quoBig(m, xb, x.form, yb, y.form, yb) && expadj > 0 {
				c.simpleReduce(z)
			}
			return z
		}
		if shift < 0 {
			if sy, ok := arith.MulPow10(y.compact, uint64(-shift)); ok {
				if z.quo(m, x.compact, x.form, sy, y.form) && expadj > 0 {
					c.simpleReduce(z)
				}
				return z
			}
			yb := z.unscaled.SetUint64(y.compact)
			yb = arith.MulBigPow10(yb, yb, uint64(-shift))
			xb := new(big.Int).SetUint64(x.compact)
			if z.quoBig(m, xb, x.form, yb, y.form, xb) && expadj > 0 {
				c.simpleReduce(z)
			}
			return z
		}
		if z.quo(m, x.compact, x.form, y.compact, y.form) && expadj > 0 {
			c.simpleReduce(z)
		}
		return z
	}

	xb, yb := &x.unscaled, &y.unscaled
	if x.isCompact() {
		xb = new(big.Int).SetUint64(x.compact)
	} else if y.isCompact() {
		yb = new(big.Int).SetUint64(y.compact)
	}

	if cmpNormBig(&z.unscaled, xb, x.Precision(), yb, yp) {
		yp--
	}

	shift := zp + yp - x.Precision()
	z.exp = (x.exp - y.exp) - shift

	var tmp *big.Int
	if shift > 0 {
		tmp = alias(&z.unscaled, yb)
		xb = arith.MulBigPow10(tmp, xb, uint64(shift))
	} else if shift < 0 {
		tmp = alias(&z.unscaled, xb)
		yb = arith.MulBigPow10(tmp, yb, uint64(-shift))
	} else {
		tmp = new(big.Int)
	}

	expadj := ideal - z.exp
	if z.quoBig(m, xb, x.form, yb, y.form, alias(tmp, &z.unscaled)) && expadj > 0 {
		c.simpleReduce(z)
	}
	return z
}

func (z *Big) quo(m RoundingMode, x uint64, xneg form, y uint64, yneg form) bool {
	z.form = xneg ^ yneg
	z.compact = x / y
	z.precision = arith.Length(z.compact)

	r := x % y
	if r == 0 {
		return true
	}

	z.Context.Conditions |= Inexact | Rounded
	if m == ToZero {
		return false
	}

	rc := 1
	if hi, lo := arith.Mul64(r, 2); hi == 0 {
		rc = arith.Cmp(lo, y)
	}

	if m == unnecessary {
		z.setNaN(InvalidOperation|InvalidContext|InsufficientStorage, qnan, quotermexp)
		return false
	}

	if m.needsInc(z.compact&1 != 0, rc, xneg == yneg) {
		z.Context.Conditions |= Rounded
		z.compact++

		// Test to see if we accidentally increased precision because of rounding.
		// For example, given n = 17 and RoundingMode = ToNearestEven, rounding
		//
		//   0.9999999999999999994284
		//
		// results in
		//
		//   0.99999999999999999 (precision = 17)
		//
		// which is rounded up to
		//
		//   1.00000000000000000 (precision = 18)
		if arith.Length(z.compact) != z.precision {
			z.compact /= 10
			z.exp++
		}
	}
	return false
}

func (z *Big) quoBig(
	m RoundingMode,
	x *big.Int, xneg form,
	y *big.Int, yneg form,
	r *big.Int,
) bool {
	z.compact = cst.Inflated
	z.form = xneg ^ yneg

	q, r := z.unscaled.QuoRem(x, y, r)
	if r.Sign() == 0 {
		z.norm()
		return true
	}

	z.Context.Conditions |= Inexact | Rounded
	if m == ToZero {
		z.norm()
		return false
	}

	var rc int
	rv := r.Uint64()
	// Drop into integers if possible.
	if r.IsUint64() && y.IsUint64() && rv <= math.MaxUint64/2 {
		rc = arith.Cmp(rv*2, y.Uint64())
	} else {
		rc = r.Mul(r, cst.TwoInt).CmpAbs(y)
	}

	if m == unnecessary {
		z.setNaN(InvalidOperation|InvalidContext|InsufficientStorage, qnan, quotermexp)
		return false
	}

	if m.needsInc(q.Bit(0) != 0, rc, xneg == yneg) {
		z.Context.Conditions |= Rounded
		z.precision = arith.BigLength(q)
		arith.Add(q, q, 1)
		if arith.BigLength(q) != z.precision {
			q.Quo(q, cst.TenInt)
			z.exp++
		}
	}
	z.norm()
	return false
}

// QuoInt sets z to x / y with the remainder truncated. See QuoRem for more
// details.
func (c Context) QuoInt(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	sign := (x.form & signbit) ^ (y.form & signbit)
	if x.IsFinite() && y.IsFinite() {
		if y.isZero() {
			if x.isZero() {
				// 0 / 0
				return z.setNaN(InvalidOperation|DivisionUndefined, qnan, quo00)
			}
			// x / 0
			z.Context.Conditions |= DivisionByZero
			return z.SetInf(sign != 0)
		}
		if x.isZero() {
			// 0 / y
			return c.fix(z.setZero(sign, 0))
		}
		z, _ = c.quorem(z, nil, x, y)
		z.exp = 0
		if z.Precision() > precision(c) {
			return z.setNaN(DivisionImpossible, qnan, quointprec)
		}
		return z
	}

	// NaN / NaN
	// NaN / y
	// x / NaN
	if z.checkNaNs(x, y, division) {
		return z
	}

	if x.form&inf != 0 {
		if y.form&inf != 0 {
			// ±Inf / ±Inf
			return z.setNaN(InvalidOperation, qnan, quoinfinf)
		}
		// ±Inf / y
		return z.SetInf(sign != 0)
	}
	// x / ±Inf
	return z.setZero(sign, 0)
}

// QuoRem sets z to the quotient x / y and r to the remainder x % y, such that
// x = z * y + r, and returns the pair (z, r).
func (c Context) QuoRem(z, x, y, r *Big) (*Big, *Big) {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		r.invalidContext(c)
		return z, r
	}

	sign := (x.form & signbit) ^ (y.form & signbit)
	if x.IsFinite() && y.IsFinite() {
		if y.isZero() {
			if x.isZero() {
				// 0 / 0
				z.setNaN(InvalidOperation|DivisionUndefined, qnan, quo00)
				r.setNaN(InvalidOperation|DivisionUndefined, qnan, quo00)
			}
			// x / 0
			z.Context.Conditions |= DivisionByZero
			r.Context.Conditions |= DivisionByZero
			return z.SetInf(sign != 0), r.SetInf(x.Signbit())
		}
		if x.isZero() {
			// 0 / y
			z.setZero((x.form^y.form)&signbit, 0)
			r.setZero(x.form, y.exp-x.exp)
			return c.fix(z), c.fix(r)
		}
		return c.quorem(z, r, x, y)
	}

	// NaN / NaN
	// NaN / y
	// x / NaN
	if z.checkNaNs(x, y, division) {
		return z, r.Set(z)
	}

	if x.form&inf != 0 {
		if y.form&inf != 0 {
			// ±Inf / ±Inf
			z.setNaN(InvalidOperation, qnan, quoinfinf)
			return z, r.Set(z)
		}
		// ±Inf / y
		return z.SetInf(sign != 0), r.SetInf(x.form&signbit != 0)
	}
	// x / ±Inf
	z.Context.Conditions |= Clamped
	z.setZero(sign, c.etiny())
	r.setZero(x.form&signbit, 0)
	return z, r
}

func (c Context) quorem(z0, z1, x, y *Big) (*Big, *Big) {
	m := c.RoundingMode
	zp := precision(c)

	if x.adjusted()-y.adjusted() > zp {
		if z0 != nil {
			z0.setNaN(DivisionImpossible, qnan, quorem_)
		}
		if z1 != nil {
			z1.setNaN(DivisionImpossible, qnan, quorem_)
		}
		return z0, z1
	}

	z := z0
	if z == nil {
		z = z1
	}

	if x.isCompact() && y.isCompact() {
		shift := x.exp - y.exp
		if shift > 0 {
			if sx, ok := arith.MulPow10(x.compact, uint64(shift)); ok {
				return m.quorem(z0, z1, sx, x.form, y.compact, y.form)
			}
			xb := z.unscaled.SetUint64(x.compact)
			xb = arith.MulBigPow10(xb, xb, uint64(shift))
			yb := new(big.Int).SetUint64(y.compact)
			return m.quoremBig(z0, z1, xb, x.form, yb, y.form)
		}
		if shift < 0 {
			if sy, ok := arith.MulPow10(y.compact, uint64(-shift)); ok {
				return m.quorem(z0, z1, x.compact, x.form, sy, y.form)
			}
			yb := z.unscaled.SetUint64(y.compact)
			yb = arith.MulBigPow10(yb, yb, uint64(-shift))
			xb := new(big.Int).SetUint64(x.compact)
			return m.quoremBig(z0, z1, xb, x.form, yb, y.form)
		}
		return m.quorem(z0, z1, x.compact, x.form, y.compact, y.form)
	}

	xb, yb := &x.unscaled, &y.unscaled
	if x.isCompact() {
		xb = new(big.Int).SetUint64(x.compact)
	} else if y.isCompact() {
		yb = new(big.Int).SetUint64(y.compact)
	}

	shift := x.exp - y.exp
	if shift > 0 {
		tmp := alias(&z.unscaled, yb)
		xb = arith.MulBigPow10(tmp, xb, uint64(shift))
	} else {
		tmp := alias(&z.unscaled, xb)
		yb = arith.MulBigPow10(tmp, yb, uint64(-shift))
	}
	return m.quoremBig(z0, z1, xb, x.form, yb, y.form)
}

// TODO(eric): quorem and quoremBig should not be methods on RoundingMode

func (m RoundingMode) quorem(
	z0, z1 *Big,
	x uint64, xneg form, y uint64, yneg form,
) (*Big, *Big) {
	if z0 != nil {
		z0.setTriple(x/y, xneg^yneg, 0)
	}
	if z1 != nil {
		z1.setTriple(x%y, xneg, 0)
	}
	return z0, z1
}

func (m RoundingMode) quoremBig(
	z0, z1 *Big,
	x *big.Int, xneg form,
	y *big.Int, yneg form,
) (*Big, *Big) {
	if z0 == nil {
		z1.unscaled.Rem(x, y)
		z1.form = xneg
		return z0, z1.norm()
	}

	if z1 != nil {
		z0.unscaled.QuoRem(x, y, &z1.unscaled)
		z1.form = xneg
		z1.norm()
	} else {
		z0.unscaled.QuoRem(x, y, new(big.Int))
	}
	z0.form = xneg ^ yneg
	return z0.norm(), z1
}

// Reduce reduces a finite z to its most simplest form.
func (c Context) Reduce(z *Big) *Big {
	if debug {
		z.validate()
	}
	c.Round(z)
	return c.simpleReduce(z)
}

// simpleReduce is the same as Reduce, but it does not round prior to reducing
// the decimal.
func (c Context) simpleReduce(z *Big) *Big {
	if z.isSpecial() {
		// Same semantics as plus(z), i.e. z + 0.
		z.checkNaNs(z, z, reduction)
		return z
	}

	if z.isZero() {
		z.exp = 0
		z.precision = 1
		return z
	}

	if z.compact == cst.Inflated {
		if z.unscaled.Bit(0) != 0 {
			return z
		}

		var r big.Int
		for z.precision >= 20 {
			z.unscaled.QuoRem(&z.unscaled, cst.OneMillionInt, &r)
			if r.Sign() != 0 {
				// TODO(eric): which is less expensive? Copying z.unscaled into
				// a temporary or reconstructing if we can't divide by N?
				z.unscaled.Mul(&z.unscaled, cst.OneMillionInt)
				z.unscaled.Add(&z.unscaled, &r)
				break
			}
			z.exp += 6
			z.precision -= 6

			// Try to avoid reconstruction for odd numbers.
			if z.unscaled.Bit(0) != 0 {
				break
			}
		}

		for z.precision >= 20 {
			z.unscaled.QuoRem(&z.unscaled, cst.TenInt, &r)
			if r.Sign() != 0 {
				z.unscaled.Mul(&z.unscaled, cst.TenInt)
				z.unscaled.Add(&z.unscaled, &r)
				break
			}
			z.exp++
			z.precision--
			if z.unscaled.Bit(0) != 0 {
				break
			}
		}

		if z.precision >= 20 {
			return z.norm()
		}
		z.compact = z.unscaled.Uint64()
	}

	for ; z.compact >= 10000 && z.compact%10000 == 0; z.precision -= 4 {
		z.compact /= 10000
		z.exp += 4
	}
	for ; z.compact%10 == 0; z.precision-- {
		z.compact /= 10
		z.exp++
	}
	return z
}

// Rem sets z to the remainder x % y. See QuoRem for more details.
func (c Context) Rem(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	if x.IsFinite() && y.IsFinite() {
		if y.isZero() {
			if x.isZero() {
				// 0 / 0
				return z.setNaN(InvalidOperation|DivisionUndefined, qnan, quo00)
			}
			// x / 0
			return z.setNaN(InvalidOperation|DivisionByZero, qnan, remx0)
		}
		if x.isZero() {
			// 0 / y
			return z.setZero(x.form&signbit, min(x.exp, y.exp))
		}
		// TODO(eric): See if we can get rid of tmp. See issue #72.
		var tmp Big
		_, z = c.quorem(&tmp, z, x, y)
		z.exp = min(x.exp, y.exp)
		tmp.exp = 0
		if tmp.Precision() > precision(c) {
			return z.setNaN(DivisionImpossible, qnan, quointprec)
		}
		return c.round(z)
	}

	// NaN / NaN
	// NaN / y
	// x / NaN
	if z.checkNaNs(x, y, division) {
		return z
	}

	if x.form&inf != 0 {
		if y.form&inf != 0 {
			// ±Inf / ±Inf
			return z.setNaN(InvalidOperation, qnan, quoinfinf)
		}
		// ±Inf / y
		return z.setNaN(InvalidOperation, qnan, reminfy)
	}
	// x / ±Inf
	return z.Set(x)
}

// Round rounds z down to the Context's precision and returns z. The result is
// undefined if z is not finite. The result of Round will always be within the
// interval [⌊10**x⌋, z] where x = the precision of z.
func (c Context) Round(z *Big) *Big {
	if debug {
		z.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	n := precision(c)
	if n == UnlimitedPrecision || z.isSpecial() {
		return z
	}

	zp := z.Precision()
	if zp <= n {
		return c.fix(z)
	}

	shift := zp - n
	if shift > c.maxScale() {
		return z.xflow(c.minScale(), false, true)
	}
	z.exp += shift

	z.Context.Conditions |= Rounded

	c.shiftr(z, uint64(shift))
	return c.fix(z)
}

func (c Context) shiftr(z *Big, n uint64) bool {
	if zp := uint64(z.Precision()); n >= zp {
		z.compact = 0
		z.precision = 1
		return n == zp
	}

	if z.isZero() {
		return false
	}

	m := c.RoundingMode
	if z.isCompact() {
		if y, ok := arith.Pow10(n); ok {
			return z.quo(m, z.compact, z.form, y, 0)
		}
		z.unscaled.SetUint64(z.compact)
		z.compact = cst.Inflated
	}
	var r big.Int
	return z.quoBig(m, &z.unscaled, z.form, arith.BigPow10(n), 0, &r)
}

func (c Context) round(z *Big) *Big {
	if c.OperatingMode == GDA {
		return c.Round(z)
	}
	return c.fix(z)
}

// RoundToInt rounds z down to an integral value.
func (c Context) RoundToInt(z *Big) *Big {
	if z.isSpecial() || z.exp >= 0 {
		return z
	}
	c.Precision = z.Precision()
	return c.Quantize(z, 0)
}

// Set sets z to x and returns z. The result might be rounded, even if z == x.
func (c Context) Set(z, x *Big) *Big {
	return c.Round(z.Copy(x))
}

// SetString sets z to the value of s, returning z and a bool indicating success.
// See Big.SetString for valid formats.
func (c Context) SetString(z *Big, s string) (*Big, bool) {
	if _, ok := z.SetString(s); !ok {
		return nil, false
	}
	return c.Round(z), true
}

// Sub sets z to x - y and returns z.
func (c Context) Sub(z, x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	if z.invalidContext(c) {
		return z
	}

	if x.IsFinite() && y.IsFinite() {
		z.form = c.add(z, x, x.form, y, y.form^signbit)
		return c.round(z)
	}

	// NaN - NaN
	// NaN - y
	// x - NaN
	if z.checkNaNs(x, y, subtraction) {
		return z
	}

	if x.form&inf != 0 {
		if y.form&inf != 0 && (x.form&signbit == y.form&signbit) {
			// -Inf - -Inf
			// -Inf - -Inf
			return z.setNaN(InvalidOperation, qnan, subinfinf)
		}
		// ±Inf - y
		// -Inf - +Inf
		// +Inf - -Inf
		return z.Set(x)
	}
	// x - ±Inf
	return z.Neg(y)
}
