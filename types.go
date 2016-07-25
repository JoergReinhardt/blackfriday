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
	"math/big"
)

// VALUE INTERFACE
//
// the base types are kept as simple as possible.
//
// NIL
//
// while returning an error and/or boolean when values don't exist, or turn out
// not to be converteable is one way of doing it, I think having a nil type to
// return instead, makes things much easyer and there are lots of additional
// reasons to have one anyway.
//
// FLAG
//
// since the type field is the essential difference to ordinary values, I
// decided to have bitflags as a type, to make the type acsess-, compare- and
// usable.
//
// INTEGER, RATIONAL & FLOAT
//
// parsing values embedded in text and perform calculations on them is the main
// goal of agiledocument, which makesnumbers are the essence of the agile
// document. they are implemented by math librarys big.int / big.float types,
// since those allready implement all nescessary type conversions, parsing from
// string included and they are highly optimized to perform well. flag type is
// another big.Int, since it also implements bitwise boolean operations.
//
// BYTES
//
// the source will be provided by blackfriday in form of a byte slice
//
type Value interface {
	Type() ValueType
	Value() Value
	ToValue(Value)
}
type Values interface {
	Value
	Values() []Value
	ToValues([]Value)
}
type KeyValue interface {
	Ident() string
	Value // Type() ValueType, Value() Value
}

// strings make the input handable as text.
//go:generate -command stringer -type ValueType
const (
	NIL  ValueType = 0
	BOOL ValueType = 1 << iota
	INTEGER
	FLOAT
	RATIONAL
	FLAG
	BYTE
	BYTES
	STRING
	VECTOR
)

type ValueType uint

// TYPES
//
type ( // are kept as close to the native types they are derived from, as possible
	emptyVal struct{} // the empty struct is taken as emptyVal
	flagVal  big.Int
	intVal   big.Int
	floatVal big.Float
	ratVal   big.Rat
	boolVal  bool
	byteVal  byte
	bytesVal []byte
	strVal   string
	vecVal   struct{ Container }
)

// funtion types that implement the interface
type (
	// for each interface method there is a function defined, representing it. By sharing it's type
	// signature and using a Value instance as its first parameter and only return type.
	//
	ValueFnT func(Value) Value // TYPE CONVERTING & RETURN METHOD
	TypeFnT  func(Value) ValueType
	IdentFnT func(KeyValue) string
	///
	/// FUNCTIONAL TYPE REPRESENTATION
	//
	// to convert from each allowed type to every type a value can be
	// converted to, without a lot of repeating, function types to
	// represent the native types are needed.
	//
	// the type system makes A function type that gets an interface type as
	// first parameter in it´s definition, acts like a method that works on
	// all instances that implement its receiving Interface parameter.
	BoolFnT     func(Value) bool
	BigIntFnT   func(Value) big.Int
	BigFloatFnT func(Value) big.Float
	BigRatFnT   func(Value) big.Rat
	ByteFnT     func(Value) byte
	BytesFnT    func(Value) []byte
	StringFnT   func(Value) string
	IntegerFnT  func(Value) int64
	FloatFnT    func(Value) float64
	ValuesFnT   func(Value) []Value
)

func RetValue(v Value) Value                      { return v }
func RetValues(v Values) Values                   { return v }
func RetInterSlice(v []interface{}) []interface{} { return v }

// INTERFACE METHODS
//
// methods that share a name, need to be implemented once per receiving type
func (emptyVal) Type() ValueType { return NIL }
func (flagVal) Type() ValueType  { return FLAG }
func (intVal) Type() ValueType   { return INTEGER }
func (floatVal) Type() ValueType { return FLOAT }
func (ratVal) Type() ValueType   { return RATIONAL }
func (boolVal) Type() ValueType  { return BOOL }
func (byteVal) Type() ValueType  { return BYTE }
func (bytesVal) Type() ValueType { return BYTES }
func (strVal) Type() ValueType   { return STRING }
func (vecVal) Type() ValueType   { return VECTOR }

func (v emptyVal) Value() Value { return v }
func (v flagVal) Value() Value  { return v }
func (v intVal) Value() Value   { return v }
func (v floatVal) Value() Value { return v }
func (v ratVal) Value() Value   { return v }
func (v boolVal) Value() Value  { return v }
func (v byteVal) Value() Value  { return v }
func (v bytesVal) Value() Value { return v }
func (v strVal) Value() Value   { return v }
func (v vecVal) Value() Value   { return v }

// tyoed return functions return values of generic type, each is implemented by
// at least the types containing that native type, and all that can be
// converted to it.
func (emptyVal) Empty() emptyVal { return emptyVal(struct{}{}) }

