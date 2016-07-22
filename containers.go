package agiledoc

import (
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
	// 	"github.com/emirpasic/gods/stacks"
	// 	"github.com/emirpasic/gods/stacks/arraystack"
	// 	"github.com/emirpasic/gods/stacks/linkedliststack"
	// 	"github.com/emirpasic/gods/trees"
	// 	"github.com/emirpasic/gods/trees/binaryheap"
	// 	"github.com/emirpasic/gods/trees/redblacktree"
	"github.com/emirpasic/gods/utils"
)

// CONTAINER INTERFACE (extends the gods container interface)
// interface to conceal god container empty interface values behind the Val
// interface, that provides a type function to inspect the nature of the
// contained value without using reflection.. Since containers themselves and
// there conained values both implement the val interface, all container types
// are fully recursive. The Values returned from and passed to mapped
// containers and sets, also implement the Var interface, KeyVal identitys are
// taken as map keys to map their values on to.
type Container interface {
	ContType() CntType
	Empty() bool
	Size() int
	Clear()
	Values() []Val // might be keyed vars/params
}

type cont struct {
	ct CntType
	containers.Container
}

func (c cont) ContType() CntType { return c.ct }

func (c *cont) Values() []Val                                   { return interfaceSlice((*c).Container.Values()).Values() }
func wrapContainer(t CntType, c containers.Container) Container { return &cont{t, c} }

// SLICE TYPES
// these types exist, so that a slice of interfaces, as well as a slice of
// Values implements a type, methods can be assigned to. That Way unwrapped
// slices can allways be converted to those types and provide either the
// Values(), or the Interfaces() method that converts them to the corredponding
// slice type.
type interfaceSlice []interface{}
type valSlice []Val

func (i interfaceSlice) Values() []Val {
	var vs = []Val{}
	for _, v := range i {
		v := v.(Val)
		vs = append(vs, v)
	}
	return vs
}
func (v valSlice) Interfaces() []interface{} {
	var is = []interface{}{}
	for _, i := range v {
		is = append(is, v)
	}
	return is
}

// LIST INTERFACE
type List interface {
	Get(index int) (Val, bool) //
	Remove(index int)
	Add(values ...Val)           //
	Contains(values ...Val) bool //
	Sort(comparator Comparator)  //
	Swap(index1, index2 int)
	Insert(index int, values ...Val) //

	Container
}

// LIST IMPLEMENTATION
type listCnt struct {
	*cont
	lists.List
}

func (l *listCnt) Add(v ...Val) {
	(*l).List.Add(valSlice(v).Interfaces())
}
func (l *listCnt) Insert(i int, v ...Val) {
	(*l).List.Insert(i, valSlice(v).Interfaces())
}
func (l *listCnt) Contains(v ...Val) bool {
	return (*l).List.Contains(valSlice(v).Interfaces())
}
func (l *listCnt) Get(i int) (Val, bool) {
	v, ok := (*l).Get(i)
	return v.(Val), ok
}

func (l *listCnt) Sort(c Comparator) {
	(*l).List.Sort(c.Convert())
}
func newListContainer(t CntType) (c Container) {
	switch t {
	case LIST_ARRAY:
		c = wrapContainer(t, arraylist.New())
	case LIST_SINGLE:
		c = wrapContainer(t, singlylinkedlist.New())
	case LIST_DOUBLE:
		c = wrapContainer(t, doublylinkedlist.New())
	}
	return c
}

// MAP INTERFACE
type Map interface {
	Put(key Val, value Val)
	Get(key Val) (value Val, found bool)
	Remove(key Val)
	Keys() []Val
}

// MAP IMPLEMENTATION
type mapCnt struct {
	*cont
	maps.Map
}

func (m *mapCnt) Values() []Val {
	return interfaceSlice((*m).Map.Values()).Values()
}

func (m *mapCnt) Put(k Val, v Val) {
	(*m).Put(k, v)
}

type BidiMap interface {
	GetKey(value Val) (key Val, found bool)
	Map
}

// trees need one, or two comparators, while maps dont. Comparators can be of
// different index type. apart from the designated container type, a variadic
// idxType can be passed. The exact ammount of comparators needed per type:
// hash, hashbidi = 0 | tree = 1 | treebidi = 2
func newMapContainer(t CntType, idxType ...ValType) (m Container) {
	switch t {
	case MAP_HASH:
		m = wrapContainer(t, hashmap.New())
	case MAP_HASHBIDI:
		m = wrapContainer(t, hashbidimap.New())
	case MAP_TREE:
		m = wrapContainer(t, treemap.NewWith(ConstructComparator(idxType[0])))
	case MAP_TREEBIDI:
		m = wrapContainer(t, treebidimap.NewWith(ConstructComparator(idxType[0]), ConstructComparator(idxType[1])))
	}
	return m
}

// SET INTERFACE
type Set interface {
	Add(elements ...Val)
	Remove(elements ...Val)
	Contains(elements ...Val) bool
}

// SET IMPLEMENTATION
type setCnt struct {
	*cont
	sets.Set
}

func (s *setCnt) Add(e ...Val) {
	(*s).Set.Add(valSlice(e).Interfaces())
}
func (s *setCnt) Remove(e ...Val) {
	(*s).Set.Remove(valSlice(e).Interfaces())
}
func (s *setCnt) Contains(e ...Val) bool {
	return (*s).Set.Contains(valSlice(e).Interfaces())
}
func newSetContainer(t CntType, idxType ValType) (c Container) {
	switch t {
	case SET_HASH:
		c = wrapContainer(t, hashset.New())
	case SET_TREE:
		c = wrapContainer(t, treeset.NewWith(ConstructComparator(idxType)))
	}
	return c
}

