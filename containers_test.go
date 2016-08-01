package agiledoc

import (
	"testing"
)

func TestContainers(t *testing.T) {

	c := newValueContainer(LIST_ARRAY)

	for _, s := range Natives.Bytes {
		s := s
		v := NativeToValue(s).Value()
		t.Log(v)
		(c.(List)).Add(v)
	}

	(*t).Log(c.ContType())
	(*t).Log(c.Container().Size())
	(*t).Log(c.Container().Empty())
	(*t).Log(c.Container().Values())

	i := c.Iterator()
	t.Log(i)

	//	for ok := i.Next(); ok; {
	//		n := i.Value()
	//		t.Log(n)
	//	}
}
