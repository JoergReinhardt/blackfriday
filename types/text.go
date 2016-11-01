package types

import (
//"math/big"
)

//"sync"
/////////////////////////////////////////////////////////////////////////
// STRING
type Text val

func (s Text) Eval() Evaluable { return s }

// text stored in the enclosed int, is retrieved, by serializing to a  Byte
// slice representation.
func (s Text) Serialize() []byte { return s().Bytes() }

// the string method builds a string representation og the contained data, by
// serializing it to bytes and representing those as a string
func (s Text) String() string  { return string(s.Serialize()) }
func (s Text) Type() ValueType { return TEXT }

// set a pre-existing Text Instance to a Value represented by the internal
// Bytes type.
func (s Text) SetBytes(x Bytes) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// SetBytes takes a public Bytes instance and sets an existing text to
	// its serialization to a byte slice
	return Value(s().SetBytes(x.Serialize())).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the native byte
// slice.
func (s Text) SetBytesNative(x []byte) Text {
	// setBytesNatice sets s to a native go byte slice. converted to Text
	// via Value
	return Value(s().SetBytes(x)).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the internal
// Text type.
func (s Text) SetText(x Text) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// setBytes with the string returned by the value converted to bytes as
	// Parameter, finaly cinverted to Text via Value
	return Value(s().SetBytes([]byte(x.String()))).(val).Text()
}

// set a pre-existing Text Instance to a Value represented by the native string.
func (s Text) SetTextNative(x string) Text {
	// setBytes with the string returned by the value converted to bytes as
	// Parameter, finaly cinverted to Text via Value
	return Value(s().SetBytes([]byte(x))).(val).Text()
}

// Append an Instance of the internal Bytes Type to a preexisting Text Instance
func (s Text) AppendBytes(x Bytes) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// uses internal append funcrion and Serialize, which must be provided
	// by all evaluable, to concatenate on a byte base
	return Value(append(s.Serialize(), x.Serialize()...)).(val).Text()
}

// Append an Instance of a native byte slice to a preexisting Text Instance
func (s Text) AppendBytesNative(x []byte) Text {
	// uses internal append funcrion and Serialize to append a native byte
	// Slice with a given text by re-valuabling using Value, asserting the
	// intermediate val type and calling the Text() method on it.
	return Value(append(s.Serialize(), x...)).(val).Text()
}

// Append an Instance of the internal Text Type to a preexisting Text Instance
func (s Text) AppendText(x Text) Text {
	// parameter instance will be reused
	defer discardInt(x())
	// uses string concatenation to append a text provided as parameter to
	// a given Text instance
	return Value(s.String() + x.String()).(val).Text()
}

// Append an Instance of a native string to a preexisting Text Instance
func (s Text) AppendTextNative(x string) Text {
	// uses gos append function and iinternal String method provided by all
	// evaluables, to concatenate annative string  to the given Text using
	// string concatenation.
	return Value(s.String() + x).(val).Text()
}
