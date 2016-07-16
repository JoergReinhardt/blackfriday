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
	i "github.com/emirpasic/gods/containers"
	l "github.com/emirpasic/gods/lists"
	ar "github.com/emirpasic/gods/lists/arraylist"
	do "github.com/emirpasic/gods/lists/doublylinkedlist"
	si "github.com/emirpasic/gods/lists/singlylinkedlist"
	m "github.com/emirpasic/gods/maps"
	hbm "github.com/emirpasic/gods/maps/hashbidimap"
	hm "github.com/emirpasic/gods/maps/hashmap"
	// tbm "github.com/emirpasic/gods/maps/treebidimap"
	// tm "github.com/emirpasic/gods/maps/treemap"
	//s "github.com/emirpasic/gods/sets"
	sh "github.com/emirpasic/gods/sets/hashset"
	//st "github.com/emirpasic/gods/sets/treeset"
	//s "github.com/emirpasic/gods/stacks"
	sa "github.com/emirpasic/gods/stacks/arraystack"
	sl "github.com/emirpasic/gods/stacks/linkedliststack"
	//t "github.com/emirpasic/gods/trees"
	//bh "github.com/emirpasic/gods/trees/binaryheap"
	//rbt "github.com/emirpasic/gods/trees/redblacktree"
	"math/big"
)

// VALUE INTERFACE
type Val interface {
	Type() ValType
	Value() Val
}

// CONTAINER INTERFACE
// interface to conceal god container values interface nature behind the val interface
type Container interface {
	// i.cont
	Empty() bool
	Size() int
	Clear()
	Values() []Val
}

// CONTAINER IMPLEMENTATION
// wraps the container in a struct
type cont struct {
	contType
	i.Container
}

// replace god containers value method through a method that encapsulates the
// empty value.
func (c cont) Values() (r []Val) {
	r = []Val{}
	for _, v := range c.Values() {
		v := v.(Val)
		r = append(r, v)
	}
	return r
}
func (c *cont) Add(v ...Val) {
	t := (*c).contType
	switch {
	case t&lists != 0:
		c.Container.(l.List).Add(v)
	case t&maps != 0:
		for _, val := range v {
			v := val.(parm)
			(*c).Container.(m.Map).Put(v.Key(), v.Value())
		}
	case t&stacks != 0:
		// c.Container.(sa.Stack).Push(v)
	case t&sets != 0:
		//(c.Container).(st.Set).Add(v)
		//	case t&trees != 0:
		//		if t&binheap != 0 {
		//			(*c).Container.(t.Heap).Push(v)
		//		} else {
		//			(*c).Container.(t.Tree).Put(v)
		//		}
	}
}

// return an empty container of a given type
func newContainer(t contType) Container {
	switch t {
	case array:
		return &cont{t, ar.New()}
	case double:
		return &cont{t, do.New()}
	case single:
		return &cont{t, si.New()}
	case hashbidi:
		return &cont{t, hbm.New()}
	case hash:
		return &cont{t, hm.New()}
	case treebidi:
		// return &cont{t,tbm.New()}
	case treemap:
		// return &cont{t,tm.New()}
	case hashset:
		return &cont{t, sh.New()}
	case treeset:
		// return &cont{t,st.New()}
	case arraystack:
		return &cont{t, sa.New()}
	case linkedstack:
		return &cont{t, sl.New()}
	case binheap:
		// return &cont{t,bh.New()}
	case redblack:
		// return &cont{t,rbt.New()}
	}
	return &cont{}
}

// container type marks the type of container taken from the god library
type contType uint16

const (
	array contType = 1 << iota
	single
	double
	hash
	hashbidi
	tree
	treebidi
	treemap
	hashset
	treeset
	arraystack
	linkedstack
	binheap
	redblack

	lists  = array | single | double
	maps   = hash | hashbidi | tree | treebidi
	sets   = hashset | treeset
	stacks = arraystack | linkedstack
	trees  = binheap | redblack
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
)

type ValType uint

// TYPES
type ( // are kept as close to the original types as possible
	emptyVal struct{}           // emptyValue
	flagVal  struct{ *big.Int } // all big based types are enveloped
	intVal   struct{ *big.Int } // by strings to encapsulate the pointer
	floatVal struct{ *big.Float }
	byteVal  []byte
	strVal   string
)

// funtion types that implement the interface
type (
	typeFunc func(Val) ValType
	valFunc  func(Val) Val
)

// INTERFACE METHODS
// methods that share a name, need to be implemented once per receiving type
func (emptyVal) Type() ValType { return NIL }
func (flagVal) Type() ValType  { return FLAG }
func (intVal) Type() ValType   { return INTEGER }
func (floatVal) Type() ValType { return FLOAT }
func (byteVal) Type() ValType  { return BYTE }
func (strVal) Type() ValType   { return STRING }

func (v emptyVal) Value() Val { return v }
func (v flagVal) Value() Val  { return v }
func (v intVal) Value() Val   { return v }
func (v floatVal) Value() Val { return v }
func (v byteVal) Value() Val  { return v }
func (v strVal) Value() Val   { return v }

// tyoed return functions return values of generic type, each is implemented by
// at least the types containing that native type, and all that can be
// converted to it.
func (emptyVal) Empty() emptyVal   { return emptyVal{} }
func (v flagVal) Flag() *big.Int   { return v.Int }
func (v intVal) Integer() *big.Int { return v.Int }
func (v floatVal) Flt() *big.Float { return v.Float }
func (v byteVal) Byte() []byte     { return v }
func (v strVal) String() string    { return string(v) }

// if the native tyoe is allready known at the time of initialization,
// reflection can be omitted.
func NewTypedVal(t ValType, i interface{}) Val {
	var v Val
	switch t {
	case NIL:
		v = emptyVal{}
	case FLAG:
		v = flagVal{big.NewInt(int64(i.(int)))}
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
