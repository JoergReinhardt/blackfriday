package agiledoc

import (
	// "fmt"
	"testing"
)

var Natives = struct {
	Bools []bool
	Ints  []int
	Rats  []float64
	Bytes []byte
	Str   []string
}{
	[]bool{false, true, false, true, false, true, false, true, false, true},
	[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	[]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
	[]byte("allet ja nuescht so einfach hier mit die viele Worte"),
	[]string{"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten", "twothousendthreehundredfiftyseven"},
}

var vals = []Value{}

/// MAIN TEST FUNCTION
func TestTypesBool(t *testing.T) {
	for _, v := range Natives.Bools {
		v := v
		val := emptyVal{}.Value()
		val = val.(emptyVal).Set(v)
		vals = append(vals, val)
	}
}
func TestTypesInt(t *testing.T) {
	for _, v := range Natives.Ints {
		v := v
		val := emptyVal{}.Value()
		val = val.(emptyVal).Set(v)
		vals = append(vals, val)
	}
}
func TestTypesRat(t *testing.T) {
	for _, v := range Natives.Rats {
		v := v
		val := emptyVal{}.Value()
		val = val.(emptyVal).Set(v)
		vals = append(vals, val)
	}
}
func TestTypesShortStrings(t *testing.T) {
	for _, v := range Natives.Str {
		v := v
		val := emptyVal{}.Value()

		val = val.(emptyVal).Set(v)
		vals = append(vals, val)
	}
}
func TestTypesLongStrings(t *testing.T) {
	for _, v := range Natives.Bytes {
		val := NativeToValue(v)
		vals = append(vals, val)
	}
	(*t).Log(vals)
}
func TestConveresion(t *testing.T) {
	for o := 0; o < 15; o++ {
		v := emptyVal{}.Set(Natives.Bytes)
		for i := 0; i < 9; i++ {
			i := i
			ty := 1 << uint(i)
			x := v.ToType(ValueType(ty))
			(*t).Log(x.Type(), x.String())
		}
	}
}
func TestCollectionList(t *testing.T) {
	col := NativeToValue(vals).(lstVal)
	(*t).Log(
		"\nType: ", col.Type(),
		"\nContainer: ", col.Container(),
		"\nContainer Values: ", col.Container().Values(),
		"\nContainer Size: ", col.Container().Size(),
		"\nEnumerable: ", col.Enumerable(),
		"\nIterator: ", col.Iterator(),
		"\nIterator Next(): ", col.Iterator().Next(),
		"\nList: ", col.List(),
	)
}
