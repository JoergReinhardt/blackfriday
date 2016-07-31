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

	KEYED   = LISTS | STACKS
	INDEXED = SETS | MAPS | TREES
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
type ValContainer interface {
	ContType() CntType
	Container
	Enumerable
	Iterator
}

// CONTAINER WRAPPER (instanciates a struct type to hold the container,
type cont struct {
	containers.Container
}

func (c cont) Slice() []interface{} { return c.Container.(containers.Container).Values() }
func (c cont) Values() []Value      { return interfaceSlice(c.Slice()).Values() }
func wrapContainer(c containers.Container) cont {
	return cont{c}
}
func newValueContainer(t CntType) ValContainer {
	var v ValContainer
	switch {
	case t&LISTS != 0:
		v = newListContainer(t)
	case t&SETS != 0:
		v = newSetContainer(t)
	case t&MAPS != 0:
		v = newMapContainer(t)
	case t&STACKS != 0:
		v = newStackContainer(t)
	case t&TREES != 0:
		v = newTreeContainer(t)
	}
	return v
}

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
type Enumerable interface {
	Each(func(koi Value, value Value))
	Any(func(koi Value, value Value) bool) bool
	All(func(koi Value, value Value) bool) bool
	Find(func(koi Value, value Value) bool) (Value, Value)
}
type enum struct {
	ContType func() CntType // CHECK TYPE FIRST!
	enum     interface{}
}

func wrapEnum(t CntType, e interface{}) Enumerable {
	return &enum{
		func() CntType { return t },
		e,
	}
}
func (rec *enum) Each(fn func(koi Value, value Value)) {
	switch {
	case (*rec).ContType()&KEYED != 0:
		(*rec).enum.(containers.EnumerableWithKey).Each(func(key interface{}, value interface{}) {
			fn(NativeToValue(key), NativeToValue(value))
		})
	case (*rec).ContType()&INDEXED != 0:
		(*rec).enum.(containers.EnumerableWithIndex).Each(func(index int, value interface{}) {
			fn(NativeToValue(index), NativeToValue(value))
		})
	}
}
func (rec *enum) Any(fn func(koi Value, value Value) bool) bool {
	var r bool
	switch {
	case (*rec).ContType()&KEYED != 0:
		r = (*rec).enum.(containers.EnumerableWithKey).Any(func(key interface{}, value interface{}) bool {
			return fn(NativeToValue(key), NativeToValue(value))
		})
	case (*rec).ContType()&INDEXED != 0:
		r = (*rec).enum.(containers.EnumerableWithIndex).Any(func(index int, value interface{}) bool {
			return fn(NativeToValue(index), NativeToValue(value))
		})
	}
	return r
}
func (rec *enum) All(fn func(koi Value, value Value) bool) bool {
	var r bool
	switch {
	case (*rec).ContType()&KEYED != 0:
		r = (*rec).enum.(containers.EnumerableWithKey).Any(func(key interface{}, value interface{}) bool {
			return fn(NativeToValue(key), NativeToValue(value))
		})
	case (*rec).ContType()&INDEXED != 0:
		r = (*rec).enum.(containers.EnumerableWithIndex).Any(func(index int, value interface{}) bool {
			return fn(NativeToValue(index), NativeToValue(value))
		})
	}
	return r
}

// FindBy<Key|Index>
// gets a function passed, that expects an arbitratry interface as key and a
// Value as it's parameters. Depending on the type of enumerator, the key
// either gets asserted as an integer, when the enumerator is indexed and left
// to be of type plain interface, when dealt with a keyed enumerator.
//
// the return values are values allready, or will be encapsulated before
// returned.
func (rec *enum) Find(fn func(index Value, value Value) bool) (Value, Value) {
	if (*rec).ContType()&INDEXED != 0 {
		var fen = func(index interface{}, value interface{}) bool {
			return fn(NativeToValue(index), NativeToValue(value))
		}
		i, v := (*rec).findByIdx(fen)
		return NativeToValue(i), NativeToValue(v)
	} else {
		var fkn = func(index int, value interface{}) bool {
			return fn(NativeToValue(index), NativeToValue(value))
		}
		i, v := (*rec).enum.(containers.EnumerableWithIndex).Find(fkn)
		return NativeToValue(i), NativeToValue(v)
	}
}
func (rec *enum) findByKey(fn func(index int, value interface{}) bool) (Value, Value) {
	var rk, rv interface{}
	rk, rv = (*rec).enum.(containers.EnumerableWithIndex).Find(fn)
	return NativeToValue(rk), NativeToValue(rv)
}
func (rec *enum) findByIdx(fn func(key interface{}, value interface{}) bool) (Value, Value) {
	var rk, rv interface{}
	rk, rv = (*rec).enum.(containers.EnumerableWithKey).Find(fn)
	return NativeToValue(rk), NativeToValue(rv)
}

