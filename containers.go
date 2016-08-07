// ENUMERATOR INTERFACE & IMPLEMENTING FUNCTIONS
//
// are defined as methods on the cont type and supposed to wrap the god
// enumerators and replace their empty interface return values by instances of
// appropriate value types.
//
// the interface conceals the differences between keyed and indexed enumerable
// containers by taking their key, respectively index as parameter of the value
// type and decide based on the value type, which enumerator, or iterator to
// choose.
//
// koi is supposed to be either a numeral type, in which case an indexed
// container is assumed, OR it can be of a container type, the bytes, or string
// value type and is assumed to be the key to map the value to.
//
// ENUMERATOR
//
//    func GetSortedValues(container Container, comparator utils.Comparator) []interface{}
//    Each calls the given function once for each element, passing that element's index and value.
//
//    Each(func(index int, value interface{}))
//	(or key interface{})
//
//    Any passes each element of the container to the given function and
//    returns true if the function ever returns true for any element.
//
//    Any(func(index int, value interface{}) bool) bool
//	(or key interface{})
//
//    All passes each element of the container to the given function and
//    returns true if the function returns true for all elements.
//
//    All(func(index int, value interface{}) bool) bool
//	(or key interface{})
//
//    Find passes each element of the container to the given function and returns
//    the first (index,value) for which the function is true or -1,nil otherwise
//    if no element matches the criteria.
//
//    Find(func(index int, value interface{}) bool) (int, interface{})
//	(or key interface{})
//
// ITERATOR
//
//    Next moves the iterator to the next element and returns true if there was a next element in the container.
//    If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
//    If Next() was called for the first time, then it will point the iterator to the first element if it exists.
//    Modifies the state of the iterator.
//
//    Next() bool
//
//    Value returns the current element's value.
//    Does not modify the state of the iterator.
// 					  ______
//    Value() interface{} 	 		\
// 						 \
//    Index returns the current element's index.  \ KoI() Value
//    Does not modify the state of the iterator.  / contains either an integer, or a type suited as a map key
// 						 /
//    Index() int 			  ______/
//
//    Begin resets the iterator to its initial state (one-before-first)
//    Call Next() to fetch the first element if any.
//
//    Begin()
//
//    First moves the iterator to the first element and returns true if there was a first element in the container.
//    If First() returns true, then first element's index and value can be retrieved by Index() and Value().
//    Modifies the state of the iterator.
//
//    First() bool
//
// REVERSE ITERATOR
//
//    Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
//    If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
//    Modifies the state of the iterator.
//
//    Prev() bool
//
//    End moves the iterator past the last element (one-past-the-end).
//    Call Prev() to fetch the last element if any.
//
//    End()
//
//    Last moves the iterator to the last element and returns true if there was a last element in the container.
//    If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
//    Modifies the state of the iterator.
//
//    Last() bool
//
// UTILS
//
//    func Sort(values []interface{}, comparator Comparator)
//    func StringComparator(a, b interface{}) int
//    type Comparator func(a, b interface{}) int
//
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
//
// LISTS
//
//    Get(index int) (interface{}, bool)
//    Remove(index int)
//    Add(values ...interface{})
//    Contains(values ...interface{}) bool
//    Sort(comparator utils.Comparator)
//    Swap(index1, index2 int)
//    Insert(index int, values ...interface{})
//
//    array 	 rev|idx
//    doublly	 rev|idx
//    singly  	 uni|idx
//
// STACKS
//
//    Push(value interface{})
//    Pop() (value interface{}, ok bool)
//    Peek() (value interface{}, ok bool)
//
//    array	 rev|idx
//    singlu  	 uni|idx
//
// MAPS
//
//    Put(key interface{}, value interface{})
//    Get(key interface{}) (value interface{}, found bool)
//    Remove(key interface{})
//    Keys() []interface{}
//
// 	BIDIMAP
//
//    	GetKey(value interface{}) (key interface{}, found bool)
//
//    hash 	 nil|key
//    hashbidi	 nil|key
//    tree  	 rev|key
//    treebidi 	 rev|key
//
//
// SETS
//
//    Add(elements ...interface{})
//    Remove(elements ...interface{})
//    Contains(elements ...interface{}) bool
//
//    hash 	 nil|nil
//    tree  	 rev|key
//
//
// TREES
//
// the two tree types have not that much in common, so tree only implements
// container, while all detail is left to the tree implementation in question.
//
//
// 	TREE (COMMON)
//
//        func NewWith(comparator utils.Comparator) *Heap |  func NewWith(comparator utils.Comparator) *Tree
//        func NewWithIntComparator() *Heap 		  |  func NewWithIntComparator() *Tree
//        func NewWithStringComparator() *Heap 		  |  func NewWithStringComparator() *Tree
//        func Clear() 			                  |  func (tree *Tree) Clear()
//        func Empty() bool 		    	          |  func (tree *Tree) Empty() bool
//        func Iterator() Iterator 	   	          |  func (tree *Tree) Iterator() Iterator
//        func Size() int 		                  |  func (tree *Tree) Size() int
//        func String() string 		 	          |  func (tree *Tree) String() string
//        func Values() []interface{} 	       	          |  func (tree *Tree) Values() []interface{}
//
//
//	    HEAP
//
//	      func (heap *Heap) Peek() (value interface{}, ok bool)
//	      func (heap *Heap) Pop() (value interface{}, ok bool)
//	      func (heap *Heap) Push(value interface{})
//
//	    RED-BLACK
//
//	      func (tree *Tree) Floor(key interface{}) (floor *Node, found bool)
//  	      func (tree *Tree) Ceiling(key interface{}) (ceiling *Node, found bool)
//	      func (tree *Tree) Put(key interface{}, value interface{})
//	      func (tree *Tree) Get(key interface{}) (value interface{}, found bool)
//	      func (tree *Tree) Left() *Node
//	      func (tree *Tree) Right() *Node
//	      func (tree *Tree) Remove(key interface{})
//	      func (tree *Tree) Keys() []interface{}
//
//
//    heap 	 rev|idx
//    heap 	 rev|key
//
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

	// sets of containers that share a more specific interface than
	// gods/containers and have other method signatures in common
	LISTS  = LIST_ARRAY | LIST_SINGLE | LIST_DOUBLE
	SETS   = SET_HASH | SET_TREE
	STACKS = STACK_ARRAY | STACK_LINKED
	MAPS   = MAP_HASH | MAP_HASHBIDI | MAP_TREE | MAP_TREEBIDI
	TREES  = TREE_BINHEAP | TREE_REDBLACK

	INDEXED = LISTS | STACKS
	KEYED   = SETS | MAPS | TREES
	REVERSE = LIST_ARRAY | LIST_DOUBLE | SET_TREE | STACK_ARRAY |
		MAP_TREEBIDI | MAP_TREE | TREE_BINHEAP | TREE_REDBLACK
)

