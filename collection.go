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
	BitFlag func() *big.Int  // implements Ranked interface
	List    func() *al.List  // implements Ranked interface
	Table   func() *al.List  // implements Tabular interface
	Map     func() *hm.Map   // imolements Mapped interface
	Stack   func() *as.Stack // imolements Mapped interface
	Set     func() *ts.Set   // imolements Mapped interface

	Iterator   func() Iterable
	Enumerator func() Enumerable
)

// lists and sublists of exactly two values length, are assumed to be either
// key/value, or index/value pairs of Pair Type, by the modules Eval function
// on first pass.
//
// All longer slices are flattened by evalCollection and refed into eval
// recursively. .  All conversions to Collected,  get instanciated as list
// type,to profit from the enumerable interface at flattening and conversion.
func evalCollection(i ...interface{}) Value {
	var v Value      // return value
	var l = al.New() // allocate intermediate list
	l.Add(i...)      // unroll outer layer of possibly nested interfaces

	//-ALL NUMERIC / NO PAIRS
	//  - MIXED NUMERIC
	//  → convert to INTEGER/<most significant>
	//  - ALL SAME NUMERIC
	//    - ALL BOOL → FLAG (index = bit position)
	//    - ALL UINT → FLAG (index = bit position)
	//    - ALL INT  → INTEGER/INTEGER
	// MIXED NUMERIC / PAIR
	//    → convert to <most significant key>/<ms value>, use string
	//	representation of iteration count to substitute for values
	//	without key.
	// ALL PAIR
	//  - MIXED PAIR
	//	→ convert to <ms key>/<ms value>, use string
	//	  representation of iteration count to substitute for values
	//	  without key.
	//  - ALL SAME PAIR
	//	NUMERIC KEY
	//    -	  FLOAT → INTEGER/RAT al.List
	//    -	  RAT	→ INTEGER/RAT al.List
	//    -	  PAIR	→ PAIR(INTEGER/VALUE) al.List
	//	MAPPED KEY
	//    -	  FLOAT → VALUE/RAT al.List
	//    -	  RAT   → VALUE/RAT al.List
	//    -	  PAIR  → PAIR(VALUE/VALUE) al.List
	// MIXED PAIR/COLLECTEDVALUE
	// MIXED COLLECTED
	//  - MIXED COLLECTED
	//  - ALL SAME COLLECTED

	return v
}

// The Type and Value methods can be pre-assigned at the level of distinct
// functional types, representing each dynamic type
// BOOLEAN VALUE (JACOBI)
type Flag Int

// The Type and Value methods can be pre-assigned at the level of distinct
// functional types, representing each dynamic type
func (f BitFlag) Eval() Value       { return Eval(f) }
func (f BitFlag) Serialize() []byte { return Int(f)().Bytes() }
func (f BitFlag) String() string    { return f().Text(2) }
func (f BitFlag) Type() ValueType   { return FLAG }

////////////////////////////////////////////////////////////////////////////////////
func (l List) Eval() Value     { return Eval(l) }
func (l List) Type() ValueType { return LIST }
func (l List) AddInterface(v ...interface{}) Ranked {
	var retval = l()
	(*retval).Add(v...)
	var retfn List = func() *al.List { return retval }
	return retfn
}
func (l List) Add(v ...Value) Ranked {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Add(value)
	}
	var retfn List = func() *al.List { return retval }
	return retfn
}
func (l List) Remove(i int) Ranked {
	var retval = l()
	(*retval).Remove(i)
	var retfn List = func() *al.List { return retval }
	return retfn
}

// LISTED IMPLEMENTING METHODS
func (l List) Empty() bool { return l().Empty() }
func (l List) Size() int   { return l().Size() }
func (l List) Clear() Ranked {
	retVal := l()
	(*retVal).Clear()
	var retFn List = func() *al.List { return retVal }
	return retFn
}

// EMPTY INTERFACE VALUE SLICE
func (l List) Values() []Value {
	var retval []Value
	// parameter function to convert slice of interfaces to slice of
	// values.
	// once.
	var fn = func(index int, value interface{}) {
		retval = append(retval, Eval(value))
	}
	l().Each(fn)
	return retval
}
func (l List) Interfaces() []interface{} {
	// call values methode of embedded container type
	return l().Values()
}

