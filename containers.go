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
)

// CONTAINER INTERFACE
type Container interface {
	ContType() CntType
	Empty() bool
	Size() int
	Clear()
}

// CONTAINER IMPLEMENTATION
type cont struct {
	ct CntType
	containers.Container
}

// returns the container type
//
func (c cont) ContType() CntType { return c.ct }

// returns a slice of Values from a container instance
//
func (c cont) Values() []Value { return interfaceSlice(c.Container.Values()).Values() }

// CONTAINER WRAPPER
//
// function to assign the nescessary interface methods to a god container, to
// make it implement the internal container interface, containing Values instances
// instead of empty interfaces.
//
func NewContainer(t CntType, c ...Comparator) (r Container) {
	switch {
	case t&LISTS != 0:
		r = newlistContainer(t)
	case t&MAPS != 0:
		r = newMapContainer(t, c...)
	case t&SETS != 0:
		r = newSetContainer(t, c...)
	case t&STACKS != 0:
		r = newStackContainer(t)
	case t&TREES != 0:
		r = newTreeContainer(t, c...)
	default:
		r = &vecVal{}
	}
	return r
}

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

// SLICE HELPER TYPES
// these types exist, so that a slice of interfaces, as well as a slice of
// Values implements a type, methods can be assigned to. That Way unwrapped
// slices can allways be converted to those types and provide either the
// Values(), or the Interfaces() method that converts them to the corredponding
// slice type.
type interfaceSlice []interface{}
type valSlice []Value
type byteSlice []byte
type boolSlice []bool

