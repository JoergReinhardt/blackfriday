package agiledoc

import (
	//"fmt"
	con "github.com/emirpasic/gods/containers"
	al "github.com/emirpasic/gods/lists/arraylist"
	hm "github.com/emirpasic/gods/maps/hashbidimap"
	ts "github.com/emirpasic/gods/sets/treeset"
	as "github.com/emirpasic/gods/stacks/arraystack"
	"math/big"
)

//////////////////////// FUNCTIONAL TYPES TO REPRESENT VALUES /////////////////////
type (
	// collections with numeric indices
	BitFlag func() *big.Int
	List    func() *al.List
	Stack   func() *as.Stack
	// collections with symbolic indices
	BidiMap func() *hm.Map
	Set     func() *ts.Set

	NumericTable  func() []List
	SymbolicTable func() []Set

	Matrix List // gonum
)

// lists and sublists of exactly two values length, are assumed to be either
// key/value, or index/value pairs of Pair Type, by the modules Eval function
// on first pass.
//
// All longer slices are flattened by evalCollection and refed into eval
// recursively. .  All conversions to Collected,  get instanciated as list
// type,to profit from the enumerable interface at flattening and conversion.
// COLLECTED IMPLEMENTING METHODS

////////////////////////////////////////////////////////////////////////////////////
//// LIST ////
//////////////
func (l List) Eval() Evaluable  { return Value(l) }
func (l List) Type() ValueType  { return LIST }
func (l List) Size() int        { return l().Size() }
func (l List) Empty() bool      { return l().Empty() }
func (l List) Clear() Collected { l().Clear(); return l }
func (l List) AddInterface(v ...interface{}) List {
	var retval = l()
	(*retval).Add(v...)
	return List(func() *al.List { return retval })
}
func (l List) Add(v ...Evaluable) List {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Add(value)
	}
	return List(func() *al.List { return retval })
}
func (l List) Remove(i int) List {
	var retval = l()
	(*retval).Remove(i)
	return List(func() *al.List { return retval })
}

func (l List) RankedValues() []Pair {
	var retval []Pair
	var fn = func(index int, value interface{}) {
		i := Value(index)
		v := Value(value)
		// pass both values as paired parameter, will trigger eval to
		// produce a key/value tuple type
		retval = append(retval, Value(i, v).(Pair))
	}
	l().Each(fn)
	return retval
}
func (l List) Interfaces() []interface{} {
	return l().Values()
}

func (l List) Values() []Evaluable {
	var retval []Evaluable
	// parameter function to convert slice of interfaces to slice of
	// values once.
	var fn = func(index int, value interface{}) {
		retval = append(retval, Value(l.Interfaces()))
	}
	// retrieve an iterator from collection and call it passing the
	// argument function, to append to the predefined slice
	(con.EnumerableWithIndex)(l()).Each(fn)
	return retval
}

func (l List) Serialize() []byte {
	// allocate return byte slice, so it can be enclosed by the parameter
	// function.
	var retval []byte

	// parameter function to pass on to internal each methode:
	var fn = func(index int, value interface{}) {
		i := Value(index).Serialize()
		v := Value(value).Serialize()

		// format each entry as one line with leading numeric index,
		// followed by a dot and blank character, the Value and a
		// newline character.
		retval = append(
			retval,
			append(
				i,
				append(
					[]byte(".) "),
					append(
						v,
						[]byte("\n")...,
					)...,
				)...,
			)...,
		)
	}
	// call function once per value, to format whole list
	l().Each(fn)
	return retval
}

// use serialization as string format base
func (l List) String() string { return string(l.Serialize()) }
func (l List) Iter() Iterable {
	return IdxIterator(func() *al.Iterator { i := l().Iterator(); return &i })
}
func (l List) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return l() }
	return r
}

// LIST FROM SLICE OF VALUES
var newList = func(v ...Evaluable) List {
	var l = al.New()
	(*l).Add(InterSlice(v)...)
	var fn List = func() *al.List { return l }
	return fn
}

//////////////////////////////////////////////////////////////////////////
// The Type and Value methods can be pre-assigned at the level of distinct
// functional types, representing each dynamic type
// BOOLEAN VALUE (JACOBI)

// The Type and Value methods can be pre-assigned at the level of distinct
// functional types, representing each dynamic type

