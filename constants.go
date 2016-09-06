package agiledoc

import (
	"math/big"
)

//go:generate stringer -type Number
type Number int64

func (n Number) Flag() *big.Int {
	return NewVal().BigInt().SetInt64(int64(n))
}

const (
	NEGATIVE_ONE Number = -1
	ZERO         Number = 0

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
	HUNDRED         = 10 * TEN
	THOUSEND        = 10 * HUNDRED
	MILLION         = 1000 * THOUSEND
	BILLION         = 1000 * MILLION
	TRILLION        = 1000 * BILLION
)

//go:generate stringer -type ValueType
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

	// SEMANTIC SETS
	NUMERIC  = FLAG | UINT | INTEGER | RATIONAL | FLOAT // int key
	SYMBOLIC = BYTES | STRING                           // map key

	// SUPER TYPES
	INT = BOOL | UINT | INTEGER | BYTES | STRING // [2]Value
	RAT = RATIONAL | FLOAT                       // [2]Value

	// COLLECTION TYPES (INCLUDES FLAG!)
	COLLECTED = FLAG | LIST | STACK |
		TABLE | MATRIX | SET // Collected

	// TYPE OF INDEX TO COLLECT BY
	// one might argue that surely floarts rationals and arguably even
	// booleans and flags are numeric by nature… this maps collection types
	// to the iteration operation there is to perform in a for loop (ether
	// they are a map, or slice/array)
	NUM_KEYS = LIST | STACK | FLAG | MATRIX | INTEGER | UINT
	SYM_KEYS = TABLE | SET | BOOL | FLAG | STRING | BYTES |
		FLOAT | RATIONAL | PAIR // takes and returns key/val pairs

	// convienient for biteise operations
	MAX_MASK = (1 << 16) - 1
)
