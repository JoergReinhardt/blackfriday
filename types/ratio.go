package types

/////////////////////////////////////////////////////////////////////////
// RATIONAL
type Ratio ratio

func (r Ratio) Eval() Evaluable   { return r }
func (r Ratio) Serialize() []byte { return []byte(ratio(r).String()) }
func (r Ratio) String() string    { return r().String() }
func (r Ratio) Type() ValueType   { return RATIONAL }

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
func (r Ratio) Add(v Ratio) Ratio {
	return wrap(ratio(r).add(r(), v())).(Ratio)
}
func (r Ratio) Cmp(v Ratio) Ratio {
	return wrap(ratio(r).cmp(v())).(Ratio)
}
func (r Ratio) Neg(v Ratio) Ratio {
	return wrap(ratio(r).neg(v())).(Ratio)
}
func (r Ratio) Mul(v Ratio) Ratio {
	return wrap(ratio(r).mul(r(), v())).(Ratio)
}
func (r Ratio) Quo(v Ratio) Ratio {
	return wrap(ratio(r).quo(r(), v())).(Ratio)
}
func (r Ratio) Sub(v Ratio) Ratio {
	return wrap(ratio(r).sub(r(), v())).(Ratio)
}
