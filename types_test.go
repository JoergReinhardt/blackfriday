package agiledoc

import (
	"testing"
)

const (
	TRUE  bool = true
	FALSE bool = false
)

var v Evaluable

func TestBoolInstanceiation(t *testing.T) {
	v = Value(TRUE)
	(*t).Log(v)

}
func TestBoolAnd(t *testing.T) {
	(*t).Log(v.Type())
	(*t).Log(v.(Bool).And(v.(Bool), Value(FALSE).(Bool)))
	(*t).Log(v.(Bool).And(v.(Bool), Value(TRUE).(Bool)))
}
func TestBoolAndNot(t *testing.T) {
	(*t).Log(v.Type())
	(*t).Log(v.(Bool).AndNot(v.(Bool), Value(FALSE).(Bool)))
	(*t).Log(v.(Bool).AndNot(v.(Bool), Value(TRUE).(Bool)))
}
