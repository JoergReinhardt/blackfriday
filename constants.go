package agiledoc

import (
	"math/big"
)

//go:generate stringer -type Number
type Number int64

func (n Number) Flag() *big.Int {
	return newVal().bigInt().SetInt64(int64(n))
}

const (
	NEGATIVE Number = -1
	ZERO     Number = 0

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

	TEEN     Number = TEN
	TY              = TEN
	TWEN            = 2 * TY
	THIR            = 3 * TY
	HUNDRED         = 10 * TEN
	THOUSEND        = 10 * HUNDRED
	MILLION         = 1000 * THOUSEND
	BILLION         = 1000 * MILLION
	TRILLION        = 1000 * BILLION
)

//go:generate stringer -type ValueType
type ValueType uint16

func (v ValueType) Uint() uint { return uint(v) }

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
	TEXT
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

	//// BIT FLAG SETS ////
	//// SEMANTIC SETS:
	SYMBOLIC = BYTES | TEXT
	NUMERIC  = BOOL | FLAG | UINT | INTEGER |
		RATIONAL | FLOAT
	// NUMERIC SUBSETS:
	NATURAL = BOOL | UINT | FLAG | INTEGER // sign for arrithmetics
	REAL    = FLOAT | RATIONAL             // both get stored as natural
	// nunber pair, quotient might be irrational nevertheless.
	//
	// When flattend to a List from a Map, Numerator is kept Denominator
	// discarded, since Denominator is often identical for all values in a
	// list:
	//
	// e.g. probabilitys get calculated by dividing occurrences per item,
	// by the number of all occurrences counted on any item, in which case
	// all probabilitys share a common denominator (sum of occurences
	// counted), while carrying individual numerators (number of
	// positive/negative occurences involving a particular item)
	//
	//// SYNTACTIC SETS
	TERMINAL = BOOL | FLAG | UINT | // ← flat base types
		INTEGER | BYTES | TEXT
	TUPLE     = FLOAT | RATIONAL | PAIR // ← two parted tyoes (FLOAT's get represented as ratio)
	COLLECTED = LIST | STACK | SET |    // ← (possibly nested) collections
		TABLE | MATRIX | FLAG //  (FLAG is implemented as big.Int but…
	// …handled like a list of bools)

	// convienient for bitwise operations
	MASK = (1 << 16) - 1
)
