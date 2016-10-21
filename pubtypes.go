package agiledoc

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	"math/rand"
	//"sync"
)

/////////////////////////////////////////////////////////////////////////
/////////// PUBLIC IMPLEMENTAIONS OF EVALUABLES /////////////////////////
///////////////// BASED ON VAL, RAT & PAIR //////////////////////////////
/////////////////////////////////////////////////////////////////////////

// BASE VALUES IMPLEMENTING FUNCTIONAL TYPES
// these functional types need to implement the absVal interface to be suitable
// base types. If called, they return their contained value. A Method set
// defined on these funtional types, implements the absVal interface, by
// manipulating the returned content. Each can implement ia couple of dynamic
// types by defining further types based on it, while overwriting and/or
// completing the method set.
//////// EMPTY //////////////////////////////////////////////////////////
type ( // functional types that form the base of all value implementations
	// empty Value
	Empty func() struct{}

	// simple type
	val func() *big.Int

	// paired types
	ratio func() *big.Rat
	pair  func() [2]Evaluable

	// collection types see collection,go
)

func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Evaluable { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }

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
func (u Bool) Type() ValueType { return BOOL }
func (u Bool) And(x, y Bool) Bool {
	defer discardInt(x(), y())
	return wrap(val(u).and(x(), y())).(val).Bool()
}
func (u Bool) AndNot(x, y Bool) Bool {
	defer discardInt(x(), y())
	return wrap(val(u).andNot(x(), y())).(val).Bool()
}
func (u Bool) Not(x Bool) Bool {
	defer discardInt(x())
	return wrap(val(u).not(x())).(val).Bool()
}
func (u Bool) Or(x, y Bool) Bool {
	defer discardInt(x(), y())
	return wrap(val(u).or(x(), y())).(val).Bool()
}
func (u Bool) Xor(x, y Bool) Bool {
	defer discardInt(x(), y())
	return wrap(val(u).xor(x(), y())).(val).Bool()
}

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
func (i Integer) Add(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).add(x(), y())).(val).Integer()
}
func (i Integer) Cmp(x Integer) Integer {
	defer discardInt(x())
	return wrap(intPool.Get().(*big.Int).Add(i(), x())).(val).Integer()
}
func (i Integer) Div(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).div(x(), y())).(val).Integer()
}
func (i Integer) DivMod(x, y, m Integer) (Integer, Integer) {
	defer discardInt(x(), y(), m())
	a, b := val(i).divMod(x(), y(), m())
	return wrap(a).(val).Integer(), wrap(b).(val).Integer()
}
func (i Integer) Exp(x, y, m Integer) Integer {
	defer discardInt(x(), y(), m())
	return wrap(val(i).exp(x(), y(), m())).(val).Integer()
}
func (i Integer) Mod(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).mod(x(), y())).(val).Integer()
}
func (i Integer) ModInverse(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).modInverse(x(), y())).(val).Integer()
}
func (i Integer) ModSqrt(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).modSqrt(x(), y())).(val).Integer()
}
func (i Integer) Mul(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).mul(x(), y())).(val).Integer()
}
func (i Integer) MulRange(a, b int64) Integer {
	return wrap(val(i).mulRange(a, b)).(val).Integer()
}
func (i Integer) Neg(x Integer) Integer {
	return wrap(val(i).neg(x())).(val).Integer()
}
func (i Integer) ProbablyPrime(n int) bool {
	return val(i).probablyPrime(n)
}
func (i Integer) Quo(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).quo(x(), y())).(val).Integer()
}
func (i Integer) QuoRem(x, y, r Integer) (Integer, Integer) {
	defer discardInt(x(), y(), r())
	a, b := val(i).quoRem(x(), y(), r())
	return wrap(a).(val).Integer(), wrap(b).(val).Integer()
}
func (i Integer) Rand(rnd *rand.Rand, x Integer) Integer {
	defer discardInt(x())
	return wrap(val(i).rand(rnd, x())).(val).Integer()
}
func (i Integer) Rem(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).rem(x(), y())).(val).Integer()
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
func (i Integer) Sub(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).sub(x(), y())).(val).Integer()
}
func (i Integer) Uint64() uint64 { return val(i).uint64() }

/////////////////////////////////////////////////////////////////////////
// BYTES
type Bytes val

func (b Bytes) Eval() Evaluable   { return b }
func (b Bytes) Serialize() []byte { return []byte(val(b)().String()) }
func (b Bytes) String() string    { return b().Text(8) }
func (b Bytes) Type() ValueType   { return BYTES }

/////////////////////////////////////////////////////////////////////////
// STRING
type Text val

func (s Text) Eval() Evaluable   { return s }
func (s Text) Serialize() []byte { return []byte(s().Bytes()) }
func (s Text) String() string    { return string(s.Serialize()) }
func (s Text) Type() ValueType   { return TEXT }

/////////////////////////////////////////////////////////////////////////
// FLOAT
type Float ratio

func (f Float) Eval() Evaluable   { return f }
func (f Float) Serialize() []byte { return []byte(f.String()) }
func (f Float) String() string    { return f().FloatString(10) }
func (f Float) Type() ValueType   { return FLOAT }

/////////////////////////////////////////////////////////////////////////
// RATIONAL
type Ratio ratio

func (r Ratio) Eval() Evaluable   { return r }
func (r Ratio) Serialize() []byte { return []byte(ratio(r).String()) }
func (r Ratio) String() string    { return r().String() }
func (r Ratio) Type() ValueType   { return RATIONAL }
