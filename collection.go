package agiledoc

import (
	//"fmt"
	con "github.com/emirpasic/gods/containers"
	cl "github.com/emirpasic/gods/lists"
	al "github.com/emirpasic/gods/lists/arraylist"
	dl "github.com/emirpasic/gods/lists/doublylinkedlist"
	sl "github.com/emirpasic/gods/lists/singlylinkedlist"
	cm "github.com/emirpasic/gods/maps"
	hbm "github.com/emirpasic/gods/maps/hashbidimap"
	hm "github.com/emirpasic/gods/maps/hashmap"
	tbm "github.com/emirpasic/gods/maps/treebidimap"
	tm "github.com/emirpasic/gods/maps/treemap"
	cs "github.com/emirpasic/gods/sets"
	hs "github.com/emirpasic/gods/sets/hashset"
	ts "github.com/emirpasic/gods/sets/treeset"
	csa "github.com/emirpasic/gods/stacks"
	as "github.com/emirpasic/gods/stacks/arraystack"
	ls "github.com/emirpasic/gods/stacks/linkedliststack"
	//ct "github.com/emirpasic/gods/trees"
	ht "github.com/emirpasic/gods/trees/binaryheap"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	"math/big"
	"sync"
)

//// ALLOCATION POOLS ////
type (
	listPool sync.Pool
	mapPool  sync.Pool
)

var (
	listGen listPool = listPool{}
	mapGen  mapPool  = mapPool{}
)

func newList() ArrayList { return func() *al.List { return listGen.New().(*al.List) } }
func newMap() HashMap    { return func() *hm.Map { return mapGen.New().(*hm.Map) } }

func init() {
	listGen.New = func() interface{} { return new(al.List) }
	mapGen.New = func() interface{} { return new(hm.Map) }
}

//////////////////////// FUNCTIONAL TYPES TO REPRESENT VALUES /////////////////////
type (
	// collections with numeric indices
	// LISTS
	BitFlag   func() *big.Int
	ArrayList func() *al.List
	DLList    func() *dl.List
	SLList    func() *sl.List
	// STACKS
	ArrayStack  func() *as.Stack
	LinkedStack func() *ls.Stack
	// collections with symbolic indices
	// MAPS
	HashMap     func() *hm.Map
	HashBidiMap func() *hbm.Map
	TreeMap     func() *tm.Map
	TreeBidiMap func() *tbm.Map
	// SETS
	TreeSet func() *ts.Set
	HashSet func() *hs.Set
	// TREES
	Heap     func() *ht.Heap
	RedBlack func() *rbt.Tree
)

////  BIT-FLAG METHODS
/// a bit-flag differs signifficantly from all other collection
// implementations. It is expressed as a big int, like the scalar types are used
// to. This type implements a collection of bools, that provides the methods
// needed to implement the stack interface. Apart from that, method are
// provided to convert to a slice of bools, slice of native bools, as well as
// one that returns an iterator over the contained boolean values. There is a
// boolean scalar, that considers the first bit for a value as well.

//// FUNCTIONS COMMON TO ALL COLLECTIONS
func evalCollection(c con.Container) Evaluable           { return Value(c) }
func collectionSize(c con.Container) int                 { return c.(con.Container).Size() }
func emptyCollection(c con.Container) bool               { return c.(con.Container).Empty() }
func collectionValues(c con.Container) []Evaluable       { return valueSlice(c.(con.Container).Values()) }
func collectionInterfaces(c con.Container) []interface{} { return c.(con.Container).Values() }
func clearCollection(c con.Container) Collected          { c.(con.Container).Clear(); return c.(Collected) }
func serializeCollection(c Collected, delims ...[]byte) []byte {
	// serialize collection expects between zero and three byte slices to keep
	// elements of the serialization seperated from one another and to seperate
	// between key, or index and the value for all tuple types,
	//
	// allocate return slice, element End marker, key value delimiter and list delimiter
	var r, elementEndMark, keyValDelim, collectionEndMark []byte

	// unroll and assign passed delimiters at there designated positions
	for i, d := range delims {
		i, d := i, d
		switch i {
		case 0: // most important, since all collections feature elements
			elementEndMark = d
		case 1: // second most important, for all types of tuple nature
			keyValDelim = d
		case 2: // to keep several lists apart (optional)
			collectionEndMark = d

		}
	}

	// prepare parameter function to pass on to internal enumerable method
	// (val integer key and plain interface expected)
	for index, value := range c.Values() {

		// convert passed plain interfaces to evaluables (which they
		// most likely allready are anyway). Serialize each passed
		// value using its types val serialize method (works
		// recursively, in case of nested collections).
		i := Value(index).Serialize()
		v := Value(value).Serialize()

		// format each collected entry, divided by the passed delimiters,
		r = append(
			r, // append to pre-assigned return value…
			append(
				i, // …current elements index (converted to an evaluable)…
				append(
					keyValDelim, // …the passed key, or index and value delimiter…
					append(
						v,                 // …the value…
						elementEndMark..., // … and the end of element marker
					)...,
				)...,
			)...,
		)
	}
	// append end of collection marker, if present
	if len(delims) >= 3 {
		r = append(r, collectionEndMark...)
	}
	return r
}

