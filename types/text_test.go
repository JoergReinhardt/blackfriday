package types

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	//	"math/big"
	//	"math/rand"
	//	"strings"
	"testing"
)

type textTestsFunc func(string, string) string

var textTests = []struct {
	a     string
	b     string
	exp   string
	opStr string
	op    textTestsFunc
}{
	{"test", "test", "testtest", "AppendText", func(a, b string) string { return a + b }},
}

func testText(t *testing.T, a string, b string, exp string, opStr string, op func(a, b string) string) {
	if op(a, b) != exp {
		(*t).Fail()
		(*t).Log("failed operation: " + fmt.Sprint(opStr) +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(op(a, b)) +
			" expected: " + fmt.Sprint(exp))
	} else {
		(*t).Log("passed operation: " + fmt.Sprint(opStr) +
			" a: " + fmt.Sprint(a) +
			" b: " + fmt.Sprint(b) +
			" got: " + fmt.Sprint(op(a, b)) +
			" expected: " + fmt.Sprint(exp))
	}
}

func TestText(t *testing.T) {

	for n, test := range textTests {

		n := n
		a := test.a
		b := test.b
		exp := test.exp
		opStr := test.opStr
		op := test.op

		t.Log(fmt.Sprintf("Test Nr. %d: ", n))

		testText(t, a, b, exp, opStr, op)

	}
}

var L string = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam convallis malesuada est eget finibus. Suspendisse suscipit tempor finibus. Cras hendrerit id lectus a pulvinar. Pellentesque a vehicula tortor, sit amet dapibus nunc. Aliquam egestas ligula lacinia odio suscipit blandit. Nulla auctor lectus dolor, a molestie lorem ultricies eget. Duis vulputate, purus et ultrices lobortis, sem tortor placerat tellus, nec venenatis enim nulla pulvinar odio. Duis lacinia faucibus mauris, at iaculis tellus dictum sed. Aliquam erat volutpat. Morbi et sollicitudin elit. Suspendisse potenti. In sed accumsan nulla, et malesuada enim. Aliquam at quam rhoncus, consequat neque et, sagittis risus. Aliquam mollis sollicitudin facilisis. Nulla sodales nibh eget viverra blandit. Pellentesque auctor porttitor consequat.\n\n Mauris et erat sed libero finibus aliquam. Proin condimentum viverra semper. Maecenas consectetur nibh id dignissim venenatis. Fusce molestie eros id mauris porta, eget semper neque vulputate. Duis tincidunt odio sed turpis facilisis euismod. Sed eget venenatis mi. Mauris efficitur orci nec ultrices aliquam. Nam lacus felis, maximus eget tincidunt sit amet, convallis eget purus. Interdum et malesuada fames ac ante ipsum primis in faucibus.\n\n Sed a imperdiet dolor, quis euismod ligula. Aenean blandit leo tortor, eget dictum libero fringilla ut. Ut a dignissim elit. Quisque elementum porta posuere. Integer rhoncus ipsum turpis, ornare egestas risus accumsan non. Quisque cursus orci ac mi auctor, in vehicula ligula accumsan. Quisque eget diam porttitor, tincidunt ex eget, iaculis leo. Nulla ultrices a neque sed feugiat. Phasellus auctor nibh eget odio tincidunt, quis commodo ipsum tempus. Suspendisse urna quam, aliquet vel risus sit amet, tempor pellentesque ligula.\n\n Nullam ac magna ac libero dapibus cursus. Phasellus vel nisl eu purus posuere aliquet. Curabitur non congue mauris, id maximus erat. Phasellus varius nisl et augue placerat, non fringilla turpis imperdiet. Mauris tellus ex, auctor vel pulvinar in, eleifend nec tortor. Vivamus sit amet nisl sit amet eros consequat pretium. In vitae leo lobortis, fringilla augue ut, facilisis eros. Donec efficitur, augue in ultrices varius, est turpis tempor enim, sed pharetra risus metus id nibh. Vivamus et sapien dictum est tempus eleifend. Nunc sed porta nibh. Duis sed felis dolor. Vestibulum sodales sagittis ex, et faucibus elit. Nam vitae felis eget neque laoreet euismod. Proin luctus efficitur lectus, non posuere dui fringilla non.\n\n Aenean quis ipsum sit amet ipsum pharetra vehicula sed porttitor erat. Integer non augue cursus erat placerat malesuada vel et libero. Aenean eu orci et augue ullamcorper malesuada. Nulla in dictum tellus. Aenean vel nisi lacus. Mauris eros lacus, mattis vel justo eu, suscipit ultricies dui. Sed porta fringilla mi vitae porttitor. Vestibulum pulvinar libero interdum eros convallis, non consequat sem placerat. Maecenas eu consequat ex. Ut mattis ut ex rutrum vestibulum.\n"
