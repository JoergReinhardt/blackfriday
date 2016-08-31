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
	INDEXED   = LISTS | STACKS | TREE_BINHEAP
	KEYMAPPED = MAPS | SET_TREE | TREE_REDBLACK

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

func wrapIterator(t CntType, c interface{}) Iterator {
	var i Iterator

	switch t {
	case LIST_ARRAY:
		iter := c.(*arraylist.List).Iterator()
		i = wrapIteratorDetailed(t, &iter)
	case LIST_SINGLE:
		iter := c.(*singlylinkedlist.List).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case LIST_DOUBLE:
		iter := c.(*doublylinkedlist.List).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case STACK_ARRAY:
		iter := c.(*arraystack.Stack).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case MAP_TREE:
		iter := c.(*treemap.Map).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case MAP_TREEBIDI:
		iter := c.(*treebidimap.Map).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case TREE_BINHEAP:
		iter := c.(*binaryheap.Heap).Iterator()
		i = wrapIteratorDetailed(t, iter)
	case TREE_REDBLACK:
		iter := c.(*redblacktree.Tree).Iterator()
		i = wrapIteratorDetailed(t, iter)
	}

	return i
}
func wrapIteratorDetailed(t CntType, i interface{}) Iterator {
	var r Iterator
	if t&INDEXED != 0 {
		if t&REVERSE != 0 {
			iter := reverseIdxIterator{i.(containers.ReverseIteratorWithIndex)}
			r = &iter
		} else {
			iter := idxIterator{i.(containers.IteratorWithIndex)}
			r = &iter
		}
	} else {
		if t&REVERSE != 0 {
			iter := reverseKeyIterator{i.(containers.ReverseIteratorWithKey)}
			r = &iter
		} else {
			iter := keyIterator{i.(containers.IteratorWithKey)}
			r = &iter
		}
	}
	return r
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
func wrapEnumerable(t CntType, c interface{}) Enumerable {
	var e Enumerable
	if t&INDEXED != 0 {
		enum := enumIdx{c.(containers.EnumerableWithIndex)}
		e = &enum
	} else {
		enum := enumKey{c.(containers.EnumerableWithKey)}
		e = &enum

	}
	return e
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

func (c container) Values() []Value      { return interfaceSlice(c.Container.Values()).Values() }
func (c container) Slice() []interface{} { return c.Container.Values() }
func wrapContainer(t CntType, c interface{}) Container {
	cnt := container{c.(containers.Container)}
	return &cnt
}

// SUB CONTAINER TYPE INTERFACES
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

//func wrapCollectionMethods()
// wraps the wrapper functions, that wrap the original iterator, enumerable and
// container interfaces.
func wrapCollectionMethods(t CntType, f func() interface{}) (func() Container, func() Iterator, func() Enumerable) {
	return func() Container { return wrapContainer(t, f()) },
		func() Iterator { return wrapIterator(t, f()) },
		func() Enumerable { return wrapEnumerable(t, f()) }

}
func newCollection(t CntType, keymapped bool, comp ...Comparator) Collection {

	// allocate empty collection struct
	var c = collection{}

	// set collection type
	c.t = t

	// instanciate native container instance and closure to return
	// reference to the concealed container.
	switch {
	case t&LIST_ARRAY != 0:

		// instanciate concealed gods container, to implement the interface
		col := *arraylist.New()

		// create a function to return a reference to the original struct
		c.ref = func() interface{} { return &col }

	case t&LIST_SINGLE != 0:
		col := *singlylinkedlist.New()
		c.ref = func() interface{} { return &col }
	case t&LIST_DOUBLE != 0:
		// instanciate concealed gods container, to implement the interface
		col := *doublylinkedlist.New()
		// create a function to return a reference to the original struct
		c.ref = func() interface{} { return &col }
	case t&STACK_ARRAY != 0:
		col := *arraystack.New()
		c.ref = func() interface{} { return &col }
	case t&STACK_LINKED != 0:
		col := *linkedliststack.New()
		c.ref = func() interface{} { return &col }
	case t&MAP_HASH != 0:
		col := *hashmap.New()
		c.ref = func() interface{} { return &col }
	case t&MAP_HASHBIDI != 0:
		col := *hashbidimap.New()
		c.ref = func() interface{} { return &col }
	case t&MAP_TREE != 0:
		col := *treemap.NewWith(wrapComparator(STRING, comp[0]))
		c.ref = func() interface{} { return &col }
	case t&MAP_TREEBIDI != 0:
		col := *treebidimap.NewWith(wrapComparator(STRING, comp[0]), wrapComparator(STRING, comp[1]))
		c.ref = func() interface{} { return &col }
	case t&SETS != 0:
		if keymapped {
			if t == SET_HASH {
				col := *hashset.New()
				c.ref = func() interface{} { return &col }
			} else {
				col := *treeset.NewWith(wrapComparator(STRING, comp[0]))
				c.ref = func() interface{} { return &col }
			}
		} else {
			if t == SET_HASH {
				col := *hashset.New()
				c.ref = func() interface{} { return &col }
			} else {
				col := *treeset.NewWith(wrapComparator(INTEGER, comp[0]))
				c.ref = func() interface{} { return &col }
			}
		}
	case t&TREES != 0:
		if keymapped {
			if t == TREE_BINHEAP {
				col := *binaryheap.NewWith(wrapComparator(STRING, comp[0]))
				c.ref = func() interface{} { return &col }
			} else {
				col := *redblacktree.NewWith(wrapComparator(STRING, comp[0]))
				c.ref = func() interface{} { return &col }
			}
		} else {
			if t == TREE_BINHEAP {
				col := *binaryheap.NewWith(wrapComparator(INTEGER, comp[0]))
				c.ref = func() interface{} { return &col }
			} else {
				col := *redblacktree.NewWith(wrapComparator(INTEGER, comp[0]))
				c.ref = func() interface{} { return &col }
			}
		}
	}

	// pass container type and the reference function to the method
	// wrapper. yields all other functions needed to assign to the
	// collection structs fields.
	c.container, c.iterator, c.enumerable = wrapCollectionMethods(t, c.ref)

	// return a reference to the collection struct, holding the closures,
	// mapped to its Collection implementing methods.
	return &c
}
func newIndexedSetContainer(t CntType, comp ...Comparator) Collection {
	var c = collection{}
	c.t = t
	switch {
	case t&SET_HASH != 0:
		col := *hashset.New()
		c.ref = func() interface{} { return &col }
	case t&SET_TREE != 0: // KEY OR INDEX
		col := *treeset.NewWith(wrapComparator(INTEGER, comp[0]))
		c.ref = func() interface{} { return &col }

	}
	c.container, c.iterator, c.enumerable = wrapCollectionMethods(t, c.ref)
	return &c
}
func newKeymappedSetContainer(t CntType, comp Comparator) Collection {
	var c = collection{}
	c.t = t
	col := *hashset.New()
	c.ref = func() interface{} { return &col }
	c.container, c.iterator, c.enumerable = wrapCollectionMethods(t, c.ref)
	return &c
}
func newIndexedTreeContainer(t CntType, comp ...Comparator) Collection {
	var c = collection{}
	c.t = t
	switch {
	case t&TREE_BINHEAP != 0:
		col := *binaryheap.NewWith(wrapComparator(INTEGER, comp[0]))
		c.ref = func() interface{} { return &col }
	case t&TREE_REDBLACK != 0:
		col := *redblacktree.NewWith(wrapComparator(INTEGER, comp[0]))
		c.ref = func() interface{} { return &col }

	}
	c.container, c.iterator, c.enumerable = wrapCollectionMethods(t, c.ref)
	return &c
}
func newKeymappedTreeContainer(t CntType, comp ...Comparator) Collection {
	var c = collection{}
	c.t = t
	switch {
	case t&TREE_BINHEAP != 0:
		col := *binaryheap.NewWith(wrapComparator(STRING, comp[0]))
		c.ref = func() interface{} { return &col }
	case t&TREE_REDBLACK != 0:
		col := *redblacktree.NewWith(wrapComparator(STRING, comp[0]))
		c.ref = func() interface{} { return &col }

	}
	c.container, c.iterator, c.enumerable = wrapCollectionMethods(t, c.ref)
	return &c
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
	CommonContainerType() CntType // combined bitflag LISTS, STACKS, SETS, MAPS, TREES
	ContainerType() CntType       // bitflag with single bit set
	Container() Container
	Iterator() Iterator
	Enumerable() Enumerable
	SubContainer() interface{} // assertable to either List, Stack, Set, Map, Tree

	contained() interface{} // either List, Stack, Set, Map, or Tree
}

// COLLECTION IMPLEMENTATION
type collection struct {
	t          CntType
	container  func() Container
	iterator   func() Iterator
	enumerable func() Enumerable
	ref        func() interface{} // contains raw struct from gods
	sub        func() interface{} // contains List, Stack, Set, Map, or Tree interface
}

func (c *collection) SubContainer() (i interface{}) {
	switch {
	case c.t&LISTS != 0:
		i = List((*c).ref().(lists.List))
	case c.t&STACKS != 0:
		i = Stack((*c).ref().(stacks.Stack))
	case c.t&SETS != 0:
		i = Set((*c).ref().(sets.Set))
	case c.t&MAPS != 0:
		i = Map((*c).ref().(maps.Map))
	case c.t&TREES != 0:
		i = Tree((*c).ref().(trees.Tree))
	}
	return i
}
func (c *collection) contained() interface{} { return (*c).ref() }

func (c *collection) Container() Container   { return (*c).container() }
func (c *collection) Iterator() Iterator     { return (*c).iterator() }
func (c *collection) Enumerable() Enumerable { return (*c).enumerable() }

func (c collection) ContainerType() CntType { return c.t }
func (c collection) CommonContainerType() CntType {
	var t = c.t
	var ct CntType = 0
	switch {
	case t&LISTS != 0:
		ct = LISTS
	case t&STACKS != 0:
		ct = LISTS
	case t&SETS != 0:
		ct = LISTS
	case t&MAPS != 0:
		ct = LISTS
	case t&TREES != 0:
		ct = LISTS
	}
	return ct
}
