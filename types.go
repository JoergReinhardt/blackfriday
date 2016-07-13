package agiledoc

import (
	// "fmt"
	"math/big"
	"strings"
	//bf "github.com/russross/blackfriday"
)

//// GENERIC VALUE TYPE AND FUNCTION DEFINITIONS// {{{
///
// represent actual values of a given type.
//
//- integers are represented by int instances from math/big
//
//- to do proper math, a float type is needed, taken from math/vbig as well
//
//- one boolean get's represented by a boolean
//
//- slice of booleans get represented by a big Int
//
//- vecConvert type is a slice of values of arbitrary kind
//
//- mtxConvert is a vector to, but also carrys to ints to define it's shape
//
type (
	/// VALUE TYPES// {{{
	// all values are wrapped in a struct. that doesent add overhead, but gives the
	// oportunity to implement the interface at the baseVar struct type, which makes
	// all structs values of type baseVar per default.
	baseVar  struct{}
	boolVal  struct{ bool }
	intVal   struct{ *big.Int }
	floatVal struct{ *big.Float }
	byteVal  struct{ byte }
	bytesVal struct{ bytes []byte }
	strVal   struct{ string }
	vecVal   struct{ vec []Value }
	mtxVal   struct {
		*vecVal
		shape [2]int
	} // }}}
	////
	/// VALUE FUNCTION TYPES// {{{
	// the base implementation of the value interface is identical for all types
	// and defined in the form of function types. Each value type implementation
	// will have to provide all of the public functions and may or may not provide
	// implementations for any, maybe even all the type functions, depending on if
	// they are convertable to the type returned by a given function, or not.
	///
	/// INTERFACE IMPLEMENTATION
	// apart from being defined as an interface, there is also a funcion
	// type signature to describes the interface. Since its taking a value
	// interface as its first parameter, it suits to fit as method for
	// every implemented tye, or instance method, for every instance
	// instanciated..
	evalFn func(Value) []byte

	/// TYPED RETURN FUNCTION TYPES
	// A type function is a function-type. It returns the value instance contained
	// in all types that implement the value interface, typed as.  given by the
	// particular funciton.  There is one type-function for each type to be
	// returned. A Type function will get instanciated for every instance of a type
	// that can be converted to the functions type.
	baseVarFn func(Value) baseVar // baseVar value to return for baseVar, or unconvertable
	boolFn    func(Value) boolVal
	intFn     func(Value) intVal
	floatFn   func(Value) floatVal
	byteFn    func(Value) byteVal
	bytesFn   func(Value) bytesVal
	strFn     func(Value) strVal
	vecFn     func(Value) vecVal
	mtxFn     func(Value) mtxVal // }}}
) // }}}

