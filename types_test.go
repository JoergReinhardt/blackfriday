package agiledoc

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

type IntGenerator int

func (i *IntGenerator) Next() int { *i = (*i) + 1; return int(*i) - 1 }
func (i *IntGenerator) NextDigit() int {
	for *i <= 9 {
		return (*i).Next()
	}
	(*i).Reset()
	return (*i).Next()
}
func (i *IntGenerator) Reset() { *i = 0 }

func NewIntGenerator() *IntGenerator { var i int = 0; return (*IntGenerator)(&i) }

var IG = NewIntGenerator()

func TestValueFromNative(t *testing.T) {
	c := IG
	for n := 0; n <= 100; n++ {
		v := Value((*c).NextDigit())
		(*t).Log(
			spew.Sprint(v.Type(), v.Serialize(), v.String(), v.Base().Bytes()),
		)
	}
}
