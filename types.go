package agiledoc

import (
	//"fmt"
	// col "github.com/emirpasic/gods/containers"
	"math/big"
	"sync"
)

type (
	intPool  sync.Pool
	ratPool  sync.Pool
	pairPool sync.Pool
)

var (
	intGen  intPool  = intPool{}
	ratGen  ratPool  = ratPool{}
	pairGen pairPool = pairPool{}
)

func newVal() val   { return func() *big.Int { return intGen.New().(*big.Int) } }
func newRat() rat   { return func() *big.Rat { return ratGen.New().(*big.Rat) } }
func newPair() pair { return func() [2]Evaluable { return pairGen.New().([2]Evaluable) } }

func init() {
	intGen.New = func() interface{} { return big.NewInt(0) }
	ratGen.New = func() interface{} { return big.NewRat(1, 1) }
	pairGen.New = func() interface{} { return [2]Evaluable{} }
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
	val func() *big.Int

	// paired types
	rat  func() *big.Rat
	pair func() [2]Evaluable

	// collection types see collection,go
)

/////////////////////////////////////////////////////////////////////////
func (b val) Eval() Evaluable   { return Value(b) }
func (b val) Serialize() []byte { return []byte(b().String()) }
func (b val) String() string    { return b().String() }
func (b val) Type() ValueType   { return TERMINAL }

////////////////////////////////////////////////////////////////
// private methods, to convert to native types
func (b val) bool() bool {
	if b().Int64() > 0 {
		return true
	} else {
		return false
	}
}
func (b val) bigInt() *big.Int { return b() }
func (b val) int() int         { return int(b().Int64()) }
func (b val) int64() int64     { return b().Int64() }
func (b val) uint64() uint64   { return b().Uint64() }
func (b val) bigRat() *big.Rat { return newRat()().SetFrac(Value(ONE).(val)(), b()) }
func (b val) float64() float64 { f, _ := b.bigRat().Float64(); return f }
func (b val) bytes() []byte    { return b().Bytes() }

// public methods to convert to other implementations of evaluable
func (b val) Rat() rat         { return Value(newRat()().SetFrac(newVal()().SetInt64(1), b.bigInt())).(rat) }
func (b val) Pair() pair       { return newPair().SetValue(b.Eval()) }
func (b val) BitFlag() BitFlag { return BitFlag(b) }
func (b val) Flag() Bool       { return Bool(b) }
func (b val) Integer() Integer { return Integer(b) }
func (b val) Text() Text       { return Text(b) }
func (b val) Bytes() Bytes     { return Bytes(b) }

/////////////////////////////////////////////////////////////////////////
func (r rat) Eval() Evaluable { return Value(r) }

// Bytes is supposed to keep as much information as possible, so this converts
// numerator and denominator to 64 bytes each, ignoring the original accuracy
// (length), to make them divideable again. Accuracys greater 64bit should not
// be serialized, but kept in absolute numbers in memoru during calculations,
func (r rat) Bytes() []byte {
	return append(
		newVal()().SetInt64(r().Num().Int64()).Bytes(),
		newVal()().SetInt64(r().Denom().Int64()).Bytes()...,
	)
}
func (r rat) Serialize() []byte { return []byte(r().String()) }
func (r rat) String() string    { return r().String() }
func (r rat) Type() ValueType   { return REAL }

////////////////////////////////////////////////////////////////
// private methods, to convert to native types
func (r rat) float64() float64 { f, _ := r().Float64(); return f }
func (r rat) bigRat() *big.Rat { return Value(r).(rat)() }

// public methods to convert to other implementations of evaluable
func (r rat) Float() Float  { return Float(r) }
func (r rat) Rational() rat { return r }
func (r rat) Pair() pair    { return pair(func() [2]Evaluable { return [2]Evaluable{r.Num(), r.Denom()} }) }

// methods that take or return the integer type, to set, or get contained values
func (r rat) Num() Integer   { return Value(r().Num()).(Integer) }
func (r rat) Denom() Integer { return Value(r().Denom()).(Integer) }
func (r rat) SetNum(v Integer) rat {
	return rat(func() *big.Rat { return newRat()().SetFrac(v(), r().Denom()) })
}
func (r rat) SetDenom(v Integer) rat {
	return rat(func() *big.Rat { return newRat()().SetFrac(r().Num(), v()) })
}
func (r rat) SetFrac(v Paired) rat {
	return rat(func() *big.Rat { return newRat()().SetFrac(Value(v.Index()).(val).bigInt(), v.Value().(val).bigInt()) })
}

// methods to implement the Paired interface
func (r rat) Index() Integer   { return r.Num() }
func (r rat) Key() Evaluable   { return val(r.Num()).Text() }
func (r rat) Value() Evaluable { return r.Denom() }

