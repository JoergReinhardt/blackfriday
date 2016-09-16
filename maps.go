package agiledoc

import (
	//"fmt"
	//con "github.com/emirpasic/gods/containers"
	hm "github.com/emirpasic/gods/maps/hashbidimap"
	//ts "github.com/emirpasic/gods/sets/treeset"
	//"math/big"
)

////////////////////////////////////////////////////////////////////////////////////
//// MAP ////
//////////////
func (m BidiMap) Eval() Evaluable  { return Value(m) }
func (m BidiMap) Type() ValueType  { return TABLE }
func (m BidiMap) Size() int        { return m().Size() }
func (m BidiMap) Empty() bool      { return m().Empty() }
func (m BidiMap) Clear() Collected { m().Clear(); return m }
func (m BidiMap) Add(v ...Evaluable) BidiMap {
	var r = m()
	for i, v := range v {
		i, v := i, v
		switch {
		case v.(Evaluable).Type()&PAIR != 0:
			(*r).Put(v.(Pair).Key(), v.(Pair).Value())

		case v.(Evaluable).Type()&RAT != 0:
			(*r).Put(Value(i), Value(v.(rat).Num(), v.(rat).Denom()).(Pair))
		default:
			(*r).Put(Value(i), v.(rat).Denom())
		}
	}
	return func() *hm.Map { return r }
}
func (m BidiMap) AddInterface(v ...interface{}) BidiMap {
	var r = m()
	for i, val := range v {
		i, val := i, val
		r.Put(Value(i), Value(val))
	}
	return func() *hm.Map { return r }
}
func (m BidiMap) Remove(i int) BidiMap {
	var retval = m()
	(*retval).Remove(i)
	return func() *hm.Map { return retval }
}
func (m BidiMap) Interfaces() []interface{} {
	return m().Values()
}

func (m BidiMap) Values() []Evaluable {
	return ValueSlice(m.Interfaces())
}

func (m BidiMap) Serialize() []byte {
	var retval []byte
	var keys = ValueSlice(m().Keys())
	var values = ValueSlice(m().Values())
	for i := len(values); i > 0; i-- {
		i := i
		retval = append(keys[i].Serialize(),
			append([]byte(": "),
				append(values[i].Serialize(),
					[]byte("\n")...,
				)...,
			)...,
		)

	}

	return retval
}

// use serialization as string format base
func (m BidiMap) String() string { return string(m.Serialize()) }
