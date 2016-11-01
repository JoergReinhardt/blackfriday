package types

import (
	"fmt"
	"testing"
)

// boolean test slug type for table driven tests
var boolTests = []struct {
	a     bool         // native boolean
	b     bool         // native boolean
	exp   string       // expection is encoded as string
	opStr string       // name of the performed operation
	op    boolTestFunc // operation to performe
}{
	// boolean arrithmetic method tests
	{true, true, "true", "and", func(a, b Bool) string { return a.And(a, b).String() }},
	{false, false, "false", "and", func(a, b Bool) string { return a.And(a, b).String() }},
	{true, false, "false", "and", func(a, b Bool) string { return a.And(a, b).String() }},
	{false, true, "false", "and", func(a, b Bool) string { return a.And(a, b).String() }},
	{true, true, "true", "andnot", func(a, b Bool) string { return a.AndNot(a, b).String() }},
	{false, false, "false", "andnot", func(a, b Bool) string { return a.AndNot(a, b).String() }},
	{true, false, "false", "andnot", func(a, b Bool) string { return a.AndNot(a, b).String() }},
	{false, true, "false", "andnot", func(a, b Bool) string { return a.AndNot(a, b).String() }},
	{true, true, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, false, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{true, false, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, true, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, true, "true", "or", func(a, b Bool) string { return a.Or(a, b).String() }},
	{true, true, "true", "or", func(a, b Bool) string { return a.Or(a, b).String() }},
	{true, false, "true", "or", func(a, b Bool) string { return a.Or(a, b).String() }},
	{false, false, "false", "or", func(a, b Bool) string { return a.Or(a, b).String() }},
	{false, true, "true", "xor", func(a, b Bool) string { return a.Xor(a, b).String() }},
	{true, true, "false", "xor", func(a, b Bool) string { return a.Xor(a, b).String() }},
	{true, false, "true", "xor", func(a, b Bool) string { return a.Xor(a, b).String() }},
	{false, false, "false", "xor", func(a, b Bool) string { return a.Xor(a, b).String() }},
}

// the passed function type
type boolTestFunc func(a, b Bool) string

// the private testBool function runs the test and formats the output as
// Log/Fail message appropriately
func testBool(a Bool, b Bool, exp string, opStr string, op boolTestFunc, t *testing.T) {

	res := op(a, b)

	if fmt.Sprint(res) != exp {
		(*t).Fail()
		(*t).Log("failed op: " + opStr +
			" x: " + fmt.Sprint(a) +
			" y: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(res) +
			" expected: " + exp)
	} else {
		(*t).Log("failed op: " + opStr +
			" x: " + fmt.Sprint(a) +
			" y: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(res) +
			" expected: " + exp)
	}
}

// public TestBool iterates over all test slugs, calling the private test
// methode once per iteration, pasing the values from the testslug iterated
// over
func TestBool(t *testing.T) {

	for n, test := range boolTests {

		n := n
		a := Value(test.a).(Bool)
		b := Value(test.b).(Bool)
		exp := test.exp
		opStr := test.opStr
		op := test.op

		t.Log(fmt.Sprintf("Test Nr. %d: ", n))

		testBool(a, b, exp, opStr, op, t)

	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// anonymois test slug type
var integerTests = []struct {
	a     int64
	b     int64
	exp   int64
	opStr string
	op    integerTestFunc
}{
	{1, 2, 3, "add", func(a, b Integer) Integer { return a.Add(b) }},
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

/////////////////////////////////////////////////////////////////////////////////////////////////////

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

/////////////////////////////////////////////////////////////////////////////////////////////////////
