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
func (m UnorderedBidiMap) Eval() Evaluable  { return Value(m) }
func (m UnorderedBidiMap) Type() ValueType  { return TABLE }
func (m UnorderedBidiMap) Size() int        { return m().Size() }
func (m UnorderedBidiMap) Empty() bool      { return m().Empty() }
func (m UnorderedBidiMap) Clear() Collected { m().Clear(); return m }
func (m UnorderedBidiMap) Add(v ...Evaluable) UnorderedBidiMap {
	var r = m()
	for i, v := range v {
		i, v := i, v
		switch {
		case v.(Evaluable).Type()&PAIR != 0:
			(*r).Put(v.(pair).Key(), v.(pair).Value())

		case v.(Evaluable).Type()&REAL != 0:
			(*r).Put(Value(i), Value(v.(rat).Num(), v.(rat).Denom()).(pair))
		default:
			(*r).Put(Value(i), v.(rat).Denom())
		}
	}
	return func() *hm.Map { return r }
}
func (m UnorderedBidiMap) AddInterface(v ...interface{}) UnorderedBidiMap {
	var r = m()
	for i, val := range v {
		i, val := i, val
		r.Put(Value(i), Value(val))
	}
	return func() *hm.Map { return r }
}
func (m UnorderedBidiMap) Remove(i int) UnorderedBidiMap {
	var retval = m()
	(*retval).Remove(i)
	return func() *hm.Map { return retval }
}
func (m UnorderedBidiMap) Interfaces() []interface{} {
	return m().Values()
}

func (m UnorderedBidiMap) Values() []Evaluable {
	return valueSlice(m.Interfaces())
}

func (m UnorderedBidiMap) Serialize() []byte {
	var retval []byte
	var keys = valueSlice(m().Keys())
	var values = valueSlice(m().Values())
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
func (m UnorderedBidiMap) String() string { return string(m.Serialize()) }
