/*
TYPE SYSTEM

To provide dunamic functional types, that auto-convert in sensible ways,
internal types are declared as closures returning an underlying native type
that decodes the raw value.

For byte slices, strings, booleans, bit-flags, integers, and unsigned integers,
the underlying base type is math/big Int. That type comes with all the
nescesssary methods to act as all those types, as well as store and manipulate
the values they encode in a higly effective manner. Ratios and Floats are
represented by math/big.Rat in arbitrary praesition. A two field struct is the
base for all tuples, key/value, value/index combinations and whatever else
comes in pairs, since you should take two of each kind and build a boat, or
something… All kinds of collections are implemented in container types taken
from github.com/emirpasic/gods

The closure based types are public and implement the evaluable interface. They
are derrived from private base types, that generalize all methods, based on the
underlying native type. Most public types method calls assert their receiver to
its base type, call those private methods and pass the return value after
asserting it to be the appropriate.type for a given context.

The only reason to convert between types (and have them in the first place), is
to treat them in different ways, depending on the semantic role they take in a
given context.  Numbers expect arithmetic operations, while strings and byte
slices get sorted, concatenated, displayed  by the operations defined up on
them. Most of the times, the value just needs to be associated with the
appropriate method set, by asserting it. No conversions needed!.

Some types can only be expressed by two distinct values, like ratios and
key/value pairs, for instance. Further on there are collections, like lists,
maps, tables, matrices and so on. Internal Methods are defined to convert
between the backend base types and types that hold one, or two big.Int values,
as well as for all kinds of collections.

All conversion needed between conceptually different types, can allmost allways
be assumed by the context, operations got called from. Operations can be
generalized for all types and determine what to do, based on the types of
arguments they got passed. Numbers for instance assume an arithmetic operation,
when the 'Add' method is called, while collections expect to be extended by a
field. The 'String' operation expects a representation of the contained value,
which is identical to the representation of that value in the source code it
came from but encoded as sting, instead of a byte slice and so on. No need for
human involvment typewise, 'the right thing™' will happen, no matter of what
rhe intention was in the first place, when calling the method. The context can
only contain instances of values, as well as provide and require operations,
that 'make sense', regarding itself.

The underlying complexity,is concealed by the public types and interfaces. All
generated and returned values, are passed as instances of the public types. All
Parameters passed are expected to implement one of the interfaces. The closure
implementing a public type can not be called directly, since interfaces are not
callable. The native type, it would return, never needs to be dealt with
directly since all needed methods are provided by the interface, or public type
of the particular instace. The overhead is kept small, since type conversion in
most cases just instanciates a new closure over the reference to the native
instance of the value. The new closure is associated to the appropriate method
set to deal with it, by its type. That happens (hopefully) entrirely on the
stack profiting from locality of reference (regarding the passed pointer) and a
highly optimized handling of the associated heap values,provided by the
math/big library (depending on escape analysis and whatnot, possibly also
located on the stack, which would propably be most effective). All operations
that generate, or manipulate instances, use the highly effective methods
provided by the math/big library and dispatch which ones to use by asserting
the private base type to do so. All Values are immutable as seen by the user of
the library, while internaly mutation is used, whenever it makes sense
performance wise. All methods that internaly manipulte, return a new instance
of themselves, by rewrapping the contained value in a new closure.  Instances
of native values are reused in sync pools, to keep allocation pressure low.
TODO: performance testing! If hypothesis turns out to be right: find ways to
lazyly pull values on stack whenever it makes sense.
*/
package types

import (
	"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	"math/rand"
	"sync"
)

//// EVALUABLE INSTANCIATION FROM ARBITRARY VALUES ////
/////////////// IN THREE SIMPLE STEPS /////////////////
////
/// to keep allocation pressure flat, cache instances of underlying native base
//  values in sync pools for instance recycling.
var (
	intPool   = sync.Pool{}
	ratPool   = sync.Pool{}
	pairPool  = sync.Pool{}
	floatPool = sync.Pool{}
)

// initializes pools with appropriate new function to return an instance of
// the type, that gets returned by this functional type
func init() {
	intPool.New = func() interface{} { return big.NewInt(0) }
	ratPool.New = func() interface{} { return big.NewRat(1, 1) }
	floatPool.New = func() interface{} { return big.NewFloat(0) }
	pairPool.New = func() interface{} { return [2]Evaluable{} }
}

