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
	NIL ValType = 0
	FLG ValType = 1 << iota
	INT
	FLT
	BYT
	STR
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

//
type typeFunc func(Val) ValType

func (emptyVal) Type() ValType { return NIL }
func (flagVal) Type() ValType  { return FLG }
func (intVal) Type() ValType   { return INT }
func (floatVal) Type() ValType { return FLT }
func (byteVal) Type() ValType  { return BYT }
func (strVal) Type() ValType   { return STR }

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
	case FLG:
		v = flagVal{big.NewInt(int64(i.(int)))}
	case INT:
		v = intVal{big.NewInt(int64(i.(int)))}
	case FLT:
		v = floatVal{big.NewFloat(i.(float64))}
	case BYT:
		v = byteVal(i.([]byte))
	case STR:
		v = strVal(i.(string))
	}
	return v
}
func NewVal(i interface{}) Val {
	var v Val
	switch i.(type) {
	case int, int8, int16, int32, int64, *big.Int:
		v = NewTypedVal(INT, i)
	case float32, float64, *big.Float:
		v = NewTypedVal(FLT, i)
	case []byte:
		v = NewTypedVal(BYT, i)
	case string:
		v = NewTypedVal(STR, i)
	}
	return v
}