func (b boolSlice) Bytes() []byte {

	// bool-slice-integer array, combines up to 8 booleans in a slice to
	// represent byte sized bitflags
	var bsi [8]bool = [8]bool{}

	// byte slice to concatenate all booleans in bitflag of arbitrary size
	var bs = []byte{}

	// split input into byte sized chunks and iterate over each of those chunks
	for o := 0; o < (len(b)/8 + 1); o++ {

		var u uint8 = 0 // allocate a new uint8 for each byte sized chunk

		if o < len(b) { // if bool slice is not yet depleted

			// iterate over each of the eight bits of the current chunk
			for i := 0; i < 8; i++ {
				i := i
				// dereference bool at the current index
				if bsi[i] { // if element is true, set a bit at the current index
					u = 1 << uint(i)
				}
			}
		} // end of chunk. since lenght check failed, another iteration is possibly needed.
		// append the last produced chunk at the byte slice intended to return
		bs = append(bs, u)

	} // either iterate on, or return byte slice
	// depending on the total number of chunks
	return bs
}
func (b boolSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range b {
		v := v
		vs = append(vs, Value().ToValue(v))
	}
	return vs
}
func (i byteSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range i {
		v := v
		vs = append(vs, byteVal(v))
	}
	return vs
}
func (v byteSlice) Interfaces() []interface{} {
	var is = []interface{}{}
	for _, i := range v {
		is = append(is, i)
	}
	return is
}
func (i interfaceSlice) Values() []Value {
	var vs = []Value{}
	for _, v := range i {
		v := v
		vs = append(vs, v.(byteVal))
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

// CUSTOMIZED CONTAINER INTERFACES AND IMPLEMENTATIONS
//
// listCnt INTERFACE
type List interface {
	Get(index int) (Value, bool)
	Remove(index int)
	Add(values ...Value)
	Contains(values ...Value) bool
	Sort(comparator Comparator)
	Swap(index1, index2 int)
	Insert(index int, values ...Value)

	Container
}

// listCnt IMPLEMENTATION
type listCnt struct {
	CntType
	lists.List
}

func (l listCnt) ContType() CntType { return l.CntType }
func (l *listCnt) Values() []Value  { s := (*l).List.Values(); return interfaceSlice(s).Values() }
func (l *listCnt) Add(v ...Value) {
	is := valSlice(v).Interfaces()
	(*l).List.Add(is)
}
func (l *listCnt) Insert(i int, v ...Value) { (*l).List.Insert(i, valSlice(v).Interfaces()) }
func (l *listCnt) Contains(v ...Value) bool { return (*l).List.Contains(valSlice(v).Interfaces()) }
func (l *listCnt) Sort(c Comparator)        { (*l).List.Sort(c.Convert()) }
func (l *listCnt) Get(i int) (Value, bool) {
	v, ok := (*l).Get(i)
	return v.(Value), ok
}

// listCnt CONSTRUCTOR
// the listCnt constructor only needs to know the dedicated type of the listCnt
// container to instanciate
func newlistContainer(t CntType) (l *listCnt) {
	switch t {
	case LIST_ARRAY:
		l = &listCnt{t, arraylist.New()}
	case LIST_SINGLE:
		l = &listCnt{t, singlylinkedlist.New()}
	case LIST_DOUBLE:
		l = &listCnt{t, doublylinkedlist.New()}
	}
	return l
}

// MAP INTERFACE
type Map interface {
	Put(key Value, value Value)
	Get(key Value) (value Value, found bool)
	Remove(key Value)
	Keys() []Value

	Container
}

// MAP IMPLEMENTATION
type mapCnt struct {
	CntType
	maps.Map
}

func (m *mapCnt) ContType() CntType { return m.CntType }
func (m *mapCnt) Values() []Value {
	return interfaceSlice((*m).Map.Values()).Values()
}
func (m *mapCnt) Get(i Value) (Value, bool) {
	v, ok := (*m).Map.Get(i)
	return v.(Value), ok
}
func (m *mapCnt) Keys() []Value {
	return interfaceSlice((*m).Map.Keys()).Values()
}
func (m *mapCnt) Put(k Value, v Value) {
	(*m).Put(k, v)
}
func (m *mapCnt) Remove(k Value) {
	(*m).Map.Remove(k)
}

type BidiMap interface {
	GetKey(value Value) (key Value, found bool)
	Map // allready contains container
}

// MAP CONSTRUCTOR
//
// trees need one, or two comparators, while maps dont. Comparators can be of
// different index type. apart from the designated container type, a variadic
// idxType can optionaly be passed. The exact ammount of comparators needed, is
// dependent on its type:
//
// | hash, hashbidi = 0 | tree = 1 | treebidi = 2 |
//
func newMapContainer(t CntType, c ...Comparator) (m Container) {
	switch t {
	case MAP_HASH:
		m = &mapCnt{t, hashmap.New()}
	case MAP_HASHBIDI:
		m = &mapCnt{t, hashbidimap.New()}
	case MAP_TREE:
		m = &mapCnt{t, treemap.NewWith(c[0].Convert())}
	case MAP_TREEBIDI:
		m = &mapCnt{t, treebidimap.NewWith(c[0].Convert(), c[1].Convert())}
	}
	return m
}
func wrapMap(t CntType, m maps.Map) (r Map) {
	return &mapCnt{t, m}
}

// SET INTERFACE
type Set interface {
	Add(elements ...Value)
	Remove(elements ...Value)
	Contains(elements ...Value) bool

	Container
}

// SET IMPLEMENTATION
type setCnt struct {
	CntType
	sets.Set
}

func (s *setCnt) ContType() CntType { return s.CntType }
func (s *setCnt) Values() []Value {
	return interfaceSlice((*s).Set.Values()).Values()
}
func (s *setCnt) Add(e ...Value) {
	i := valSlice(e).Interfaces()
	(*s).Set.Add(i...)
}
func (s *setCnt) Remove(e ...Value) {
	(*s).Set.Remove(valSlice(e).Interfaces())
}
func (s *setCnt) Contains(e ...Value) bool {
	return (*s).Set.Contains(valSlice(e).Interfaces())
}

// SET CONSTRUCTOR
// the treeset needs a comparator closure
func newSetContainer(t CntType, c ...Comparator) (r *setCnt) {
	switch t {
	case SET_HASH:
		r = &setCnt{t, hashset.New()}
	case SET_TREE:
		r = &setCnt{t, treeset.NewWith(c[0].Convert())}
	}
	return r
}

// STACK INTERFACE
type Stack interface {
	Push(value Value)
	Pop() (value Value, ok bool)
	Peek() (value Value, ok bool)

	Container
}

// STACK IMPLEMENTATION
type stackCnt struct {
	CntType
	stacks.Stack
}

func (s *stackCnt) ContType() CntType { return s.CntType }
func (s *stackCnt) Values() []Value {
	return interfaceSlice((*s).Stack.Values()).Values()
}
func (s *stackCnt) Peek() (Value, bool) {
	v, ok := (*s).Stack.Peek()
	return v.(Value), ok
}
func (s *stackCnt) Pop() (Value, bool) {
	v, ok := (*s).Stack.Pop()
	return v.(Value), ok
}
func (s *stackCnt) Push(v Value) { (*s).Stack.Push(v) }

// STACK CONSTRUCTOR
func newStackContainer(t CntType) (s Stack) {
	switch t {
	case STACK_ARRAY:
		s = &stackCnt{t, arraystack.New()}
	case STACK_LINKED:
		s = &stackCnt{t, linkedliststack.New()}
	}
	return s
}

// TREE INTERFACE
type Tree interface {
	Container
}

// TREE IMPLEMENTATION
type treeCnt struct {
	CntType
	trees.Tree
}

func (t *treeCnt) ContType() CntType { return (*t).CntType }
func (t *treeCnt) Values() []Value {
	return interfaceSlice((*t).Tree.Values()).Values()
}

// TREE CONSTRUCTOR
func newTreeContainer(t CntType, c ...Comparator) (r Container) {
	switch t {
	case TREE_REDBLACK:
		r = &treeCnt{t, redblacktree.NewWith(c[0].Convert())}
	case TREE_BINHEAP:
		r = &treeCnt{t, binaryheap.NewWith(c[0].Convert())}
	}
	return r
}
