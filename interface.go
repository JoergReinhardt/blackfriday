package agiledoc

type Evaluable interface {
	Type() ValueType   // return designated dynamic type
	Eval() Evaluable   // produce a value from contained data
	Serialize() []byte // return evaluated content in a byte slice
	String() string    // return serialized evaluated content
}
type Tupled interface {
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

///////////////////////////////////////////////////
//// FLAT COLLECTIONS OF ELEMENTS
///
// Ranked is a list of values addressed by index.
type Listed interface {
	Collected
	Add(...Evaluable) Listed
	Remove(int) Listed
	//Join() Mapped // join indices and values to slice of pairs
	//Split() (indices Listed, values Listed)
}

// is a list of values accessed sequencialy
type Stacked interface {
	Collected
	Add(...Evaluable) Stacked
	Push(Evaluable) Stacked
	Pop() (Evaluable, bool, Stacked)
	Peek() (Evaluable, bool)
}

// FLAGGED INTERFACE
// defines common methods of  all types that can be represented as boolean
// value and encodes as big.Int:using its bitwise methods.
//
//  - bool, []bool [true,false]
//  - int, []int [-2,0,1]
//  - uint, []uint [boolean bitwise]
type Flagged interface {
	//
	// Bool(BoolType) Evaluable (return values type will match the BITWISE set)
	//
	// bool method returns the contained bool values in sequence encoded
	// according to the passed BoolType, which is one of either NATIVE |
	// LISTED | SIGNED | BIT-FLAG | UINT_FLAG | VAL_TYPE | TOKEN_TYPE |
	// NODE_TYPE, all of which encode the contained list of booleans in
	// different ways.
	Bool(BoolType) Evaluable
	//
	// Shift(boolVal uint, digit int)
	//
	// shift either one, or zero digit (according to passed uint) by int
	// digits, to the left, if int is positive, or to the right if int
	// turns out to have a negative negative.
	//  0	,   0	→   val 0 << 0 (not legal)
	//  1	,   0	→   val 1 << 0 (not legal)
	//  0	,   1	→   val 0 << 1	→ Value = 000-
	//  1	,   1	→   val 1 << 1	→ Value = 0001
	//  0	,   2	→   val 0 << 2	→ Value = 00-0
	//  1	,   2	→   val 1 << 2	→ Value = 0010
	//  0	,   3	→   val 0 << 3	→ Value = 0-00
	//  1	,   3	→   val 1 << 3	→ Value = 0100
	Shift(boolVal uint, digitToSet int) Flagged
	//
	// Add(...Evaluable) !!! Expected to be part of BITWISE set !!!
	// adds booleans at there appropriate index. Current position is based
	// on the bit length of the contained flag, or derrived by the most
	// signifficants bits index within the passed value, If particular bit
	// happens to be set allready. it keeps its state, once set (strictly
	// additive)
	Add(...Evaluable) Flagged
	//
	// Set(...Evaluable) !!! Expected to be part of BITWISE set !!!
	//
	// Remove()
	// sets bit at passed index position to false/0/-1 (according to the
	// implementations bool type
	Remove(int) Flagged
	//
	// Match()
	// determine if passed flaggable value is within the set of contained
	// values
	Match(Flagged) bool
}

///////////////////////////////////////////////////
//// ELEMENT COLLECTIONS MAPPED ON TO KEYS
// Mapped is a collection, with elements mapped on to keys
type Mapped interface {
	Collected
	Add(...Evaluable) Mapped
	Put(k, v Evaluable) Mapped
	Remove(Evaluable) Mapped
	Keys() []Evaluable
}

// DeDublicated is a Set of unique keys with one or more values mapped up on.
type DeDublicated interface {
	Collected
	Add(...Evaluable) DeDublicated
	Remove(int) DeDublicated
	Contains(v ...Evaluable) bool
}

///////////////////////////////////////////////////
//// ELEMENT COLLECTIONS OF PREDEFINED DIMENSION
// pre-defined dimensions make a collection tabular
type Tabular interface {
	Listed
	Shape() []int // [len(x-vector/row),len(y-vector/column),len(z-vector),...]
	Dim() int     // == len(Shape)
}

// matrix containing numeric elements exclusively
type NumericTabular interface {
	Listed
	Element(x int, y int) Evaluable
	Column(i int) Listed
	Row(i int) Listed
}

// table containing elements of symbolic value
type SymbolicTabular interface {
	abs() []Mapped
	Element(Evaluable, Evaluable) Evaluable
	Column(Evaluable) Mapped
	Row(Evaluable) Mapped
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
// parameters and return values encapsulated in pair, or Value instances
// respectively. All enumerables take two parameters, so one of the types
// implementing the pair interface is expected as a Parameter. The Key of that
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
