package agiledoc

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	"sync"
)

type (
	intPool sync.Pool
	ratPool sync.Pool
)

func NewVal() Val { return intGen.New().(Val) }
func NewRat() Rat { return ratGen.New().(Rat) }

var (
	intGen intPool = intPool{}
	ratGen ratPool = ratPool{}
)

func init() {
	intGen.New = func() interface{} { return big.NewInt(0) }
	ratGen.New = func() interface{} { return big.NewRat(1, 1) }
}

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
	Val func() *big.Int

	// paired types
	Rat  func() *big.Rat
	Pair func() [2]Evaluable

	// collection types see collection,go
)

func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Evaluable { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }

/////////////////////////////////////////////////////////////////////////
func (b Val) Eval() Evaluable   { return Value(b) }
func (b Val) Serialize() []byte { return []byte(b().String()) }
func (b Val) String() string    { return b().String() }
func (b Val) Type() ValueType   { return INT }

////////////////////////////////////////////////////////////////
func (b Val) BigInt() *big.Int { return b() }
func (b Val) Integer() Integer { return Integer(b) }
func (b Val) Bool() bool {
	if b().Int64() > 0 {
		return true
	} else {
		return false
	}
}
func (b Val) IntUntyped() int    { return int(b().Int64()) }
func (b Val) Int() int64         { return b().Int64() }
func (b Val) Unsigned() Unsigned { return Unsigned(b) }
func (b Val) Uint() uint64       { return b().Uint64() }
func (b Val) BigRat() *big.Rat   { return new(big.Rat).SetFrac(Value(ONE).(Val)(), b()) }
func (b Val) Flt() float64       { f, _ := b.BigRat().Float64(); return f }
func (b Val) Rat() Rat           { return Value(b.BigRat()).(Rat) }
func (b Val) Pair() [2]Evaluable { return [2]Evaluable{Value(), b} } // negative == index not set
func (b Val) Bytes() []byte      { return b().Bytes() }

/////////////////////////////////////////////////////////////////////////
func (b Pair) Eval() Evaluable { return Value(b) }

func (b Pair) Key() Evaluable   { return b()[0].Eval() }
func (b Pair) Value() Evaluable { return b()[1].Eval() }
func (b Pair) Index() int {
	if i, ok := b.Key().(Val); ok {
		return i.IntUntyped()
	} else {
		return -1 // negative → not set
	}
}
func (p Pair) Serialize() []byte {
	return append(
		p()[0].Serialize(),
		append(
			[]byte(": "),
			append(
				p()[1].Serialize(),
				[]byte("\n")...,
			)...,
		)...,
	)
}

func (b Pair) String() string  { return string(b.Serialize()) }
func (b Pair) Type() ValueType { return RAT }

/////////////////////////////////////////////////////////////////////////
func (r Rat) Eval() Evaluable { return Value(r) }

// Bytes is supposed to keep as much information as possible, so this converts
// numerator and denominator to 64 bytes each, ignoring the original accuracy
// (length), to make them divideable again. Accuracys greater 64bit should not
// be serialized, but kept in absolute numbers in memoru during calculations,
func (r Rat) Bytes() []byte {
	return append(
		new(big.Int).SetInt64(r().Num().Int64()).Bytes(),
		new(big.Int).SetInt64(r().Denom().Int64()).Bytes()...,
	)
}
func (r Rat) Serialize() []byte { return []byte(r().String()) }
func (r Rat) String() string    { return r().String() }
func (r Rat) Type() ValueType   { return RAT }
func (r Rat) Num() Evaluable    { return Value(r().Num()) }
func (r Rat) Denom() Evaluable  { return Value(r().Denom()) }

func (r Rat) BigInt() *big.Int { return Value(r).(Val)() }
func (r Rat) Float() Float     { return Float(r) }
func (r Rat) Flt() float64     { f, _ := r().Float64(); return f }

/////////////////////////////////////////////////////////////////////////

// INTEGER
type Integer Val

func (i Integer) Eval() Evaluable   { return Val(i).Eval() }
func (i Integer) Serialize() []byte { return []byte(Val(i)().String()) }
func (i Integer) String() string    { return i().Text(10) }
func (i Integer) Type() ValueType   { return INTEGER }

// BYTES
type Bytes Val

func (b Bytes) Eval() Evaluable   { return Val(b).Eval() }
func (b Bytes) Serialize() []byte { return []byte(Val(b)().String()) }
func (b Bytes) String() string    { return b().Text(8) }
func (b Bytes) Type() ValueType   { return BYTES }

