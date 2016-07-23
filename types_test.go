package agiledoc

import (
	"fmt"
	"testing"
)

var TestVals = []interface{}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9,
	"zero", "one", "two", "three", "four", "five",
	"six", "seven", "eight", "nine", []byte("allet jar nuescht so einfach"),
}

func TestTypes(t *testing.T) {
	var l []Val
	var out string = ""
	for _, v := range TestVals {
		v := v
		val := NewVal(v)
		l = append(l, val)
		out = out + " " + fmt.Sprint(val.(Val).Type())
	}
	t.Log(out)
	t.Log(l)
	for i, v := range l {
		v := v
		if v.Type() == FLOAT {
			n, e := v.(*floatVal).Rat.Float64()
			t.Log("Float Number ", i, ": ", n, " the TestVals exactnes is: ", e)
		}
	}
}
