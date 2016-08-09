package agiledoc

import (
	//"fmt"
	"github.com/emirpasic/gods/containers"
	"github.com/emirpasic/gods/lists"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/lists/doublylinkedlist"
	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/hashbidimap"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/maps/treebidimap"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/emirpasic/gods/trees"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
	//"os"
)

// CONTAINER INTERFACE (extends the gods container interface)
// interface to conceal god container empty interface values behind the Values
// interface, that provides a type function to inspect the nature of the
// contained value without using reflection.. Since containers themselves and
// there conained values both implement the val interface, all container types
// are fully recursive. The Values returned from and passed to mapped
// containers and sets, also implement the Var interface, KeyValues identitys are
// taken as map keys to map their values on to.

// CONTAINER TYPE
// bitflag to indicate the contained gods containers exact type
type CntType uint16

//go:generate -command stringer -type CntType
const (
	///////////// listCnts
	LIST_ARRAY CntType = 1 << iota
	LIST_SINGLE
	LIST_DOUBLE
	///////////// sets
	SET_HASH
	SET_TREE
	///////////// stacks
	STACK_LINKED
	STACK_ARRAY
	///////////// maps
	MAP_HASH
	MAP_HASHBIDI
	MAP_TREE
	MAP_TREEBIDI
	///////////// trees
	TREE_REDBLACK
	TREE_BINHEAP
	///////////// utils

	// generic container is the nil type of container types, it wont get
	// initialized much, since every container implementation also
	// implements one of the more specific interfaces, that additionaly
	// contain the container interface.
	CONTAINER CntType = 0

	// COMMON CONTAINER TYPES
	// sets of containers that share a more specific interface than
	// gods/containers, like the lists, set, stack, map, or tree interface
	LISTS  = LIST_ARRAY | LIST_SINGLE | LIST_DOUBLE
	SETS   = SET_HASH | SET_TREE
	STACKS = STACK_ARRAY | STACK_LINKED
	MAPS   = MAP_HASH | MAP_HASHBIDI | MAP_TREE | MAP_TREEBIDI
	TREES  = TREE_BINHEAP | TREE_REDBLACK

	// ENUMERABLE/ITERATOR/COMPARATOR MAPPING TYPE
	// sets of containers, which have their contents either mapped by an
	// integer index, or by a key of arbitrary type.
	INDEXED = LISTS | STACKS
	KEYED   = SETS | MAPS | TREES

	// REVERSEABLE ITERATOR INDICATOR
	// reverseable iterators allways also implement the basic iterator of
	// the same mapping type.
	REVERSE = LIST_ARRAY | LIST_DOUBLE | SET_TREE | STACK_ARRAY |
		MAP_TREEBIDI | MAP_TREE | TREE_BINHEAP | TREE_REDBLACK
)

// COMPARATOR INTERFACE
// replaces the empty interface parameters expected by the gods library, with
// parameters that implement the Value interface.  Integer and String
// Comparator wrap the Implementations from gods, the untyped comparator uses
// the Compare operation to compare arbitrary values and return an integer.
type Comparator func(a, b Value) int
type IntComparator func(a, b Value) int
type StringComparator func(a, b Value) int

// wrap comparator exoecting values as parameters, to return one implementing gods utils.comparator.
func wrapComparator(t ValueType, f func(a, b Value) int) (c utils.Comparator) {
	switch t {
	case STRING, BYTES:
		c = func(a, b interface{}) int { return f(NativeToValue(a), NativeToValue(b)) }
	case INTEGER, FLOAT, RATIONAL:
		c = func(a, b interface{}) int {
			return f(NativeToValue(a), NativeToValue(b))
		}
	default:
		c = func(a, b interface{}) int {
			return f(NativeToValue(a), NativeToValue(b))
		}
	}
	return c
}

// SORT INTERFACE
// inplace sort, wraps gods utils.Sort() function. If the Type of all contained values is either integer, or string, the gods integer or string comparators will be instanciated. For all other, or mixed Values, the compare operation is to be used (operations.go).
func Sort(values []Value, comparator Comparator) {
	var t ValueType = 0
	// concatenate all embedded types
	for _, v := range values {
		val := v
		t = t | val.Type()
	}
	utils.Sort(valSlice(values).Interfaces(), wrapComparator(t, comparator))
}

