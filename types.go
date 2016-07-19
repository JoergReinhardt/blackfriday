// TYPE SYSTEM
//
// after lots of experimenting, I decided the best way is to keep the interface
// as simple as possible. One Method interfaces for the win but you might want
// to add one just to get shure. The main reason to reimplement a typesystem on
// top of gos existing type system, is due to the fact, that the reflection
// capabilitys are far to complex in features. far to complicated to use and
// don't perform that well either.
//
// That makes a type marker that performs well in comparisions, the essence of
// a type. The Value needs to be returnable in a generalized form, which is
// perfomed by the Value() mehode. Every value can additionaly be returned as
// its contained native type, but those methods differ by signature and can't
// be defined without the use of an empty interface. Everything else is up for
// grabs,
package agiledoc

import (
	// "fmt"
	"github.com/emirpasic/gods/containers"
	"math/big"
)

// VALUE INTERFACE
// exposes a type flag and it's content embedded in an instance of a value
// type.
type Val interface {
	Type() ValType
	Value() Val
}

// KEYVAL INTERFACE
// variable interface combines a value with an identifiyer
type Var interface {
	Val
	Key() string
}

// CONTAINER INTERFACE (extends the gods container interface)
// interface to conceal god container empty interface values behind the Val
// interface, that provides a type function to inspect the nature of the contained value without using reflection.. Since containers themselves and there conained values both
// implement the val interface, all container types are fully recursive. The Values returned from and passed to mapped containers and sets, also implement the Var interface, KeyVal identitys are taken as map keys to map their values on to.
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

// the base types are kept as simple as possible.
//
// NIL
//
// while returning an error and/or boolean when values don't exist, or turn out
// not to be converteable is one way of doing it, I think having a nil type to
// return instead, makes things much easyer and there are lots of additional
// reasons to have one anyway.
//
// FLAG
//
// since the type field is the essential difference to ordinary values, I
// decided to have bitflags as a type, to make the type acsess-, compare- and
// usable.
//
// INTEGER & FLOAT
//
// parsing values embedded in text and perform calculations on them is the main
// goal of agiledocument, which makesnumbers are the essence of the agile
// document. they are implemented by math librarys big.int / big.float types,
// since those allready implement all nescessary type conversions, parsing from
// string included and they are highly optimized to perform well. flag type is
// another big.Int, since it also implements bitwise boolean operations.
//
// BYTE
//
// the source will be provided by blackfriday in form of a byte slice
//
// strings make the input handable as text.
//go:generate -command stringer -type ValType
const (
	NIL  ValType = 0
	FLAG ValType = 1 << iota
	INTEGER
	FLOAT
	BYTE
	STRING
	KEYVAL    // type to hold key value pairs like variables and parameters
	CONTAINER // holds lists, sets, stacks, maps and trees
)

type ValType uint

// TYPES
type ( // are kept as close to the original types as possible
	empty    struct{}
	flagVal  struct{ *big.Int } // all big based types are enveloped
	intVal   struct{ *big.Int } // by strings to encapsulate the pointer
	floatVal struct{ *big.Float }
	byteVal  []byte
	strVal   string
	keyVal   struct { // variable, map member and parameter type
		id string
		Val
	}
	cntVal struct {
		CntType
		containers.Container
	}
)

// INTERFACE METHODS
// methods that share a name, need to be implemented once per receiving type
func (empty) Type() ValType    { return NIL }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (floatVal) Type() ValType { return FLOAT }
func (byteVal) Type() ValType  { return BYTE }
func (strVal) Type() ValType   { return STRING }
func (keyVal) Type() ValType   { return KEYVAL }
func (cntVal) Type() ValType   { return CONTAINER }

func (v empty) Value() Val    { return empty{} }
func (v flagVal) Value() Val  { return v }
func (v intVal) Value() Val   { return v }
func (v floatVal) Value() Val { return v }
func (v byteVal) Value() Val  { return v }
func (v strVal) Value() Val   { return v }
func (v keyVal) Value() Val   { return v }
func (v cntVal) Value() Val   { return v }

// native type return functions
func (v flagVal) Uint() uint          { return uint(v.Int.Uint64()) }
func (v flagVal) Flag() *big.Int      { return v.Int }
func (v intVal) BigInt() *big.Int     { return v.Int }
func (v floatVal) BigFlt() *big.Float { return v.Float }
func (v byteVal) Byte() []byte        { return v }
func (v strVal) String() string       { return string(v) }

// if the native tyoe is allready known at the time of initialization,
// reflection can be omitted.
func NewTypedVal(t ValType, i interface{}) Val {
	var v Val
	switch t {
	case NIL:
		v = empty{}
	case FLAG:
		v = flagVal{big.NewInt(int64(i.(uint)))}
	case INTEGER:
		v = intVal{big.NewInt(int64(i.(int)))}
	case FLOAT:
		v = floatVal{big.NewFloat(i.(float64))}
	case BYTE:
		v = byteVal(i.([]byte))
	case STRING:
		v = strVal(i.(string))
	}
	return v
}

// arbitratry values will be performed to the appropriate type, or an empty
// value will be returned.
func NewVal(i interface{}) Val {
	var v Val
	switch i.(type) {
	case uint, uint8, uint16, uint32, uint64:
		v = NewTypedVal(FLAG, i)
	case int, int8, int16, int32, int64, *big.Int:
		v = NewTypedVal(INTEGER, i)
	case float32, float64, *big.Float:
		v = NewTypedVal(FLOAT, i)
	case []byte:
		v = NewTypedVal(BYTE, i)
	case string:
		v = NewTypedVal(STRING, i)
	}
	return v
}
