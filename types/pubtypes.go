package types

import (
	"math/big"
)

/////////////////////////////////////////////////////////////////////////
/////////// PUBLIC IMPLEMENTAIONS OF EVALUABLES /////////////////////////
///////////////// BASED ON VAL, RAT & PAIR //////////////////////////////
/////////////////////////////////////////////////////////////////////////

// BASE VALUES IMPLEMENTING FUNCTIONAL TYPES
// these functional types need to implement the absVal interface to be suitable
// base types. If called, they return their contained value. A Method set
// defined on these funtional types, implements the absVal interface, by
// manipulating the returned content. Each can implement ia couple of dynamic
// types by defining further types based on it, while overwriting and/or
// completing the method set.
//////// EMPTY //////////////////////////////////////////////////////////
type ( // functional types that form the base of all value implementations
	// empty Value
	Empty func() struct{}

	// simple type
	val func() *big.Int

	// paired types
	ratio func() *big.Rat
	Pair  func() [2]Evaluable

	// collection types see collection,go
)

func (Empty) Type() ValueType   { return EMPTY }
func (e Empty) Eval() Evaluable { return Empty(func() struct{} { return struct{}{} }) }
func (Empty) Serialize() []byte { return []byte{0} }
func (e Empty) String() string  { return e.Type().String() }
