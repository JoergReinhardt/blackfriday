package agiledoc

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

// constants that denote the type of variable, so that the marshal function can
// assert it.

const (
	Bool int = 1 + iota
	INt
	Flt
	Byt
	Str
	Num
	Wor
	Sce
)

var NativeValues = struct {
	Bool []bool
	INt  []int
	Flt  []float64
	Byt  []byte
	Str  []string
	Num  []string
	Wor  []string
	Sce  []string
}{
	[]bool{false, true, false, true, false, true, false, true, false, true},
	[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	[]float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
	[]byte("allet ja nuescht so einfach hier mit die viele Worte"),
	[]string{"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten"},
	[]string{"fourtytwothousendonehundredtwentythree", "sevenhundredthirteen"},
	[]string{"Word", "I", "am", "slim", "shady"},
	[]string{"This is a scentence, exemplary for something parts of speech can be extracted from"},
}

func TestNativeToBig(t *testing.T) {
	slice := []Value{}
	for _, v := range NativeValues.INt {
		val := nativeToValue(v)
		spew.Print(val.(Integer)())
		spew.Print(val.Type())
		slice = append(slice, val)
	}
	(*t).Log(
		"value slice after conversion and append",
		spew.Sprint(slice),
	)
}
