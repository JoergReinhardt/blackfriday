package types

/////////////////////////////////////////////////////////////////////////
// RATIONAL
type Ratio ratio

func (r Ratio) Eval() Evaluable   { return r }
func (r Ratio) Serialize() []byte { return []byte(ratio(r).String()) }
func (r Ratio) String() string    { return r().String() }
func (r Ratio) Type() ValueType   { return RATIONAL }
