package agiledoc

import (
	"testing"
)

func wrap(ct CntType, v [][]interface{}, t *testing.T) {

	cal := NewContainer(LIST_ARRAY).(*listCnt)
	st := NewVal("test addition")
	(*cal).Add(st)

	for o, s := range v {
		o := o
		s := s
		for _, i := range s {
			i := i
			v := NewVal(i)
			(*cal).Add(v)
		}
		(*t).Log("Nr.", o, "  ┄>  ", cal.Values())
	}
}

func TestContainers(t *testing.T) { wrap(LIST_ARRAY, InitValues(), t) }
