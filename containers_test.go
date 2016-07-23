package agiledoc

import (
	"testing"
)

var ValList = []Val{}

func TestContainers(t *testing.T) {
	for _, v := range TestVals {
		val := NewVal(v)
		ValList = append(ValList, val)
	}

	t.Log("list of values: ", ValList)

	cal := NewContainer(LIST_ARRAY).(List)
	for _, v := range ValList {
		cal.(List).Add(v)
	}
	t.Log("list Container after adding Values: ", cal, cal.Values(), "Size: ", cal.Size())
}