/////////////////////////////////////////////////////////////////////////
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
		ret = Value(-1).(val).Integer() // negative → not set
	} else { // NUMERIC
		// if natural number, return as interger
		if b.Key().Type()&NATURAL != 0 {
			ret = b.Key().(val).Integer()
		}
		// if real number, return numerator as interger
		if b.Key().Type()&REAL != 0 {
			ret = b.Key().(rat).Num()
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
func (b pair) Type() ValueType { return PAIR }

/////////////////////////////////////////////////////////////////////////
// EMPTY
func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Evaluable { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }

/////////////////////////////////////////////////////////////////////////
// BOOLS ENCODED AS UNSIGNED INTEGER
type Bool val

func (u Bool) Eval() Evaluable   { return val(u).Eval().(Bool) }
func (u Bool) Serialize() []byte { return []byte(u().String()) }
func (u Bool) String() string    { return u().Text(2) }
func (u Bool) Type() ValueType   { return BOOL }

/////////////////////////////////////////////////////////////////////////
// UNSIGNED
type Unsigned val

func (u Unsigned) Eval() Evaluable   { return val(u).Eval().(Bool) }
func (u Unsigned) Serialize() []byte { return []byte(u().String()) }
func (u Unsigned) String() string    { return u().Text(2) }
func (u Unsigned) Type() ValueType   { return UINT }

/////////////////////////////////////////////////////////////////////////
// INTEGER
type Integer val

func (i Integer) Eval() Evaluable   { return val(i).Eval().(val).Integer() }
func (i Integer) Serialize() []byte { return []byte(val(i)().String()) }
func (i Integer) String() string    { return i().Text(10) }
func (i Integer) Type() ValueType   { return INTEGER }

/////////////////////////////////////////////////////////////////////////
// BYTES
type Bytes val

func (b Bytes) Eval() Evaluable   { return val(b).Eval().(val).Bytes() }
func (b Bytes) Serialize() []byte { return []byte(val(b)().String()) }
func (b Bytes) String() string    { return b().Text(8) }
func (b Bytes) Type() ValueType   { return BYTES }

/////////////////////////////////////////////////////////////////////////
// STRING
type Text val

func (s Text) Eval() Evaluable   { return val(s).Eval().(val).Text() }
func (s Text) Serialize() []byte { return []byte(s().Bytes()) }
func (s Text) String() string    { return string(s.Serialize()) }
func (s Text) Type() ValueType   { return TEXT }

/////////////////////////////////////////////////////////////////////////
// FLOAT
type Float rat

func (f Float) Eval() Evaluable   { return rat((f).Eval().(rat).Float()) }
func (f Float) Serialize() []byte { return []byte(f.String()) }
func (f Float) String() string    { return f().FloatString(10) }
func (f Float) Type() ValueType   { return FLOAT }

/////////////////////////////////////////////////////////////////////////
// RATIONAL
type Ratio rat

func (r Ratio) Eval() Evaluable   { return rat((r).Eval().(rat).Rational()) }
func (r Ratio) Serialize() []byte { return []byte(rat(r).String()) }
func (r Ratio) String() string    { return r().String() }
func (r Ratio) Type() ValueType   { return RATIONAL }

/////////////////////////////////////////////////////////////////////////
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
		v = pair(func() [2]Evaluable { return [2]Evaluable{Value(i[0]), Value(i[1])} })
	}
	if len(i) > 2 { // if more than two values are passed, we assume an indexed list of values. Should they turn out to be key Value Pairs, they will be converted to a list of maps, due to recursion.
		v = newOrderedList()
		v.(OrderedList).Add(valueSlice(i)...)
	}
	return v
}
func nativeToValue(i interface{}) Evaluable {

	var retFn Evaluable

	switch i.(type) {
	case bool: // a boolean returns a flag with the first bit set
		v := newVal()
		if i.(bool) {
			retFn = val(func() *big.Int { return v().SetInt64(int64(1)) })
		} else {
			retFn = val(func() *big.Int { return v().SetInt64(int64(0)) })
		}
	case []bool: // slice of bools gets spooled to a bitflag
		v := newVal()()
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
	case ValueType: // the underlying uint marks the type and is used as a bit f;ag
		v := (newVal()().SetUint64((uint64(i.(ValueType).Uint()))))
		retFn = BitFlag(func() *big.Int { return v })

	case uint: // a uint is assumed to be a single byte
		v := (newVal()().SetUint64((uint64(i.(uint)))))
		retFn = BitFlag(func() *big.Int { return v })

	case uint16: // a uint is assumed to be a single byte
		v := (newVal()().SetUint64((uint64(i.(uint16)))))
		retFn = BitFlag(func() *big.Int { return v })

	case uint64: // a uint is assumed to be a single byte
		v := (newVal()().SetUint64((i.(uint64))))
		retFn = BitFlag(func() *big.Int { return v })

	case int: // integers are integer
		v := (newVal()().SetInt64((int64(i.(int)))))
		retFn = Integer(func() *big.Int { return v })

	case int16: // integers are integer
		v := (newVal()().SetInt64((int64(i.(int16)))))
		retFn = Integer(func() *big.Int { return v })

	case int32: // integers are integer
		v := (newVal()().SetInt64((int64(i.(int32)))))
		retFn = Integer(func() *big.Int { return v })

	case int64: // integers are integer
		v := (newVal()().SetInt64(i.(int64)))
		retFn = Integer(func() *big.Int { return v })

	case float32: // floating point values get assigned to rationals
		v := (new(big.Rat).SetFloat64(float64(i.(float32))))
		retFn = Float(func() *big.Rat { return v })

	case float64: // floating point values get assigned to rationals
		v := (new(big.Rat).SetFloat64(i.(float64)))
		retFn = Float(func() *big.Rat { return v })

	case byte: // == uint8
		v := new(big.Int).SetBytes(i.([]byte))
		retFn = Text(func() *big.Int { return v })

	case []byte: // a bytes slice gets assigned as bytes
		v := (new(big.Int).SetBytes(i.([]byte)))
		retFn = Bytes(func() *big.Int { return v })

	case uint32: // == uint32
		v := (newVal()().SetUint64((uint64(i.(uint32)))))
		retFn = Text(func() *big.Int { return v })

	case string: // a string gets assigned by its bislice as well
		v := (new(big.Int).SetBytes([]byte(i.(string))))
		retFn = Text(func() *big.Int { return v })
	}
	return retFn
}