//// GENERIC FUNCTION IMPLEMENTATIONS// {{{
///
// value functions implement methods on the funcrion level, that are supposed
// to exist on all imlementations in one or another way. The first argument of
// a value function is a value. The function calls its type function to
// determin if and how to convert it. There is onw typed value return function
// per tyoe that can be returned,
var (
	// NEW VALUE FUNCTION (instance from a previously unknown type)// {{{
	NewVal = func(v interface{}) Value {
		switch v.(type) {
		case bool:
			return boolVal{
				v.(bool),
			}
		case int, int8, int16, int32, int64:
			return strVal{
				v.(string),
			}
		case float32, float64:
			return floatVal{
				big.NewFloat(v.(float64)),
			}
		case byte:
			return byteVal{
				byte(v.(float64)),
			}
		case []byte:
			return bytesVal{
				v.([]byte),
			}
		case string:
			return strVal{
				v.(string),
			}
			// CASES TO BE READ AS MATRIX
			// cpt. obvious
		case vecVal:
			return v.(vecVal)
			// slice of values, or interfaces
		case []Value, []interface{}:
			return vecVal{
				v.([]Value),
			}

			// CASES TO BE READ AS MATRIX
			// cpt. obvious
		case mtxVal:
			return v.(mtxVal)
			// flatten slice of vector values
		case []vecVal:
			vec := vecVal{
				[]Value{},
			}
			ret := mtxVal{
				&vec,
				[2]int{0, 0},
			}
			// range over contained vector instances
			for n, i := range v.([]vecVal) {
				// interface is allready converted to vecVal
				// due to the range statement assertion
				v := i
				// get length of this row (column count)
				m := Length(v)
				// append this vectors fields to the new vector
				vec.vec = append(vec.vec, v.vec...)
				// if this row happens to be longer then the
				// longest row, update column count
				if m > ret.shape[0] {
					ret.shape[0] = m
				}
				// update row count
				ret.shape[1] = n
			}
			return ret
			// flatten slice of slices of values or interfaces
		case [][]Value, [][]interface{}:
			vec := vecVal{
				[]Value{},
			}
			ret := mtxVal{
				&vec,
				// set length of outer slice as row count
				[2]int{0, len(v.([][]Value))},
			}
			for _, vals := range v.([][]Value) {

				for n, v := range vals {
					v := v
					n := n
					// append this vectors fields to the new vector
					vec.vec = append(vec.vec, v)
					// if this row happens to be longer then the
					// longest row, update column count
					if n > ret.shape[0] {
						ret.shape[0] = n
					}
				}
			}
			return ret
		}
		return baseVar{}
	} // }}}

	// NEW EMPTY VALUE FUNCTION (returns empty value of passed type)// {{{
	// returnes an empty instance of a designated type.
	NewemptyConvertVal = func(v valueType) Value {
		switch v {
		case BOOL:
			return boolVal{false}
		case INTEGER:
			return strVal{""}
		case FLOAT:
			return floatVal{big.NewFloat(0)}
		case BYTE:
			return byteVal{byte(0)}
		case BYTES:
			return bytesVal{[]byte{}}
		case STRING:
			return strVal{""}
		case VECTOR:
			return vecVal{[]Value{}}
		case MATRIX:
			return mtxVal{
				&vecVal{[]Value{}},
				[2]int{0, 0},
			}
		default:
			return baseVar{}
		}
	} // }}}

	// HELPER FUNCTIONS (length, width, heigth)// {{{
	// length function to return total length of vector and/or matrix type
	Length = func(v Value) int {
		if v.Type()&VECTOR != 0 {
			return len(v.(vecVal).vec)
		}
		if v.Type()&MATRIX != 0 {
			return v.(mtxVal).shape[0] * v.(mtxVal).shape[1]
		}
		return -1
	}

	// WIDTH returns the number of columns a matrix contains
	Width = func(v Value) int {
		if v.Type()&MATRIX != 0 {
			return v.(mtxVal).shape[0]
		}
		return -1
	}

	// HEIGTH returns the number of rows a matrix contains
	Heigth = func(v Value) int {
		if v.Type()&MATRIX != 0 {
			return v.(mtxVal).shape[1]
		}
		return -1
	} // }}}

	//// GENERIC CONVERSION FUNCTIONS (implements all possible return types, per type)// {{{
	///
	// functions to return arbitrary type as the designated return type.
	// the typed return functions are defined on the value interface and
	// therefore apply to all possible values. Each of the typed return
	// functions implements a type switch, to split by passed type
	// initially and then defines the appropriate conversion function to
	// the passed type within the switches cases.
	//
	// the designated returned value is initiated as value of the
	// appropriate contained type right above the type switch. The
	// conversion function assigns the appropriate value to that
	// preinitialized return variable.
	//
	// The switch cases may call other typed return functions to perform
	// the designated conversion.

	// RETURN EMPTY TYPE// {{{
	emptyConvert = func(Value) baseVar { return baseVar{} } // now for the tricky part...// }}}

	// RETURN BOOL TYPE// {{{
	// (if value set, or in case of string, if its reading "true" return true)
	boolConvert = func(v Value) boolVal {
		ret := false
		switch v.Type() {
		case BOOL:
			// just pass on, contained value as a new instance
			ret = v.(boolVal).bool
		case INTEGER:
			// if contained integer is greater than zero, return true,
			// false for zero and all negative values
			if v.(intVal).Int64() > 0 {
				ret = true
			} else {
				ret = false
			}
		case FLOAT:
			// if contained float is greater than zero, return true,
			// false for zero and all negative values
			val, _ := v.(floatVal).Float64()
			if val > 0 {
				ret = true
			} else {
				ret = false
			}
		case BYTE:
			// if zero, return false, for each other value return true
			if v.(byteVal).byte&byte(0) != 0 {
				ret = false
			} else {
				ret = true
			}
		case BYTES:
			// if len is zero, then "bytes" returns false. If there
			// are any bytes, then its return value is true
			if len(v.(bytesVal).bytes) > 0 {
				ret = false
			} else {
				ret = true
			}
		case STRING:
			// when the contained value turns out to be a string, try to
			// parse it, by comparing it to the word "true"
			cmp := v.(strVal).string
			if strings.Compare(cmp, "true") == 1 || strings.Compare(cmp, "TRUE") == 1 || strings.Compare(cmp, "True") == 1 {
				ret = true
			} else {
				ret = false
			}
		case VECTOR:
			// if vecConvert Values are set -> true, else false
			if Length(v) > 0 {
				ret = true
			} else {
				ret = false
			}
		case MATRIX:
			// if mtxConvert Fields are set -> true, else false
			if Length(v) > 0 {
				ret = true
			} else {
				ret = false
			}
		}
		return boolVal{
			ret,
		}
	} // }}}

	// RETURN INTEGER TYPE// {{{
	intConvert = func(v Value) intVal {
		var ret int64 = 0
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return intVal{
			big.NewInt(ret),
		}
	} // }}}

	// RETURN FLOAT TYPE// {{{
	floatConvert = func(v Value) floatVal {
		var ret float64 = 0
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return floatVal{
			big.NewFloat(ret),
		}
	} // }}}

	// RETURN BYTE TYPE// {{{
	byteConvert = func(v Value) byteVal {
		var ret byte = 0
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return byteVal{
			ret,
		}
	} // }}}

	// RETURN BYTES TYPE// {{{
	bytesConvert = func(v Value) bytesVal {
		var ret []byte = []byte{}
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return bytesVal{ret}
	} // }}}

	// RETURN STRING TYPE// {{{
	strConvert = func(v Value) strVal {
		var ret string = ""
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return strVal{ret}
	} // }}}

	// RETURN VECTOR TYPE// {{{
	vecConvert = func(v Value) vecVal {
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return NewVal(VECTOR).(vecVal)
	} // }}}

	// RETURN MATRIX TYPE// {{{
	mtxConvert = func(v Value) mtxVal {
		switch v.Type() {
		case BOOL:
		case INTEGER:
		case FLOAT:
		case BYTE:
		case BYTES:
		case STRING:
		case VECTOR:
		case MATRIX:
		}
		return NewVal(MATRIX).(mtxVal)
	} // }}}
	// }}}

	//// GENERIC TO-TYPE (one ring to rule them all)// {{{
	///
	//  to-type is defined at the value interface level to apply on all
	//  value implementations. It calls the appropriate conversion function
	//  and passes it's receiver and the passed designated type as its
	//  arguments.
	//
	// if the value fails to have a type, a new value instance will be defined
	toType = func(t valueType, v Value) Value {
		switch t {
		case EMPTY:
			return emptyConvert(v)
		case BOOL:
			return boolConvert(v)
		case INTEGER:
			return intConvert(v)
		case FLOAT:
			return floatConvert(v)
		case BYTE:
			return byteConvert(v)
		case BYTES:
			return bytesConvert(v)
		case STRING:
			return strConvert(v)
		case VECTOR:
			return vecConvert(v)
		case MATRIX:
			return mtxConvert(v)
		}
		// define new value and convert to given type
		return NewVal(v).ToType(t)
	} // }}}

	// GENERIC EVAL (second ring to rule them all... wait a minute!)// {{{
	eval = func(v Value) []byte {
		return bytesConvert(v).bytes

	} // }}}
) // }}}