//// FUNCTIONS COMMON TO ALL LISTS
func listToString(l Listed) string                                 { return string(serializeCollection(l, []byte("\n"))) }
func serializeList(l Listed) []byte                                { return serializeCollection(l, []byte("\n")) }
func getFromList(l Listed, i int) (Evaluable, bool)                { v, ok := l.(cl.List).Get(i); return Value(v), ok }
func removeFromList(l Listed, i int) Listed                        { l.(cl.List).Remove(i); return l }
func addToList(l Listed, v ...Evaluable) Listed                    { l.(cl.List).Add(interfaceSlice(v)...); return l }
func addSliceOfInterfacesToList(l Listed, v ...interface{}) Listed { l.(cl.List).Add(v...); return l }
func listContains(l Listed, v ...Evaluable) bool                   { return l.(cl.List).Contains(interfaceSlice(v)...) }
func sortList(l Listed, c Compareable) Listed                      { l.(cl.List).Sort(c.InterfaceComparator()); return l }
func swapList(l Listed, idx int, idy int) Listed                   { l.(cl.List).Swap(idx, idy); return l }
func insertList(l Listed, i int, v ...Evaluable) Listed {
	l.(cl.List).Insert(i, interfaceSlice(v)...)
	return l
}

// LIST GENERATORS
func newArrayList() (r ArrayList) {
	l := al.New()
	r = func() *al.List { return l }
	return r
}
func newDLLList() (r DLList) {
	l := dl.New()
	r = func() *dl.List { return l }
	return r
}
func newSLLList() (r SLList) {
	l := sl.New()
	r = func() *sl.List { return l }
	return r
}

//// FUNCTIONS COMMON TO All STACKS
func pushToStack(s Stacked, v Evaluable) Stacked { s.Push(v); return s }
func popFromStack(s Stacked) (Evaluable, bool, Stacked) {
	v, ok := s.(csa.Stack).Pop()
	return Value(v), ok, s
}
func peekOnStack(s Stacked) (Evaluable, bool) {
	v, ok := s.(csa.Stack).Peek()
	return Value(v), ok
}
func newArraystack() (r ArrayStack) {
	s := as.New()
	r = func() *as.Stack { return s }
	return r
}
func newLinkedStack() (r LinkedStack) {
	s := ls.New()
	r = func() *ls.Stack { return s }
	return r
}

//// FUNCTIONS COMMON TO All MAPS
func Add(m Mapped, k Evaluable, v Evaluable) Mapped      { return putToMap(m, k, v) }
func putToMap(m Mapped, k Evaluable, v Evaluable) Mapped { m.(cm.Map).Put(m, v); return m }
func getFromMap(m Mapped, v Evaluable) (Evaluable, bool) {
	val, ok := m.(cm.Map).Get(v)
	return Value(val), ok
}
func removeFromMap(m Mapped, v Evaluable) Mapped { m.(cm.Map).Remove(v); return m }
func keysOfMap(m Mapped) []Evaluable             { return valueSlice(m.(cm.Map).Keys()) }
func valuesFromMap(m Mapped) []Evaluable         { return valueSlice(m.(con.Container).Values()) }

