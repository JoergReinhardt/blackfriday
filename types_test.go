package agiledoc

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"math/big"
	"math/rand"
	"strings"
	"testing"
)

var L string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam convallis malesuada est eget finibus. Suspendisse suscipit tempor finibus. Cras hendrerit id lectus a pulvinar. Pellentesque a vehicula tortor, sit amet dapibus nunc. Aliquam egestas ligula lacinia odio suscipit blandit. Nulla auctor lectus dolor, a molestie lorem ultricies eget. Duis vulputate, purus et ultrices lobortis, sem tortor placerat tellus, nec venenatis enim nulla pulvinar odio. Duis lacinia faucibus mauris, at iaculis tellus dictum sed. Aliquam erat volutpat. Morbi et sollicitudin elit. Suspendisse potenti. In sed accumsan nulla, et malesuada enim. Aliquam at quam rhoncus, consequat neque et, sagittis risus. Aliquam mollis sollicitudin facilisis. Nulla sodales nibh eget viverra blandit. Pellentesque auctor porttitor consequat.\n\n Mauris et erat sed libero finibus aliquam. Proin condimentum viverra semper. Maecenas consectetur nibh id dignissim venenatis. Fusce molestie eros id mauris porta, eget semper neque vulputate. Duis tincidunt odio sed turpis facilisis euismod. Sed eget venenatis mi. Mauris efficitur orci nec ultrices aliquam. Nam lacus felis, maximus eget tincidunt sit amet, convallis eget purus. Interdum et malesuada fames ac ante ipsum primis in faucibus.\n\n Sed a imperdiet dolor, quis euismod ligula. Aenean blandit leo tortor, eget dictum libero fringilla ut. Ut a dignissim elit. Quisque elementum porta posuere. Integer rhoncus ipsum turpis, ornare egestas risus accumsan non. Quisque cursus orci ac mi auctor, in vehicula ligula accumsan. Quisque eget diam porttitor, tincidunt ex eget, iaculis leo. Nulla ultrices a neque sed feugiat. Phasellus auctor nibh eget odio tincidunt, quis commodo ipsum tempus. Suspendisse urna quam, aliquet vel risus sit amet, tempor pellentesque ligula.\n\n Nullam ac magna ac libero dapibus cursus. Phasellus vel nisl eu purus posuere aliquet. Curabitur non congue mauris, id maximus erat. Phasellus varius nisl et augue placerat, non fringilla turpis imperdiet. Mauris tellus ex, auctor vel pulvinar in, eleifend nec tortor. Vivamus sit amet nisl sit amet eros consequat pretium. In vitae leo lobortis, fringilla augue ut, facilisis eros. Donec efficitur, augue in ultrices varius, est turpis tempor enim, sed pharetra risus metus id nibh. Vivamus et sapien dictum est tempus eleifend. Nunc sed porta nibh. Duis sed felis dolor. Vestibulum sodales sagittis ex, et faucibus elit. Nam vitae felis eget neque laoreet euismod. Proin luctus efficitur lectus, non posuere dui fringilla non.\n\n Aenean quis ipsum sit amet ipsum pharetra vehicula sed porttitor erat. Integer non augue cursus erat placerat malesuada vel et libero. Aenean eu orci et augue ullamcorper malesuada. Nulla in dictum tellus. Aenean vel nisi lacus. Mauris eros lacus, mattis vel justo eu, suscipit ultricies dui. Sed porta fringilla mi vitae porttitor. Vestibulum pulvinar libero interdum eros convallis, non consequat sem placerat. Maecenas eu consequat ex. Ut mattis ut ex rutrum vestibulum.\n"

type Generator int

func (i *Generator) NextInt() int { *i = (*i) + 1; return int(*i) - 1 }
func (i *Generator) NextRat() float64 {
	*i = (*i) + 1
	j := new(val).BigInt().Rand(&rand.Rand{}, new(big.Int).SetInt64(int64(*i)))
	return float64(*i) * float64(j.Int64())
}
func (i *Generator) NextChar() string {
	*i = (*i) + 1
	if int(*i) >= len(strings.Split(L, "")) {
		*i = 0
	}
	return strings.Split(L, "")[*i]
}
func (i *Generator) NextWord() string {
	*i = (*i) + 1
	if int(*i) >= len(strings.Split(L, " ")) {
		*i = 0
	}
	return strings.Split(L, " ")[*i]
}
func (i *Generator) NextLine() string {
	*i = (*i) + 1
	if int(*i) >= len(strings.Split(L, "\n")) {
		*i = 0
	}
	return strings.Split(L, "\n")[*i]
}
func (i *Generator) NextParagraph() string {
	*i = (*i) + 1
	if int(*i) >= len(strings.Split(L, "\n\n")) {
		*i = 0
	}
	return strings.Split(L, "\n\n")[*i]
}
func (i *Generator) NextDigit() int {
	for int(*i) <= 9 {
		return (*i).NextInt()
	}
	(*i).Reset()
	return (*i).NextInt()
}
func (i *Generator) Reset() { *i = 0 }

func NewGenerator() *Generator { var i int = 0; return (*Generator)(&i) }

var G = NewGenerator()

func TestValueFromNative(t *testing.T) {
	c := G
	for n := 0; n <= 10; n++ {
		v := Value((*c).NextInt())
		(*t).Log(
			spew.Sprint("Type: ", v.Type(), " serialized: ", v.Serialize(), " string: ", v.String()),
		)
		a := Value((*c).NextChar())
		(*t).Log(
			spew.Sprint("Chars: ", a.Type(), " serialized: ", a.Serialize(), " string: ", a.String()),
		)
		s := Value((*c).NextWord())
		(*t).Log(
			spew.Sprint("Words: ", s.Type(), " serialized: ", s.Serialize(), " string: ", s.String()),
		)
		l := Value((*c).NextLine())
		(*t).Log(
			spew.Sprint("Lines: ", l.Type(), " serialized: ", l.Serialize(), " string: ", l.String()),
		)
		p := Value((*c).NextParagraph())
		(*t).Log(
			spew.Sprint("Lines: ", p.Type(), " serialized: ", p.Serialize(), " string: ", p.String()),
		)
	}
}