type Iterator interface {
	Next() bool
	Value() Value
	KoI() Value
	Begin()
	First() bool
}

// ITERATOR WRAPPER
type iterator struct {
	ContType func() CntType
	i        interface{}
	next     func() bool
	value    func() Value
	koi      func() Value
	begin    func()
	first    func() bool
	prev     func() bool
	end      func()
	last     func() bool
}

func (i iterator) Next() bool   { return i.next() }
func (i iterator) Value() Value { return i.value() }
func (i iterator) KoI() Value   { return i.value() }
func (i iterator) Begin()       { i.begin() }
func (i iterator) First() bool  { return i.first() }
func (i iterator) Prev() bool   { return i.prev() }
func (i iterator) End()         { i.prev() }
func (i iterator) Last() bool   { return i.prev() }

func wrapIterator(t CntType, i interface{}) Iterator {
	var r = iterator{}
	if t&INDEXED != 0 {
		r.next = i.(containers.IteratorWithIndex).Next
		r.value = NativeToValue(i.(containers.IteratorWithIndex).Value).Value
		r.koi = NativeToValue(i.(containers.IteratorWithIndex).Index).Value
		r.begin = i.(containers.IteratorWithIndex).Begin
		r.first = i.(containers.IteratorWithIndex).First
		if t&REVERSE != 0 {
			r.prev = i.(containers.ReverseIteratorWithIndex).Prev
			r.end = i.(containers.ReverseIteratorWithIndex).End
			r.last = i.(containers.ReverseIteratorWithIndex).Last
		}
	} else {
		r.next = i.(containers.IteratorWithKey).Next
		r.value = NativeToValue(i.(containers.IteratorWithKey).Value).Value
		r.koi = NativeToValue(i.(containers.IteratorWithKey).Key).Value
		r.begin = i.(containers.IteratorWithKey).Begin
		r.first = i.(containers.IteratorWithKey).First
		if t&REVERSE != 0 {
			r.prev = i.(containers.ReverseIteratorWithKey).Prev
			r.end = i.(containers.ReverseIteratorWithKey).End
			r.last = i.(containers.ReverseIteratorWithKey).Last
		}
	}
	return r
}

//var ( // functions to reveal original underlaying gods container type
//	arraylistFn        = func(l *listCnt) lists.List { return (*l).list.(arraylist.List) }
//	singlylinkedlistFn = func(l *listCnt) lists.List { return singlylinkedlist.List((*l).list) }
//	doublylinkedlistFn = func(l *listCnt) lists.List { return doublylinkedlist.List((*l).list) }
//
//	hashsetFn = func(s *setCnt) sets.Set { return hashset.Set((*s).Set) }
//	treesetFn = func(s *setCnt) sets.Set { return treeset.Set((*s).Set) }
//
//	hashbidimapFn = func(m *mapCnt) maps.Map { return hashbidimap.Map((*m).Map) }
//	hashmapFn     = func(m *mapCnt) maps.Map { return hashmap.Map((*m).Map) }
//	treebidimapFn = func(m *mapCnt) maps.Map { return treebidimap.Map((*m).Map) }
//	treemapFn     = func(m *mapCnt) maps.Map { return treemap.Map((*m).Map) }
//
//	arraystackFn      = func(s *stackCnt) stacks.Stack { return arraystack.Stack((*s).Stack) }
//	linkedliststackFn = func(s *stackCnt) stacks.Stack { return linkedliststack.Stack((*s).Stack) }
//
//	binaryheapFn   = func(t *treeCnt) trees.Tree { return binaryheap.Tree((*t).Tree) }
//	redblacktreeFn = func(t *treeCnt) trees.Tree { return redblacktree.Tree((*t).Tree) }
//)

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

