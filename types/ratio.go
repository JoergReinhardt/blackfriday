package types

/////////////////////////////////////////////////////////////////////////
// RATIONAL
type Ratio ratio

func (r Ratio) Eval() Evaluable   { return r }
func (r Ratio) Serialize() []byte { return []byte(ratio(r).String()) }
func (r Ratio) String() string    { return r().String() }
func (r Ratio) Type() ValueType   { return RATIONAL }

func (r Ratio) Integer() Integer { return ratPool.Get().(Ratio).SetValue(r).Integer() }
func (r Ratio) Float() Float     { return wrap(r()).(Float) }

func (r Ratio) SetKey(v Evaluable) Ratio   { return wrap(pairPool.Get().(ratio).Ratio()).(Ratio) }
func (r Ratio) SetValue(v Evaluable) Ratio { return wrap(pairPool.Get().(ratio).Ratio()).(Ratio) }
func (r Ratio) SetNum(v Evaluable) Ratio   { return r.SetValue(v) }
func (r Ratio) SetDenom(v Evaluable) Ratio { return r.SetKey(v) }
func (r Ratio) SetFrac(a, b Evaluable) Ratio {
	ret := wrap(pairPool.Get().(ratio).Ratio()).(Ratio)
	ret.SetKey(a)
	ret.SetValue(b)
	return ret
}
func (r Ratio) Abs(v Ratio) Ratio {
	return wrap(ratPool.Get().(ratio).abs(v())).(Ratio)
}
