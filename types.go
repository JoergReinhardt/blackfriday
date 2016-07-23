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
type Val interface {
	Type() ValType
	Value() Val
}
type KeyVal interface {
	Val
	Key() string
}

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
// INTEGER & FLOAT
//
// parsing values embedded in text and perform calculations on them is the main
// goal of agiledocument, which makesnumbers are the essence of the agile
// document. they are implemented by math librarys big.int / big.float types,
// since those allready implement all nescessary type conversions, parsing from
// string included and they are highly optimized to perform well. flag type is
// another big.Int, since it also implements bitwise boolean operations.
//
// BYTE
//
// the source will be provided by blackfriday in form of a byte slice
//
// strings make the input handable as text.
//go:generate -command stringer -type ValType
const (
	NIL  ValType = 0
	FLAG ValType = 1 << iota
	INTEGER
	FLOAT
	BYTE
	STRING
)

type ValType uint

// TYPES
type ( // are kept as close to the original types as possible
	emptyVal struct{}           // emptyValue
	flagVal  struct{ *big.Int } // all big based types are enveloped
	intVal   struct{ *big.Int } // by strings to encapsulate the pointer
	floatVal struct{ *big.Rat }
	byteVal  []byte
	strVal   string
)

// funtion types that implement the interface
type (
	typeFunc func(Val) ValType
	valFunc  func(Val) Val
)

// INTERFACE METHODS
// methods that share a name, need to be implemented once per receiving type
func (emptyVal) Type() ValType { return NIL }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (floatVal) Type() ValType { return FLOAT }
func (byteVal) Type() ValType  { return BYTE }
func (strVal) Type() ValType   { return STRING }

func (v emptyVal) Value() Val { return v }
func (v flagVal) Value() Val  { return v }
func (v intVal) Value() Val   { return v }
func (v floatVal) Value() Val { return v }
func (v byteVal) Value() Val  { return v }
func (v strVal) Value() Val   { return v }

// tyoed return functions return values of generic type, each is implemented by
// at least the types containing that native type, and all that can be
// converted to it.
func (emptyVal) Empty() emptyVal   { return emptyVal{} }
func (v flagVal) Flag() *big.Int   { return v.Int }
func (v intVal) Integer() *big.Int { return v.Int }
func (v floatVal) Flt() *big.Rat   { return v.Rat }
func (v byteVal) Byte() []byte     { return v }
func (v strVal) String() string    { return string(v) }

// if the native tyoe is allready known at the time of initialization,
// reflection can be omitted.
func NewTypedVal(t ValType, i interface{}) Val {
	var v Val
	switch t {
	case NIL:
		v = emptyVal{}
	case FLAG:
		v = flagVal{big.NewInt(int64(i.(int)))}
	case INTEGER:
		v = intVal{big.NewInt(int64(i.(int)))}
	case FLOAT:
		v = floatVal{big.NewRat(0, 0).SetFloat64(i.(float64))}
	case BYTE:
		v = byteVal(i.([]byte))
	case STRING:
		v = strVal(i.(string))
	}
	return v
}

// arbitratry values will be performed to the appropriate type, or an empty
// value will be returned.
func NewVal(i interface{}) Val {
	var v Val
	switch i.(type) {
	case int, int8, int16, int32, int64, *big.Int:
		v = NewTypedVal(INTEGER, i)
	case float32, float64, *big.Float:
		v = NewTypedVal(FLOAT, i)
	case []byte:
		v = NewTypedVal(BYTE, i)
	case string:
		v = NewTypedVal(STRING, i)
	}
	return v
}
