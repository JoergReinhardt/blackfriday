package types

import (
	"math/big"
	"math/rand"
)

/////////////////////////////////////////////////////////////////////////
// INTEGER
type Integer val

func (i Integer) Eval() Evaluable { return i }
func (i Integer) Serialize() []byte {
	defer discardInt(i())
	return []byte(val(i)().String())
}
func (i Integer) String() string  { return val(i).text(10) }
func (i Integer) Type() ValueType { return INTEGER }
func (i Integer) Int64() int64    { return val(i).int64() }
func (i Integer) Add(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).add(i(), y())).(val).Integer()
}
func (i Integer) Sub(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).sub(i(), y())).(val).Integer()
}
func (i Integer) Cmp(x Integer) int {
	defer discardInt(x())
	a := wrap(intPool.Get().(*big.Int).Set(i())).(val).Integer()
	return a().Cmp(x())
}
func (i Integer) Div(y Integer) Integer {
	defer discardInt(y())
	return wrap(val(i).div(i(), y())).(val).Integer()
}
func (i Integer) DivMod(y Integer) Pair {
	// assume base ten arrithmetic
	m := Value(10).(val).Integer()
	defer discardInt(i(), y(), m())
	return func() [2]Evaluable { return [2]Evaluable{wrap(i), wrap(m)} }
}
func (i Integer) Exp(y Integer) Integer {
	m := Value(10).(val).Integer()
	defer discardInt(i(), y(), m())
	return wrap(val(i).exp(i(), y(), m())).(val).Integer()
}
func (i Integer) Mod(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).mod(i(), y())).(val).Integer()
}

func (i Integer) ModInverse(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).modInverse(i(), y())).(val).Integer()
}

func (i Integer) ModSqrt(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).modSqrt(i(), y())).(val).Integer()
}
func (i Integer) Mul(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).mul(i(), y())).(val).Integer()
}

//func (i Integer) MulRange(a, b int64) Integer {
//	return wrap(val(i).mulRange(a, b)).(val).Integer()
//}
func (i Integer) Neg(x Integer) Integer {
	return wrap(val(i).neg(x())).(val).Integer()
}

func (i Integer) ProbablyPrime(n Integer) Bool {
	return Value(val(i).probablyPrime(int(n.Int64()))).(val).Bool()
}

func (i Integer) Quo(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).quo(i(), y())).(val).Integer()
}

func (i Integer) QuoRem(y Integer) Pair {
	r := intPool.Get().(val).Integer()
	defer discardInt(i(), y(), r())
	a, b := val(i).quoRem(i(), y(), r())
	ret := pairPool.Get().(Pair).SetKey(wrap(a).(val).Integer())
	ret = ret.SetValue(wrap(b).(val).Integer())
	return ret
}

func (i Integer) Rand() Integer {
	var rnd = rand.New(rand.NewSource(i().Int64()))
	defer discardInt(i())
	return wrap(val(i).rand(rnd, i())).(val).Integer()
}
func (i Integer) Rem(y Integer) Integer {
	defer discardInt(i(), y())
	return wrap(val(i).rem(i(), y())).(val).Integer()
}
func (i Integer) Set(x Integer) Integer {
	defer discardInt(x())
	return wrap(val(i).set(x())).(val).Integer()
}
func (i Integer) SetInt64(x int64) Integer {
	return wrap(val(i).setInt64(x)).(val).Integer()
}
func (i Integer) SetUint64(x uint64) Integer {
	return wrap(val(i).setUint64(x)).(val).Integer()
}
func (i Integer) SetString(s string, b int) (Integer, bool) {
	x, y := val(i).setString(s, b)
	return wrap(x).(val).Integer(), y
}
func (i Integer) Uint64() uint64 { return val(i).uint64() }