// container type marks the type of container taken from the god library
type CntType uint16

const (
	// liat of all container tupes imported from gods
	// every type implements at least the container interface and one more
	// specific interface that combines all types that share a common kind
	// of data structure: lists, sets, maps, stacks and trees.
	//
	// some of those data structures can exist in indexed and/or mapped
	// versions. Dependend from the type, they may implement additional
	// interfaces, like iteratorWithKey, IteratorWithIndex... Comparators
	// and so on (see gods documentation)
	//
	///////////// lists
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

	// sets of containers that share a more specific interface than
	// gods/containers and have other method signatures in common
	LISTS  = LIST_ARRAY | LIST_SINGLE | LIST_DOUBLE
	SETS   = SET_HASH | SET_TREE
	STACKS = STACK_ARRAY | STACK_LINKED
	MAPS   = MAP_HASH | MAP_HASHBIDI | MAP_TREE | MAP_TREEBIDI
	TREES  = TREE_BINHEAP | TREE_REDBLACK
)

var ContainerConstructors = []func() Container{}

func InitializeContainers() {
	// 	ContainerConstructors[LIST_ARRAY] = func() Container {}
	// 	ContainerConstructors[LIST_SINGLE]
	// 	ContainerConstructors[LIST_DOUBLE]
	// 	ContainerConstructors[SET_HASH]
	// 	ContainerConstructors[SET_TREE]
	// 	ContainerConstructors[STACK_LINKED]
	// 	ContainerConstructors[STACK_ARRAY]
	// 	ContainerConstructors[MAP_HASH]
	// 	ContainerConstructors[MAP_HASHBIDI]
	// 	ContainerConstructors[MAP_TREE]
	// 	ContainerConstructors[MAP_TREEBIDI]
	// 	ContainerConstructors[TREE_REDBLACK]
	// 	ContainerConstructors[TREE_BINHEAP]
}

// COMPARATOR INTERFACES
type Comparator func(a, b Val) int

func (c Comparator) Convert() utils.Comparator {
	var r utils.Comparator = c.Convert()
	return r
}

type IntComparator func(a, b Val) int
type StringComparator func(a, b Val) int

func ConstructComparator(t ValType) utils.Comparator {
	var f utils.Comparator
	switch t {
	case STRING:
		f = utils.StringComparator
	case INTEGER:
		f = utils.IntComparator
	}
	return f
}

// ENUMERABLE & ITERATOR INTERFACES
type EnumerableWithIndex interface {
	// Each calls the given function once for each element, passing that element's index and value.
	Each(func(index int, value Val))

	// Any passes each element of the container to the given function and
	// returns true if the function ever returns true for any element.
	Any(func(index int, value Val) bool) bool

	// All passes each element of the container to the given function and
	// returns true if the function returns true for all elements.
	All(func(index int, value Val) bool) bool

	// Find passes each element of the container to the given function and returns
	// the first (index,value) for which the function is true or -1,nil otherwise
	// if no element matches the criteria.
	Find(func(index int, value Val) bool) (int, Val)
}
type EnumerableWithKey interface {
	// Each calls the given function once for each element, passing that element's key and value.
	Each(func(key Val, value Val))

	// Any passes each element of the container to the given function and
	// returns true if the function ever returns true for any element.
	Any(func(key Val, value Val) bool) bool

	// All passes each element of the container to the given function and
	// returns true if the function returns true for all elements.
	All(func(key Val, value Val) bool) bool

	// Find passes each element of the container to the given function and returns
	// the first (key,value) for which the function is true or nil,nil otherwise if no element
	// matches the criteria.
	Find(func(key Val, value Val) bool) (Val, Val)
}
type IteratorWithIndex interface {
	// Next moves the iterator to the next element and returns true if there was a next element in the container.
	// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
	// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Value returns the current element's value.
	// Does not modify the state of the iterator.
	Value() Val

	// Index returns the current element's index.
	// Does not modify the state of the iterator.
	Index() int

	// Begin resets the iterator to its initial state (one-before-first)
	// Call Next() to fetch the first element if any.
	Begin()

	// First moves the iterator to the first element and returns true if there was a first element in the container.
	// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
	// Modifies the state of the iterator.
	First() bool
}
type IteratorWithKey interface {
	// Next moves the iterator to the next element and returns true if there was a next element in the container.
	// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
	// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
	// Modifies the state of the iterator.
	Next() bool

	// Value returns the current element's value.
	// Does not modify the state of the iterator.
	Value() Val

	// Key returns the current element's key.
	// Does not modify the state of the iterator.
	Key() Val

	// Begin resets the iterator to its initial state (one-before-first)
	// Call Next() to fetch the first element if any.
	Begin()

	// First moves the iterator to the first element and returns true if there was a first element in the container.
	// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	First() bool
}
type ReverseIteratorWithIndex interface {
	// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
	// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
	// Modifies the state of the iterator.
	Prev() bool

	// End moves the iterator past the last element (one-past-the-end).
	// Call Prev() to fetch the last element if any.
	End()

	// Last moves the iterator to the last element and returns true if there was a last element in the container.
	// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
	// Modifies the state of the iterator.
	Last() bool

	IteratorWithIndex
}
type ReverseIteratorWithKey interface {
	// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
	// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	Prev() bool

	// End moves the iterator past the last element (one-past-the-end).
	// Call Prev() to fetch the last element if any.
	End()

	// Last moves the iterator to the last element and returns true if there was a last element in the container.
	// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
	// Modifies the state of the iterator.
	Last() bool

	IteratorWithKey
}
