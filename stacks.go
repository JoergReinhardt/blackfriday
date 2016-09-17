package agiledoc

import (
	as "github.com/emirpasic/gods/stacks/arraystack"
	ls "github.com/emirpasic/gods/stacks/linkedliststack"
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
//// STACKS ////
//////////////
//// ITERABLE STACK ////
// wraps the array-stack
func (l UnorderedStack) Eval() Evaluable  { return Value(l) }
func (l UnorderedStack) Type() ValueType  { return LIST }
func (l UnorderedStack) Size() int        { return l().Size() }
func (l UnorderedStack) Empty() bool      { return l().Empty() }
func (l UnorderedStack) Clear() Collected { l().Clear(); return l }
func (l UnorderedStack) AddInterface(v ...interface{}) UnorderedStack {
	var retval = l()
	for _, val := range v {
		(*retval).Push(val)
	}
	return UnorderedStack(func() *as.Stack { return retval })
}
func (l UnorderedStack) Add(v ...Evaluable) UnorderedStack {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Push(value)
	}
	return UnorderedStack(func() *as.Stack { return retval })
}
func (l UnorderedStack) Pop() (Evaluable, bool, UnorderedStack) {
	v, ok := l().Pop()
	return Value(v), ok, l
}
func (l UnorderedStack) Peek() (Evaluable, bool, UnorderedStack) {
	v, ok := l().Peek()
	return Value(v), ok, l
}
func (l UnorderedStack) Push(v Evaluable) UnorderedStack { l().Push(v); return l }
func (l UnorderedStack) Interfaces() []interface{} {
	return l().Values()
}

func (l UnorderedStack) Values() []Evaluable {
	return valueSlice(l.Values())
}

func (l UnorderedStack) Serialize() []byte {
	// allocate return byte slice, so it can be enclosed by the parameter
	// function.
	var retval []byte

	// parameter function to pass on to internal each methode:
	for index, value := range l.Values() {
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
	return retval
}

//// ITERABLE STACK ////
// wraps the linked-list-stack
// use serialization as string format base
func (l UnorderedStack) String() string { return string(l.Serialize()) }

func (l IterableStack) Eval() Evaluable                { return Value(l) }
func (l IterableStack) Type() ValueType                { return LIST }
func (l IterableStack) Size() int                      { return l().Size() }
func (l IterableStack) Empty() bool                    { return l().Empty() }
func (l IterableStack) Clear() Collected               { l().Clear(); return l }
func (l IterableStack) Push(v Evaluable) IterableStack { l().Push(v); return l }
func (l IterableStack) Pop() (Evaluable, bool, IterableStack) {
	v, ok := l().Pop()
	return Value(v), ok, l
}
func (l IterableStack) Peek() (Evaluable, bool, IterableStack) {
	v, ok := l().Peek()
	return Value(v), ok, l
}
func (l IterableStack) AddInterface(v ...interface{}) IterableStack {
	var retval = l()
	for _, val := range v {
		val := val
		(*retval).Push(val)
	}
	return IterableStack(func() *ls.Stack { return retval })
}
func (l IterableStack) Add(v ...Evaluable) IterableStack {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Push(value)
	}
	return IterableStack(func() *ls.Stack { return retval })
}
func (l IterableStack) Interfaces() []interface{} {
	return l().Values()
}

func (l IterableStack) Values() []Evaluable {
	var retval []Evaluable
	// parameter function to convert slice of interfaces to slice of
	// values once.
	for _, val := range l().Values() {
		val := Value(val)
		retval = append(retval, val)
	}
	return retval
}

func (l IterableStack) Serialize() []byte {
	// allocate return byte slice, so it can be enclosed by the parameter
	// function.
	var retval []byte

	// parameter function to pass on to internal each methode:
	for index, value := range l().Values() {
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
	return retval
}

// use serialization as string format base
func (l IterableStack) String() string { return string(l.Serialize()) }
func (l IterableStack) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
