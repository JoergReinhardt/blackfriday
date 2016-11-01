package agiledoc

import (
	"fmt"
	"testing"
)

var tests = []struct {
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
	for _, te := range tests {
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

var tests = []struct {
	x   int64
	y   int64
	exp string
	op  string
}{
	{1, 10, "11", "add"},
}
func testfunc(x Evaluable, y Evaluable, exp string, op string, t *testing.T) {
		  if x(),Add(x(),y()).String() != exp {
				(*t).Fail()
				(*t).Log("failed op: " + op +
					" x: " + fmt.Sprint(x) +
					" y: " + fmt.Sprint(y) +
					" expected: " + exp)
		  } else {
		(*t).Log(
				(*t).Log("passed op: " + op +
					" x: " + fmt.Sprint(x) +
					" y: " + fmt.Sprint(y) +
					" expected: " + exp))
		  }
}

func TestInteger(t *testing.T) {
	for _, te := range tests {
		x := Value(te.x).(Integer)
		y := Value(te.y).(Integer)
		exp := te.exp
		switch te.op {
		case "add":
		  testfunc(x,y,exp,op,t)
		}
	}
}
