package agiledoc

//go:generate -comand stringer -type Number
type Number int64

const (
	ONE Number = 1 + iota
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	ELEVEN
	TWELVE
	THIRTEEN

	NEGATIVE_ONE Number = -1

	TEEN     Number = TEN
	HUNDRED         = 10 * TEN
	THOUSEND        = 10 * HUNDRED
	MILLION         = 1000 * THOUSEND
	BILLION         = 1000 * MILLION
	TRILLION        = 1000 * BILLION
)

//go:generate -command stringer -type ValueType
type ValueType uint16

// DYNAMIC TYPES
// All dynamic types are defined as function types. Methods are defined on
// those functional types, To implement the variable and all higher level
// interfaces. The typefunction itsef returns a variable of it's type.
const ( // NUMERIC TYPES
	// encoded as *big.Int
	EMPTY ValueType = 0
	BOOL  ValueType = 1 << iota
	UINT
	INTEGER
	BYTES
	STRING
	// FLOAT/PAIR TYPES
	// encoded as *big.RAT
	FLOAT
	RATIONAL
	PAIR
	// COLLECTION TYPES
	FLAG   // *big.Int
	LIST   // *arraylist.List
	STACK  // *arraystack.Stack
	TABLE  // *hashbidimap.Map
	MATRIX // *arraylist.List
	SET    // *treeset.Set
	// FUNCTIONAL SETS
	NUMERIC  = FLOAT | UINT | INTEGER // int key
	SYMBOLIC = BYTES | STRING         // map key
	PAIRED   = PAIR | RATIONAL        // [2]Value
	// collected types
	COLLECTED = FLAG | LIST | STACK |
		TABLE | MATRIX | SET // Collected
	INDEXED = LIST | STACK | FLAG |
		MATRIX | INTEGER | UINT
	MAPPED = TABLE | SET | BOOL | STRING |
		BYTES | FLOAT | RATIONAL | PAIR // takes and returns key/val pairs

	MASK = (1 << 16) - 1
)

type Value interface {
	Type() ValueType   // return designated dynamic type
	Eval() Value       // produce a value from contained data
	Serialize() []byte // return evaluated content in a byte slice
	String() string    // return serialized evaluated content
}

// COLLECTED META INTERFACE
// Collected is a super interface, that defines common functionality of all
// types that define collections of values
type Collected interface {
	Value
	Empty() bool
	Size() int
	Interfaces() []interface{}
	Values() []Value
	Iterator() Iterable
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
	Add(...Value) Flagged
	Remove(int) Flagged
	Clear() Flagged
}

// is a list of values accessed sequencialy
type Stacked interface {
	Collected
	Enumerable
	Add(...Value) Stacked
	Remove(int) Stacked
	Clear() Stacked
}

// Ranked is a list of values addressed by index.
type Ranked interface {
	Collected
	Add(...Value) Ranked
	Remove(int) Ranked
	Clear() Ranked
	RankedValues() []Pair
}

// Contained is intended to be tested against if it contains a certain value,
// or not.  (Set)
type Contained interface {
	Collected
	Contains(...Value) bool
	Add(...interface{}) Contained
	AddValue(...Value) Contained
	Remove(...Value) Contained
	Clear() Contained
}

// Map is a collection, with values mapped on to keys
type Mapped interface {
	Collected
	Add(...Pair) Mapped
	Remove(...Pair) Mapped
	Clear() Mapped
	Keys() []Value
	KeyValues() []Pair
}

// matrices and tables are tabular
type Tabular interface {
	Ranked
	Shape() []int // [len(x-vector/row),len(y-vector/column),len(z-vector),...]
	Dim() int     // == len(Shape)
}

// numeric matrix
type NumericMatrix interface {
	Ranked
	Element(x int, y int) Value
	Column(i int) Ranked
	Row(i int) Ranked
}

// symbolic table
type SymbolicTable interface {
	Mapped
	Element(Value, Value) Value
	Column(Value) Mapped
	Row(Value) Mapped
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
	Each(func(Value, Value)) Enumerable

	// key:int/value ← val:val ←|→ bool
	Any(func(Value, Value) bool) (Enumerable, bool)

	// key:int/value ← val:val ←|→ bool
	All(func(Value, Value) bool) (Enumerable, bool)

	// key:int/value ← val:val ←|→ pair(value (index|key), value)
	Find(func(Value, Value) bool) (Enumerable, Pair)
}

// Iterables provide a Rev methode, that returns a boolean to indicate wether
// or not the reversable methodes exist and an instance of Reverse so that the
// caller can call them.
type Iterable interface {
	Next() (bool, Iterable)
	Value() (Value, Iterable)
	Index() (int, Iterable)
	Begin() Iterable
	First() (bool, Iterable)
}
type Reverse interface {
	Iterable
	Prev() (bool, Reverse)
	End() Reverse
	Last() (bool, Reverse)
}
