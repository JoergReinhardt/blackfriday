package types

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	//"math/rand"
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
func (u Bool) Native() bool {
	if u().Int64() > int64(0) {
		return true
	} else {
		return false
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

// sets a Bools value to the value of a passed Bool
func (u Bool) SetBool(x Bool) Bool {
	// discard parameter and old version
	defer discardInt(x(), u())
	// pre allocate return value
	var res Bool
	// return either positive, or negative one, Based on truthyness of
	// the passed parameters, Where negativety, or abscense of a value, is
	// considered false, while all positive values will be considdered
	// true. by rewriting to positive, or negative one, Value will be
	// normalized.
	if u.And(u, x)().Int64() > 0 {
		res = wrap(intPool.Get().(val).setInt64(1)).(val).Bool()
	} else {
		res = wrap(intPool.Get().(val).setInt64(-1)).(val).Bool()
	}
	return res
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
		r = wrap(intPool.Get().(val).setInt64(1)).(val).Bool()
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
func (i Integer) Sub(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).sub(x(), y())).(val).Integer()
}
func (i Integer) Cmp(x Integer) int {
	defer discardInt(x())
	a := wrap(intPool.Get().(*big.Int).Set(i())).(val).Integer()
	return a().Cmp(x())
}
func (i Integer) Div(y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).div(x(), i())).(val).Integer()
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

//func (i Integer) ModInverse(x, y Integer) Integer {
//	defer discardInt(x(), y())
//	return wrap(val(i).modInverse(x(), y())).(val).Integer()
//}
//func (i Integer) ModSqrt(x, y Integer) Integer {
//	defer discardInt(x(), y())
//	return wrap(val(i).modSqrt(x(), y())).(val).Integer()
//}
func (i Integer) Mul(x, y Integer) Integer {
	defer discardInt(x(), y())
	return wrap(val(i).mul(x(), y())).(val).Integer()
}

//func (i Integer) MulRange(a, b int64) Integer {
//	return wrap(val(i).mulRange(a, b)).(val).Integer()
//}
func (i Integer) Neg(x Integer) Integer {
	return wrap(val(i).neg(x())).(val).Integer()
}

//func (i Integer) ProbablyPrime(n int) bool {
//	return val(i).probablyPrime(n)
//}
//func (i Integer) Quo(x, y Integer) Integer {
//	defer discardInt(x(), y())
//	return wrap(val(i).quo(x(), y())).(val).Integer()
//}
//func (i Integer) QuoRem(x, y, r Integer) (Integer, Integer) {
//	defer discardInt(x(), y(), r())
//	a, b := val(i).quoRem(x(), y(), r())
//	return wrap(a).(val).Integer(), wrap(b).(val).Integer()
//}
//func (i Integer) Rand(rnd *rand.Rand, x Integer) Integer {
//	defer discardInt(x())
//	return wrap(val(i).rand(rnd, x())).(val).Integer()
//}
//func (i Integer) Rem(x, y Integer) Integer {
//	defer discardInt(x(), y())
//	return wrap(val(i).rem(x(), y())).(val).Integer()
//}
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

/////////////////////////////////////////////////////////////////////////
// BYTES
type Bytes val

func (b Bytes) Eval() Evaluable { return b }

// the string representation is provided by serializing the integer to a slice
// of bytes and converting that to a string, that way preserving all contained
// information. Lower order types are supposed to be stored in a more
// appropriate internal type, like Bool, or Integer, and otherwise need to be
// re-parsed to regain arithmetic, or boolean functionality.
func (b Bytes) String() string { return string(b().Bytes()) }

// if serialized, the string representation is converted to a slice of bytes,
// in order to not use any valid information. In case a 'lower' type was stored
// by the Bytes instance, it must br reparsed at a later point to convert it to
// the appropriate internal type
func (b Bytes) Serialize() []byte { return []byte(b.String()) }
func (b Bytes) Type() ValueType   { return BYTES }

