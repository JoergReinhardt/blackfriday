package types

import (
	con "github.com/emirpasic/gods/containers"
)

func (s HashSet) Eval() Evaluable                 { return evalCollection(s()) }
func (s HashSet) Type() ValueType                 { return SET }
func (s HashSet) Size() int                       { return collectionSize(s()) }
func (s HashSet) Empty() bool                     { return emptyCollection(s()) }
func (s HashSet) Clear() Collected                { return clearCollection(s()) }
func (s HashSet) Contains(v ...Evaluable) bool    { return setContains(s, v...) }
func (s HashSet) Add(v ...Evaluable) DeDublicated { return addToSet(s, v...) }
func (s HashSet) Remove(i int) DeDublicated       { return removeFromSet(s, i) }
func (s HashSet) Interfaces() []interface{}       { return interfacesFromSet(s) }
func (s HashSet) String() string                  { return s().String() }
func (s HashSet) Serialize() []byte               { return []byte(s().String()) }
func (s HashSet) Values() []Evaluable             { return valueSlice(s().Values()) }

func (s TreeSet) Eval() Evaluable                 { return evalCollection(s()) }
func (s TreeSet) Type() ValueType                 { return SET }
func (s TreeSet) Size() int                       { return collectionSize(s()) }
func (s TreeSet) Empty() bool                     { return emptyCollection(s()) }
func (s TreeSet) Clear() Collected                { return clearCollection(s()) }
func (s TreeSet) Contains(v ...Evaluable) bool    { return setContains(s, v...) }
func (s TreeSet) Add(v ...Evaluable) DeDublicated { return addToSet(s, v...) }
func (s TreeSet) Remove(i int) DeDublicated       { return removeFromSet(s, i) }
func (s TreeSet) Interfaces() []interface{}       { return interfacesFromSet(s) }
func (s TreeSet) String() string                  { return s().String() }
func (s TreeSet) Serialize() []byte               { return []byte(s().String()) }
func (s TreeSet) Values() []Evaluable             { return valueSlice(s().Values()) }
func (t TreeSet) Iter() Iterable {
	iter := t().Iterator()
	return IdxIterator{&iter}
}
func (t TreeSet) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return t() }
	return r
}
