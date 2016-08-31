package agiledoc

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
)

// ABSOLUTE VALUE IMPLEMENTING FUNCTIONAL TYPES
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
	Int func() *big.Int
	// paired types
	Rat  func() *big.Rat
	Pair func() [2]Value
	// collection types see collection,go
)

func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Value     { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }

/////////////////////////////////////////////////////////////////////////
func (b Int) Eval() Value       { return Eval(b) }
func (b Int) Serialize() []byte { return []byte(b().String()) }
func (b Int) String() string    { return b().String() }
func (b Int) Type() ValueType   { return NUMERIC }

////////////////////////////////////////////////////////////////
func (b Int) BigInt() *big.Int { return b() }
func (b Int) Integer() Integer { return Integer(b) }
func (b Int) Bool() bool {
	if b().Int64() > 0 {
		return true
	} else {
		return false
	}
}
func (b Int) IntUntyped() int    { return int(b().Int64()) }
func (b Int) Int() int64         { return b().Int64() }
func (b Int) Unsigned() Unsigned { return Unsigned(b) }
func (b Int) Uint() uint64       { return b().Uint64() }
func (b Int) BigRat() *big.Rat   { return new(big.Rat).SetFrac(Eval(ONE).(Int)(), b()) }
func (b Int) Flt() float64       { f, _ := b.BigRat().Float64(); return f }
func (b Int) Rat() Rat           { return Eval(b.BigRat()).(Rat) }
func (b Int) Pair() [2]Value     { return [2]Value{Eval(), b} } // negative == index not set

/////////////////////////////////////////////////////////////////////////
func (b Pair) Eval() Value { return Eval(b) }

func (b Pair) Serialize() []byte {
	return append(
		b.Key().Serialize(),
		append(
			[]byte(": "),
			append(
				b.Value().Serialize(),
				[]byte("\n")...)...)...)
}
func (b Pair) Key() Value   { return b()[0].Eval() }
func (b Pair) Value() Value { return b()[1].Eval() }
func (b Pair) Index() int {
	if i, ok := b.Key().(Int); ok {
		return i.IntUntyped()
	} else {
		return -1 // negative → not set
	}
}

func (b Pair) String() string  { return string(b.Serialize()) }
func (b Pair) Type() ValueType { return PAIRED }

/////////////////////////////////////////////////////////////////////////
func (r Rat) Eval() Value { return Eval(r) }

func (r Rat) Serialize() []byte { return []byte(r().String()) }
func (r Rat) String() string    { return r().String() }
func (r Rat) Type() ValueType   { return PAIRED }
func (r Rat) Num() Value        { return Eval(r().Num()) }
func (r Rat) Denom() Value      { return Eval(r().Denom()) }

func (r Rat) BigInt() *big.Int { return Eval(r).(Int)() }
func (r Rat) Float() Float     { return Float(r) }
func (r Rat) Flt() float64     { f, _ := r().Float64(); return f }

/////////////////////////////////////////////////////////////////////////

// INTEGER
type Integer Int

func (i Integer) Eval() Value       { return Int(i).Eval() }
func (i Integer) Serialize() []byte { return []byte(Int(i)().String()) }
func (i Integer) String() string    { return i().Text(10) }
func (i Integer) Type() ValueType   { return INTEGER }

// BYTES
type Bytes Int

func (b Bytes) Eval() Value       { return Int(b).Eval() }
func (b Bytes) Serialize() []byte { return []byte(Int(b)().String()) }
func (b Bytes) String() string    { return b().Text(8) }
func (b Bytes) Type() ValueType   { return BYTES }

// STRING
type String Int

func (s String) Eval() Value       { return Int(s).Eval() }
func (s String) Serialize() []byte { return Int(s)().Bytes() }
func (s String) String() string    { return s().String() }
func (s String) Type() ValueType   { return BYTES }

// UNSIGNED INTEGER
type Unsigned Int

func (u Unsigned) Eval() Value       { return Int(u).Eval() }
func (u Unsigned) Serialize() []byte { return Int(u)().Bytes() }
func (u Unsigned) String() string    { return u().Text(2) }
func (u Unsigned) Type() ValueType   { return UINT }

// FLOAT
type Float Rat

func (f Float) Eval() Value       { return Rat(f).Eval() }
func (f Float) Serialize() []byte { return []byte(f.String()) }
func (f Float) String() string    { return Eval(f()).(Rat).String() }
func (f Float) Type() ValueType   { return FLOAT }

// PAIREDIONAL
type Ratio Rat

func (r Ratio) Eval() Value       { return Rat(r).Eval() }
func (r Ratio) Serialize() []byte { return []byte(Rat(r).String()) }
func (r Ratio) String() string    { return r.String() }
func (r Ratio) Type() ValueType   { return RATIONAL }

// PAIR
// implements KeyVal interface
type KeyValue Pair

func (k KeyValue) Eval() Value       { return Pair(k).Eval() }
func (k KeyValue) Serialize() []byte { return Pair(k).Serialize() }
func (k KeyValue) Type() ValueType {
	if k()[0].Eval().Type()&INDEXED != 0 {
		return NUMERIC
	} else {
		return SYMBOLIC
	}
}
func (t KeyValue) Key() Value   { return t()[0].Eval() }
func (t KeyValue) Value() Value { return t()[1].Eval() }
func (t KeyValue) String() string {
	return string(t.Key().Eval().String() + ": " + t.Value().Eval().String())
}

//func newMatrix() Matrix {
//	var ret Matrix
//	ret = func() Tab { return func() (c Collection) { return c } }
//	return ret
//}
//}

// INSTANCIATE A NEW VALUE
//
// values are represented internaly by either a Big, Rat, or Col type instance,
// each of which implement the absVal interface. Implemented as functional
// types, that return a value of destince type, either *big.Int, *big.Rat, or
// an Instance of the Collection interface. A Method set defined on the
// function type implements the absVal interface, by manipulating the main
// return value.
func Eval(i ...interface{}) Value {
	for _, e := range i {
		// if allready a value return
		if val, ok := e.(Value); ok {
			return val
		}
		// if not defined yet, allocate and return empty value
		if len(i) == 0 {
			return Empty(func() struct{} { return struct{}{} })
		}
		// if it is a pair of values
		if len(i) == 2 { // → assume key/value pair
			return Pair(func() [2]Value { return [2]Value{Eval(i[0]), Eval(i[1])} })
		}
		// if assertable to slice of interfaces
		if len(i) > 2 { // → pass to evalCollection
			return evalCollection(i...)
		}
	}
	// everything else is converted to value from its native tupe.
	return nativeToValue(i)
}
func nativeToValue(i interface{}) Value {

	var retFn Value

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
		var Fn Integer = func() *big.Int { return v }
		retFn = Fn

	case []byte: // a bytes slice gets assigned as bytes
		v := (new(big.Int).SetBytes(i.([]byte)))
		retFn = Bytes(func() *big.Int { return v })

	case uint8, uint32:
		v := (new(big.Int).SetBytes(i.([]byte)))
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
