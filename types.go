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
package agiledoc

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	"sync"
)

//// EVALUABLE INSTANCIATION FROM ARBITRARY VALUES ////
/////////////// IN THREE SIMPLE STEPS /////////////////
////
/// to keep allocation pressure flat, cache instances of underlying native base
//values in sync pools for instance recycling.
var (
	intGen  = sync.Pool{}
	ratGen  = sync.Pool{}
	pairGen = sync.Pool{}
)

// initializes pools with appropriate new function to return an instance of
// the type, that gets returned by this functional type
func init() {
	intGen.New = func() interface{} { return big.NewInt(0) }
	ratGen.New = func() interface{} { return big.NewRat(1, 1) }
	pairGen.New = func() interface{} { return [2]Evaluable{} }
}

// retrieve recycled empty instance of return type, for each of the functional
// base types
func newVal() val   { return intGen.New().(val) }
func newRat() ratio { return ratGen.New().(ratio) }
func newPair() pair { return pairGen.New().(pair) }

// finally, after possibly setting, or mutating the raw instance, returned by
// functional type, it needs to be wrapped in the appropriate type of closure
// again, to implement evaluable
func valWrap(v *big.Int) val               { return func() *big.Int { return v } }
func ratioWrap(v *big.Rat) func() *big.Rat { return func() *big.Rat { return v } }
func pairWrap(k, v Evaluable) pair         { return func() [2]Evaluable { return [2]Evaluable{k, v} } }

///// VALUE RECYCLING /////
///
// puts evaluables enclosed return value back in the appropriate pool for later
// reuse
func DiscardEvaluable(v Evaluable) {
	switch { //…discard each in appropriate pool
	case v.Type()&NATURAL != 0:
		discardInt(v.(val)())
	case v.Type()&REAL != 0:
		discardRat(v.(ratio)())
	case v.Type()&PAIR != 0:
		discardPair(v.(pair)())
	}
}

// typed discard functions
func discardInt(v *big.Int)      { intGen.Put(v) }
func discardRat(v *big.Rat)      { ratGen.Put(v) }
func discardPair(v [2]Evaluable) { pairGen.Put(v) }

// BASE VALUES IMPLEMENTING FUNCTIONAL TYPES
// these functional types need to implement the absVal interface to be suitable
// base types. If called, they return their contained value. A Method set
// defined on these funtional types, implements the absVal interface, by
// manipulating the returned content. Each can implement ia couple of dynamic
// types by defining further types based on it, while overwriting and/or
// completing the method set.
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

/////////////////////////////////////////////////////////////////////////
/////// VAL /////////////////////////////////////////////////////////////
func (b val) bool() bool {
	if b().Int64() > 0 {
		return true
	} else {
		return false
	}
}
func (b val) set(i int64) val   { return valWrap(b().SetInt64(i)) }
func (b val) Type() ValueType   { return EMPTY }
func (b val) Eval() Evaluable   { return Value(b) }
func (b val) Serialize() []byte { return []byte(b().String()) }
func (b val) String() string    { return b().String() }

func (b val) bigInt() *big.Int { return b() }
func (b val) BigRat() *big.Rat { ; return newRat()().SetFloat64(float64(b().Int64())) }
func (b val) Int() int         { return int(b().Int64()) }
func (b val) Int64() int64     { return b().Int64() }
func (b val) Uint64() uint64   { return b().Uint64() }
func (b val) Bytes() []byte    { return b().Bytes() }

// methods to convert to other type that share common base value
func (b val) toBitFlag() BitFlag { return BitFlag(b) }
func (b val) toFlag() Bool       { return Bool(b) }
func (b val) toInteger() Integer { return Integer(b) }
func (b val) toText() Text       { return Text(b) }

// assign receiver value as returnvalue, set key to zero
func (b val) toPair() pair {
	var r = newPair()
	r.SetValue(b)
	return r
}

// set receivers value as numerator and denomintor to one (don't devide by zero).
func (b val) toRational() ratio {
	var r = newRat()
	r.SetNum(b.toInteger())
	return r
}

/////////////////////////////////////////////////////////////////////////
//////// RATIONAL ///////////////////////////////////////////////////////
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
func (r ratio) Serialize() []byte { return []byte(r().String()) }
func (r ratio) String() string    { return r().String() }
func (r ratio) Type() ValueType   { return REAL }

////////////////////////////////////////////////////////////////
// private methods, to convert to native types
func (r ratio) float64() float64 { f, _ := r().Float64(); return f }
func (r ratio) bigRat() *big.Rat { return Value(r).(ratio)() }

// public methods to convert to other implementations of evaluable
func (r ratio) Float() Float    { return Float(r) }
func (r ratio) Rational() ratio { return r }
func (r ratio) Pair() pair {
	return pair(func() [2]Evaluable { return [2]Evaluable{r.Num(), r.Denom()} })
}