// STRING
type String Val

func (s String) Eval() Evaluable   { return Val(s).Eval() }
func (s String) Serialize() []byte { return []byte(Val(s)().Bytes()) }
func (s String) String() string    { return string(s.Serialize()) }
func (s String) Type() ValueType   { return STRING }

// UNSIGNED INTEGER
type Unsigned Val

func (u Unsigned) Eval() Evaluable   { return Val(u).Eval() }
func (u Unsigned) Serialize() []byte { return Val(u)().Bytes() }
func (u Unsigned) String() string    { return u().Text(2) }
func (u Unsigned) Type() ValueType   { return UINT }

// FLOAT
type Float Rat

func (f Float) Eval() Evaluable   { return Rat(f).Eval() }
func (f Float) Serialize() []byte { return []byte(f.String()) }
func (f Float) String() string    { return Value(f()).(Rat).String() }
func (f Float) Type() ValueType   { return FLOAT }

// PAIREDIONAL
type Ratio Rat

func (r Ratio) Eval() Evaluable   { return Rat(r).Eval() }
func (r Ratio) Serialize() []byte { return []byte(Rat(r).String()) }
func (r Ratio) String() string    { return r.String() }
func (r Ratio) Type() ValueType   { return RATIONAL }

// INSTANCIATE A NEW VALUE
//
// values are represented internaly by either a Big, Rat, or Col type instance,
// each of which implement the absVal interface. Implemented as functional
// types, that return a value of destince type, either *big.Int, *big.Rat, or
// an Instance of the Collection interface. A Method set defined on the
// function type implements the absVal interface, by manipulating the main
// return value.
func Value(i ...interface{}) Evaluable {

	var v Evaluable
	// if one Element only
	if len(i) == 1 {
		// if allready a value, return that immedeately
		if v, ok := i[0].(Evaluable); ok {
			return v
		}
		// otherwise convert native to value
		v = nativeToValue(i[0])
	}
	// if exactly two elements, assume a pair of key/value as element for a map
	if len(i) == 2 { // convert key and value to be shure they implement value
		v = newMap(Pair(func() [2]Evaluable { return [2]Evaluable{Value(i[0]), Value(i[1])} }))
	}
	if len(i) > 2 { // if more than two values are passed, we assume an indexed list of values. Should they turn out to be key Value Pairs, they will be converted to a list of maps, due to recursion.
		var vals = []Evaluable{}
		for _, v := range i {
			v := v
			vals = append(vals, Value(v))
		}
		v = newList(vals...)
	}
	return v
}
func nativeToValue(i interface{}) Evaluable {

	var retFn Evaluable

	switch i.(type) {
	case bool: // a boolean returns a flag with the first bit set
		v := new(big.Int)
		if i.(bool) {
			v = v.SetInt64(int64(1))
		} else {
			v = v.SetInt64(int64(0))
		}
		retFn = BitFlag(func() *big.Int { return v })
	case []bool: // slice of bools gets spooled to a bitflag
		v := new(big.Int)
		for n, val := range i.([]bool) {
			var u uint
			val := val
			n := n
			if val { // set true
				u = 1
			} else { // or false
				u = 0
			} // at appropriate place
			v = (*v).SetBit(v, n, u)
		}
		retFn = BitFlag(func() *big.Int { return v })
	case uint, uint16, uint64: // a uint is assumed to be a single byte
		v := (new(big.Int).SetUint64((uint64(i.(uint8)))))
		retFn = BitFlag(func() *big.Int { return v })

	case int, int8, int16, int32, int64: // integers are integer
		v := (new(big.Int).SetInt64(int64(i.(int))))
		var Fn Integer = Integer(func() *big.Int { return v })
		retFn = Fn

	case []byte: // a bytes slice gets assigned as bytes
		v := (new(big.Int).SetBytes(i.([]byte)))
		retFn = Bytes(func() *big.Int { return v })

	case uint8, uint32:
		v := new(big.Int).SetBytes(i.([]byte))
		retFn = String(func() *big.Int { return v })

	case string: // a string gets assigned by its bislice as well
		v := (new(big.Int).SetBytes([]byte(i.(string))))
		retFn = String(func() *big.Int { return v })

	case float32, float64: // floating point values get assigned to rationals
		v := (new(big.Rat).SetFloat64(i.(float64)))
		retFn = Float(func() *big.Rat { return v })
	}
	return retFn
}
