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
	"github.com/emirpasic/gods/stacks"
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/emirpasic/gods/stacks/linkedliststack"
	"github.com/emirpasic/gods/trees"
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/emirpasic/gods/trees/redblacktree"
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
		is = append(is, i)
	}
	return is
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

// STACK INTERFACE
type Stack interface {
	Push(value Val)
	Pop() (value Val, ok bool)
	Peek() (value Val, ok bool)
}

// STACK IMPLEMENTATION
type stackCnt struct {
	*cont
	stacks.Stack
}

func newStackContainer(t CntType) (c Container) {
	switch t {
	case STACK_ARRAY:
		c = wrapContainer(t, arraystack.New())
	case STACK_LINKED:
		c = wrapContainer(t, linkedliststack.New())
	}
	return c
}

// STACK INTERFACE
type Tree interface {
	Container
}

// STACK IMPLEMENTATION
type treeCnt struct {
	*cont
	trees.Tree
}

func newTreeContainer(t CntType, comp Comparator) (c Container) {
	switch t {
	case TREE_REDBLACK:
		c = wrapContainer(t, redblacktree.NewWith(comp.Convert()))
	case TREE_BINHEAP:
		c = wrapContainer(t, binaryheap.NewWith(comp.Convert()))
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
