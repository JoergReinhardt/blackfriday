package agiledoc

import (
// "github.com/emirpasic/gods/containers"
//"github.com/emirpasic/gods/lists"
// 	"github.com/emirpasic/gods/lists/arraylist"
// 	"github.com/emirpasic/gods/lists/doublylinkedlist"
// 	"github.com/emirpasic/gods/lists/singlylinkedlist"
//"github.com/emirpasic/gods/maps"
// 	 "github.com/emirpasic/gods/maps/hashbidimap"
// 	 "github.com/emirpasic/gods/maps/hashmap"
// 	"github.com/emirpasic/gods/maps/treebidimap"
// 	"github.com/emirpasic/gods/maps/treemap"
//"github.com/emirpasic/gods/sets"
// 	"github.com/emirpasic/gods/sets/hashset"
// 	"github.com/emirpasic/gods/sets/treeset"
//"github.com/emirpasic/gods/stacks"
// 	"github.com/emirpasic/gods/stacks/arraystack"
// 	"github.com/emirpasic/gods/stacks/linkedliststack"
// "github.com/emirpasic/gods/trees"
// 	"github.com/emirpasic/gods/trees/binaryheap"
// 	"github.com/emirpasic/gods/trees/redblacktree"
// 	"github.com/emirpasic/gods/utils"
)

// KEYVAL INTERFACE
// variable interface combines a value with an identifiyer
type Var interface {
	Val
	Key() string
}

func (v keyVal) Key() string { return v().id }

// CONTAINER INTERFACE (extends the gods container interface)
// interface to conceal god container empty interface values behind the Val
// interface, that provides a type function to inspect the nature of the
// contained value without using reflection.. Since containers themselves and
// there conained values both implement the val interface, all container types
// are fully recursive. The Values returned from and passed to mapped
// containers and sets, also implement the Var interface, KeyVal identitys are
// taken as map keys to map their values on to.
type Container interface {
	Val // contVal imlements the value interface
	Empty() bool
	Size() int
	Clear()
	Values() []Val // might be keyed vars/params
	ContType() CntType
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
	ARRAYLIST CntType = 1 << iota
	SINGLELIST
	DOUBLELIST
	///////////// sets
	HASHSET
	TREESET
	///////////// stacks
	LINKEDSTACK
	ARRAYSTACK
	///////////// maps
	HASHMAP
	HASHBIDIMAP
	TREEMAP
	TREEBIDIMAP
	///////////// trees
	REDBLACK
	BINHEAP

	// sets of containers that share a more specific interface than
	// gods/containers and have other method signatures in common
	LISTS  = ARRAYLIST | SINGLELIST | DOUBLELIST
	SETS   = HASHSET | TREESET
	STACKS = ARRAYSTACK | LINKEDSTACK
	MAPS   = HASHMAP | HASHBIDIMAP | TREEMAP | TREEBIDIMAP
	TREES  = BINHEAP | REDBLACK

	INDEXED = LISTS | STACKS
	MAPPED  = MAPS | TREES | SETS

	ADD = LISTS | SETS

	REVERSEABLE = DOUBLELIST | HASHBIDIMAP | TREEBIDIMAP
)

func (v cntVal) ContType() CntType { return v().CntType }
func (v *cntVal) Values() (r []Val) {
	r = []Val{}
	for _, val := range (*v).Values() {
		val := val.Value()
		r = append(r, val)
	}
	return r
}

// func (v *cntVal) Add(vals ...Val) {
// 	switch {
// 	case v.ContType()&LISTS != 0:
// 		for _, val := range vals {
// 			val := val.Value()
// 			(*v).Container.(lists.List).Add(val)
// 		}
// 	case v.ContType()&SETS != 0:
// 		for _, val := range vals {
// 			val := val.Value()
// 			(*v).Container.(sets.Set).Add(val)
// 		}
// 	case v.ContType()&STACKS != 0:
// 		for _, val := range vals {
// 			val := val.Value()
// 			(*v).Container.(stacks.Stack).Push(val)
// 		}
// 	case v.ContType()&MAPS != 0:
// 		for i, kv := range vals {
// 			i := i
// 			kv := kv.Value()
// 			var key string
// 			// if the value is of type keyVal, it will implement
// 			// the Var interface and have a key method, which is
// 			// used to set the map key, otherwise take string
// 			// representation of the intefer index of current
// 			// element.
// 			if kv.Type()&KEYVAL != 0 {
// 				key = kv.(Var).Key()
// 			} else {
// 				key = string(i)
// 			}
// 			(*v).Container.(maps.Map).Put(key, kv.Value())
// 		}
// 		// 	case v.ContType()&REDBLACK != 0:
// 		// 	case v.ContType()&BINHEAP != 0:
// 	}
// }

// ITERATOR INTERFACES (wrapper for Emir Pasic's gods Iterators, that return
// and take Val instances instead of empty interfaces)
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
