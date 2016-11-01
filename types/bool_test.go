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
	{true, true, "true", "and", func(a, b Bool) string { return a.And(b).String() }},
	{false, false, "false", "and", func(a, b Bool) string { return a.And(b).String() }},
	{true, false, "false", "and", func(a, b Bool) string { return a.And(b).String() }},
	{false, true, "false", "and", func(a, b Bool) string { return a.And(b).String() }},
	{true, true, "true", "andnot", func(a, b Bool) string { return a.AndNot(b).String() }},
	{false, false, "false", "andnot", func(a, b Bool) string { return a.AndNot(b).String() }},
	{true, false, "false", "andnot", func(a, b Bool) string { return a.AndNot(b).String() }},
	{false, true, "false", "andnot", func(a, b Bool) string { return a.AndNot(b).String() }},
	{true, true, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, false, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{true, false, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, true, "false", "not", func(a, b Bool) string { return a.Not(b).String() }},
	{false, true, "true", "or", func(a, b Bool) string { return a.Or(b).String() }},
	{true, true, "true", "or", func(a, b Bool) string { return a.Or(b).String() }},
	{true, false, "true", "or", func(a, b Bool) string { return a.Or(b).String() }},
	{false, false, "false", "or", func(a, b Bool) string { return a.Or(b).String() }},
	{false, true, "true", "xor", func(a, b Bool) string { return a.Xor(b).String() }},
	{true, true, "false", "xor", func(a, b Bool) string { return a.Xor(b).String() }},
	{true, false, "true", "xor", func(a, b Bool) string { return a.Xor(b).String() }},
	{false, false, "false", "xor", func(a, b Bool) string { return a.Xor(b).String() }},
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
		(*t).Log("passed op: " + opStr +
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
