package types

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
func (a ArrayStack) Eval() Evaluable                 { return evalCollection(a()) }
func (a ArrayStack) Type() ValueType                 { return STACK }
func (a ArrayStack) Size() int                       { return collectionSize(a()) }
func (a ArrayStack) Empty() bool                     { return emptyCollection(a()) }
func (a ArrayStack) Clear() Collected                { return clearCollection(a()) }
func (a ArrayStack) Push(v Evaluable) Stacked        { return pushToStack(a, v) }
func (a ArrayStack) Pop() (Evaluable, bool, Stacked) { return popFromStack(a) }
func (a ArrayStack) Peek() (Evaluable, bool)         { return peekOnStack(a) }

func (l ArrayStack) AddInterface(v ...interface{}) ArrayStack {
	var retval = l()
	for _, val := range v {
		(*retval).Push(val)
	}
	return ArrayStack(func() *as.Stack { return retval })
}
func (l ArrayStack) Add(v ...Evaluable) Stacked {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Push(value)
	}
	return ArrayStack(func() *as.Stack { return retval })
}
func (l ArrayStack) Interfaces() []interface{} {
	return l().Values()
}

func (l ArrayStack) Values() []Evaluable {
	return valueSlice(l.Values())
}
func (l ArrayStack) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
func (l ArrayStack) String() string { return string(l.Serialize()) }
func (l ArrayStack) Serialize() []byte {
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

func (l LinkedStack) Eval() Evaluable                 { return evalCollection(l()) }
func (l LinkedStack) Type() ValueType                 { return STACK }
func (l LinkedStack) Size() int                       { return collectionSize(l()) }
func (l LinkedStack) Empty() bool                     { return emptyCollection(l()) }
func (l LinkedStack) Clear() Collected                { return clearCollection(l()) }
func (l LinkedStack) Push(v Evaluable) Stacked        { return pushToStack(l, v) }
func (l LinkedStack) Pop() (Evaluable, bool, Stacked) { return popFromStack(l) }
func (l LinkedStack) Peek() (Evaluable, bool)         { return peekOnStack(l) }

func (l LinkedStack) AddInterface(v ...interface{}) Stacked {
	var retval = l()
	for _, val := range v {
		val := val
		(*retval).Push(val)
	}
	return LinkedStack(func() *ls.Stack { return retval })
}
func (l LinkedStack) Add(v ...Evaluable) Stacked {
	var retval = l()
	for _, value := range v {
		value := value
		(*retval).Push(value)
	}
	return LinkedStack(func() *ls.Stack { return retval })
}
func (l LinkedStack) Interfaces() []interface{} {
	return l().Values()
}

func (l LinkedStack) Values() []Evaluable {
	var retval []Evaluable
	// parameter function to convert slice of interfaces to slice of
	// values once.
	for _, val := range l().Values() {
		val := Value(val)
		retval = append(retval, val)
	}
	return retval
}

func (l LinkedStack) Serialize() []byte {
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
func (l LinkedStack) String() string { return string(l.Serialize()) }
func (l LinkedStack) Iter() Iterable {
	iter := l().Iterator()
	return IdxIterator{&iter}
}
