package agiledoc

import (
	"math/big"
)

//go:generate stringer -type Number
type Number int64

func (n Number) Flag() *big.Int {
	return new(big.Int).SetInt64(int64(n))
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
	EMPTY   ValueType = 0
	BOOL    ValueType = 1 << iota
	UINT              //	    ← 2	    <<  2
	INTEGER           //	    ← 4	    <<  3
	BYTES             //	    ← 8	    <<  4
	TEXT              //	    ← 16    <<  5
	// FLOAT/PAIR TYPES
	// encoded as *big.RAT
	FLOAT    //		    ← 32    <<  6
	RATIONAL //		    ← 64    <<  7
	PAIR     //		    ← 128   <<  8
	// COLLECTION TYPES
	FLAG   // *big.Int	    ← 256   <<  9
	LIST   // *arraylist.List   ← 512   <<  9
	STACK  // *arraystack.Stack ← 1024  <<  9
	TABLE  // *hashbidimap.Map  ← 2048  << 10
	MATRIX // *arraylist.List   ← 4096  << 11
	SET    // *treeset.Set	    ← 8192  << 12
	MAP    // *maps.Map	    ← 16384  << 12
	//	    ← 32768  << 12

	//////////// BIT FLAG SETS /////////////
	/////////////////
	// SEMANTIC SETS:
	BITWISE = UINT | FLAG | BOOL             // bitwise operable
	NUMERIC = BOOL | FLAG | UINT | INTEGER | // arithmetic operable…
		RATIONAL | FLOAT // …possibly not enumerable
	SYMBOLIC = BYTES | TEXT // syntacticly operable
	////////////////////////////////
	//// NUMERIC (SEMANTIC SUBSETS):
	///
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
	REAL = FLOAT | RATIONAL // both get stored as natural
	// nunber pair, quotient might be irrational nevertheless.
	NATURAL = BOOL | UINT | FLAG | INTEGER // sign relevant for arrithmetic
	// use only
	////////////////////
	//// SYNTACTIC SETS:
	TERMINAL = BOOL | FLAG | UINT | // ← flat base types
		INTEGER | BYTES | TEXT
	TUPLE = FLOAT | RATIONAL | PAIR // ← two parted tyoes (FLOAT's get
	// represented as ratio)
	COLLECTED = LIST | STACK | SET | // ← (possibly nested) collections
		TABLE | MATRIX | FLAG //  (FLAG is implemented as big.Int but…
	// …handled like a list of bools)

	// convienience mask with all bits set for bitwise operations
	MASK = (1 << 16) - 1
)

//go:generate stringer -type BoolType
type BoolType uint8

const (
	NATIVE     BoolType = 0         // returns bool: true/false
	LISTED     BoolType = 1 << iota // returns []bool
	SIGNED                          // returns  int: -1, 0, +1
	BIT_FLAG                        // returns BitFlag
	UINT_FLAG                       // returns uint Flag
	VAL_TYPE                        // returns ValueType Flag
	TOKEN_TYPE                      // returns TokenType Flag
	NODE_TYPE                       // returns NodeType Flag
)

//defined by runes, stored as uint32
//go:generate stringer -type Reserved
type Reserved uint32 // ← identical with rune type

const (
	BLANK          Reserved = ' '  //  //
	TAB                     = '	'  //
	LINEBREAK               = '\n' //  //
	FULL_STOP               = '.'
	QUESTION_MARK           = '?'
	ATTENTION_MARK          = '!'
	COMMA                   = ','
	COLON                   = ':'
	SEMICOLON               = ';'
	HYPHEN                  = '-'
	ELLIPSIS                = '…'
)
