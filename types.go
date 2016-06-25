package agiledoc

import (
	"bytes"
	// "fmt"
	"math/big"
	//bf "github.com/russross/blackfriday"
)

type Value interface {
	Type() valueType
	ToType(valueType) (ok bool, val Value)
	Eval() []byte
}

type (
	emptyValue  struct{ int } // special value -1
	boolValue   struct{ bool }
	intValue    struct{ *big.Int }
	floatValue  struct{ *big.Float }
	stringValue struct{ s []byte }
	Vector      struct{ v []Value }
	Matrix      struct{ Vector }
)

func (emptyValue) Type() valueType  { return EMPTY }
func (boolValue) Type() valueType   { return BOOL }
func (intValue) Type() valueType    { return INTEGER }
func (floatValue) Type() valueType  { return FLOAT }
func (stringValue) Type() valueType { return STRING }
func (Vector) Type() valueType      { return VECTOR }
func (Matrix) Type() valueType      { return MATRIX }

func (emptyValue) ToType(t valueType) (ok bool, val Value)  { return false, nil }
func (boolValue) ToType(t valueType) (ok bool, val Value)   { return false, nil }
func (intValue) ToType(t valueType) (ok bool, val Value)    { return false, nil }
func (floatValue) ToType(t valueType) (ok bool, val Value)  { return false, nil }
func (stringValue) ToType(t valueType) (ok bool, val Value) { return false, nil }
func (Vector) ToType(t valueType) (ok bool, val Value)      { return false, nil }
func (Matrix) ToType(t valueType) (ok bool, val Value)      { return false, nil }

func (v emptyValue) Eval() []byte { return []byte{} }
func (v boolValue) Eval() []byte {
	if v.bool {
		return []byte{1}
	} else {
		return []byte{}
	}
}
func (v intValue) Eval() []byte    { return v.Bytes() }
func (v floatValue) Eval() []byte  { return []byte(v.String()) }
func (v stringValue) Eval() []byte { return v.s }
func (v Vector) Eval() []byte {

	ret := make([]byte, len(v.v))

	for _, val := range v.v {
		val := val
		ret = append(ret, val.Eval()...)
	}
	return ret
}

func (v Matrix) Eval() []byte {
	return v.Eval()
}

func newValue(t valueType, v interface{}) Value {

	var ret Value

	switch t {
	case EMPTY:
		ret = emptyValue{}

	case BOOL:
		if v.(bool) {
			ret = boolValue{true}
		} else {
			ret = boolValue{false}
		}

	case INTEGER:
		ret = intValue{big.NewInt(v.(int64))}

	case FLOAT:
		ret = floatValue{big.NewFloat(v.(float64))}

	case STRING:

	case VECTOR:

	case MATRIX:
	}

	return ret
}

type pos [2]int

type Token struct {
	ttype     tokenType
	position  pos
	content   []byte
	flags     int
	parameter []param
}

type (
	tokenType uint32
	valueType uint8
)

const (
	EMPTY valueType = 0
	BOOL            = 1 << iota
	INTEGER
	FLOAT
	STRING
	VECTOR
	MATRIX

	termVTypes = BOOL | INTEGER | FLOAT | STRING
	combVTypes = VECTOR | MATRIX

	// BLOCK LEVEL
	DOCUMENT tokenType = 0
	D_HEADER           = 1 << iota // header
	D_FOOTER
	SECTION
	TITLE
	PARAGRAPH
	CODE
	QUOTE
	HTML
	HRULE
	LIST
	L_ITEM
	TABLE
	T_HEADER_CELL
	T_ROW
	T_CELL
	FOOTNOTES
	F_ITEM
	//SPAN_LEVEL
	AUTO_LINK
	CODE_SPAN
	LINE_BREAK
	EMPHASIS
	DOUBLE_EMPHASIS
	TRIPLE_EMPHASIS
	STRIKE_THROUGHT
	RAW_HTML_TAG
	LINK
	IMAGE
	F_REF // footnote reference
	// LOW LEVEL
	ENTITY
	TEXT

	blockElements = DOCUMENT | D_HEADER | D_FOOTER | SECTION | TITLE | PARAGRAPH | CODE | QUOTE | HTML | HRULE | LIST | L_ITEM | TABLE | T_HEADER_CELL | T_ROW | T_CELL | FOOTNOTES | F_ITEM

	spanElements = AUTO_LINK | CODE_SPAN | LINE_BREAK | EMPHASIS | DOUBLE_EMPHASIS | TRIPLE_EMPHASIS | STRIKE_THROUGHT | RAW_HTML_TAG | LINK | IMAGE | F_REF

	lowLevelElements = ENTITY | TEXT
)

