package agiledoc

import (
	"testing"
)

var vals = struct{ v []interface{} }{
	[]interface{}{
		struct{}{},
		true,
		false,
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
		8,
		9,
		10,
		0.0,
		1.1,
		2.2,
		3.3,
		4.4,
		5.5,
		6.6,
		7.7,
		8.8,
		9.9,
		10.1,
		"zero",
		"one",
		"one",
		"two",
		"three",
		"four",
		"five",
		"six (you propbably guessed that by now!?)",
		"seven (up)",
		"eight",
		"nine",
	},
}

func TestTypes(t *testing.T) {
	l := []Value{}
	for _, v := range vals.v {
		v := NewVal(v)
		s := v.Type()
		(*t).Log(s.String())
		l = append(l, v)
	}
	(*t).Log("slice of values: "+" ", l)

	// generate a vector from the valueInt64s
	vec := NewVal(l)
	(*t).Log("vector Value generated from slice of values ", vec.Type().String(), vec)
	(*t).Log("evaluating vector value: ", vec.Eval())
}
