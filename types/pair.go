package types

import ()

/////////////////////////////////////////////////
/////// PAIR ////////////////////////////////////
/////////////////////////////////////////////////
func (b Pair) Eval() Evaluable { return Value(b) }

func (b Pair) Value() Evaluable { return b()[1].Eval() }

// a pair allways provides a key, which can be of any given base type
func (b Pair) Key() Evaluable { return b()[0].Eval() }

// Index() int
// returns the key of the element as native integger, if it turns out to be
// convertable, otherwise return a negative integer to indicate that the key is
// not convertable to a Number
func (b Pair) Index() Integer {
	var ret Integer
	if b.Key().Type()&SYMBOLIC != 0 {
		ret = Value(-1).(Integer) // negative → not set
	} else { // NUMERIC
		// if natural number, return as interger
		if b.Key().Type()&NATURAL != 0 {
			ret = b.Key().(Integer)
		}
		// if real number, return numerator as interger
		if b.Key().Type()&REAL != 0 {
			ret = b.Key().(ratio).Num()
		}
	}
	return ret
}
func (b Pair) SetKey(v Evaluable) Pair {
	return func() [2]Evaluable { return [2]Evaluable{v, b.Value()} }
}
func (b Pair) SetValue(v Evaluable) Pair {
	return func() [2]Evaluable { return [2]Evaluable{b.Key(), v} }
}
func (b Pair) SetBoth(k Evaluable, v Evaluable) Pair {
	return func() [2]Evaluable { return [2]Evaluable{k, v} }
}
func (p Pair) Serialize() []byte {
	var delim = []byte{}
	if p.Index()().Int64() == -1 {
		delim = []byte(": ")
	} else {
		delim = []byte(".) ")
	}
	return append(
		p()[0].Serialize(),
		append(
			delim,
			append(
				p()[1].Serialize(),
				[]byte("\n")...,
			)...,
		)...,
	)
}
func (b Pair) String() string  { return string(b.Serialize()) }
func (b Pair) Type() ValueType { return TUPLE }

// generate pair from evaluables
func pairFromValues(k, v Evaluable) (r Pair) { return r }
