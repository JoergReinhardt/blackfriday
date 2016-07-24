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
// INTEGER & FLOAT
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
type Val interface {
	Type() ValType
	Value() Val
}
type KeyVal interface {
	Val
	Key() string
}

// strings make the input handable as text.
//go:generate -command stringer -type ValType
const (
	NIL  ValType = 0
	BOOL ValType = 1 << iota
	INTEGER
	FLOAT
	FLAG
	BYTE
	BYTES
	STRING
	VECTOR
)

type ValType uint

// TYPES
//
type ( // are kept as close to the original types as possible
	emptyVal struct{}           // emptyValue
	flagVal  struct{ *big.Int } // all big based types are enveloped
	intVal   struct{ *big.Int } // by strings to encapsulate the pointer
	floatVal struct{ *big.Rat } // nature of the embedded type (otherwise
	// methods get hard to define in a way that respects their signature)
	boolVal  bool
	byteVal  byte
	bytesVal []byte
	strVal   string
	vecVal   struct{ Container }
)

// funtion types that implement the interface
type (
	typeFunc func(Val) ValType
	valFunc  func(Val) Val
)

// INTERFACE METHODS
//
// methods that share a name, need to be implemented once per receiving type
func (emptyVal) Type() ValType { return NIL }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (floatVal) Type() ValType { return FLOAT }
func (boolVal) Type() ValType  { return BOOL }
func (byteVal) Type() ValType  { return BYTE }
func (bytesVal) Type() ValType { return BYTES }
func (strVal) Type() ValType   { return STRING }
func (vecVal) Type() ValType   { return VECTOR }

func (v emptyVal) Value() Val { return v }
func (v flagVal) Value() Val  { return v }
func (v intVal) Value() Val   { return v }
func (v floatVal) Value() Val { return v }
func (v boolVal) Value() Val  { return v }
func (v byteVal) Value() Val  { return v }
func (v bytesVal) Value() Val { return v }
func (v strVal) Value() Val   { return v }
func (v vecVal) Value() Val   { return v }

// tyoed return functions return values of generic type, each is implemented by
// at least the types containing that native type, and all that can be
// converted to it.
func (*emptyVal) Empty() emptyVal { return emptyVal(struct{}{}) }

func (b *boolVal) Boolean() bool { return bool(*b) }
func (b *boolVal) Integer() int {
	if *b {
		return 1
	} else {
		return 0
	}
}

func (v *flagVal) Flag() *flagVal   { return &flagVal{(*v).Int} }
func (v *flagVal) Uint() uint64     { return (*v.Int).Uint64() }
func (v *flagVal) Integer() int64   { return (*v.Int).Int64() }
func (v *flagVal) BigInt() *big.Int { return v.Int }
func (v *flagVal) Boolean() bool {
	if v.Int64() > 0.0 {
		return true
	} else {
		return false
	}
}

func (v *intVal) Flag() *flagVal       { return &flagVal{(*v).BigInt()} }
func (v *intVal) Uint() uint64         { return (*v.Int).Uint64() }
func (v *intVal) Integer() int64       { return (*v.Int).Int64() }
func (v *intVal) Float() float64       { return float64((*v.Int).Int64()) }
func (v *intVal) BigInt() *big.Int     { return (*v).Int }
func (v *intVal) BigRat() *big.Rat     { return big.NewRat(v.Integer(), 1) }
func (v *intVal) BigFloat() *big.Float { return big.NewFloat(float64((*v).Int64())) }
func (v *intVal) Boolean() bool {
	if (*v).Int64() > 0 {
		return true
	} else {
		return false
	}
}

func (v *floatVal) Float() float64       { r, _ := (*v.Rat).Float64(); return r }
func (v *floatVal) Integer() int64       { r, _ := (*v.Rat).Float64(); return int64(r) }
func (v *floatVal) BigRat() *big.Rat     { return v.Rat }
func (v *floatVal) BigInt() *big.Int     { f, _ := (*v.Rat).Float64(); return big.NewInt(int64(f)) }
func (v *floatVal) BigFloat() *big.Float { f, _ := (*v.Rat).Float64(); return big.NewFloat(f) }
func (v *floatVal) Boolean() bool {
	f, _ := (*v).Float64()
	if f > 0.0 {
		return true
	} else {
		return false
	}
}

func (v *bytesVal) Bytes() []byte    { return []byte(*v) }
func (v *bytesVal) String() string   { return string(*v) }
func (v *bytesVal) Uint() uint64     { return big.NewInt(0).SetBytes(*v).Uint64() }
func (v *bytesVal) Flag() *flagVal   { return &flagVal{big.NewInt(0).SetBytes((*v).Bytes())} }
func (v *bytesVal) BigInt() *big.Int { return big.NewInt(0).SetBytes((*v).Bytes()) }
func (v *bytesVal) Uints() (r []uint8) {
	r = []uint8{}
	for _, val := range *v {
		b := uint8(val)
		r = append(r, b)
	}
	return r
}
func (v *bytesVal) Boolean() bool {
	// see if any bit is set and return
	if len(*v) > 1 { // true, when set
		return true
	} else { // false, when no bits are set
		return false
	}
}

func (v *strVal) Bytes() []byte  { return []byte(*v) }
func (v *strVal) String() string { return string(*v) }

// convenience conversion functions
func (v *vecVal) Vector() Val { return v }

// TO-TYPE | FROM-TYPE
//
// if the native tyoe is allready known at the time of initialization,
// reflection can be omitted.
func toType(i Val, t ValType) (v Val) {
	switch t {
	case NIL:
		v = emptyVal{}
	case FLAG:
		v = &flagVal{big.NewInt(i.(*flagVal).Integer())}
	case INTEGER:
		v = &intVal{big.NewInt(i.(*intVal).Int64())}
	case FLOAT:
		f, _ := i.(*floatVal).Float64()
		v = &floatVal{big.NewRat(1, 1).SetFloat64(f)}
	case BYTES:
		v = bytesVal(i.(*bytesVal).Bytes())
	case BYTE:
		v = i.(byteVal)
	case STRING:
		v = strVal(i.(*strVal).String())
	default:
		v = emptyVal{}
	}
	return v
}
func fromType(r Val, v Val) {}

// INSTANCIATE A NEW VALUE
//
// instanciate arbitratry values of native type and return instances
// implementing the value interface.
func NewVal(i interface{}) Val {
	var v Val
	switch i.(type) {
	case bool:
		if i.(bool) {
			return &flagVal{big.NewInt(1)}
		} else {
			return &flagVal{big.NewInt(1)}
		}
	case uint, uint16, uint32, uint64:
		v = &flagVal{big.NewInt(int64(i.(uint)))}
	case int, int8, int16, int32, int64:
		v = &intVal{big.NewInt(int64(i.(int)))}
	case float32, float64:
		v = &floatVal{big.NewRat(1, 1).SetFloat64(i.(float64))}
	case byte:
		v = byteVal(i.(byte))
	case []byte:
		v = bytesVal(i.([]byte))
	case string:
		v = strVal(i.(string))
	}
	return v
}
