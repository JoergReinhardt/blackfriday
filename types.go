// TYPE SYSTEM
//
// after lots of experimenting, I decided the best way is to keep the interface
// as simple as possible. One Method interfaces for the win but you might want
// to add one just to get shure. The main reason to reimplement a typesystem on
// top of gos existing type system, is due to the fact, that the reflection
// capabilitys are far to complex in features. far to complicated to use and
// don't perform that well either.
//
// That makes a type marker that performs well in comparisions, the essence of
// a type. The Value needs to be returnable in a generalized form, which is
// perfomed by the Value() mehode. Every value can additionaly be returned as
// its contained native type, but those methods differ by signature and can't
// be defined without the use of an empty interface. Everything else is up for
// grabs,
package agiledoc

import (
	"errors"
	"fmt"
	//"strconv"
	//"github.com/emirpasic/gods/containers"
	"math/big"
	//"strconv"
)

// VALUE INTERFACE
//
// to implement the interface a type must provide a type method to return a
// uint based Type Flag of type ValueType.
//
// A Value function thet returns the contained value in a form that implements
// the Val interface has to be provided.
//
// Finaly a method needs to be provided, that returns the contained type
// converted to a given internal return type and throws an error in cases where
// that's impossible.
type Val interface {
	Type() ValType
	Value() Val
}

// KEY-VAL INTERFACE
// values that are part of a finite set, mapped to, or identifyed by a key,
// need a type to store them including the key
type KVal interface {
	Val
	Ident() string
}

// type to return as error, when type casting failed
type ConversionError error

// generate a new type conversion error, providing information about what
// actually went wrong
func NewConvError(t ValType, v Val) ConversionError {
	// pregenerate formated string
	str := "Type Conversion Error:\n"
	str = str + "tryed to convert Value\n"
	str = str + "from Type %s to Type %s\n"
	// instanciate error at runtime, substituting values by calling sprintf
	err := errors.New(fmt.Sprintf(str, v.Type(), t))
	// return as conversion error
	return ConversionError(err)
}

// internal value type identifyer is a uint bitflag
type ValType uint

//go:generate -command stringer -type ValType
const (
	NIL     ValType = 0
	BOOLEAN ValType = 1 << iota
	FLAG
	INTEGER
	RATIONAL
	BYTES
	STRING
	KEYVAL    // type to hold key value pairs like variables and parameters
	CONTAINER // holds lists, sets, stacks, maps and trees
)

// FUNCTIONAL TYPES
//
// all internal types are defined as function types that return a natural type
// representing the contained value(s). Function types are immutable and get
// constructed dynamicaly by a constructor. All Methods except the type and val
// methods are defined on the function instance. Type and Val Method of each
// functional type are defined as methods on the function type.
type (
	empty    func() empty // the empty value is represented by the empty struct
	boolVal  func() bool     // single boolen. Sets will be expressed as intVal
	intVal   func() *big.Int // all integers are represented by big ints as well
	ratVal   func() *big.Rat // floats will be represented as big rational
	bytesVal func() []byte   // most input streams will be byte buffers
	strVal   func() string   // since we're dealing with text, strings need to be represented
	keyVal   func() struct { // variable, map member and parameter type
		id  string // provides an string identifyer
		Val        // wraps a typable value implementation
	}
	cntVal func() struct { // implements internal container, iterator
		// and comparator interface by wrapping gods containers, iterators and
		// comparators in its fields and methods.
		CntType
		Vals []Val
	}
	flagVal intVal // a bitflag stored in a big int
)

// COMBINED TYPES METHODS
// methods for combined tupes keyval and container
func (k keyVal) Identity() string       { return k().id }
func (k keyVal) Get() Val               { return k().Val }
func (c cntVal) ContainerType() CntType { return c().CntType }
func (c cntVal) Get() []Val             { return c().Vals }

// INTERFACE TYPE METHODS
func (empty) Type() ValType    { return NIL }
func (boolVal) Type() ValType  { return BOOLEAN }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (ratVal) Type() ValType   { return RATIONAL }
func (bytesVal) Type() ValType { return BYTES }
func (strVal) Type() ValType   { return STRING }
func (keyVal) Type() ValType   { return KEYVAL }
func (cntVal) Type() ValType   { return CONTAINER }

