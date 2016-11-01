package types

import (
	"fmt"
	"testing"
)

var bytesTests = []struct {
	a     []byte
	b     []byte
	ex    string
	opStr string
	op    func(Evaluable, Evaluable) string
}{
	{[]byte("a"), []byte("b"), "ab", "AppendBytes", func(a, b Evaluable) string { return a.(Bytes).AppendBytes(b.(Bytes)).String() }},
	{[]byte(""), []byte(""), "0", "BitLen", func(a, b Evaluable) string { return fmt.Sprint(a.(Bytes).BitLen()) }},
}

func testBytesFunc(t *testing.T, a Evaluable, b Evaluable, ex string, opStr string, op func(a, b Evaluable) string) {
	if op(a, b) != ex {
		(*t).Fail()
		(*t).Log("failed operation: " + fmt.Sprint(opStr) +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(op(a, b)) +
			" expected: " + fmt.Sprint(ex))
	} else {
		(*t).Log("passed operation: " + fmt.Sprint(opStr) +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(op(a, b)) +
			" expected: " + fmt.Sprint(ex))
	}
}

func TestBytes(t *testing.T) {
	for _, test := range bytesTests {
		a := Value(test.a).(val).Bytes()
		b := Value(test.b).(val).Bytes()
		ex := test.ex
		opStr := test.opStr
		op := test.op

		testBytesFunc(t, a, b, ex, opStr, op)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
