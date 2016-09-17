package agiledoc

import (
	//"fmt"
	con "github.com/emirpasic/gods/containers"
	al "github.com/emirpasic/gods/lists/arraylist"
	dl "github.com/emirpasic/gods/lists/doublylinkedlist"
	hm "github.com/emirpasic/gods/maps/hashbidimap"
	ts "github.com/emirpasic/gods/sets/treeset"
	as "github.com/emirpasic/gods/stacks/arraystack"
	ls "github.com/emirpasic/gods/stacks/linkedliststack"
	"math/big"
)

//////////////////////// FUNCTIONAL TYPES TO REPRESENT VALUES /////////////////////
type (
	// collections with numeric indices
	BitFlag        func() *big.Int
	OrderedList    func() *al.List
	LinkedList     func() *dl.List
	UnorderedStack func() *as.Stack
	IterableStack  func() *ls.Stack
	// collections with symbolic indices
	UnorderedBidiMap func() *hm.Map
	TreeSet          func() *ts.Set

	NumericTable  func() []OrderedList
	SymbolicTable func() []TreeSet

	Matrix OrderedList // gonum
)

//////////////////////////////////////////////////////////////////////////
//
// ITERATOR IMPLEMENTING TYPES (to wrap different iterator implementations)
// the iterator embedded in a arraylist is a struct, of type
// arraylist.Iterator. the list has a method to generate it. Iterators with
// index differ from iterators with key in the expected parameters, not in the
// type of returnvalues they generate. It alters it's state and needs to be
// returned each time.
//// IDX ITERATOR ////
type IdxIterator struct {
	con.IteratorWithIndex
}

func (l IdxIterator) Index() Integer   { return Value((l.IteratorWithIndex).Index()).(Integer) }
func (l IdxIterator) Value() Evaluable { return Value((l.IteratorWithIndex).Value()) }
func (l IdxIterator) Next() bool       { return (l.IteratorWithIndex).Next() }
func (l IdxIterator) First() bool      { return (l.IteratorWithIndex).First() }
func (l IdxIterator) Begin()           { (l.IteratorWithIndex).Begin() }

//// KEY ITERATOR ////
type KeyIterator struct {
	con.IteratorWithKey
	Integer
}

func (k KeyIterator) Index() Integer   { return k.Integer }
func (k KeyIterator) Value() Evaluable { return Value(k.IteratorWithKey.Value()) }
func (k KeyIterator) Next() bool {
	k.Integer = Value(k.Integer().Add(k.Integer(), Value(1).(Integer)())).(Integer)
	return k.IteratorWithKey.Next()
}
func (k KeyIterator) First() bool { return k.IteratorWithKey.First() }
func (k KeyIterator) Begin()      { k.IteratorWithKey.Begin() }

type IdxRevIterator struct {
	con.ReverseIteratorWithIndex
}

// reverse iterator interface (works for indexed as well as key mapped iterables)
func (l IdxRevIterator) End()       { l.ReverseIteratorWithIndex.End() }
func (l IdxRevIterator) Prev() bool { return l.ReverseIteratorWithIndex.Prev() }
func (l IdxRevIterator) Last() bool { return l.ReverseIteratorWithIndex.Last() }

type KeyRevIterator struct {
	con.ReverseIteratorWithKey
}

// reverse iterator interface (works for indexed as well as key mapped iterables)
func (l KeyRevIterator) End()       { l.ReverseIteratorWithKey.End() }
func (l KeyRevIterator) Prev() bool { return l.ReverseIteratorWithKey.Prev() }
func (l KeyRevIterator) Last() bool { return l.ReverseIteratorWithKey.Last() }

// ENUMERABLE IMPLEMENTING TYPE
// the enumerator is imolemented by the list itself and alters it's State. Two
// types of enumerable interfaces exist, different in parameters and different
// regarding the type of return values one of its methods returns. the possible
// return types of the differing find method are either int index and interface
// value, or both of the value type. the internal interface returns a value of
// the Pair type instead
//// IDX ENUMERABLE ////
type IdxEnumerable func() con.EnumerableWithIndex

func (e IdxEnumerable) Each(pf func(index Evaluable, value Evaluable)) Enumerable {
	e().Each(
		func(index int, value interface{}) {
			pf(Value(index), Value(value)) // each does not return a boolean
		})
	return IdxEnumerable(e)
}

func (e IdxEnumerable) Any(pf func(index Evaluable, value Evaluable) bool) (bool, Enumerable) {
	ok := e().Any(
		func(index int, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return ok, e
}
func (e IdxEnumerable) All(pf func(index Evaluable, value Evaluable) bool) (bool, Enumerable) {
	ok := e().All(
		func(index int, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return ok, e
}
func (e IdxEnumerable) Find(pf func(index Evaluable, value Evaluable) bool) (pair, Enumerable) {
	i, v := e().Find(
		func(index int, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return Value(i, v).(pair), e
}

//// KEY ENUMERABLE ////
type KeyEnumerable func() con.EnumerableWithKey

func (e KeyEnumerable) Each(pf func(index Evaluable, value Evaluable)) Enumerable {
	e().Each(
		func(index interface{}, value interface{}) {
			pf(Value(index), Value(value)) // each does not return a boolean
		})
	return e
}

func (e KeyEnumerable) Any(pf func(index Evaluable, value Evaluable) bool) (bool, Enumerable) {
	ok := e().Any(
		func(index interface{}, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return ok, e
}
func (e KeyEnumerable) All(pf func(index Evaluable, value Evaluable) bool) (bool, Enumerable) {
	ok := e().All(
		func(index interface{}, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return ok, e
}
func (e KeyEnumerable) Find(pf func(index Evaluable, value Evaluable) bool) (pair, Enumerable) {
	i, v := e().Find(
		func(index interface{}, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return Value(i, v).(pair), e
}

type EnumParameter func(index, value Evaluable) bool

//// SLICE ////
// helper type to convert between slices of interfaces and slices of value
func interfaceSlice(i interface{}) []interface{} {
	if s, ok := i.([]interface{}); ok {
		return s
	} else {
		var e = []interface{}{}
		for _, v := range i.([]Evaluable) {
			v := v
			e = append(e, v)
		}
		return e
	}
}

func valueSlice(i interface{}) []Evaluable {
	if v, ok := i.([]Evaluable); ok {
		return v
	} else {
		var e = []Evaluable{}
		for _, val := range i.([]interface{}) {
			val := Value(val)
			e = append(e, val)
		}
		return e
	}
}

// LIST FROM SLICE OF VALUES
func newOrderedList(v ...Evaluable) OrderedList {
	var l = al.New()
	(*l).Add(interfaceSlice(v)...)
	return func() *al.List { return l }
}

// MAP FROM PAIRS OF VALUES
func unorderedBidiMapFromPairs(v ...pair) UnorderedBidiMap {
	var r = hm.New()
	for _, v := range v {
		k := v.Key()
		v := v.Value()
		switch {
		case k.Type()&SYMBOLIC != 0:
			(*r).Put(k.Serialize(), v)
		case k.Type()&NATURAL != 0:
			(*r).Put(k.(val).int(), v)
		case k.Type()&REAL != 0:
			(*r).Put(val(k.(rat).Num()).int(), v)
		}
	}
	return func() *hm.Map { return r }
}