// CUSTOMIZED CONTAINER INTERFACES AND IMPLEMENTATIONS
//
//
//
//    type Iterator
//        func (iterator *Iterator) Begin()
//        func (iterator *Iterator) End()
//        func (iterator *Iterator) First() bool
//        func (iterator *Iterator) Index() int
//        func (iterator *Iterator) Last() bool
//        func (iterator *Iterator) Next() bool
//        func (iterator *Iterator) Prev() bool
//        func (iterator *Iterator) Value() interface{}
//    type List
//        func New() *List
//        func (list *List) Add(values ...interface{})
//        func (list *List) All(f func(index int, value interface{}) bool) bool
//        func (list *List) Any(f func(index int, value interface{}) bool) bool
//        func (list *List) Clear()
//        func (list *List) Contains(values ...interface{}) bool
//        func (list *List) Each(f func(index int, value interface{}))
//        func (list *List) Empty() bool
//        func (list *List) Find(f func(index int, value interface{}) bool) (int, interface{})
//        func (list *List) Get(index int) (interface{}, bool)
//        func (list *List) Insert(index int, values ...interface{})
//        func (list *List) Iterator() Iterator
//        func (list *List) Map(f func(index int, value interface{}) interface{}) *List
//        func (list *List) Remove(index int)
//        func (list *List) Select(f func(index int, value interface{}) bool) *List
//        func (list *List) Size() int
//        func (list *List) Sort(comparator utils.Comparator)
//        func (list *List) String() string
//        func (list *List) Swap(i, j int)
//        func (list *List) Values() []interface{}
//
// listCnt INTERFACE
type List interface {
	// List
	Get(index int) (Value, bool)
	Remove(index int)
	Add(values ...Value)
	Contains(values ...Value) bool
	Sort(comparator Comparator)
	Swap(index1, index2 int)
	Insert(index int, values ...Value)
}

// listCnt IMPLEMENTATION
//
// type of the container needs to be embedded in the encapsulating struct.
// Agiledocs own Container, List, Enumerable, Iterator and Comparator
// interfaces will be implemented in the structs methods to encapsulate all
// those typeless interface types into nice little values.
//
// the list interface, as well as all container interfaces will be implemented
// by one instance of a god container type, that will be instanciated and
// embedded by the constructor. Depending on gods constructor type, the
// returned type is different for each container type, but allways implements a
// more gegenral interface, like map, set, stack... which will be used as the
// fields name to store it in. that list/map/set/stack/tree implementation then
// also implements a couple of container interfaces, including container
// itself.
type ListValue interface {
	List
	Container
	Enumerable
	Iterator
}
type listCnt struct {
	CntType
	list lists.List
	Container
	Enumerable
	Iterator
}

// VALUE INTERFACE
func (l listCnt) ContType() CntType    { return l.CntType }
func (l listCnt) Slice() []interface{} { return valSlice(l.Container.Values()).Interfaces() }

// LIST INTERFACE
func (l *listCnt) Get(i int) (Value, bool) {
	v, ok := (*l).Get(i)
	return v.(Value), ok
}
func (l *listCnt) Remove(index int) {}
func (l *listCnt) Add(v ...Value) {
	is := valSlice(v).Interfaces()
	(*l).list.Add(is)
}
func (l listCnt) Contains(v ...Value) bool  { return l.list.Contains(valSlice(v).Interfaces()) }
func (l *listCnt) Sort(c Comparator)        { (*l).list.Sort(c.Convert()) }
func (l *listCnt) Swap(index1, index2 int)  { (*l).list.Swap(index1, index2) }
func (l *listCnt) Insert(i int, v ...Value) { (*l).list.Insert(i, valSlice(v).Interfaces()) }

