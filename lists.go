package agiledoc

import (
	con "github.com/emirpasic/gods/containers"
	al "github.com/emirpasic/gods/lists/arraylist"
	//as "github.com/emirpasic/gods/stacks/arraystack"
	"math/big"
)

// lists and sublists of exactly two values length, are assumed to be either
// key/value, or index/value pairs of Pair Type, by the modules Eval function
// on first pass.
//
// All longer slices are flattened by evalCollection and refed into eval
// recursively. .  All conversions to Collected,  get instanciated as list
// type,to profit from the enumerable interface at flattening and conversion.
// COLLECTED IMPLEMENTING METHODS

////////////////////////////////////////////////////////////////////////////////////
//// LISTS ////
//////////////
//// ORDERED LIST ////
// wraps the array-list
func (l OrderedList) Eval() Evaluable  { return Value(l) }
func (l OrderedList) Type() ValueType  { return LIST }
func (l OrderedList) Size() int        { return l().Size() }
func (l OrderedList) Empty() bool      { return l().Empty() }
func (l OrderedList) Clear() Collected { l().Clear(); return l }
func (l OrderedList) AddInterface(v ...interface{}) OrderedList {
	var retval = l()
	(*retval).Add(v...)
	return OrderedList(func() *al.List { return retval })
}
func (l OrderedList) Add(v ...Evaluable) OrderedList {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Add(value)
	}
	return OrderedList(func() *al.List { return retval })
}
func (l OrderedList) Remove(i int) OrderedList {
	var retval = l()
	(*retval).Remove(i)
	return OrderedList(func() *al.List { return retval })
}

func (l OrderedList) RankedValues() []pair {
	var retval []pair
	var fn = func(index int, value interface{}) {
		i := Value(index)
		v := Value(value)
		// pass both values as paired parameter, will trigger eval to
		// produce a key/value tuple type
		retval = append(retval, Value(i, v).(pair))
	}
	l().Each(fn)
	return retval
}
func (l OrderedList) Interfaces() []interface{} {
	return l().Values()
}

func (l OrderedList) Values() []Evaluable {
	var retval []Evaluable
	// parameter function to convert slice of interfaces to slice of
	// values once.
	var fn = func(index int, value interface{}) {
		retval = append(retval, Value(l.Interfaces()))
	}
	// retrieve an iterator from collection and call it passing the
	// argument function, to append to the predefined slice
	(con.EnumerableWithIndex)(l()).Each(fn)
	return retval
}

func (l OrderedList) Serialize() []byte {
	// allocate return byte slice, so it can be enclosed by the parameter
	// function.
	var retval []byte

	// parameter function to pass on to internal each methode:
	var fn = func(index int, value interface{}) {
		i := Value(index).Serialize()
		v := Value(value).Serialize()

		// format each entry as one line with leading numeric index,
		// followed by a dot and blank character, the Value and a
		// newline character.
		retval = append(
			retval,
			append(
				i,
				append(
					[]byte(".) "),
					append(
						v,
						[]byte("\n")...,
					)...,
				)...,
			)...,
		)
	}
	// call function once per value, to format whole list
	l().Each(fn)
	return retval
}

// use serialization as string format base
func (l OrderedList) String() string { return string(l.Serialize()) }
func (l OrderedList) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
func (l OrderedList) Enum() Enumerable {
	var r IdxEnumerable = func() con.EnumerableWithIndex { return l() }
	return r
}

//////////////////////////////////////////////////////////////////////////
func (f BitFlag) Eval() Evaluable { return Value(f) }

// uses byte method of contained big int
func (f BitFlag) Serialize() []byte { return f().Bytes() }

// returns Flag converted to string on base two
func (f BitFlag) String() string { return f().Text(2) }

// returns pure type Flag
func (f BitFlag) Type() ValueType { return FLAG }

func (f BitFlag) Empty() bool {
	if f().Cmp(ZERO.Flag()) > 0 {
		return false
	} else {
		return true
	}
}
func (f BitFlag) Size() int { return f().BitLen() }
func (f BitFlag) Clear() Collected {
	var r *big.Int = f()
	r.Set(ZERO.Flag())
	return BitFlag(func() *big.Int { return r })
}
func (f BitFlag) Values() []Evaluable {
	return valueSlice(f.Values())
}
func (f BitFlag) Interfaces() []interface{} {
	var v []interface{}
	for _, val := range f().Bits() {
		val := val
		v = append(v, val)
	}
	return v
}