func (b Bytes) Bit(n int) uint {
	return b().Bit(n)
}
func (b Bytes) BitLen() int {
	return b().BitLen()
}
func (b Bytes) Bytes() Bytes {
	//the returned byte slice neds to be represented by a big Int, which is
	//provided by the modules public Value funcitons slice
	return Value(b.Serialize()).(val).Bytes()
}
func (b Bytes) SetBytes(x Bytes) Bytes {
	// parameter instance will be reused
	defer discardInt(x())
	// returns a big Int, which only needs to be enclosed in a fresh
	// closure, provided by the modules wrap function.
	return wrap(b().SetBytes(x.Serialize())).(val).Bytes()
}
func (b Bytes) SetBytesNative(x []byte) Bytes {
	// the wrapper encloses the returned big Int in a fresh closure for
	// return.
	return wrap(b().SetBytes(x)).(val).Bytes()
}
func (b Bytes) AppendBytes(x Bytes) Bytes {
	defer discardInt(x(), b())
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(append(b.Serialize(), x.Serialize()...)).(val).Bytes()
}
func (b Bytes) AppendBytesNative(x []byte) Bytes {
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(append(b.Serialize(), x...)).(val).Bytes()
}

// set text is allmost identical to set bytes, since text is stored in the same
// way as a byte slice and features a Serialize method just like it, since its
// an implementation of an Evaluable just like the Bytes type.
func (b Bytes) SetText(x Text) Bytes {
	// parameter instance will be reused
	defer discardInt(x(), b())
	// the wrapper encloses the big Int, representing the string
	// representation of the passed Text instance in a fresh closure for
	// return.
	return wrap(b().SetBytes([]byte(x.String()))).(val).Bytes()
}
func (b Bytes) SetTextNative(x string) Bytes {
	// after returning the new instance, the old one is designated for
	// reuse.
	defer discardInt(b())
	// b is set to a native string by replacing it vit a new value instance
	return Value(x).(val).Bytes()
}
func (b Bytes) AppendText(x Text) Bytes {
	defer discardInt(x())
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(b.String() + x.String()).(val).Bytes()
}
func (b Bytes) AppendTextNative(x string) Bytes {
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(b.String() + x).(val).Bytes()
}

/////////////////////////////////////////////////////////////////////////
// STRING
type Text val

func (s Text) Eval() Evaluable { return s }

// text stored in the enclosed int, is retrieved, by serializing to a  Byte
// slice representation.
func (s Text) Serialize() []byte { return s().Bytes() }

// the string method builds a string representation og the contained data, by
// serializing it to bytes and representing those as a string
func (s Text) String() string  { return string(s.Serialize()) }
func (s Text) Type() ValueType { return TEXT }

// set a pre-existing Text Instance to a Value represented by the internal
// Bytes type.
func (s Text) SetBytes(x Bytes) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// SetBytes takes a public Bytes instance and sets an existing text to
	// its serialization to a byte slice
	return Value(s().SetBytes(x.Serialize())).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the native byte
// slice.
func (s Text) SetBytesNative(x []byte) Text {
	// setBytesNatice sets s to a native go byte slice. converted to Text
	// via Value
	return Value(s().SetBytes(x)).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the internal
// Text type.
func (s Text) SetText(x Text) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// setBytes with the string returned by the value converted to bytes as
	// Parameter, finaly cinverted to Text via Value
	return Value(s().SetBytes([]byte(x.String()))).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the native string.
func (s Text) SetTextNative(x string) Text {
	// setBytes with the string returned by the value converted to bytes as
	// Parameter, finaly cinverted to Text via Value
	return Value(s().SetBytes([]byte(x))).(val).Text()
}

// Append an Instance of the internal Bytes Type to a preexisting Text Instance
func (s Text) AppendBytes(x Bytes) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// uses internal append funcrion and Serialize, which must be provided
	// by all evaluable, to concatenate on a byte base
	return Value(append(s.Serialize(), x.Serialize()...)).(val).Text()
}

// Append an Instance of a native byte slice to a preexisting Text Instance
func (s Text) AppendBytesNative(x []byte) Text {
	// uses internal append funcrion and Serialize to append a native byte
	// Slice with a given text by re-valuabling using Value, asserting the
	// intermediate val type and calling the Text() method on it.
	return Value(append(s.Serialize(), x...)).(val).Text()
}

// Append an Instance of the internal Text Type to a preexisting Text Instance
func (s Text) AppendText(x Text) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// uses string concatenation to append a text provided as parameter to
	// a given Text instance
	return Value(s.String() + x.String()).(val).Text()
}

// Append an Instance of a native string to a preexisting Text Instance
func (s Text) AppendTextNative(x string) Text {
	// uses gos append function and iinternal String method provided by all
	// evaluables, to concatenate annative string  to the given Text using
	// string concatenation.
	return Value(s.String() + x).(val).Text()
}

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
