package agiledoc

import (
	"fmt"
	"testing"
)

var values = []interface{}{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9,
	"zero", "one", "two", "three", "four", "five",
	"six", "seven", "eight", "nine",
}

func TestTypes(t *testing.T) {
	var out string = ""
	for _, v := range values {
		v := v
		val := NewVal(v)
		out = out + " " + fmt.Sprint(val.(Val).Type())
	}
	t.Log(out)
}
