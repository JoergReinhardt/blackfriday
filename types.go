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
	"fmt"
	"math/big"
	"strconv"
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
	Bytes() []byte
	String() string
	ToType(ValueType) Value
}
type Values interface {
	Value
	Container
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
	ratVal   big.Rat
	boolVal  bool
	byteVal  uint8
	bytesVal []byte
	strVal   string
	vecVal   struct {
		CntType
		Container
	}
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

// INTERFACE METHODS
//
// methods that share a name, need to be implemented once per receiving type
func (flagVal) Type() ValueType  { return FLAG }
func (intVal) Type() ValueType   { return INTEGER }
func (ratVal) Type() ValueType   { return RATIONAL }
func (boolVal) Type() ValueType  { return BOOL }
func (byteVal) Type() ValueType  { return BYTE }
func (bytesVal) Type() ValueType { return BYTES }
func (strVal) Type() ValueType   { return STRING }
func (vecVal) Type() ValueType   { return VECTOR }

func (v flagVal) Value() Value  { return v }
func (v intVal) Value() Value   { return v }
func (v ratVal) Value() Value   { return v }
func (v boolVal) Value() Value  { return v }
func (v byteVal) Value() Value  { return v }
func (v bytesVal) Value() Value { return v }
func (v strVal) Value() Value   { return v }
func (v vecVal) Value() Value   { return v }

// EMPTY TYPE
// the implementation of the empty type, will be instanciateable from thin air,
// used as base value type, that can convert a native value into the
// appropriate Value instance. It´s also used as a return type, whenever a type
// conversion can not be performed, or a variable lookup fails.
func (emptyVal) Type() ValueType         { return NIL }
func (emptyVal) Empty() emptyVal         { return emptyVal(struct{}{}) }
func (e emptyVal) Value() Value          { return emptyVal(struct{}{}) }
func (e emptyVal) Bytes() []byte         { return []byte{} }
func (e emptyVal) String() string        { return "" }
func (emptyVal) Set(v interface{}) Value { return NativeToValue(v) }
func (v emptyVal) ToType(t ValueType) Value {
	return emptyVal{}
}

// BOOL TYPE
func (b boolVal) Boolean() bool { return bool(b) }
func (b boolVal) Native() bool  { return b.Boolean() }
func (b boolVal) Integer() int {
	if b {
		return 1
	} else {
		return 0
	}
}
func (b boolVal) Bytes() []byte { return []byte(b.String()) }
func (b boolVal) String() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
func (v boolVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = flagVal(*big.NewInt(int64(v.Integer())))
	case INTEGER:
		val = intVal(*big.NewInt(int64(v.Integer())))
	case FLOAT:
		val = ratVal(*big.NewRat(1, 1).SetFloat64(float64(v.Integer())))
	case BYTES:
		val = bytesVal(v.Bytes())
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// FLAG TYPE
func (v flagVal) Flag() flagVal    { return v }
func (v flagVal) BigInt() *big.Int { i := big.Int(v); return &i }
func (v flagVal) Integer() int64   { return v.BigInt().Int64() }
func (v flagVal) Uint() uint64     { return v.BigInt().Uint64() }
func (v flagVal) Native() uint64   { return v.Uint() }
func (v flagVal) Bytes() []byte    { return []byte(v.String()) }
func (v flagVal) String() string   { return fmt.Sprint((v.Boolean())) }
func (v flagVal) Boolean() bool {
	if v.Integer() > 0 {
		return true
	} else {
		return false
	}
}
func (v flagVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = flagVal(big.Int(v))
	case INTEGER:
		val = intVal(big.Int(v))
	case FLOAT:
		i := big.Int(v)
		val = ratVal(*big.NewRat(1, 1).SetInt64((&i).Int64()))
	case BYTES:
		i := big.Int(v)
		val = bytesVal((&i).Bytes())
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = strVal(string(v.Bytes()))
	default:
		val = emptyVal{}
	}
	return val
}

// INTEGER TYPE
func (v intVal) BigInt() *big.Int     { r := big.Int(v); return &r }
func (v intVal) BigRat() *big.Rat     { return big.NewRat(v.Integer(), 1) }
func (v intVal) BigFloat() *big.Float { return big.NewFloat(v.Float()) }
func (v intVal) Flag() flagVal        { return flagVal(*v.BigInt()) }
func (v intVal) Uint() uint64         { return (v.BigInt()).Uint64() }
func (v intVal) Integer() int64       { return v.BigInt().Int64() }
func (v intVal) Native() int64        { return v.Integer() }
func (v intVal) Float() float64       { return float64(v.BigInt().Int64()) }
func (v intVal) Bytes() []byte        { return []byte(v.String()) }
func (v intVal) String() string       { return fmt.Sprint((v.BigInt().Int64())) }
func (v intVal) Boolean() bool {
	if v.Integer() > 0 {
		return true
	} else {
		return false
	}
}
func (v intVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = v.Flag()
	case INTEGER:
		val = intVal(*v.BigInt())
	case FLOAT:
		val = ratVal(*v.BigRat())
	case BYTES:
		val = bytesVal(v.Bytes())
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = strVal(string(v.Bytes()))
	default:
		val = emptyVal{}
	}
	return val
}

