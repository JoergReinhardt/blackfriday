package agiledoc

type Evaluable interface {
	Type() ValueType   // return designated dynamic type
	Eval() Evaluable   // produce a value from contained data
	Serialize() []byte // return evaluated content in a byte slice
	String() string    // return serialized evaluated content
}
type Paired interface {
	Index() int
	Key() Evaluable
	Value() Evaluable
}

// COLLECTED META INTERFACE
// Collected is a super interface, that defines common functionality of all
// types that define collections of values
type Collected interface {
	Evaluable
	Clear() Collected
	Empty() bool
	Size() int
	Interfaces() []interface{}
	Values() []Evaluable
}

// TYPED INTERFACES
// applys to types with common abitlitys. there are often several
// implementations, but one designated implementation, such functionality is
// best stored in and also expected at, by the user.
//
// FLAGGED INTERFACE
// defines common methods of  all types that can be represented as boolean
// value and encodes as big.Int:using its bitwise methods.
//
//  - bool, []bool [true,false]
//  - int, []int [-2,0,1]
//  - uint, []uint [boolean bitwise]
type Flagged interface {
	Match(Flagged) bool
	// shift either one, or zero uint by int digits, to the left, if int is
	// positive, or right if it's negative.
	Shift(uint, int) Flagged
	Add(...Evaluable) Flagged
	Remove(int) Flagged
}

// is a list of values accessed sequencialy
type Stacked interface {
	Collected
	Add(...Evaluable) Stacked
	Remove(int) Stacked
}

// Ranked is a list of values addressed by index.
type Listed interface {
	Collected
	Add(...Evaluable) Listed
	Remove(int) Listed
	Ordered() []pair
}

// Contained is intended to be tested against if it contains a certain value,
// or not.  (Set)
//type Appendable interface {
//	Collected
//	Contains(...Evaluator) bool
//	Add(...interface{}) Appendable
//	AddValue(...Evaluator) Appendable
//	Remove(...Evaluator) Appendable
//}

// Map is a collection, with values mapped on to keys
type Mapped interface {
	Collected
	Add(...pair) Mapped
	Remove(...pair) Mapped
	Keys() []Evaluable
	KeyValues() []pair
}

// matrices and tables are tabular
type Tabular interface {
	Listed
	Shape() []int // [len(x-vector/row),len(y-vector/column),len(z-vector),...]
	Dim() int     // == len(Shape)
}

// numeric matrix
type NumericMatrickable interface {
	Listed
	Element(x int, y int) Evaluable
	Column(i int) Listed
	Row(i int) Listed
}

// symbolic table
type SymbolicTabular interface {
	abs() []Mapped
	Element(Evaluable, Evaluable) Evaluable
	Column(Evaluable) Mapped
	Row(Evaluable) Mapped
}

// symbolic table
type NumericTabular interface {
	abs() []Listed
	Element(int, Evaluable) Evaluable
	Column(int) Listed
	Row(int) Listed
}

// ENUMERABLE & ITERATOR INTERFACES
//
// Enumerable & Iteraor wrap gods collection enumerables and iterators, that
// come in different types. Depending on being indexed, or mapped to keys and
// if they are reverseable or not, parameters and return values vary. To make
// gods collections more dynamic, parameters and return values implement the
// value interface.
//
// iterator type to call internaly will be determined dvnamicly based on its
// parameters and return values encapsulated in Pair, or Value instances
// respectively. All enumerables take two parameters, so one of the types
// implementing the Pair interface is expected as a Parameter. The Key of that
// pair will either be of numeric kind, which leads to the declaration of an
// enumerable with integer index, or of some other kind (most likely a string)
// and the declaration of a maped collection.
//
// enumerablesare stateful. To keep values safe from beeing mutated by multiple
// refering callers Enumerable. enum() replicates the list and returns the
// altered version of the list as return value after each mutation.
type Enumerable interface {
	// key:int/value ← val:val ←|→ empty
	Each(func(Evaluable, Evaluable)) Enumerable

	// key:int/value ← val:val ←|→ bool
	Any(func(Evaluable, Evaluable) bool) (bool, Enumerable)

	// key:int/value ← val:val ←|→ bool
	All(func(Evaluable, Evaluable) bool) (bool, Enumerable)

	// key:int/value ← val:val ←|→ pair(value (index|key), value)
	Find(func(Evaluable, Evaluable) bool) (pair, Enumerable)
}

// Iterables provide a Rev methode, that returns a boolean to indicate wether
// or not the reversable methodes exist and an instance of Reverse so that the
// caller can call them.
type Iterable interface {
	Next() bool
	Value() Evaluable
	Index() Integer
	Begin()
	First() bool
}
type Reverse interface {
	Iterable
	Prev() bool
	End()
	Last() bool
}