// ITERATOR INTERFACE
type Iterator interface {
	Value() Value
	Begin()
	Next() bool
	First() bool
}

//	REVERSE ITERATOR INTERFACE
type Reverse interface {
	Iterator
	End()
	Preview() bool
	Last() bool
}
type idxIterator struct {
	containers.IteratorWithIndex
}

func (i idxIterator) Value() Value { return NativeToValue(i.Value()) }

type keyIterator struct {
	containers.IteratorWithKey
}

func (i keyIterator) Value() Value { return NativeToValue(i.Value()) }

type reverseIdxIterator struct {
	containers.ReverseIteratorWithIndex
}

func (i reverseIdxIterator) Value() Value { return NativeToValue(i.Value()) }

type reverseKeyIterator struct {
	containers.ReverseIteratorWithKey
}

func (i reverseKeyIterator) Value() Value { return NativeToValue(i.Value()) }

func wrapIter(t CntType, c interface{}) Iterator {
	var ret Iterator
	if t&INDEXED != 0 {
		if t&REVERSE != 0 {
			ret = reverseIdxIterator{c.(containers.ReverseIteratorWithIndex)}
		} else {
			ret = idxIterator{c.(containers.IteratorWithIndex)}
		}
	} else {
		if t&REVERSE != 0 {
			ret = reverseKeyIterator{c.(containers.ReverseIteratorWithKey)}
		} else {
			ret = keyIterator{c.(containers.IteratorWithKey)}
		}
	}
	return ret
}

// ENUMERABLE INTERFACE
//
type Enumerable interface {
	HasKey() bool
}

// ENUMERABLE RETURN TYPE IMPLEMENTATIONS
//
type enumIdx struct {
	containers.EnumerableWithIndex
}

func (e enumIdx) HasKey() bool { return false }

type enumKey struct {
	containers.EnumerableWithKey
}

func (e enumKey) HasKey() bool { return true }

// WRAP_ENUM
// chooses which enumerable interface implementation to instanciate and maps
// the appropriate functions to the instance.
func wrapEnum(t CntType, c interface{}) Enumerable {
	if t&INDEXED != 0 {
		return enumIdx{c.(containers.EnumerableWithIndex)}
	} else {
		return enumKey{c.(containers.EnumerableWithKey)}

	}
}

// CONTAINER INTERFACE
//
type Container interface {
	Values() []Value
	Slice() []interface{}
	Empty() bool
	Size() int
	Clear()
}

// CONTAINER IMPLEMENTATION
type container struct {
	containers.Container
}

func (c container) Slice() []interface{} { return c.Container.Values() }

// CONTAINER COMMON TYPE INTERFACES
type List interface {
	lists.List
}
type Stack interface {
	stacks.Stack
}
type Map interface {
	maps.Map
}
type Set interface {
	sets.Set
}
type Tree interface {
	trees.Tree
}
type EnumerableWithIndex interface {
	containers.EnumerableWithIndex
}
type EnumerableWithKey interface {
	containers.EnumerableWithKey
}

func (c container) Values() []Value { return interfaceSlice(c.Container.Values()).Values() }