// INTERFACE VALUE METHOD
// the value method of a functional type calls its receiver to instanciate (a)
// copy(s) of the returned value(s) and and wraps a closure literal around that
// copy to get a new functional instance to return to the caller.
func (v empty) Value() Val    { return func() Val { return v } }
func (v boolVal) Value() Val  { return func() Val { return v } }
func (v intVal) Value() Val   { return func() Val { return v } }
func (v ratVal) Value() Val   { return func() Val { return v } }
func (v bytesVal) Value() Val { return func() Val { return v } }
func (v strVal) Value() Val   { return func() Val { return v } }
func (v keyVal) Value() Val {
	return func() Val {
		return struct {
			string
			Val
		}{v.Identity(), v.Get()}
	}
}
func (v cntVal) Value() Val {
	v = v()
	return func() Val {
		return struct {
			VecType
			Vals []Val
		}{v.ContainerType(), v.Get()}
	}
}
func (f flagVal) Value() Val { return func() Val { return *big.Int(v()) } }

// GENERALIZED CONVERSION METHODS
//
// the conversion functions are defined as methods on the value function type
// and can therefore be called up on any Val implementing types value method
// implementation.
//
// the to function type takes instances of Val and converts to arbitrary types
type ToTypeFn func(Val) interface{}

// the from type takes arbitrary types and returns Val instances
type FromTypeFn func(Val interface{})

// the Conversion function returns all conversion functions a value method provides
type ConversionsFn func(Val) ([]FromTypeFn, []ToTypeFn)

// DEDICATED TYPE METHOD TYPES
// these types define functions that all implement the contained type, but
// return a different type instead of the empty interface the contained type
// defines. Which type to assert, is easy to recocnize given the function
// names.
type (
	EmptyFn     func() empty
	BoolValFn   func(Val) boolVal
	IntValFn    func(Val) intVal
	RatValFn    func(Val) ratVal
	BytesValFn  func(Val) bytesVal
	StrValFn    func(Val) strVal
	KeyValFn    func(Val) (string, Val)
	VecValFn    func(Val) (CntType, []Val)
	BoolFn      func(Val) bool
	IntFn       func(Val) int
	BigIntFn    func(Val) *big.Int
	RatFn       func(Val) (*big.Int, *big.Int)
	BigRatFn    func(Val) *big.Rat
	FloatFn     func(Val) float64
	ByteSliceFn func(Val) []byte
	StrFn       func(Val) string
	FlagFn      BigIntFn
	///////////////////////////
	ValEmpty     func() Val
	ValBoolVal   func(boolVal) Val
	ValIntVal    func(intVal) Val
	ValRatVal    func(ratVal) Val
	ValBytesVal  func(bytesVal) Val
	ValStrVal    func(strVal) Val
	ValKeyVal    func(string Val) Val
	ValVecVal    func(VecType ...Val) Val
	ValBool      func(bool) Val
	ValInt       func(int) Val
	ValBigInt    func(*big.Int) Val
	ValRat       func(counter *big.Int, denominator *big.Int) Val
	ValBigRat    func(*big.Rat) Val
	ValFloat     func(float64) Val
	ValByteSlice func([]byte) Val
	ValStr       func(string) Val
	ValFlag      ValBigInt
)

const ( // all type identyfiers as unsigned integer constant
	// internal types that implement the Val interface
	Empty  uint = 0
	IntVal      = iota << 1
	RatVal
	StrVal
	ByteVal
	KeyVal
	VecVal
	// native types returned by Val instances
	Bool
	Uint
	Int
	BigInt
	Float
	BigFloat
	BigRat
	ByteSlice
	String
	ValSlice

	FlagVal  = IntVal
	FloatVal = RatVal

	Native = Bool | Uint | Int | BigInt | Float |
		BigFloat | Bytes | String | ValSlice

	Internal = Empty | FlagVal | IntVal | FloatVal |
		RatVal | StrVal | ByteVal | KeyVal | VecVal
)

type TypeConversions struct {
	From []FromTypeFn
	To   []ToTypeFn
}

var Conversions []TypeConversions

func init() {
	Conversions[Empty] = TypeConversions{
		[]FromTypeFn{
		  func(),
		},
		  []ToTypeFn{
		  func() empty,
		},
	}
	Conversions[IntVal] = func(v Val) { return v.(IntVal) }
}
