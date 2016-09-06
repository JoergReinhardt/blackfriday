package agiledoc

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

type NumGenerator int

func (i *NumGenerator) NextInt() int { *i = (*i) + 1; return int(*i) - 1 }
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
	for n := 0; n <= 100; n++ {
		v := Value((*c).NextInt())

		//		if Int(v.Base()).Integer() != n {
		//			(*t).Fail("Error: expected ", n, " got, ", v.Base().Integer(), " instead\n")
		//		}
		//		if string(v.Base().Bytes()) != string(n) {
		//			(*t).Fail("Error: expected ", string(n), " got, ", string(v.Base().Bytes()), " instead\n")
		//		}
		//		if v.Type() != INTEGER {
		//			(*t).Fail("Error: expected ", INTEGER, " got, ", v.Type(), " instead\n")
		//		}
		//		if v.String() != string(n) {
		//			(*t).Fail("Error: expected ", string(n), " got, ", v.String(), " instead\n")
		//		}
		//		if v.Serialize() != []byte(string(n)) {
		//			(*t).Fail("Error: expected ", []byte(string{n}), " got, ", v.Serialize(), " instead\n")
		//		}

		(*t).Log(
			spew.Sprint("Type: ", v.Type(), " serialized: ", v.Serialize(), " string: ", v.String()),
		)
	}
}