// CONTAINER INTERFACE
//
// enumerator and iterator interfaces, provided by the concealed container
// type, as well as the Interface defining the parent type of the container
// (list, Set, Stack, Map, Tree).
//
// function to assign the nescessary interface methods to a god container, to
// make it implement the internal container interface, containing Values instances
// instead of empty interfaces.
//
type Container interface {
	Values() []Value
	Slice() []interface{}
	Empty() bool
	Size() int
	Clear()
}
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

// CONTAINER WRAPPER (instanciates a struct type to hold the container,
type cont struct {
	containers.Container
}

// CONVERSIONS
func (c cont) Slice() []interface{} { return c.Container.(containers.Container).Values() }
func (c cont) Values() []Value      { return interfaceSlice(c.Slice()).Values() }
func wrapContainer(c containers.Container) cont {
	return cont{c}
}
func newContainer(t CntType, comp ...Comparator) Container {
	var v Container
	switch {
	case t&LISTS != 0:
		v = newListContainer(t)
	case t&SETS != 0:
		c := newSetContainer(t)
	case t&MAPS != 0:
		c := newMapContainer(t, comp[0])
	case t&STACKS != 0:
		v = newStackContainer(t)
	case t&TREES != 0:
		v = newTreeContainer(t)
	}
	return v
}
func newListContainer(t CntType) Container {
	var c lists.List
	switch {
	case t&LIST_ARRAY != 0:
		c = arraylist.New()
	case t&LIST_SINGLE != 0:
		c = singlylinkedlist.New()
	case t&LIST_DOUBLE != 0:
		c = doublylinkedlist.New()
	}
	return wrapContainer(c)
}
func newStackContainer(t CntType) Container {
	var c stacks.Stack
	switch {
	case t&STACK_ARRAY != 0:
		c = arraystack.New()
	case t&STACK_LINKED != 0:
		c = linkedliststack.New()
	}
	return wrapContainer(c)
}
func newMapContainer(t CntType, comp ...Comparator) Container {
	var c maps.Map
	switch {
	case t&MAP_HASH != 0:
		c = hashmap.New()
	case t&MAP_HASHBIDI != 0:
		c = hashbidimap.New()
	case t&MAP_TREE != 0:
		c = treemap.NewWith(comp[0].Convert())
	case t&MAP_TREEBIDI != 0:
		c = treebidimap.NewWith(comp[0].Convert(), comp[1].Convert())
	}
	return wrapContainer(c)
}
func newSetContainer(t CntType, comp ...Comparator) Container {
	var c sets.Set
	switch {
	case t&SET_HASH != 0:
		c = hashset.New()
	case t&SET_TREE != 0:
		c := treeset.NewWith(comp[0].Convert())

	}
	return wrapContainer(c)
}
func newTreeContainer(t CntType, comp ...Comparator) Container {
	var c trees.Tree
	switch {
	case t&TREE_BINHEAP != 0:
		c = binaryheap.NewWith(comp[0].Convert())
	case t&TREE_REDBLACK != 0:
		c = redblacktree.NewWith(comp[0].Convert())

	}
	return wrapContainer(c)
}

