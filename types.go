package agiledoc

import (
	"math/big"
)

type Val interface {
	Type() ValType
}

type ValType uint

//go:generate -command stringer -type ValType
const (
	NIL  ValType = 0
	FLAG ValType = 1 << iota
	INTEGER
	FLOAT
	BYTE
	STRING
)

//
type (
	emptyVal struct{} // emptyValue
	flagVal  struct{ *big.Int }
	intVal   struct{ *big.Int }
	floatVal struct{ *big.Float }
	byteVal  []byte
	strVal   string
)

type typeFunc func(Val) ValType
type valFunc func(Val) Val

func (emptyVal) Type() ValType { return NIL }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (floatVal) Type() ValType { return FLOAT }
func (byteVal) Type() ValType  { return BYTE }
func (strVal) Type() ValType   { return STRING }

func (v flagVal) Value() Val  { return v }
func (v intVal) Value() Val   { return v }
func (v floatVal) Value() Val { return v }
func (v byteVal) Value() Val  { return v }
func (v strVal) Value() Val   { return v }

func (emptyVal) Empty() emptyVal { return emptyVal{} }
func (v flagVal) Flag() uint     { return uint(v.Int64()) }
func (v intVal) Integer() int    { return int(v.Int64()) }
func (v floatVal) Flt() float64  { f, _ := v.Float64(); return float64(f) }
func (v byteVal) Byte() []byte   { return v }
func (v strVal) String() string  { return string(v) }

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
		v = floatVal{big.NewFloat(i.(float64))}
	case BYTE:
		v = byteVal(i.([]byte))
	case STRING:
		v = strVal(i.(string))
	}
	return v
}
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
