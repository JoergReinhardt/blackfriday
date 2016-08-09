package agiledoc

import (
	"fmt"
	"math/big"
)

// VALUE INTERFACE
//
type Value interface {
	Type() ValueType
	Value() Value
	Eval() []byte
	String() string
	Native() interface{}
	ToType(ValueType) Value
}
type KeyValue interface {
	Value // Type() ValueType, Value() Value
	Key() Value
	Ident() string
}
type Values interface {
	Value
	Container
}

func Eval(v []Value) []byte {
	var bytes = []byte{}
	var vals = v
	for _, v := range vals {
		v := v
		b := v.Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}

type KeyValues interface {
	Key() Value
	Ident() string
	Value() Value
}

// SLICE HELPER TYPES
// these types exist, so that a slice of interfaces, as well as a slice of
// Values implements a type, methods can be assigned to. That Way unwrapped
// slices can allways be converted to those types and provide either the
// Values(), or the Interfaces() method that converts them to the corredponding
// slice type.
type interfaceSlice []interface{}

func (i interfaceSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range i {
		v := v
		val := NativeToValue(v)
		vs = append(vs, val)
	}
	return vs
}

type valSlice []Value

func (v valSlice) Interfaces() []interface{} {
	var is = []interface{}{}
	for _, i := range v {
		is = append(is, i)
	}
	return is
}

type keyValSlice []KeyValue
type byteSlice []byte
type boolSlice []bool

func (b boolSlice) Eval() []byte {

	// bool-slice-integer array, combines up to 8 booleans in a slice to
	// represent byte sized bitflags
	var bsi [8]bool = [8]bool{}

	// byte slice to concatenate all booleans in bitflag of arbitrary size
	var bs = []byte{}

	// split input into byte sized chunks and iterate over each of those chunks
	for o := 0; o < (len(b)/8 + 1); o++ {

		var u uint8 = 0 // allocate a new uint8 for each byte sized chunk

		if o < len(b) { // if bool slice is not yet depleted

			// iterate over each of the eight bits of the current chunk
			for i := 0; i < 8; i++ {
				i := i
				// dereference bool at the current index
				if bsi[i] { // if element is true, set a bit at the current index
					u = 1 << uint(i)
				}
			}
		} // end of chunk. since lenght check failed, another iteration is possibly needed.
		// append the last produced chunk at the byte slice intended to return
		bs = append(bs, u)

	} // either iterate on, or return byte slice
	// depending on the total number of chunks
	return bs
}
func (b boolSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range b {
		v := v
		vs = append(vs, NativeToValue(v))
	}
	return vs
}
func (i byteSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range i {
		v := v
		vs = append(vs, NativeToValue(v))
	}
	return vs
}
func (v byteSlice) Interfaces() []interface{} {
	var is = []interface{}{}
	for _, i := range v {
		is = append(is, i)
	}
	return is
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
	KEYVAL
	LIST
	STACK
	MAP
	SET
	TREE
)

type ValueType uint16

// TYPES
//
type ( // are kept as close to gos native types they are derived from, as possible
	emptyVal struct{} // the empty struct is taken as emptyVal
	flagVal  big.Int
	intVal   big.Int
	ratVal   big.Rat
	boolVal  bool
	byteVal  uint8
	bytesVal []byte
	strVal   string
	keyVal   struct {
		key Value
		val Value
	}
	lstVal struct {
		Collection
		list func() List
	}
	stkVal struct {
		Collection
		stack func() Stack
	}
	mapVal struct {
		Collection
		getmap func() Map
	}
	setVal struct {
		Collection
		set func() Set
	}
	treeVal struct {
		Collection
		tree func() Tree
	}
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
func (keyVal) Type() ValueType   { return KEYVAL }
func (lstVal) Type() ValueType   { return LIST }
func (stkVal) Type() ValueType   { return STACK }
func (mapVal) Type() ValueType   { return MAP }
func (setVal) Type() ValueType   { return SET }
func (treeVal) Type() ValueType  { return TREE }

func (v flagVal) Value() Value  { return v }
func (v intVal) Value() Value   { return v }
func (v ratVal) Value() Value   { return v }
func (v boolVal) Value() Value  { return v }
func (v byteVal) Value() Value  { return v }
func (v bytesVal) Value() Value { return v }
func (v strVal) Value() Value   { return v }
func (v keyVal) Value() Value   { return v }
func (v lstVal) Value() Value   { return v }
func (v stkVal) Value() Value   { return v }
func (v mapVal) Value() Value   { return v }
func (v setVal) Value() Value   { return v }
func (v treeVal) Value() Value  { return v }

func (v flagVal) Native() interface{}  { return v.Uint() }
func (v intVal) Native() interface{}   { return v.Integer() }
func (v ratVal) Native() interface{}   { return v.Float() }
func (v boolVal) Native() interface{}  { return v.Boolean() }
func (v byteVal) Native() interface{}  { return v.Byte() }
func (v bytesVal) Native() interface{} { return v.Eval() }
func (v strVal) Native() interface{}   { return v.String() }
func (v keyVal) Native() interface{}   { return v.val }
func (v lstVal) Native() interface{}   { return v.list() }
func (v stkVal) Native() interface{}   { return v.stack() }
func (v mapVal) Native() interface{}   { return v.getmap() }
func (v setVal) Native() interface{}   { return v.set() }
func (v treeVal) Native() interface{}  { return v.tree() }

// EMPTY TYPE
// the implementation of the empty type, will be instanciateable from thin air,
// used as base value type, that can convert a native value into the
// appropriate Value instance. It´s also used as a return type, whenever a type
// conversion can not be performed, or a variable lookup fails.
func (emptyVal) Type() ValueType         { return NIL }
func (emptyVal) Empty() emptyVal         { return emptyVal(struct{}{}) }
func (e emptyVal) Native() interface{}   { return e.Empty() }
func (e emptyVal) Value() Value          { return emptyVal(struct{}{}) }
func (e emptyVal) Eval() []byte          { return []byte{} }
func (e emptyVal) String() string        { return "" }
func (emptyVal) Set(v interface{}) Value { return NativeToValue(v) }
func (v emptyVal) ToType(t ValueType) Value {
	return emptyVal{}
}

// BOOL TYPE
func (b boolVal) Boolean() bool { return bool(b) }
func (b boolVal) Integer() int {
	if b {
		return 1
	} else {
		return 0
	}
}
func (b boolVal) Eval() []byte { return []byte(b.String()) }
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
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
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
func (v flagVal) Eval() []byte     { return []byte(v.String()) }
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
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(string(v.Eval()))
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
func (v intVal) Float() float64       { return float64(v.BigInt().Int64()) }
func (v intVal) Eval() []byte         { return []byte(v.String()) }
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
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(string(v.Eval()))
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
func (v ratVal) Integer() int64       { return v.BigInt().Int64() }
func (v ratVal) String() string       { r, _ := v.BigRat().Float64(); return fmt.Sprint(r) }
func (v ratVal) Eval() []byte         { return []byte(v.String()) }
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
func (v byteVal) String() string   { return string(v.Eval()) }
func (v byteVal) Eval() []byte     { return []byte{v.Byte()} }
func (v byteVal) Byte() byte       { return byte(v) }
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
func (v bytesVal) Eval() []byte     { return []byte(v) }
func (v bytesVal) String() string   { return string(v) }
func (v bytesVal) Uint() uint64     { return big.NewInt(0).SetBytes(v).Uint64() }
func (v bytesVal) Integer() int64   { return big.NewInt(0).SetBytes(v).Int64() }
func (v bytesVal) Flag() flagVal    { return flagVal(*big.NewInt(0).SetBytes(v.Eval())) }
func (v bytesVal) BigInt() *big.Int { return big.NewInt(0).SetBytes(v.Eval()) }
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
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// STRING TYPE
func (v strVal) Eval() []byte   { return []byte(v) }
func (v strVal) String() string { return string(v) }
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
		val = byteVal(v.Eval()[0])
	case STRING:
		val = v
	default:
		val = emptyVal{}
	}
	return val
}