// CONTAINER INTERFACE
func (l listCnt) Empty() bool      { return l.list.Empty() }
func (l listCnt) Size() int        { return l.list.Size() }
func (l *listCnt) Clear()          { (*l).list.Clear() }
func (l *listCnt) Values() []Value { return interfaceSlice((*l).list.Values()).Values() }

// LIST-CNT CONSTRUCTOR
// the listCnt constructor only needs to know the dedicated type of the listCnt
// container to instanciate
func newListContainer(t CntType) ValContainer {
	var l = listCnt{}
	switch t {
	case LIST_ARRAY:
		// list.List, implementing enum, container and returning
		// iterator from its method
		al := *arraylist.New()
		l.CntType = t
		l.list = &al
		l.Container = wrapContainer(&al)
		l.Enumerable = wrapEnum(t, &al)
		l.Iterator = wrapIterator(t, (&al).Iterator())
	case LIST_SINGLE:
		sl := *singlylinkedlist.New()
		l.CntType = t
		l.list = &sl
		l.Enumerable = wrapEnum(t, &sl)
		l.Iterator = wrapIterator(t, (&sl).Iterator())
	case LIST_DOUBLE:
		dl := *doublylinkedlist.New()
		l.CntType = t
		l.list = &dl
		l.Enumerable = wrapEnum(t, &dl)
		l.Iterator = wrapIterator(t, (&dl).Iterator())
	}
	return ValContainer(&l)
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
	Container
	Enumerable
	Iterator
}

type MapValue interface {
	Map
	Enumerable
	Iterator
}

func (m *mapCnt) ContType() CntType    { return m.CntType }
func (l *mapCnt) Slice() []interface{} { return (*l).Map.Values() }
func (m *mapCnt) Values() []Value      { return interfaceSlice((*m).Map.Values()).Values() }
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

// CONTAINER INTERFACE
func (m mapCnt) Empty() bool { return m.Map.Empty() }
func (m mapCnt) Size() int   { return m.Map.Size() }
func (m *mapCnt) Clear()     { (*m).Map.Clear() }

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
func newMapContainer(t CntType, c ...Comparator) ValContainer {
	var m mapCnt
	switch t {
	case MAP_HASH:
		am := *hashmap.New()
		m.CntType = t
		m.Map = &am
		m.Container = wrapContainer(&am)
		m.Enumerable = wrapEnum(t, &am)
	case MAP_HASHBIDI:
		hbm := *hashbidimap.New()
		m.CntType = t
		m.Map = &hbm
		m.Container = wrapContainer(&hbm)
		m.Enumerable = wrapEnum(t, &hbm)
	case MAP_TREE:
		tm := *treemap.NewWith(c[0].Convert())
		m.CntType = t
		m.Map = &tm
		m.Container = wrapContainer(&tm)
		m.Enumerable = wrapEnum(t, &tm)
	case MAP_TREEBIDI:
		tbm := *treebidimap.NewWith(c[0].Convert(), c[1].Convert())
		m.CntType = t
		m.Map = &tbm
		m.Container = wrapContainer(&tbm)
		m.Enumerable = wrapEnum(t, &tbm)
	}
	return ValContainer(&m)
}

// SET INTERFACE
type Set interface {
	Add(elements ...Value)
	Remove(elements ...Value)
	Contains(elements ...Value) bool

	Container
}
type SetValue interface {
	Set
	Enumerable
	Iterator
}

// SET IMPLEMENTATION
type setCnt struct {
	CntType
	sets.Set
	Container
	Enumerable
	Iterator
}

func (s *setCnt) ContType() CntType    { return s.CntType }
func (s *setCnt) Values() []Value      { return interfaceSlice((*s).Set.Values()).Values() }
func (l *setCnt) Slice() []interface{} { return (*l).Set.Values() }
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

// CONTAINER INTERFACE
func (s setCnt) Empty() bool { return s.Set.Empty() }
func (s setCnt) Size() int   { return s.Set.Size() }
func (s *setCnt) Clear()     { (*s).Set.Clear() }