// FLOAT TYPE
func (v ratVal) BigRat() *big.Rat     { r := big.Rat(v); return &r }
func (v ratVal) BigInt() *big.Int     { return big.NewInt(int64(v.Float())) }
func (v ratVal) BigFloat() *big.Float { return big.NewFloat((v).Float()) }
func (v ratVal) Float() float64       { f, _ := (v).BigRat().Float64(); return f }
func (v ratVal) Native() float64      { return v.Float() }
func (v ratVal) Integer() int64       { return v.BigInt().Int64() }
func (v ratVal) String() string       { r, _ := v.BigRat().Float64(); return fmt.Sprint(r) }
func (v ratVal) Bytes() []byte        { return []byte(v.String()) }
func (v ratVal) Boolean() bool {
	if v.Float() > 0.0 {
		return true
	} else {
		return false
	}
}
func (v ratVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = ratVal(*v.BigRat())
	case INTEGER:
		val = intVal(*v.BigInt())
	case FLOAT:
		val = ratVal(*v.BigRat())
	case BYTES:
		val = bytesVal(v.BigInt().Bytes())
	case BYTE:
		if len(v.BigInt().Bytes()) > 0 {
			val = byteVal(v.BigInt().Bytes()[0])
		} else {
			val = byteVal(byte(0))
		}
	case STRING:
		val = strVal(string(v.BigInt().Bytes()))
	default:
		val = emptyVal{}
	}
	return val
}

// SINGLE BYTE TYPE
func (v byteVal) String() string   { return string(v.Bytes()) }
func (v byteVal) Bytes() []byte    { return []byte{v.Byte()} }
func (v byteVal) Byte() byte       { return byte(v) }
func (v byteVal) Native() uint8    { return v.Uint() }
func (v byteVal) Uint() uint8      { return uint8(v) }
func (v byteVal) Integer() int64   { return int64(v) }
func (v byteVal) BigInt() *big.Int { return big.NewInt(v.Integer()) }
func (v byteVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = flagVal(*big.NewInt(int64(v.Uint())))
	case INTEGER:
		val = intVal(*big.NewInt(int64(v.Uint())))
	case FLOAT:
		val = ratVal(*big.NewRat(1, 1).SetInt64(int64(v.Uint())))
	case BYTES:
		val = bytesVal([]byte{v.Byte()})
	case BYTE:
		val = byteVal([]byte{v.Byte()}[0])
	case STRING:
		val = strVal(string(v.Byte()))
	default:
		val = emptyVal{}
	}
	return val
}

