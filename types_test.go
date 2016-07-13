package agiledoc

import (
	"testing"
)

var vals = struct{ []interface{} }{
	var v = []interface{
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
      }
}

func TestTypes(t *testing.T) {
	l := []Value{}
	for _, v := range vals {
		v := newUntypedValue{v}
		l = append(l, v)
	}
	(*t).Log(l)
}