// SET CONSTRUCTOR
// the treeset needs a comparator closure
func newSetContainer(t CntType, c ...Comparator) ValContainer {
	var s setCnt
	switch t {
	case SET_HASH:
		hs := *hashset.New()
		s.CntType = t
		s.Set = &hs
		s.Container = wrapContainer(&hs)
		s.Enumerable = wrapEnum(t, &hs)
	case SET_TREE:
		ts := *treeset.NewWith(c[0].Convert())
		s.CntType = t
		s.Set = &ts
		s.Container = wrapContainer(&ts)
		s.Enumerable = wrapEnum(t, &ts)
		s.Iterator = wrapIterator(t, (&ts).Iterator())
	}
	return ValContainer(&s)
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
	Container
	Enumerable
	Iterator
}
type StackValue interface {
	Stack
	Enumerable
	Iterator
}

func (s *stackCnt) ContType() CntType    { return s.CntType }
func (l *stackCnt) Slice() []interface{} { return (*l).Stack.Values() }

//func (s *stackCnt) Values() []Value {
//	return interfaceSlice((*s).Stack.Values()).Values()
//}
func (s *stackCnt) Peek() (Value, bool) {
	v, ok := (*s).Stack.Peek()
	return v.(Value), ok
}
func (s *stackCnt) Pop() (Value, bool) {
	v, ok := (*s).Stack.Pop()
	return v.(Value), ok
}
func (s *stackCnt) Push(v Value) { (*s).Stack.Push(v) }

// CONTAINER INTERFACE
func (s stackCnt) Empty() bool      { return s.Stack.Empty() }
func (s stackCnt) Size() int        { return s.Stack.Size() }
func (s *stackCnt) Clear()          { (*s).Stack.Clear() }
func (s *stackCnt) Values() []Value { return interfaceSlice((*s).Stack.Values()).Values() }

// STACK CONSTRUCTOR
func newStackContainer(t CntType) ValContainer {
	var s stackCnt
	switch t {
	case STACK_ARRAY:
		as := *arraystack.New()
		s.CntType = t
		s.Stack = &as
		s.Container = wrapContainer(&as)
		s.Enumerable = wrapEnum(t, &as)
		s.Iterator = wrapIterator(t, &as)
	case STACK_LINKED:
		ll := *linkedliststack.New()
		s.CntType = t
		s.Stack = &ll
		s.Container = wrapContainer(&ll)
		s.Enumerable = wrapEnum(t, &ll)
		s.Iterator = wrapIterator(t, &ll)
	}
	return ValContainer(&s)
}

// TREE INTERFACE
type Tree interface {
	Container
}
type TreeValue interface {
	Tree
	Enumerable
	Iterator
}

// TREE IMPLEMENTATION
type treeCnt struct {
	CntType
	trees.Tree
	Container
	Enumerable
	Iterator
}

func (t *treeCnt) ContType() CntType    { return (*t).CntType }
func (l *treeCnt) Slice() []interface{} { return (*l).Tree.Values() }
func (t *treeCnt) Values() []Value {
	return interfaceSlice((*t).Tree.Values()).Values()
}

// CONTAINER INTERFACE
func (t treeCnt) Empty() bool { return t.Tree.Empty() }
func (t treeCnt) Size() int   { return t.Tree.Size() }
func (t *treeCnt) Clear()     { (*t).Tree.Clear() }

// TREE CONSTRUCTOR
func newTreeContainer(ct CntType, c ...Comparator) *treeCnt {
	var t treeCnt
	switch ct {
	case TREE_REDBLACK:
		tr := *redblacktree.NewWith(c[0].Convert())
		t.CntType = ct
		t.Tree = &tr
		t.Container = wrapContainer(&tr)
		t.Enumerable = wrapEnum(ct, &tr)
		t.Iterator = wrapIterator(ct, &tr)
	case TREE_BINHEAP:
		bm := *binaryheap.NewWith(c[0].Convert())
		t.CntType = ct
		t.Tree = &bm
		t.Container = wrapContainer(&bm)
		t.Enumerable = wrapEnum(ct, &bm)
		t.Iterator = wrapIterator(ct, &bm)
	}
	return &t
}