///// VALUE RECYCLING /////
///
// puts evaluables enclosed return value back in the appropriate pool for later
// reuse
func discard(v Evaluable) {
	switch { //…discard each in appropriate pool
	case v.Type()&NATURAL != 0:
		discardInt(v.(val)())
	case v.Type()&REAL != 0:
		discardRat(v.(Ratio)())
	case v.Type()&FLOAT != 0:
		discardRat(v.(Float)())
	case v.Type()&TUPLE != 0:
		discardPair(v.(Pair)())
	}
}

// typed discard functions
func discardInt(v ...*big.Int) {
	for n := 0; n < len(v); n++ {
		n := n
		intPool.Put(v[n])
	}
}
func discardRat(v ...*big.Rat) {
	for n := 0; n < len(v); n++ {
		n := n
		ratPool.Put(v[n])
	}
}
func discardFloat(v ...*big.Float) {
	for n := 0; n < len(v); n++ {
		n := n
		floatPool.Put(v[n])
	}
}
func discardPair(v ...[2]Evaluable) {
	for n := 0; n < len(v); n++ {
		n := n
		pairPool.Put(v[n])
	}
}

func wrap(i interface{}) (r Evaluable) {
	switch i.(type) {
	case *big.Int:
		r = val(func() *big.Int { return i.(*big.Int) })
	case *big.Rat:
		r = Ratio(func() *big.Rat { return i.(*big.Rat) })
	case *big.Float:
		r = Float(func() *big.Rat { return i.(*big.Rat) })
	case Pair:
		// set key to zero and value to passed interface
		r = wrap([2]Evaluable{i.(Pair).Key(), i.(Pair).Value()})
	}
	return r
}