// KEYVALUE
// a keyvalue is a valu identifyable by a key. It represents values, contained
// in keymapped containers.
func (v keyVal) Key() Value    { return v.key }
func (v keyVal) Ident() string { return string(v.Key().String()) }
func (v keyVal) Eval() []byte {
	return []byte(fmt.Sprintf("%s: %s", v.Key().String(), v.Value().String()))
}
func (v keyVal) String() string { return string(v.Eval()) }

func (v keyVal) ToType(t ValueType) Value {
	var val Value
	//	switch t {
	//	case FLAG:
	//		val = emptyVal{}
	//	case INTEGER:
	//		val = emptyVal{}
	//	case FLOAT:
	//		val = emptyVal{}
	//	case BYTES:
	//		val = bytesVal(v.Eval())
	//	case BYTE:
	//		val = byteVal(v.Eval()[0])
	//	case STRING:
	//		val = strVal(v.String())
	//	default:
	//		val = emptyVal{}
	//	}
	return val
}

// convenience conversion functions

// convenience conversion functions
func (v lstVal) List() List { return v.list() }
func (v lstVal) Eval() []byte {
	var bytes = []byte{}
	var vals = v.Container().Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v lstVal) String() string { return string(v.Eval()) }

func (v lstVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = emptyVal{}
	case INTEGER:
		val = emptyVal{}
	case FLOAT:
		val = emptyVal{}
	case BYTES:
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// STACK VALUE
func (v stkVal) Stack() Stack { return v.stack() }
func (v stkVal) Eval() []byte {
	var bytes = []byte{}
	var vals = v.Container().Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v stkVal) String() string { return string(v.Eval()) }

func (v stkVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = emptyVal{}
	case INTEGER:
		val = emptyVal{}
	case FLOAT:
		val = emptyVal{}
	case BYTES:
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// SET VALUE
func (v setVal) Set() Set { return v.set() }
func (v setVal) Eval() []byte {
	var bytes = []byte{}
	var vals = v.Container().Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v setVal) KeyValues() []KeyValue {
	var retv []KeyValue
	var vals = v.Container().Values()
	for _, val := range vals {
		retv = append(retv, keyVal{val.(keyVal).Key(), val.(keyVal).Value()})
	}
	return retv
}
func (v setVal) String() string { return string(v.Eval()) }

func (v setVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = emptyVal{}
	case INTEGER:
		val = emptyVal{}
	case FLOAT:
		val = emptyVal{}
	case BYTES:
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

// convenience conversion functions
func (v mapVal) Map() Map { return v.getmap() }
func (v mapVal) Eval() []byte {
	var bytes = []byte{}
	var vals = v.Container().Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v mapVal) String() string { return string(v.Eval()) }
func (v mapVal) KeyValues() []KeyValue {
	var retv []KeyValue
	var vals = v.Container().Values()
	for _, val := range vals {
		retv = append(retv, keyVal{val.(keyVal).Key(), val.(keyVal).Value()})
	}
	return retv
}

func (v mapVal) ToType(t ValueType) Value {
	var val Value
	//	switch t {
	//	case FLAG:
	//		val = emptyVal{}
	//	case INTEGER:
	//		val = emptyVal{}
	//	case FLOAT:
	//		val = emptyVal{}
	//	case BYTES:
	//		val = bytesVal(v.Eval())
	//	case BYTE:
	//		val = byteVal(v.Eval()[0])
	//	case STRING:
	//		val = strVal(v.String())
	//	default:
	//		val = emptyVal{}
	//	}
	return val
}

// TREE VALUE
func (v treeVal) Tree() Tree { return v.tree() }
func (v treeVal) Eval() []byte {
	var bytes = []byte{}
	var vals = v.Container().Values()
	for _, v := range vals {
		v := v
		b := NativeToValue(v).Eval()
		bytes = append(bytes, b...)
	}
	return bytes
}
func (v treeVal) String() string { return string(v.Eval()) }
func (v treeVal) KeyValues() []KeyValue {
	var retv []KeyValue
	var vals = v.Container().Values()
	for _, val := range vals {
		retv = append(retv, keyVal{val.(keyVal).Key(), val.(keyVal).Value()})
	}
	return retv
}

func (v treeVal) ToType(t ValueType) Value {
	var val Value
	switch t {
	case FLAG:
		val = emptyVal{}
	case INTEGER:
		val = emptyVal{}
	case FLOAT:
		val = emptyVal{}
	case BYTES:
		val = bytesVal(v.Eval())
	case BYTE:
		val = byteVal(v.Eval()[0])
	case STRING:
		val = strVal(v.String())
	default:
		val = emptyVal{}
	}
	return val
}

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
		col := newCollection(LIST_ARRAY, false)
		v = lstVal{col, func() List { return List(col.native().(List)) }}
	case []KeyValue:
		col := newCollection(MAP_HASHBIDI, true)
		v = mapVal{col, func() Map { return Map(col.native().(Map)) }}
		for _, kv := range i.([]KeyValue) {
			key := kv.Key()
			value := kv.Value()
			v.(mapVal).Map().Put(key, value)
		}
	}
	return v
}