// wrap flag in a fresh closure and return that.
// TODO: chaeck if this pull's parameters on the stack that evaluation time
func (f BitFlag) Eval() Evaluable { return Value(f) }

// uses byte method of contained big int
func (f BitFlag) Serialize() []byte { return f().Bytes() }

// returns Flag converted to string on base two
func (f BitFlag) String() string { return f().Text(2) }

// returns pure type Flag
func (f BitFlag) Type() ValueType { return FLAG }

func (f BitFlag) Empty() bool {
	if f().Cmp(ZERO.Flag()) > 0 {
		return false
	} else {
		return true
	}
}
func (f BitFlag) Size() int { return f().BitLen() }
func (f BitFlag) Clear() Collected {
	var r *big.Int = f()
	r.Set(ZERO.Flag())
	return BitFlag(func() *big.Int { return r })
}
func (f BitFlag) Values() []Evaluable {
	return ValueSlice(f.Values())
}
func (f BitFlag) Interfaces() []interface{} {
	var v []interface{}
	for _, val := range f().Bits() {
		val := val
		v = append(v, val)
	}
	return v
}

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

// MAP FROM PAIRS OF VALUES
var newMap = func(v ...Pair) BidiMap {
	var r = hm.New()
	for _, v := range v {
		v := v
		(*r).Put(v.Key(), v.Value())
	}
	return func() *hm.Map { return r }
}

//////////////////////////////////////////////////////////////////////////
//
// ITERATOR IMPLEMENTING TYPES (to wrap different iterator implementations)
// the iterator embedded in a arraylist is a struct, of type
// arraylist.Iterator. the list has a method to generate it. Iterators with
// index differ from iterators with key in the expected parameters, not in the
// type of returnvalues they generate. It alters it's state and needs to be
// returned each time.
type IdxIterator func() *al.Iterator

func (l IdxIterator) Index() (Evaluable, Iterable) { return Value(l().Index()), l }
func (l IdxIterator) Value() (Evaluable, Iterable) { return Value(l().Value()), l }
func (l IdxIterator) Next() (bool, Iterable)       { return l().Next(), l }
func (l IdxIterator) First() (bool, Iterable)      { return l().First(), l }
func (l IdxIterator) Begin() Iterable              { l.Begin(); return l }

// reverse iterator interface
func (l IdxIterator) End() Iterable          { l().End(); return l }
func (l IdxIterator) Prev() (bool, Iterable) { return l().Prev(), l }
func (l IdxIterator) Last() (bool, Iterable) { return l().Last(), l }

//// KEY ITERATOR ////
type KeyIterator func() con.IteratorWithKey

func (k KeyIterator) Index() (Evaluable, Iterable) { return Value(k().Key()), k }
func (k KeyIterator) Value() (Evaluable, Iterable) { return Value(k().Value()), k }
func (k KeyIterator) Next() (bool, Iterable)       { return k().Next(), k }
func (k KeyIterator) First() (bool, Iterable)      { return k().First(), k }
func (k KeyIterator) Begin() Iterable              { k().Begin(); return k }

// ENUMERABLE IMPLEMENTING TYPE
// the enumerator is imolemented by the list itself and alters it's State. Two
// types of enumerable interfaces exist, different in parameters and different
// regarding the type of return values one of its methods returns. the possible
// return types of the differing find method are either int index and interface
// value, or both of the value type. the internal interface returns a value of
// the Pair type instead
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
func (e IdxEnumerable) Find(pf func(index Evaluable, value Evaluable) bool) (Pair, Enumerable) {
	i, v := e().Find(
		func(index int, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return Value(i, v).(Pair), e
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
func (e KeyEnumerable) Find(pf func(index Evaluable, value Evaluable) bool) (Pair, Enumerable) {
	i, v := e().Find(
		func(index interface{}, value interface{}) bool {
			return pf(Value(index), Value(value))
		})
	return Value(i, v).(Pair), e
}

//// SLICE ////
// helper type to convert between slices of interfaces and value slices
func InterSlice(i interface{}) []interface{} {
	if s, ok := i.([]interface{}); ok {
		return s
	} else {
		var e = []interface{}{}
		for _, v := range i.([]interface{}) {
			v := v
			e = append(e, v)
		}
		return e
	}
}

func ValueSlice(i interface{}) []Evaluable {
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