/////////////////////////////////////////////////////////////////////////
/////// VAL /////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
/// wrapper methods for each native big.Int method
func (v val) abs(x *big.Int) *big.Int                       { return v().Abs(x) }
func (v val) add(x, y *big.Int) *big.Int                    { return v().Add(x, y) }
func (v val) and(x, y *big.Int) *big.Int                    { return v().And(x, y) }
func (v val) andNot(x, y *big.Int) *big.Int                 { return v().And(x, y) }
func (v val) append(buf []byte, base int) []byte            { return v().Append(buf, base) }
func (v val) binomial(n, k int64) *big.Int                  { return v().Binomial(n, k) }
func (v val) bit(i int) uint                                { return v().Bit(i) }
func (v val) bitLen() int                                   { return v().BitLen() }
func (v val) bits() []big.Word                              { return v().Bits() }
func (v val) bytes() []byte                                 { return v().Bytes() }
func (v val) cmp(y *big.Int) (r int)                        { return v().Cmp(y) }
func (v val) div(x, y *big.Int) *big.Int                    { return v().Div(x, y) }
func (v val) divMod(x, y, m *big.Int) (*big.Int, *big.Int)  { return v().DivMod(x, y, m) }
func (v val) exp(x, y, m *big.Int) *big.Int                 { return v().Exp(x, y, m) }
func (v val) format(s fmt.State, ch rune)                   { v().Format(s, ch) }
func (v val) gCD(x, y, a, b *big.Int) *big.Int              { return v().GCD(x, y, a, b) }
func (v val) gobDecode(buf []byte) error                    { return v().GobDecode(buf) }
func (v val) gobEncode() ([]byte, error)                    { return v().GobEncode() }
func (v val) int64() int64                                  { return v().Int64() }
func (v val) uint64() uint64                                { return v().Uint64() }
func (v val) bitFlag() BitFlag                              { return BitFlag(wrap(v()).(val)) }
func (v val) lsh(x *big.Int, n uint) *big.Int               { return v().Lsh(x, n) }
func (v val) marshalJSON() ([]byte, error)                  { return v().MarshalJSON() }
func (v val) marshalText() (text []byte, err error)         { return v().MarshalText() }
func (v val) mod(x, y *big.Int) *big.Int                    { return v().Mod(x, y) }
func (v val) modInverse(g, n *big.Int) *big.Int             { return v().ModInverse(g, n) }
func (v val) modSqrt(x, p *big.Int) *big.Int                { return v().ModSqrt(x, p) }
func (v val) mul(x, y *big.Int) *big.Int                    { return v().Mul(x, y) }
func (v val) mulRange(a, b int64) *big.Int                  { return v().MulRange(a, b) }
func (v val) neg(x *big.Int) *big.Int                       { return v().Neg(x) }
func (v val) not(x *big.Int) *big.Int                       { return v().Not(x) }
func (v val) or(x, y *big.Int) *big.Int                     { return v().Or(x, y) }
func (v val) probablyPrime(n int) bool                      { return v().ProbablyPrime(n) }
func (v val) quo(x, y *big.Int) *big.Int                    { return v().Quo(x, y) }
func (v val) quoRem(x, y, r *big.Int) (*big.Int, *big.Int)  { return v().QuoRem(x, y, r) }
func (v val) rand(rnd *rand.Rand, n *big.Int) *big.Int      { return v().Rand(rnd, n) }
func (v val) rem(x, y *big.Int) *big.Int                    { return v().Rem(x, y) }
func (v val) rsh(x *big.Int, n uint) *big.Int               { return v().Rsh(x, n) }
func (v val) scan(s fmt.ScanState, ch rune) error           { return v().Scan(s, ch) }
func (v val) set(x *big.Int) *big.Int                       { return v().Set(x) }
func (v val) setBit(x *big.Int, i int, b uint) *big.Int     { return v().SetBit(x, i, b) }
func (v val) setBits(abs []big.Word) *big.Int               { return v().SetBits(abs) }
func (v val) setBytes(buf []byte) *big.Int                  { return v().SetBytes(buf) }
func (v val) setInt64(x int64) *big.Int                     { return v().SetInt64(x) }
func (v val) setString(s string, base int) (*big.Int, bool) { return v().SetString(s, base) }
func (v val) setUint64(x uint64) *big.Int                   { return v().SetUint64(x) }
func (v val) sign() int                                     { return v().Sign() }
func (v val) string() string                                { return v().String() }
func (v val) sub(x, y *big.Int) *big.Int                    { return v().Sub(x, y) }
func (v val) text(base int) string                          { return v().Text(base) }
func (v val) unmarshalJSON(text []byte) error               { return v().UnmarshalJSON(text) }
func (v val) unmarshalText(text []byte) error               { return v().UnmarshalText(text) }
func (v val) xor(x, y *big.Int) *big.Int                    { return v().Xor(x, y) }

///////////////////////////////////////////////////
// VAL METHODS TO CONVERT TO DIFFERENT EVALUABLE //
///////////////////////////////////////////////////

func (v val) Integer() Integer { return Integer(v) }
func (v val) Bool() Bool       { return Bool(v) }
func (v val) Bytes() Bytes     { return Bytes(v) }
func (v val) Text() Text       { return Text(v) }

/////////////////////////////////////////////////
////// VAL METHODS TO IMPLEMENT EVALUABLE ///////
/////////////////////////////////////////////////

func (b val) bool() bool {
	if b().Int64() > 0 {
		return true
	} else {
		return false
	}
}
func (b val) Type() ValueType   { return EMPTY }
func (b val) Eval() Evaluable   { return Value(b) }
func (b val) Serialize() []byte { return []byte(b().String()) }
func (b val) String() string    { return b().String() }

func (b val) bigInt() *big.Int { return b() }
func (b val) Int() int         { return int(b().Int64()) }
func (b val) Int64() int64     { return b().Int64() }
func (b val) Uint64() uint64   { return b().Uint64() }

// methods to convert to other type that share common base value
func (b val) toBitFlag() BitFlag { return BitFlag(b) }
func (b val) toFlag() Bool       { return Bool(b) }
func (b val) toInteger() Integer { return Integer(b) }
func (b val) toText() Text       { return Text(b) }

// assign receiver value as returnvalue, set key to zero
func (b val) toPair() Pair {
	var r Pair
	return r
}

// set receivers value as numerator and denomintor to one (don't devide by zero).
func (b val) toRational() (r ratio) {
	return r
}