//// MAPPING OF VALUE METHODS TO GENERIC FUNCTIONS// {{{
///
// while data-/ and function-type definitions and generic implementations on
// base of the value interface describe the general implementation, each
// data-type needs to provide its customized implementation of the interface
// method types.
//
/// TYPE FUNCTIONS (one per type)// {{{
//
// One method per type to return the Type of this particular methods receiver value
func (baseVar) Type() valueType  { return EMPTY }
func (boolVal) Type() valueType  { return BOOL }
func (intVal) Type() valueType   { return INTEGER }
func (floatVal) Type() valueType { return FLOAT }
func (byteVal) Type() valueType  { return BYTE }
func (bytesVal) Type() valueType { return BYTES }
func (strVal) Type() valueType   { return STRING }
func (vecVal) Type() valueType   { return VECTOR }
func (mtxVal) Type() valueType   { return MATRIX } // }}}

/// TO-TYPE FUNCTIONS (one per type)// {{{
// map one to-type function per type to the generic implementation
func (v baseVar) ToType(t valueType) Value  { return toType(t, v) }
func (v boolVal) ToType(t valueType) Value  { return toType(t, v) }
func (v intVal) ToType(t valueType) Value   { return toType(t, v) }
func (v floatVal) ToType(t valueType) Value { return toType(t, v) }
func (v byteVal) ToType(t valueType) Value  { return toType(t, v) }
func (v bytesVal) ToType(t valueType) Value { return toType(t, v) }
func (v strVal) ToType(t valueType) Value   { return toType(t, v) }
func (v vecVal) ToType(t valueType) Value   { return toType(t, v) }
func (v mtxVal) ToType(t valueType) Value   { return toType(t, v) } // }}}

/// EVAL FUNCTIONS (one per type)// {{{
// map one eval function per type to the generic implementation
func (v baseVar) Eval() []byte  { return eval(v) }
func (v boolVal) Eval() []byte  { return eval(v) }
func (v intVal) Eval() []byte   { return eval(v) }
func (v floatVal) Eval() []byte { return eval(v) }
func (v byteVal) Eval() []byte  { return eval(v) }
func (v bytesVal) Eval() []byte { return eval(v) }
func (v strVal) Eval() []byte   { return eval(v) }
func (v vecVal) Eval() []byte   { return eval(v) }
func (v mtxVal) Eval() []byte   { return eval(v) } // }}}
// }}}