// SLICE OF RANKED VALUES
func (l List) RankedValues() []Pair {
	var retval []Pair
	var fn = func(index int, value interface{}) {
		i := Eval(index)
		v := Eval(value)
		// pass both values as paired parameter, will trigger eval to
		// produce a key/value tuple type
		retval = append(retval, Eval(i, v).(Pair))
	}
	l().Each(fn)
	return retval
}
func (l List) Serialize() []byte {
	// allocate return byte slice, so it can be enclosed by the parameter
	// function.
	var retval []byte

	// parameter function to pass on to internal each methode:
	var fn = func(index int, value interface{}) {
		i := Eval(index).Serialize()
		v := Eval(value).Serialize()

		// format each entry as one line with leading numeric index,
		// followed by a dot and blank character, the Value and a
		// newline character.
		retval = append(
			retval,
			append(
				i,
				append(
					[]byte(". "),
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
func (l List) Iterator() Iterable {
	return intIterator(func() *al.Iterator { i := l().Iterator(); return &i })
}

// ITERATOR IMPLEMENTING METHODS
// the iterator embedded in a arraylist is a struct, of type
// arraylist.Iterator. the list has a method to generate it. Iterators with
// index differ from iterators with key in the expected parameters, not in the
// type of returnvalues they generate. It alters it's state and needs to be
// returned each time.
type intIterator func() *al.Iterator

func (l intIterator) Index() (int, Iterable)   { return l().Index(), l }
func (l intIterator) Value() (Value, Iterable) { return Eval(l().Value()), l }
func (l intIterator) Next() (bool, Iterable)   { return l().Next(), l }
func (l intIterator) First() (bool, Iterable)  { return l().First(), l }
func (l intIterator) Begin() Iterable          { l.Begin(); return l }

// reverse iterator interface
func (l intIterator) End() Iterable          { l().End(); return l }
func (l intIterator) Prev() (bool, Iterable) { return l().Prev(), l }
func (l intIterator) Last() (bool, Iterable) { return l().Last(), l }

//////////////////////// ENUM /////////////////////////////////////////////////////
// ENUMERABLE IMPLEMENTING METHODS
// the enumerator is imolemented by the list itself and alters it's State. Two
// types of enumerable interfaces exist, different in parameters and different
// regarding return values of one of its methods. the return types of the
// differing find method are either int index and interface value, or both of
// the value type. the internal interface returns a value of the Pair type
// instead
type intEnumerable func() con.EnumerableWithIndex

func (e intEnumerable) Each(pf func(index Value, value Value)) Enumerable {
	e().Each(
		func(index int, value interface{}) {
			pf(Eval(index), Eval(value)) // each does not return a boolean
		})
	return intEnumerable(e)
}

func (e intEnumerable) Any(pf func(value Value, index Value) bool) (Enumerable, bool) {
	// calls each method of contained container, with demanded parameter types
	ok := e().Any(
		func(index int, value interface{}) bool {
			return pf(Eval(index), Eval(value))
		})
	return e, ok
}
func (e intEnumerable) All(pf func(value Value, index Value) bool) (Enumerable, bool) {
	ok := e().All(
		func(index int, value interface{}) bool {
			return pf(Eval(index), Eval(value))
		})
	return e, ok
}
func (e intEnumerable) Find(pf func(index Value, value Value) bool) (Enumerable, Pair) {
	i, v := e().Find(
		func(index int, value interface{}) bool {
			return pf(Eval(index), Eval(value))
		})
	return e, Eval(i, v).(Pair)
}

type Slice []Value

func (s Slice) Interfaces() (i []interface{}) {
	for _, val := range s {
		i = append(i, val)
	}
	return s.Interfaces()
}
func (s Slice) Values() (i []Value) {
	for _, val := range s {
		i = append(i, val.(Value))
	}
	return i
}

///////////////////////////////////////////////////////////////////////////////////////////
// LIST FROM NATIVE VALUES
var EvalList = func(v ...Value) Ranked {
	var l = al.New()
	(*l).Add(Slice(v).Interfaces()...)
	var fn List = func() *al.List { return l }
	return fn
}
