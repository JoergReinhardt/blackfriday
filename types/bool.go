package types

import (
	"math/big"
)

/////// BOOL ////////////////////////////////////////////////////////////
// booleans allways come in slices (encoded as big.int, handled bitwise using
// uint representation)
type Bool val

func (u Bool) Eval() Evaluable   { return u }
func (u Bool) Serialize() []byte { return val(u).bytes() }
func (u Bool) String() string {
	if u().Int64() > int64(0) {
		return "true"
	} else {
		return "false"
	}
}
func (u Bool) Native() bool {
	if u().Int64() > int64(0) {
		return true
	} else {
		return false
	}
}
func (u Bool) Type() ValueType { return BOOL }
func (u Bool) And(y Bool) Bool {
	defer discardInt(u(), y())
	return wrap(val(u).and(u(), y())).(val).Bool()
}
func (u Bool) AndNot(y Bool) Bool {
	defer discardInt(u(), y())
	return wrap(val(u).andNot(u(), y())).(val).Bool()
}
func (u Bool) Not(x Bool) Bool {
	defer discardInt(x())
	return wrap(val(u).not(x())).(val).Bool()
}
func (u Bool) Or(y Bool) Bool {
	defer discardInt(u(), y())
	return wrap(val(u).or(u(), y())).(val).Bool()
}
func (u Bool) Xor(y Bool) Bool {
	defer discardInt(u(), y())
	return wrap(val(u).xor(u(), y())).(val).Bool()
}

// sets a Bools value to the value of a passed Bool
func (u Bool) SetBool(x Bool) Bool {
	// discard parameter and old version
	defer discardInt(u())
	// pre allocate return value

	if u.And(x)().Int64() > 0 {
		return Value(+1).(val).Bool()
	} else {
		return Value(-1).(val).Bool()
	}
}
func (u Bool) SetBoolSlice(x ...Bool) (r BitFlag) {
	var res *big.Int
	// range over slice of Bools
	for _, i := range x {
		// use native to get a true/false value as the if condition
		if i.Native() { // left shift either uint one…
			res = val(r).lsh(r(), 1)
		} else { // or uint zero to preexisting uint and overwrite it
			// with the result
			res = val(r).lsh(r(), 0)
		}
	}
	return wrap(res).(val).bitFlag()
}
func (u Bool) SetBoolNative(x bool) (r Bool) {
	if x {
		r = wrap(intPool.Get().(val).setInt64(+1)).(val).Bool()
	} else {
		r = wrap(intPool.Get().(val).setInt64(-1)).(val).Bool()
	}
	return r
}
func (u Bool) SetBoolSliceNative(x ...bool) (r BitFlag) {
	var res *big.Int
	// range over slice of Bools
	for _, i := range x {
		// use native to get a true/false value as the if condition
		if i { // left shift either uint one…
			res = val(r).lsh(r(), 1)
		} else { // or uint zero to preexisting uint and overwrite it
			// with the result
			res = val(r).lsh(r(), 0)
		}
	}
	return wrap(res).(val).bitFlag()
}
func (u Bool) SetInteger(x Integer) Bool {
	// discard parameter and old version
	defer discardInt(x(), u())
	// pre allocate return value
	var res Bool
	if x().Int64() > 0 {
		res = wrap(intPool.Get().(val).setInt64(1)).(val).Bool()
	} else {
		res = wrap(intPool.Get().(val).setInt64(-1)).(val).Bool()
	}
	return res
}
func (u Bool) SetIntegerNative(x int64) Bool {
	// discard parameter and old version
	defer discardInt(u())
	// pre allocate return value
	var res Bool
	if x > 0 {
		res = wrap(intPool.Get().(val).setInt64(1)).(val).Bool()
	} else {
		res = wrap(intPool.Get().(val).setInt64(-1)).(val).Bool()
	}
	return res
}

// if a uint iis passed, a bit-flag will be returned. Bit Flag is considered a
// list type and implemented among the collections.
func (u Bool) SetUintNative(x uint64) BitFlag {
	// discard parameter and old version
	defer discardInt(u())
	// pre allocate return value
	return wrap(intPool.Get().(val).setUint64(x)).(val).bitFlag()
}
