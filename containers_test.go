package agiledoc

import (
	"testing"
)

func TestContainers(t *testing.T) {

	c := NewContainer(LIST_ARRAY)

	for _, s := range Natives.Bytes {
		s := s
		v := NativeToValue(s)
		c.(List).Add(v)
	}

	v := NativeToValue(c.Values())

	(*t).Log(v.(vecVal).Size())
	(*t).Log(v.(vecVal).ContType())
	(*t).Log(v.(vecVal).Empty())
	(*t).Log(v.(vecVal).Values())
	(*t).Log(v.(vecVal).Values())
}
