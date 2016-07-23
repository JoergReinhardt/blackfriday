package agiledoc

import (
	// "fmt"
	"testing"
)

func InitValues() [][]interface{} {

	return [][]interface{}{
		{false, true, false, true, false, true, false, true, false, true},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
		{"allet ja nuescht so einfach", "hier mit die viele Worte"},
		{"zero", "one", "two", "three", "four", "five",
			"six", "seven", "eight", "nine", "ten", "twothousendthreehundredfiftyseven"},
	}
}
func convertSliceOfSlices(s [][]interface{}, t *testing.T) (l []Container) {

	l = []Container{}

	for o, s := range InitValues() {

		o := o // outer loop counter
		s := s // slice of values of the same native type

		l = append(l, NewContainer(LIST_ARRAY))

		for i, v := range s {

			i := i // inner Loop counter
			v := v // single value
			val := NewVal(v).(Val)
			(*t).Log("\nvalue Nr.:", i, "\nType: ", val.Type(), "\nValue:", val.Value())
			l[o].(List).Add(val)
		}
	}
	(*t).Log("\n container   ┄>  ", l)
	return l
}

// pass the testing struct as parameter to the log function, to call it from
// within the called function.
func logValueSlice(vals []Val, t *testing.T) {

	for i, v := range vals {

		i := i       // counter
		v := v.(Val) // value

		switch v.Type() {
		case FLAG:
			n := v.(*intVal).Uint64()
			(*t).Log("Nr.", i, ", of Type ", v.Type(), "┄>", n)
		case INTEGER:
			n := v.(*intVal).Int64()
			(*t).Log("Nr.", i, ", of Type ", v.Type(), "┄>", n)
		case FLOAT:
			n, e := v.(*floatVal).Float64()
			(*t).Log("Nr.", i, ", of Type ", v.Type(), "┄>", n, " Exactness: ", e)
		}

	}
}

/// MAIN TEST FUNCTION
func TestTypes(t *testing.T) { convertSliceOfSlices(InitValues(), t) }
