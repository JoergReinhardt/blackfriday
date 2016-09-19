package agiledoc

import (
	con "github.com/emirpasic/gods/containers"
	"math/big"
)

// lists and sublists of exactly two values length, are assumed to be either
// key/value, or index/value pairs of pair Type, by the modules Eval function
// on first pass.
//
// All longer slices are flattened by evalCollection and refed into eval
// recursively. .  All conversions to Collected,  get instanciated as list
// type,to profit from the enumerable interface at flattening and conversion.
// COLLECTED IMPLEMENTING METHODS

////////////////////////////////////////////////////////////////////////////////////
//////// LISTS /////////
//// COMMON METHODS ///
//////////////////////
////////////////////////////////////////////////////////////////////////////////////
// To prevent code repetition, most methods wrap private methods defined at
// module level, that take implementations of public interfaces defined by the
// package. as their arguments and return values.
func (l ArrayList) Eval() Evaluable  { return evalCollection(l()) }
func (l ArrayList) Size() int        { return collectionSize(l()) }
func (l ArrayList) Empty() bool      { return emptyCollection(l()) }
func (l ArrayList) Clear() Collected { return clearCollection(l()) }

func (l ArrayList) Get(i int) (Evaluable, bool) { return getFromList(l, i) }
func (l ArrayList) Remove(i int) Listed         { return removeFromList(l, i) }
func (l ArrayList) Add(v ...Evaluable) Listed   { return addToList(l, v...) }
func (l ArrayList) AddInterface(v ...interface{}) Listed {
	return addSliceOfInterfacesToList(l, v...)
}
func (l ArrayList) Contains(v ...Evaluable) bool        { return listContains(l, v...) }
func (l ArrayList) Sort(c Compareable) Listed           { return sortList(l, c) }
func (l ArrayList) Swap(idx int, idy int) Listed        { return swapList(l, idx, idy) }
func (l ArrayList) Insert(i int, v ...Evaluable) Listed { return insertList(l, i, v...) }
func (l ArrayList) Values() []Evaluable                 { return collectionValues(l()) }
func (l ArrayList) Interfaces() []interface{}           { return collectionInterfaces(l()) }
func (l ArrayList) Serialize() []byte                   { return serializeList(l) }
func (l ArrayList) String() string                      { return listToString(l) }

//////////////////////////////////////
//// ARRAY LIST CUSTOM METHODS ////
/// methods that are custom to this particular implementation of listed
// Cpt. Obvious:
func (l ArrayList) Type() ValueType { return LIST }

// join yields a slice of pairs with key set to the index of an element while
// value is set to the value holding the index
func (l ArrayList) Join() Collected {
	var retval ArrayList = newList()

	// prepare enumerable parameter func to pass to lists each enumerable
	var fn = func(index int, value interface{}) {
		i := Value(index) // iteration count converted to evaluable integer
		v := Value(value) // value propbably allready a value, otherwise converted
		// pass both values as paired parameter, will trigger eval to
		// produce a key/value tuple type.
		retval.Add(Value(i, v).(pair))
	}

	// pass prepared closure to each, appends all keys and values to the
	// predefined return value
	l().Each(fn)

	// return value got all pairs attendet by each enumerable and is ready
	// to return
	return retval
}

// ENUMERABLE & ITERABLE GENERATORS
//
// the type of iterator depends on the implementation and/or contained value
// types the iter method wraps the underlying containers Iterable() method, to
// zield an iterable, that takes evaluable implementations as arguments, hiding
// the empty interface expected by the embedded container
func (l ArrayList) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}

func (l ArrayList) RevIter() Reverse { rev := l().Iterator(); return IdxRevIterator{&rev} }

// the type of enumerable depends on the implementation and/or contained value types
// the enum method wraps list types in the appropriate enumerable according to
// key types provided by implementation
func (l ArrayList) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return l() }
	return r
}

//////////////////////////////////////////////////////////////////////////
func (l SLList) Eval() Evaluable  { return evalCollection(l()) }
func (l SLList) Size() int        { return collectionSize(l()) }
func (l SLList) Empty() bool      { return emptyCollection(l()) }
func (l SLList) Clear() Collected { return clearCollection(l()) }

func (l SLList) Get(i int) (Evaluable, bool) { return getFromList(l, i) }
func (l SLList) Remove(i int) Listed         { return removeFromList(l, i) }
func (l SLList) Add(v ...Evaluable) Listed   { return addToList(l, v...) }
func (l SLList) AddInterface(v ...interface{}) Listed {
	return addSliceOfInterfacesToList(l, v...)
}
func (l SLList) Contains(v ...Evaluable) bool        { return listContains(l, v...) }
func (l SLList) Sort(c Compareable) Listed           { return sortList(l, c) }
func (l SLList) Swap(idx int, idy int) Listed        { return swapList(l, idx, idy) }
func (l SLList) Insert(i int, v ...Evaluable) Listed { return insertList(l, i, v...) }
func (l SLList) Values() []Evaluable                 { return collectionValues(l()) }
func (l SLList) Interfaces() []interface{}           { return collectionInterfaces(l()) }
func (l SLList) Serialize() []byte                   { return serializeList(l) }
func (l SLList) String() string                      { return listToString(l) }
func (l SLList) Type() ValueType                     { return LIST }
func (l SLList) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
func (l SLList) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return l() }
	return r
}

////////////////////////////////////////////////////////////////////////////////////
func (l DLList) Eval() Evaluable  { return evalCollection(l()) }
func (l DLList) Size() int        { return collectionSize(l()) }
func (l DLList) Empty() bool      { return emptyCollection(l()) }
func (l DLList) Clear() Collected { return clearCollection(l()) }

func (l DLList) Get(i int) (Evaluable, bool) { return getFromList(l, i) }
func (l DLList) Remove(i int) Listed         { return removeFromList(l, i) }
func (l DLList) Add(v ...Evaluable) Listed   { return addToList(l, v...) }
func (l DLList) AddInterface(v ...interface{}) Listed {
	return addSliceOfInterfacesToList(l, v...)
}
func (l DLList) Contains(v ...Evaluable) bool        { return listContains(l, v...) }
func (l DLList) Sort(c Compareable) Listed           { return sortList(l, c) }
func (l DLList) Swap(idx int, idy int) Listed        { return swapList(l, idx, idy) }
func (l DLList) Insert(i int, v ...Evaluable) Listed { return insertList(l, i, v...) }
func (l DLList) Values() []Evaluable                 { return collectionValues(l()) }
func (l DLList) Interfaces() []interface{}           { return collectionInterfaces(l()) }
func (l DLList) Serialize() []byte                   { return serializeList(l) }
func (l DLList) String() string                      { return listToString(l) }
func (l DLList) Type() ValueType                     { return LIST }
func (l DLList) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
func (l DLList) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return l() }
	return r
}

//////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////
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
	return valueSlice(f.Values())
}
func (f BitFlag) Interfaces() []interface{} {
	var v []interface{}
	for _, val := range f().Bits() {
		val := val
		v = append(v, val)
	}
	return v
}
