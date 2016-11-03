package tokens

import (
	"bytes"
	t "github.com/JoergReinhardt/blackfriday/types"
	//b "github.com/russross/blackfriday"
)

type Token struct {
	val  t.Bytes
	parm []t.Pair
}

type Renderer struct {
	tok chan<- Token
}

func (r Renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	r.tok <- Token{
		t.Value(text).(t.Bytes),
		[]t.Pair{
			t.Value("lang", lang).(t.Pair),
		},
	}
}