func (b boolVal) Boolean() bool { return bool(b) }
func (b boolVal) Integer() int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (v flagVal) Flag() flagVal   { return v }
func (v flagVal) BigInt() big.Int { return big.Int(v) }
func (v flagVal) Uint() uint64    { return v.Uint() }
func (v flagVal) Integer() int64  { return v.Integer() }
func (v flagVal) Boolean() bool {
	if v.Integer() > 0 {
		return true
	} else {
		return false
	}
}

func (v intVal) BigInt() *big.Int     { r := big.Int(v); return &r }
func (v intVal) BigRat() *big.Rat     { return big.NewRat(v.Integer(), 1) }
func (v intVal) BigFloat() *big.Float { return big.NewFloat(v.Float()) }
func (v intVal) Flag() flagVal        { return flagVal(*v.BigInt()) }
func (v intVal) Uint() uint64         { return (v.BigInt()).Uint64() }
func (v intVal) Integer() int64       { return v.BigInt().Int64() }
func (v intVal) Float() float64       { return float64(v.BigInt().Int64()) }
func (v intVal) Bytes() []byte        { return []byte(v.Bytes()) }
func (v intVal) Boolean() bool {
	if v.Integer() > 0 {
		return true
	} else {
		return false
	}
}

func (v ratVal) BigRat() *big.Rat     { r := big.Rat(v); return &r }
func (v ratVal) BigInt() *big.Int     { return big.NewInt(v.Integer()) }
func (v ratVal) BigFloat() *big.Float { return (v).BigFloat() }
func (v ratVal) Float() float64       { f, _ := (v).BigFloat().Float64(); return f }
func (v ratVal) Integer() int64       { return v.BigInt().Int64() }
func (v ratVal) Boolean() bool {
	if v.Float() > 0.0 {
		return true
	} else {
		return false
	}
}

func (v bytesVal) Bytes() []byte    { return []byte(v) }
func (v bytesVal) String() string   { return string(v) }
func (v bytesVal) Uint() uint64     { return big.NewInt(0).SetBytes(v).Uint64() }
func (v bytesVal) Flag() flagVal    { return flagVal(*big.NewInt(0).SetBytes(v.Bytes())) }
func (v bytesVal) BigInt() *big.Int { return big.NewInt(0).SetBytes(v.Bytes()) }
func (v bytesVal) Uints() (r []uint8) {
	r = []uint8{}
	for _, val := range v {
		b := uint8(val)
		r = append(r, b)
	}
	return r
}
func (v bytesVal) Boolean() bool {
	// see if any bit is set and return
	if len(v) > 1 { // true, when set
		return true
	} else { // false, when no bits are set
		return false
	}
}

func (v strVal) Bytes() []byte  { return []byte(v) }
func (v strVal) String() string { return string(v) }

// convenience conversion functions
func (v vecVal) Vector() Value { return v }

// TO-TYPE | FROM-TYPE
//
// if the native tyoe is allready known at the time of initialization,
// reflection can be omitted.
//func toType(i Value, t ValueType) (v Value) {
//	switch t {
//	case NIL:
//		v = emptyVal{}
//	case FLAG:
//		v = &flagVal{big.NewInt(i.(*flagVal).Integer())}
//	case INTEGER:
//		v = &intVal{big.NewInt(i.(*intVal).Int64())}
//	case FLOAT:
//		f, _ := i.(*floatVal).Float64()
//		v = &floatVal{big.NewRat(1, 1).SetFloat64(f)}
//	case BYTES:
//		v = bytesVal(i.(*bytesVal).Bytes())
//	case BYTE:
//		v = i.(byteValue)
//	case STRING:
//		v = strValue(i.(*strValue).String())
//	default:
//		v = emptyVal{}
//	}
//	return v
//}
func fromType(r Value, v Value) {}

// INSTANCIATE A NEW VALUE
//
// instanciate arbitratry values of native type and return instances
// implementing the value interface.
//func NewValue(i interface{}) Value {
//	var v Value
//	switch i.(type) {
//	case bool:
//		if i.(bool) {
//			return &flagVal{big.NewInt(1)}
//		} else {
//			return &flagVal{big.NewInt(1)}
//		}
//	case uint, uint16, uint32, uint64:
//		v = &flagVal{big.NewInt(int64(i.(uint)))}
//	case int, int8, int16, int32, int64:
//		v = &intVal{big.NewInt(int64(i.(int)))}
//	case float32, float64:
//		v = &floatVal{big.NewRat(1, 1).SetFloat64(i.(float64))}
//	case byte:
//		v = byteValue(i.(byte))
//	case []byte:
//		v = bytesVal(i.([]byte))
//	case string:
//		v = strValue(i.(string))
//	}
//	return v
//}