// ENUMERABLE INTERFACE
//
type Enumerable interface {
	Each(func(koi Value, value Value))
	Any(func(koi Value, value Value) bool) bool
	All(func(koi Value, value Value) bool) bool
	Find(func(koi Value, value Value) bool) (Value, Value)
}

// ENUMERABLE RETURN TYPE IMPLEMENTATIONS
//
type enumIdx struct {
	ContType CntType // CHECK TYPE FIRST!
	enum     containers.EnumerableWithIndex
}
type enumKey struct {
	ContType CntType // CHECK TYPE FIRST!
	enum     containers.EnumerableWithKey
}

// WRAP_ENUM
// chooses which enumerable interface implementation to instanciate and maps
// the appropriate functions to the instance.
//func wrapEnum(t CntType, koi Value, value Value) Enumerable {
//	var ret utils.Comparator
//	if t&INDEXED != 0 {
//		ret = func(koi, value interface{}) int {
//			idx := koi.(intVal).Integer()
//			return utils.IntComparator(idx, value)
//		}
//	} else {
//		ret = func(koi, value interface{}) int {
//			key := koi.(strVal).String()
//			return utils.StringComparator(key, value)
//		}
//	}
//}

// ENUMERABLE FUNCTIONS
// enumerable interface implementing functions
//
func Each(e containers.Enumerable, t CntType, koi Value, value Value) {
	var comp utils.Comparator
	if t&INDEXED != 0 {
		comp = func(koi, value interface{}) int {
			idx := koi.(intVal).Integer()
			return utils.IntComparator(idx, value)
		}
	} else {
		comp = func(koi, value interface{}) int {
			key := koi.(strVal).String()
			return utils.StringComparator(key, value)
		}
	}
	e.Each(comp)
}
func Any(e Enumerable, fn func(koi Value, value Value) bool) bool                             {}
func All(e Enumerable, fn func(koi Value, value Value) bool) bool                             {}
func Find(e Enumerable, fn func(index Value, value Value) bool) (Value, Value)                {}
func findByKey(e Enumerable, fn func(index int, value interface{}) bool) (Value, Value)       {}
func findByIdx(e Enumerable, fn func(key interface{}, value interface{}) bool) (Value, Value) {}

// ITERATOR INTERFACE
type Iterator interface {
	Value() Value
	Begin()
	Next() bool
	First() bool
}

//	REVERSE
//
type Reverse interface {
	End()
	Preview() bool
	Last() bool
}

//			r.i
//			r.next
//			r.value
//			r.koi
//			r.begin
//			r.first
//			r.prev
//			r.end
//			r.last
//		r.next
//		r.value
//		r.koi
//		r.begin
//		r.first
//			r.next
//			r.value
//			r.koi
//			r.begin
//			r.first
//			r.prev
//			r.end
//			r.last
//		r.next
//		r.value
//		r.koi
//		r.begin
//		r.first
//
//	arraylistFn
//	singlylinkedlistFn
//	doublylinkedlistFn
//
//	hashsetFn
//	treesetFn
//
//	hashbidimapFn
//	hashmapFn
//	treebidimapFn
//	treemapFn
//
//	arraystackFn
//	linkedliststackFn
//
//	binaryheapFn
//	redblacktreeFn

// COMPARATOR INTERFACES
// gods comparators expect empty interfaces that are assertable to either
// string, or int. The Values interface allows for much more complex types. the
// comparator function can be set up on arbitratry types, methods, or fields of
// complex types, as long as it is converted to the correct signature before
// passed to god.
type Comparator func(a, b Value) int
type IntComparator func(a, b Value) int
type StringComparator func(a, b Value) int

func (c Comparator) Convert() utils.Comparator {
	var r utils.Comparator = c.Convert()
	return r
}

func ConstructComparator(t ValueType) utils.Comparator {
	var f utils.Comparator
	switch t {
	case STRING:
		f = utils.StringComparator
	case INTEGER:
		f = utils.IntComparator
	}
	return f
}