// BYTE SLICE TYPE
func (v bytesVal) Bytes() []byte    { return []byte(v) }
func (v bytesVal) Native() []byte   { return v.Bytes() }
func (v bytesVal) String() string   { return string(v) }
func (v bytesVal) Uint() uint64     { return big.NewInt(0).SetBytes(v).Uint64() }
func (v bytesVal) Integer() int64   { return big.NewInt(0).SetBytes(v).Int64() }
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
func (v bytesVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = flagVal(*big.NewInt(0).SetBytes([]byte(v)))
	case INTEGER:
		val = intVal(*big.NewInt(0).SetBytes([]byte(v)))
	case FLOAT:
		val = ratVal(*big.NewRat(1, 1).SetInt64(int64(v.Integer())))
	case BYTES:
		val = bytesVal([]byte(v))
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// STRING TYPE
func (v strVal) Bytes() []byte  { return []byte(v) }
func (v strVal) String() string { return string(v) }
func (v strVal) Native() string { return v.String() }
func (v strVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case BOOL:
		vari, err := parseBool(v)
		if err != nil {
			val = emptyVal{}
		}
		val = NativeToValue(vari)
	case FLAG:
		vari, err := parseInt(v)
		if err != nil {
			val = emptyVal{}
		}
		val = NativeToValue(vari)
	case INTEGER:
		vari, err := parseInt(v)
		if err != nil {
			val = emptyVal{}
		}
		val = NativeToValue(vari)
	case FLOAT:
		vari, err := parseBool(v)
		if err != nil {
			val = emptyVal{}
		}
		val = NativeToValue(vari)
	case BYTES:
		val = bytesVal([]byte(string(v)))
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = v
	default:
		val = emptyVal{}
	}
	return val
}

// PARSE STRINGS OR BYTE SLICES TO NUMERALS
func parseInt(v Value) (Value, error) {
	i, err := strconv.Atoi(v.String())
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}
func parseUint(v Value) (Value, error) {
	i, err := strconv.ParseUint(v.String(), 2, 8)
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}
func parseFloat(v Value) (Value, error) {
	f, err := strconv.ParseFloat(v.String(), 10)
	if err != nil {
		return nil, err
	}
	return NativeToValue(f), nil
}
func parseBool(v Value) (Value, error) {
	i, err := strconv.ParseBool(v.String())
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}

// convenience conversion functions
func (v vecVal) Vector() Value   { return v }
func (v vecVal) Native() []Value { return v.Container.Values() }
func (v vecVal) Bytes() []byte {
	var bytes = []byte{}
	var vals = v.Container.Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Bytes()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v vecVal) String() string { return string(v.Bytes()) }

func (v vecVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = emptyVal{}
	case INTEGER:
		val = emptyVal{}
	case FLOAT:
		val = emptyVal{}
	case BYTES:
		val = bytesVal(v.Bytes())
	case BYTE:
		val = byteVal(v.Bytes()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

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
func NativeToValue(i interface{}) (v Value) {
	switch i.(type) {
	// if its allready a value, just return it
	case Value:
		v = i.(Value).Value()
		// convert booleans to big.Int
	case bool:
		if i.(bool) {
			v = flagVal(*big.NewInt(1))
		} else {
			v = flagVal(*big.NewInt(0))
		}
		// convert all native unsigned integer types to big.Int flags
		// for bitwise comparison
	case uint, uint16, uint32, uint64:
		v = flagVal(*big.NewInt(int64(i.(uint))))

		// convert all native integer types to big.Int
	case int, int8, int16, int32, int64:
		v = intVal(*big.NewInt(int64(i.(int))))
		// floats get stored as big.Rat rational types, to not loose
		// praesition where not nescessary
	case float32, float64:
		v = ratVal(*big.NewRat(1, 1).SetFloat64(i.(float64)))
		// convert a single byte to a byte val
	case byte:
		v = byteVal(i.(uint8))
		// convert byte slice to bytes value
	case []byte:
		v = bytesVal(i.([]byte))
		// while being technicaly the same as the byte slice type, a
		// string type helps to keep tagged and/or formated text apart
		// from not yet parsed input.
	case string:
		v = strVal(i.(string))
		// when a slice of values is passed, it is intended to be
		// encapsulated in a container.
	case []Value:
		val := wrapContainer(LIST_ARRAY, nil)
		val.(List).Add(i.([]Value)...)
		v = vecVal{LIST_ARRAY, val}
	}
	return v
}