/////////////////////////////////////////////////////////////////////////
//////// RATIONAL ///////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
// private wrapper methods for big Rat methods
func (r ratio) abs(x *big.Rat) *big.Rat               { return r().Abs(x) }
func (r ratio) add(x, y *big.Rat) *big.Rat            { return r().Add(x, y) }
func (r ratio) cmp(y *big.Rat) int                    { return r().Cmp(y) }
func (r ratio) denom() *big.Int                       { return r().Denom() }
func (r ratio) float32() (f float32, exact bool)      { return r().Float32() }
func (r ratio) float64() (f float64, exact bool)      { return r().Float64() }
func (r ratio) floatString(prec int) string           { return r().FloatString(prec) }
func (r ratio) gobDecode(buf []byte) error            { return r().GobDecode(buf) }
func (r ratio) gobEncode() ([]byte, error)            { return r().GobEncode() }
func (r ratio) inv(x *big.Rat) *big.Rat               { return r().Inv(x) }
func (r ratio) isInt() bool                           { return r().IsInt() }
func (r ratio) marshalText() (text []byte, err error) { return r().MarshalText() }
func (r ratio) mul(x, y *big.Rat) *big.Rat            { return r().Mul(x, y) }
func (r ratio) neg(x *big.Rat) *big.Rat               { return r().Neg(x) }
func (r ratio) num() *big.Int                         { return r().Num() }
func (r ratio) quo(x, y *big.Rat) *big.Rat            { return r().Quo(x, y) }
func (r ratio) ratString() string                     { return r().RatString() }
func (r ratio) scan(s fmt.ScanState, ch rune) error   { return r().Scan(s, ch) }
func (r ratio) set(x *big.Rat) *big.Rat               { return r().Set(x) }
func (r ratio) setFloat64(f float64) *big.Rat         { return r().SetFloat64(f) }
func (r ratio) setFrac(a, b *big.Int) *big.Rat        { return r().SetFrac(a, b) }
func (r ratio) setFrac64(a, b int64) *big.Rat         { return r().SetFrac64(a, b) }
func (r ratio) setInt(x *big.Int) *big.Rat            { return r().SetInt(x) }
func (r ratio) setInt64(x int64) *big.Rat             { return r().SetInt64(x) }
func (r ratio) setString(s string) (*big.Rat, bool)   { return r().SetString(s) }
func (r ratio) sign() int                             { return r().Sign() }
func (r ratio) string() string                        { return r().String() }
func (r ratio) sub(x, y *big.Rat) *big.Rat            { return r().Sub(x, y) }
func (r ratio) unmarshalText(text []byte) error       { return r().UnmarshalText(text) }

func (r ratio) Float() Float { return Float(r) }
func (r ratio) Ratio() Ratio { return Ratio(r) }

/////////////////////////////////////////////////
////// METHODS TO IMPLEMENT EVALUABLE ///////////
/////////////////////////////////////////////////
func (r ratio) Eval() Evaluable { return Value(r) }

// Bytes is supposed to keep as much information as possible, so this converts
// numerator and denominator to 64 bytes each, ignoring the original accuracy
// (length), to make them divideable again. Accuracys greater 64bit should not
// be serialized, but kept in absolute numbers in memoru during calculations,
func (r ratio) Bytes() []byte {
	return append(
		r().Num().Bytes(),
		r().Denom().Bytes()...,
	)
}
func (r ratio) Type() ValueType { return REAL }

///////////////////////////////////////////////////
////// METHODS TO CONVERT TO DIFFERENT EVALUABLE //
///////////////////////////////////////////////////
func (r ratio) Serialize() []byte { return []byte(r().String()) }
func (r ratio) String() string    { return r().String() }

////////////////////////////////////////////////////////////////
// private methods, to convert to native types
func (r ratio) bigRat() *big.Rat { return Value(r).(ratio)() }

// public methods to convert to other implementations of evaluable
func (r ratio) Rational() ratio { return r }
func (r ratio) Pair() Pair {
	return Pair(func() [2]Evaluable { return [2]Evaluable{r.Num(), r.Denom()} })
}

// methods that take or return the integer type, to set, or get contained values
func (r ratio) Num() Integer   { return Value(r().Num()).(Integer) }
func (r ratio) Denom() Integer { return Value(r().Denom()).(Integer) }
