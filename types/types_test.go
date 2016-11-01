package types

import (
	"fmt"
	"testing"
)

var boolTests = []struct {
	x   bool
	y   bool
	exp string
	op  string
}{
	{true, true, "true", "and"},
	{false, false, "false", "and"},
	{true, false, "false", "and"},
	{false, true, "false", "and"},
	{true, true, "true", "andnot"},
	{false, false, "false", "andnot"},
	{true, false, "false", "andnot"},
	{false, true, "false", "andnot"},
	{true, true, "false", "not"},
	{false, false, "false", "not"},
	{true, false, "false", "not"},
	{false, true, "false", "not"},
	{false, true, "true", "or"},
	{true, true, "true", "or"},
	{true, false, "true", "or"},
	{false, false, "false", "or"},
	{false, true, "true", "xor"},
	{true, true, "false", "xor"},
	{true, false, "true", "xor"},
	{false, false, "false", "xor"},
}

func TestBool(t *testing.T) {
	for _, te := range boolTests {
		x := Value(te.x).(Bool)
		y := Value(te.y).(Bool)
		exp := te.exp
		switch te.op {
		case "and":
			if x.And(x, y).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + te.op +
					" x: " + fmt.Sprint(te.x) +
					" y: " + fmt.Sprint(te.y) +
					" expected: " + te.exp)
			}
		case "andnot":
			if x.And(x, y).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + te.op +
					" x: " + fmt.Sprint(te.x) +
					" y: " + fmt.Sprint(te.y) +
					" expected: " + te.exp)
			}
		case "or":
			if x.Or(x, y).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + te.op +
					" x: " + fmt.Sprint(te.x) +
					" y: " + fmt.Sprint(te.y) +
					" expected: " + te.exp)
			}
		case "xor":
			if x.Xor(x, y).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + te.op +
					" x: " + fmt.Sprint(te.x) +
					" y: " + fmt.Sprint(te.y) +
					" expected: " + te.exp)
			}
		case "not":
			if x.Not(x).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + te.op +
					" x: " + fmt.Sprint(te.x) +
					" y: " + fmt.Sprint(te.y) +
					" expected: " + te.exp)
			}
		}
		(*t).Log(te)
	}
}

// anonymois test slug type
var integerTests = []struct {
	a     int64
	b     int64
	exp   int64
	opStr string
	op    integerTestFunc
}{
	{1, 2, 3, "add", func(a, b Integer) Integer { return a.Add(a, b) }},
	{3, 2, 1, "sub", func(a, b Integer) Integer { return a.Sub(a, b) }},
	{3, 3, 0, "cmp", func(a, b Integer) Integer { return Value(a.Cmp(b)).(val).Integer() }},
	{3, 22, -1, "cmp", func(a, b Integer) Integer { return Value(a.Cmp(b)).(val).Integer() }},
	{3, 2, 1, "div", func(a, b Integer) Integer { return a.Div(a, b) }},
	{10, 2, 5, "div", func(a, b Integer) Integer { return a.Div(a, b) }},
	{3, 2, 1, "divmod", func(a, b Integer) Integer { return a.Div(a, b) }},
	{4, 2, 2, "divmod", func(a, b Integer) Integer { return a.Div(a, b) }},
}

// function type, to be passed in the test slug
type integerTestFunc func(a, b Integer) Integer

// testInteger runs the actual testing by calling op with its parameters and
// reporting Log and Error Messages through the passed testing.T
func testInteger(t *testing.T, a Integer, b Integer, exp int64, opStr string, op integerTestFunc) {
	// run operation with parameters a and b
	res := op(a, b).Int64()
	if res != exp { // if result of the operation differs from the expected result. fail.
		(*t).Fail()
		(*t).Log("failed operation: " + opStr +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(res) +
			" expected: " + fmt.Sprint(exp))
	} else { // if result and epexted integer are identical, return log
		(*t).Log("passed operation: " + opStr +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(res) +
			" expected: " + fmt.Sprint(exp))
	}
}

// TestInteger iterates over all test slugs, calling testInteger for each slug
// and passing it's parameters.
func TestInteger(t *testing.T) {

	for n, test := range integerTests {

		n := n
		a := Value(test.a).(val).Integer()
		b := Value(test.b).(val).Integer()
		exp := test.exp
		opStr := test.opStr
		op := test.op

		t.Log(fmt.Sprintf("Test Nr. %d: ", n))

		testInteger(t, a, b, exp, opStr, op)

	}
}

var bytesTests = []struct {
	a  []byte
	b  []byte
	ex string
	op func(Evaluable, Evaluable) string
}{
	{[]byte("a"), []byte("b"), "ab", func(a, b Evaluable) string { return a.(Bytes).AppendBytes(b.(Bytes)).String() }},
	{[]byte(""), []byte(""), "0", func(a, b Evaluable) string { return fmt.Sprint(a.(Bytes).BitLen()) }},
}

func testBytesFunc(t *testing.T, a Evaluable, b Evaluable, ex string, op func(a, b Evaluable) string) {
	if op(a, b) != ex {
		(*t).Fail()
		(*t).Log("failed operation: " + fmt.Sprint(op) +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(op(a, b)) +
			" expected: " + fmt.Sprint(ex))
	} else {
		(*t).Log("passed operation: " + fmt.Sprint(op) +
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
		op := test.op

		testBytesFunc(t, a, b, ex, op)
	}
}