type param struct {
	name  string
	value Value
}

type tokenizer struct {
	name    string
	flags   int
	lastPos int
	queue   chan Token
}

// newToken is called by all methods implementing the blackfriday renderer
// interface, to generate a token and propagate it to the caller
func (t *tokenizer) newToken(typ tokenType, content []byte, flags int, parms ...param) {

	// instanciate a new token
	tok := Token{
		ttype: typ, // determined by each callback
		position: pos{ // instanciate pos from last position to lp plus elements length
			(*t).lastPos + 1,                // starts one byte after end of last tag
			(*t).lastPos + 1 + len(content), // ends at (start + length of content)
		},
		content:   content, // byte slice
		flags:     flags,   // OR concatenated ints
		parameter: parms,   // slice of name/value pairs
	}

	// update tokeniers last position of to end of new token
	(*t).lastPos = tok.position[1]

	// send token to queue
	(*t).queue <- tok
}

// Header and footer
func (t *tokenizer) DocumentHeader(out *bytes.Buffer) {
	(*t).newToken(D_HEADER, out.Bytes(), 0)
}

func (t *tokenizer) DocumentFooter(out *bytes.Buffer) {
	(*t).newToken(D_FOOTER, out.Bytes(), 0)
}

// Document Blocks
func (t *tokenizer) Header(out *bytes.Buffer, text func() bool, level int, id string) { // header as in headline as in section
	var parms []param = []param{}
	parms = append(parms, param{"level", newValue(INTEGER, level)})
	parms = append(parms, param{"id", newValue(STRING, id)})

	if text() { // call callback to generate byteslice and writes it to the buffer
		(*t).newToken( // generates a new token and sends it to the callers queue
			SECTION,
			out.Bytes(),
			0,
			parms...,
		)
	}
}
func (t *tokenizer) BlockCode(out *bytes.Buffer, text []byte, lang string)                 {}
func (t *tokenizer) BlockQuote(out *bytes.Buffer, text []byte)                             {}
func (t *tokenizer) BlockHtml(out *bytes.Buffer, text []byte)                              {}
func (t *tokenizer) HRule(out *bytes.Buffer)                                               {}
func (t *tokenizer) List(out *bytes.Buffer, text func() bool, flags int)                   {}
func (t *tokenizer) ListItem(out *bytes.Buffer, text []byte, flags int)                    {}
func (t *tokenizer) Paragraph(out *bytes.Buffer, text func() bool)                         {}
func (t *tokenizer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (t *tokenizer) TableRow(out *bytes.Buffer, text []byte)                               {}
func (t *tokenizer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             {}
func (t *tokenizer) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (t *tokenizer) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (t *tokenizer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}
func (t *tokenizer) TitleBlock(out *bytes.Buffer, text []byte)                             {}

// Span-level callbacks
func (t *tokenizer) AutoLink(out *bytes.Buffer, link []byte, kind int)                 {}
func (t *tokenizer) CodeSpan(out *bytes.Buffer, text []byte)                           {}
func (t *tokenizer) DoubleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (t *tokenizer) Emphasis(out *bytes.Buffer, text []byte)                           {}
func (t *tokenizer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)    {}
func (t *tokenizer) LineBreak(out *bytes.Buffer)                                       {}
func (t *tokenizer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {}
func (t *tokenizer) RawHtmlTag(out *bytes.Buffer, tag []byte)                          {}
func (t *tokenizer) TripleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (t *tokenizer) StrikeThrough(out *bytes.Buffer, text []byte)                      {}
func (t *tokenizer) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                 {}

// Low-level callbacks
func (t *tokenizer) Entity(out *bytes.Buffer, entity []byte)   {}
func (t *tokenizer) NormalText(out *bytes.Buffer, text []byte) {}

func (t *tokenizer) GetFlags() int { return (*t).flags }

func newTokenizer(queue chan Token, name string, flags ...int) *tokenizer {

	var nf int = 0

	for _, f := range flags {
		f := f
		nf = nf | f
	}

	return &tokenizer{
		name:    name,
		flags:   nf,
		lastPos: 0,
		queue:   queue,
	}
}