// methods that take or return the integer type, to set, or get contained values
func (r ratio) Num() Integer               { return Value(r().Num()).(Integer) }
func (r ratio) Denom() Integer             { return Value(r().Denom()).(Integer) }
func (r ratio) SetNum(v Integer) ratio     { r().SetFrac(v(), newVal()().SetInt64(1)); return r }
func (r ratio) SetDenom(v Integer) ratio   { r().SetFrac(newVal()().SetInt64(1), v()); return r }
func (r ratio) SetFrac(n, d Integer) ratio { r().SetFrac(n(), d()); return r }

/////////////////////////////////////////////////////////////////////////
/////// PAIR ////////////////////////////////////////////////////////////
func (b pair) Eval() Evaluable { return Value(b) }

func (b pair) Value() Evaluable { return b()[1].Eval() }

// a pair allways provides a key, which can be of any given base type
func (b pair) Key() Evaluable { return b()[0].Eval() }

// Index() int
// returns the key of the element as native integger, if it turns out to be
// convertable, otherwise return a negative integer to indicate that the key is
// not convertable to a Number
func (b pair) Index() Integer {
	var ret Integer
	if b.Key().Type()&SYMBOLIC != 0 {
		ret = Value(-1).(Integer) // negative → not set
	} else { // NUMERIC
		// if natural number, return as interger
		if b.Key().Type()&NATURAL != 0 {
			ret = b.Key().(Integer)
		}
		// if real number, return numerator as interger
		if b.Key().Type()&REAL != 0 {
			ret = b.Key().(ratio).Num()
		}
	}
	return ret
}
func (b pair) SetKey(v Evaluable) pair {
	return func() [2]Evaluable { return [2]Evaluable{v, b.Value()} }
}
func (b pair) SetValue(v Evaluable) pair {
	return func() [2]Evaluable { return [2]Evaluable{b.Key(), v} }
}
func (b pair) SetBoth(k Evaluable, v Evaluable) pair {
	return func() [2]Evaluable { return [2]Evaluable{k, v} }
}
func (p pair) Serialize() []byte {
	var delim = []byte{}
	if p.Index()().Int64() == -1 {
		delim = []byte(": ")
	} else {
		delim = []byte(".) ")
	}
	return append(
		p()[0].Serialize(),
		append(
			delim,
			append(
				p()[1].Serialize(),
				[]byte("\n")...,
			)...,
		)...,
	)
}
func (b pair) String() string  { return string(b.Serialize()) }
func (b pair) Type() ValueType { return TUPLE }

// generate pair from evaluables
func pairFromValues(k, v Evaluable) (r pair) { return r }

/////////////////////////////////////////////////////////////////////////
/////////// PUBLIC IMPLEMENTAIONS OF EVALUABLES /////////////////////////
///////////////// BASED ON VAL, RAT & PAIR //////////////////////////////
/////////////////////////////////////////////////////////////////////////

//////// EMPTY //////////////////////////////////////////////////////////
func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Evaluable { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }

/////// BOOL ////////////////////////////////////////////////////////////
// booleans allways come in slices (encoded as big.int, handled bitwise using
// uint representation)
type Bool val

func (u Bool) Eval() Evaluable   { return u }
func (u Bool) Serialize() []byte { return []byte(u().String()) }
func (u Bool) String() string    { return u().Text(2) }
func (u Bool) Type() ValueType   { return BOOL }

/////////////////////////////////////////////////////////////////////////
// INTEGER
type Integer val

func (i Integer) Eval() Evaluable   { return i }
func (i Integer) Serialize() []byte { return []byte(val(i)().String()) }
func (i Integer) String() string    { return i().Text(10) }
func (i Integer) Type() ValueType   { return INTEGER }
func (i Integer) Int64() int64      { return i().Int64() }

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

/////////////////////////////////////////////////////////////////////////////
// INSTANCIATE NEW VALUE(S) FROM GOLANG NATIVE VALUES
//
// 1.) chack number of passed values:
//	- one: pass on to convert from native type
//	- two: pass on to create a pair of values
//	- > two:  pass on to create a collection
func Value(i ...interface{}) (v Evaluable) {

	// IF SINGLE ELEMENT GOT PASSED
	//
	//// TEST IF ALLREADY EVALUABLE ////
	if len(i) == 1 { // value generation is indempotent and just ommitted,
		// if parameter is allready evaluable.
		if v, ok := i[0].(Evaluable); ok {
			// !!! EARLY BIRD RETURN SPECIAL !!!
			return v
		}

		// NATIVE INTENDED FOR CONVERSION TO EVALUABLE
		v = nativeToValue(i[0])
	}

	// IF TWO ELEMENTS GOT PASSED
	//
	// if exactly two elements, assume a pair of key/value as element for a map
	if len(i) == 2 { // convert key and value recursively to make shure
		// they implement evaluate
		v = pairFromValues(Value(i[0]), Value(i[1]))
	}

	// MORE THAN TWO ELEMENTS GOT PASSED
	//
	// if more than two values are passed, we assume an
	// slice of values to be converted to some kind of collection.
	if len(i) > 2 {
		v = Collect(valueSlice(i)...)
	}
	return v
}
