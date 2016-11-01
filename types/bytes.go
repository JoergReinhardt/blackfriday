package types

import (
//"math/big"
)

/////////////////////////////////////////////////////////////////////////
// BYTES
type Bytes val

func (b Bytes) Eval() Evaluable { return b }

// the string representation is provided by serializing the integer to a slice
// of bytes and converting that to a string, that way preserving all contained
// information. Lower order types are supposed to be stored in a more
// appropriate internal type, like Bool, or Integer, and otherwise need to be
// re-parsed to regain arithmetic, or boolean functionality.
func (b Bytes) String() string { return string(b().Bytes()) }

// if serialized, the string representation is converted to a slice of bytes,
// in order to not use any valid information. In case a 'lower' type was stored
// by the Bytes instance, it must br reparsed at a later point to convert it to
// the appropriate internal type
func (b Bytes) Serialize() []byte { return []byte(b.String()) }
func (b Bytes) Type() ValueType   { return BYTES }

func (b Bytes) Bit(n int) uint {
	return b().Bit(n)
}
func (b Bytes) BitLen() int {
	return b().BitLen()
}
func (b Bytes) Bytes() Bytes {
	//the returned byte slice neds to be represented by a big Int, which is
	//provided by the modules public Value funcitons slice
	return Value(b.Serialize()).(val).Bytes()
}
func (b Bytes) SetBytes(x Bytes) Bytes {
	// parameter instance will be reused
	defer discardInt(x())
	// returns a big Int, which only needs to be enclosed in a fresh
	// closure, provided by the modules wrap function.
	return wrap(b().SetBytes(x.Serialize())).(val).Bytes()
}
func (b Bytes) SetBytesNative(x []byte) Bytes {
	// the wrapper encloses the returned big Int in a fresh closure for
	// return.
	return wrap(b().SetBytes(x)).(val).Bytes()
}
func (b Bytes) AppendBytes(x Bytes) Bytes {
	defer discardInt(x(), b())
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(append(b.Serialize(), x.Serialize()...)).(val).Bytes()
}
func (b Bytes) AppendBytesNative(x []byte) Bytes {
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(append(b.Serialize(), x...)).(val).Bytes()
}

// set text is allmost identical to set bytes, since text is stored in the same
// way as a byte slice and features a Serialize method just like it, since its
// an implementation of an Evaluable just like the Bytes type.
func (b Bytes) SetText(x Text) Bytes {
	// parameter instance will be reused
	defer discardInt(x(), b())
	// the wrapper encloses the big Int, representing the string
	// representation of the passed Text instance in a fresh closure for
	// return.
	return wrap(b().SetBytes([]byte(x.String()))).(val).Bytes()
}
func (b Bytes) SetTextNative(x string) Bytes {
	// after returning the new instance, the old one is designated for
	// reuse.
	defer discardInt(b())
	// b is set to a native string by replacing it vit a new value instance
	return Value(x).(val).Bytes()
}
func (b Bytes) AppendText(x Text) Bytes {
	defer discardInt(x())
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(b.String() + x.String()).(val).Bytes()
}
func (b Bytes) AppendTextNative(x string) Bytes {
	// since big Ints Append returns a byte slice, we need to allocate a
	// complete new instance of an Evaluable using Value
	return Value(b.String() + x).(val).Bytes()
}