func newCollection(t CntType, mapped bool, comp ...Comparator) Collection {
	var v Collection
	switch {
	case t&LISTS != 0:
		v = newListContainer(t)
	case t&MAPS != 0:
		v = newMapContainer(t, comp[0])
	case t&STACKS != 0:
		v = newStackContainer(t)
	case t&SETS != 0:
		if mapped {
			v = newKeymappedSetContainer(comp[0])
		} else {
			v = newIndexedSetContainer(t)
		}
	case t&TREES != 0:
		if mapped {
			v = newKeymappedTreeContainer(t)
		} else {
			v = newIndexedTreeContainer(t)
		}
	}
	return v
}
func newListContainer(t CntType) Collection {
	var c lists.List
	switch {
	case t&LIST_ARRAY != 0:
		c = arraylist.New()
	case t&LIST_SINGLE != 0:
		c = singlylinkedlist.New()
	case t&LIST_DOUBLE != 0:
		c = doublylinkedlist.New()
	}
	return collection{t, c}
}
func newStackContainer(t CntType) Collection {
	var c stacks.Stack
	switch {
	case t&STACK_ARRAY != 0:
		c = arraystack.New()
	case t&STACK_LINKED != 0:
		c = linkedliststack.New()
	}
	return collection{t, c}
}
func newMapContainer(t CntType, comp ...Comparator) Collection {
	var c maps.Map
	switch {
	case t&MAP_HASH != 0:
		c = hashmap.New()
	case t&MAP_HASHBIDI != 0:
		c = hashbidimap.New()
	case t&MAP_TREE != 0:
		c = treemap.NewWith(wrapComparator(STRING, comp[0]))
	case t&MAP_TREEBIDI != 0:
		c = treebidimap.NewWith(wrapComparator(STRING, comp[0]), wrapComparator(STRING, comp[1]))
	}
	return collection{t, c}
}
func newIndexedSetContainer(t CntType, comp ...Comparator) Collection {
	var c sets.Set
	switch {
	case t&SET_HASH != 0:
		c = hashset.New()
	case t&SET_TREE != 0: // KEY OR INDEX
		c = treeset.NewWith(wrapComparator(INTEGER, comp[0]))

	}
	return collection{t, c}
}
func newKeymappedSetContainer(comp Comparator) Collection {
	var c = treeset.NewWith(wrapComparator(STRING, comp))
	return collection{SET_TREE, c}
}
func newIndexedTreeContainer(t CntType, comp ...Comparator) Collection {
	var c trees.Tree
	switch {
	case t&TREE_BINHEAP != 0:
		c = binaryheap.NewWith(wrapComparator(INTEGER, comp[0]))
	case t&TREE_REDBLACK != 0:
		c = redblacktree.NewWith(wrapComparator(INTEGER, comp[0]))

	}
	return collection{t, c}
}
func newKeymappedTreeContainer(t CntType, comp ...Comparator) Collection {
	var c trees.Tree
	switch {
	case t&TREE_BINHEAP != 0:
		c = binaryheap.NewWith(wrapComparator(STRING, comp[0]))
	case t&TREE_REDBLACK != 0:
		c = redblacktree.NewWith(wrapComparator(STRING, comp[0]))

	}
	return collection{t, c}
}

// COLLECTION INTERFACE
//
// collection provides a common Interface, to be implemented by all collection
// types alike. A Collection has a common type, to be recognizeable as either
// list, stack, set, map, or tree. Depending on the common type, the Native
// method yields the appropriate type.
//
// the exact type gives further information about methods to expect being
// present, or if and which type of comparator function needs to be passed to
// access the enumerable methods.
//
// All collections implement the container, iterator, and enumerable methods,
// yielding the appropriate interface instances, depending on their common
// type.
//
// methods provided by the contained type, that are not part of any interface
// will be accessable via additional mathods mapped to the instance in question
// directly.
type Collection interface {
	CommonType() CntType // combined bitflag LISTS, STACKS, SETS, MAPS, TREES
	ExactType() CntType  // bitflag with single bit set
	Container() Container
	Iterator() Iterator
	Enumerable() Enumerable
	native() interface{} // either List, Stack, Set, Map, or Tree
}

// COLLECTION IMPLEMENTATION
type collection struct {
	t CntType
	c interface{}
}

func (c collection) ExactType() CntType { return c.t }
func (c collection) CommonType() CntType {
	var t = c.t
	var r CntType
	switch {
	case t&LISTS != 0:
		r = LISTS
	case t&STACKS != 0:
		r = LISTS
	case t&SETS != 0:
		r = LISTS
	case t&MAPS != 0:
		r = LISTS
	case t&TREES != 0:
		r = LISTS
	}
	return r
}
func (c collection) Container() Container {
	return container{c.c.(containers.Container)}
}
func (c collection) Iterator() Iterator {
	return wrapIter(c.t, c.c)
}
func (c collection) Enumerable() Enumerable {
	return wrapEnum(c.t, c.c)
}
func (c collection) native() interface{} { return c.c }
