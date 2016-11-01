package types

import (
	"math/big"
)

/////////////////////////////////////////////////////////////////////////
// FLOAT
type Float ratio

func (f Float) Eval() Evaluable   { return f }
func (f Float) Serialize() []byte { return []byte(f.String()) }
func (f Float) String() string    { return f().FloatString(10) }
func (f Float) Type() ValueType   { return FLOAT }

func (f Float) Float() Evaluable {
	f = wrap(floatPool.Get().(*big.Float)).(Float)
	defer discardRat(f())
	return wrap(f().Set(f()))
}
func (f Float) Ratio() Evaluable {
	f = wrap(floatPool.Get().(*big.Rat)).(Float)
	defer discardRat(f())
	return wrap(f().Set(f()))
}