func serializeMap(m Mapped) []byte {
	var retval []byte
	var keys = valueSlice(m.Keys())
	var values = valueSlice(m.Values())
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
func interfacesFromMap(m Mapped) []interface{} { return m.(cm.Map).Values() }
func mapToString(m Mapped) string              { return string(m.Serialize()) }

//// FUNCTIONS COMMON TO All BIDIRECTIONAL MAPS
func getKeyFromMap(m Mapped, v Evaluable) (Evaluable, bool) {
	r, ok := m.(cm.BidiMap).GetKey(v)
	return Value(r), ok
}
func newHashMap() (r HashMap) {
	m := hm.New()
	r = func() *hm.Map { return m }
	return r
}
func newHashBidiMap() (r HashBidiMap) {
	m := hbm.New()
	r = func() *hbm.Map { return m }
	return r
}
func newTreeMapNumeric() (r TreeMap) {
	m := tm.NewWithIntComparator()
	r = func() *tm.Map { return m }
	return r
}
func newTreeMapSymbolic() (r TreeMap) {
	m := tm.NewWithStringComparator()
	r = func() *tm.Map { return m }
	return r
}
func newTreeBidiMapNumeric() (r TreeBidiMap) {
	m := tbm.NewWithIntComparators()
	r = func() *tbm.Map { return m }
	return r
}
func newTreeBidiMapSymbolic() (r TreeBidiMap) {
	m := tbm.NewWithStringComparators()
	r = func() *tbm.Map { return m }
	return r
}

//// FUNCTIONS COMMON TO All SETS OF UNIQUE ELEMENTS
func removeFromSet(u DeDublicated, i int) DeDublicated { u.(cs.Set).Remove(i); return u }
func addToSet(u DeDublicated, v ...Evaluable) DeDublicated {
	u.(cs.Set).Add(interfaceSlice(v)...)
	return u
}
func setContains(u DeDublicated, v ...Evaluable) bool {
	return u.(cs.Set).Contains(interfaceSlice(v)...)
}
func interfacesFromSet(u DeDublicated) []interface{} { return interfaceSlice(u.Values()) }
func newHashSet(v ...Evaluable) (r HashSet) {
	m := hs.New()
	r = func() *hs.Set { return m }
	return r
}
func newTreeSetNumeric() (r TreeSet) {
	m := ts.NewWithIntComparator()
	r = func() *ts.Set { return m }
	return r
}
func newTreeSetSymbolic() (r TreeSet) {
	m := ts.NewWithStringComparator()
	r = func() *ts.Set { return m }
	return r
}

/////////////// TREE /////////////////
func newHeap() (r Heap) {
	m := ht.NewWithIntComparator()
	r = func() *ht.Heap { return m }
	return r
}
func newRedBlack() (r RedBlack) {
	m := rbt.NewWithStringComparator()
	r = func() *rbt.Tree { return m }
	return r
}

////////////////////////////////////////////////
//// COMPARARATOR IMPLEMENTED BY A FUNCTION TYPE
type Compareable func(a, b Evaluable) int

/// method to cast evaluable parameters as interfaces with empty method set, as
// expected by the underlying collection implementation by a comparator
// implementation
func (c Compareable) InterfaceComparator() utils.Comparator {
	return func(a, b interface{}) int { return c(Value(a), Value(b)) }
}

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
func (l IdxRevIterator) End()             { l.ReverseIteratorWithIndex.End() }
func (l IdxRevIterator) Prev() bool       { return l.ReverseIteratorWithIndex.Prev() }
func (l IdxRevIterator) Last() bool       { return l.ReverseIteratorWithIndex.Last() }
func (l IdxRevIterator) Index() Integer   { return Value(l.ReverseIteratorWithIndex.Index()).(Integer) }
func (l IdxRevIterator) Value() Evaluable { return Value(l.ReverseIteratorWithIndex.Value()) }

type KeyRevIterator struct {
	con.ReverseIteratorWithKey
}

// reverse iterator interface (works for indexed as well as key mapped iterables)
func (l KeyRevIterator) End()             { l.ReverseIteratorWithKey.End() }
func (l KeyRevIterator) Prev() bool       { return l.ReverseIteratorWithKey.Prev() }
func (l KeyRevIterator) Last() bool       { return l.ReverseIteratorWithKey.Last() }
func (l KeyRevIterator) Key() Evaluable   { return Value(l.ReverseIteratorWithKey.Key()) }
func (l KeyRevIterator) Value() Evaluable { return Value(l.ReverseIteratorWithKey.Value()) }

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
func newOrderedList(v ...Evaluable) ArrayList {
	var l = al.New()
	(*l).Add(interfaceSlice(v)...)
	return func() *al.List { return l }
}

// MAP FROM PAIRS OF VALUES
func unorderedBidiMapFromPairs(v ...pair) HashBidiMap {
	var r = hbm.New()
	for _, v := range v {
		k := v.Key()
		v := v.Value()
		switch {
		case k.Type()&SYMBOLIC != 0:
			(*r).Put(k.Serialize(), v)
		case k.Type()&NATURAL != 0:
			(*r).Put(k.(val).Int(), v)
		case k.Type()&REAL != 0:
			// both parts of ratio are stored in the value field,
			// numerator is taken as the index
			(*r).Put(val(k.(ratio).Num()).Int(), v)
		}
	}
	return func() *hbm.Map { return r }
}
