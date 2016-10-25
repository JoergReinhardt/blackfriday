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

var integerTests = []struct {
	a  int64
	b  int64
	ex int64
	op string
}{
	{1, 2, 3, "add"},
	{3, 2, 1, "sub"},
	{-2, 2, -2, "neg"},
	{3, 3, 0, "cmp"},
	{10, 5, 2, "div"},
	{10, 5, 50, "mul"},
	{12, 5, 2, "divmod"},
	{12, 5, 2, "mod"},
	{10, 3, 1000, "exp"},
	{0, 3, 3, "set"},
	{0, 3, 3, "setInt64"},
	{0, 3, 3, "setUint64"},
	{0, 3, 3, "setString"},
}

func TestInteger(t *testing.T) {
	for _, test := range integerTests {
		a := Integer(Value(test.a).(val))
		b := Integer(Value(test.b).(val))
		ex := test.ex
		op := test.op
		switch op {
		case "add":
			if a.Add(a, b)().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "sub":
			if a.Sub(a, b)().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "cmp":
			if a.Cmp(b) != int(ex) {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "neg":
			if a.Neg(b).Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "div":
			if a.Div(a, b)().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "divmod":
			_, mod := a.DivMod(a, b, Integer(Value(ex).(val)))
			if mod().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" got: " + fmt.Sprint(mod) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "mod":
			_, mod := a.DivMod(a, b, Integer(Value(ex).(val)))
			if mod().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" got: " + fmt.Sprint(mod) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "exp":
			res := a.Exp(a, b, Integer(Value(0).(val)))
			if res.Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex) +
					" got: " + fmt.Sprint(res.Int64()))
			}
		case "mul":
			if a.Mul(a, b)().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "set":
			if a.Set(b)().Int64() != ex {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "setUint64":
			if a.SetUint64(b().Uint64())().Uint64() != uint64(ex) {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		case "setString":
			str, _ := a.SetString(b.String(), 10)
			if str().String() != fmt.Sprint(ex) {
				(*t).Fail()
				(*t).Log("failed operation: " + test.op +
					" a: " + fmt.Sprint(test.a) +
					" b: " + fmt.Sprint(test.b) +
					" expected: " + fmt.Sprint(test.ex))
			}
		}
		(*t).Log(test)
	}
}
