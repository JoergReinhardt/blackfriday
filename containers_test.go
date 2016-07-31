package agiledoc

import (
	"testing"
)

func TestContainers(t *testing.T) {

	c := wrapContainer(LIST_ARRAY)

	for _, s := range Natives.Bytes {
		s := s
		v := NativeToValue(s)
		c.(List).Add(v)
	}

	(*t).Log(c.Size())
	(*t).Log(c.ContType())
	(*t).Log(c.Empty())
	(*t).Log(c.Values())
}
