package agiledoc

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

type NumGenerator int

func (i *NumGenerator) NextInt() int    { *i = (*i) + 1; return int(*i) - 1 }
func (i *NumGenerator) NextLit() string { *i = (*i) + 1; return Number(*i - 1).String() }
func (i *NumGenerator) NextDigit() int {
	for *i <= 9 {
		return (*i).NextInt()
	}
	(*i).Reset()
	return (*i).NextInt()
}
func (i *NumGenerator) Reset() { *i = 0 }

func NewNumGenerator() *NumGenerator { var i int = 0; return (*NumGenerator)(&i) }

var UG = NewNumGenerator()

func TestValueFromNative(t *testing.T) {
	c := UG
	for n := 0; n <= 10; n++ {
		v := Value((*c).NextInt())
		(*t).Log(
			spew.Sprint("Type: ", v.Type(), " serialized: ", v.Serialize(), " string: ", v.String()),
		)
		s := Value((*c).NextLit())
		(*t).Log(
			spew.Sprint("Type: ", s.Type(), " serialized: ", s.Serialize(), " string: ", s.String()),
		)
	}
}
