package agiledoc

import (
	"testing"
)

func TestContainers(t *testing.T) {
	cnt := NewContainer(LIST_ARRAY)
	vals := []Value{}

	for _, v := range Natives.Bytes {
		v := NativeToValue(v)
		vals = append(vals, v)
	}

	cnt.(List).Add(vals...)

	(*t).Log(cnt.Size())

	for _, v := range vals {
		v := v.Value()
		(*t).Log(v.String())
	}
}
