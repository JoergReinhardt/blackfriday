package types

import (
	"fmt"
	"testing"
)

// anonymois test slug type
var integerTests = []struct {
	a     int64
	b     int64
	exp   int64
	opStr string
	op    integerTestFunc
}{
	{1, 2, 3, "add", func(a, b Integer) Evaluable { return a.Add(b) }},
	{3, 2, 1, "sub", func(a, b Integer) Evaluable { return a.Sub(b) }},
	{3, 3, 0, "cmp", func(a, b Integer) Evaluable { return Value(a.Cmp(b)).(val).Integer() }},
	{3, 22, -1, "cmp", func(a, b Integer) Evaluable { return Value(a.Cmp(b)).(val).Integer() }},
	{3, 2, 1, "div", func(a, b Integer) Evaluable { return a.Div(b) }},
	{10, 2, 5, "div", func(a, b Integer) Evaluable { return a.Div(b) }},
	{11, 5, 1, "mod", func(a, b Integer) Evaluable { return a.Mod(b) }},
	{11, 5, 1, "modInverse", func(a, b Integer) Evaluable { return a.ModInverse(b) }},
	{11, 5, 1, "ModSqrt", func(a, b Integer) Evaluable { return a.ModSqrt(b) }},
	{11, 5, 55, "mul", func(a, b Integer) Evaluable { return a.Mul(b) }},
	{11, 5, 2, "quo", func(a, b Integer) Evaluable { return a.Quo(b) }},
	{11, 5, 2, "quoRem", func(a, b Integer) Evaluable { return a.Quo(b) }},
	{11, 5, 2, "rem", func(a, b Integer) Evaluable { return a.Quo(b) }},
	// {5, 0, 0, "rand", func(a, b Integer) Evaluable { return a.Rand() }},
}

// function type, to be passed in the test slug
type integerTestFunc func(a, b Integer) Evaluable

// testInteger runs the actual testing by calling op with its parameters and
// reporting Log and Error Messages through the passed testing.T
func testInteger(t *testing.T, a Integer, b Integer, exp int64, opStr string, op integerTestFunc) {
	// run operation with parameters a and b
	res := op(a, b).(Integer).Int64()
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
