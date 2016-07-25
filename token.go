// TOKENS
//
// agiledoc emulates the token format of bleve search, to enshure compability
// and ability to make use of bleves analysing and indexing capabilitys.
//
// type Token struct {
// 	// Start specifies the byte offset of the beginning of the term in the
// 	// field.
// 	Start int `json:"start"`
//
// 	// End specifies the byte offset of the end of the term in the field.
// 	End  int    `json:"end"`
// 	Term []byte `json:"term"`
//
// 	// Position specifies the 1-based index of the token in the sequence of
// 	// occurrences of its term in the field.
// 	Position int       `json:"position"`
// 	Type     TokenType `json:"type"`
// 	KeyWord  bool      `json:"keyword"`
// }
//
// func (t *Token) String() string {
// 	return fmt.Sprintf("Start: %d  End: %d  Position: %d  Token: %s  Type:
// 	%d", t.Start, t.End, t.Position, string(t.Term), t.Type)
// }
//
// type TokenStream []*Token
//
// since the approach is different, position is stored in a more convenient way
// but can be resolved to a start and endposition like expected and provided by
// the bleve library.
//
// Tokens map metainformation to the corresponding piece of text. It is
// mandatory to know the position relative to the containing context, as well
// as the text itself, or a refer to it. further meta information can be in the
// form of tags to mark all sets the content is part of, or arbitrary other
// parameters, to express parsed values and other information extracted.
package agiledoc

import (
	"bytes"
	// "github.com/emirpasic/gods/maps/hashmap"
)

// the position type provides integer indices that reference the byte indecx of
// the start and end point of a token relative to its containing context.
// input
type pos [2]int

func (p pos) Start() int { return p[0] }
func (p pos) End() int   { return p[1] }

// Position instances calculate and return the length of the part they reference.
func (p pos) SetLength(l int) { p[1] = p[0] + l }

// An instance of position can allways return a new instance of position with
// the next byte behind its end index as the start position for the new
// position.
func (p pos) genNextPos() pos { return pos{p.End() + 1, p.End() + 1} }
func (p pos) GenNextToken(t TType, raw []byte, flags uint, parms ...Value) Token {
	var pc Container
	return Token{
		t,
		p.genNextPos(),
		bytesVal(raw),
		flags,
		pc,
	}
}

// The token type combines a type with a position marker and a list of
// parameters. The parameter list is implemented using the gods library to
// profit from its enumerators and iterators. the list of parameters will be
// implemented bu gods hashmap. While encapsulating the empty interface the god
// interface has to expose being universal.
type Token struct {
	ttype TType
	pos
	rawTxt Value
	flags  uint
	params Container // contains god hashmap
}

type TType uint32

//go:generate -command stringer -type TType ./token.go
const (
	MD_AutoLink TType = 1 << iota
	MD_BlockCode
	MD_BlockHtml
	MD_BlockQuote
	MD_CodeSpan
	MD_DocumentFooter
	MD_DocumentHeader
	MD_DoubleEmphasis
	MD_Emphasis
	MD_Entity
	MD_FootnoteItem
	MD_FootnoteRef
	MD_Footnotes
	MD_GetFlags
	MD_HRule
	MD_Header
	MD_Image
	MD_LineBreak
	MD_Link
	MD_List
	MD_ListItem
	MD_NormalText
	MD_Paragraph
	MD_RawHtmlTag
	MD_StrikeThrough
	MD_Table
	MD_TableCell
	MD_TableHeaderCell
	MD_TableRow
	MD_TitleBlock
	MD_TripleEmphasis
)

//// Tokenizer
/// TODO: evaluate if queue/semafore performs any better
// Tokenizer implements blackfriday renderer. token content is written to
// semaBuf, token instance is appended to semaQueue of the contained parser
type Tokenizer struct {
	// options blackfriday.Options // parameters
	flags intVal     // options
	cur   Token      // current position
	out   chan Token // returns tokens to the caller
}

// func (tkz *Tokenizer) newToken(t TType, raw []byte, flags uint, parms ...keyValue) {
// 	var c Container
// 	// copy parameters, only if there are any
// 	if len(parms) > 0 {
// 		// allocate array to take parameters
// 		m := hashmap.New()
//
// 		for _, v := range parms {
// 			(*m).Put(v.Key(), v.Valueue())
// 		}
//
// 		c = NewTypedValue(CONTAINER, m).(intVal)
// 	}
// 	// calculate new position
// 	ns := (*tkz).cur.End() + 1
// 	pos := pos{ns, ns + len(raw)}
//
// 	// put token into channel
// 	(*tkz).out <- (*tkz).cur
//
// 	// generate token of designated type, with calculated position and the options container
// 	(*tkz).cur = Token{t, pos, raw, flags, c}
//
// }

// func NewTokenizer(flags ...uint) (*Tokenizer, chan Token) {
// 	f := NewTypedValue(FLAG, 0)
// 	// XOR flags
// 	for _, flag := range flags {
// 		cmp := NewTypedValue(FLAG, flag)
// 		f.Valueue.Flag().Xor(f.Valueue.Flag(), cmp.Valueue.Flag())
// 	}
// 	o := make(chan Token, 1)
// 	return &Tokenizer{f, Token{}, o}, o
// }

//// BLACK FRIDAY INTERFACE IMPLEMENTATION
///
// implement the black friday renderer to act as receiver for the parsed parts
// of a markdown document. Generates tokens of the token type, that take all
// provided metadata, parsed and raw data.
//
// DOCUMENT METAINFO HEADER AND FOOTER
// func (t *Tokenizer) DocumentHeader(out *bytes.Buffer) {
// 	raw := out.Bytes()
// 	(*t).newToken(DocumentHeader, raw, 0)
// }
//
// func (t *Tokenizer) DocumentFooter(out *bytes.Buffer) {
// 	raw := out.Bytes()
// 	(*t).newToken(DocumentFooter, raw, 0)
// }

// DOCUMENT BLOCKS
func (t *Tokenizer) Header(out *bytes.Buffer, text func() bool, level int, id string) { // header as in headline of a section
}
func (t *Tokenizer) BlockCode(out *bytes.Buffer, text []byte, lang string)                 {}
func (t *Tokenizer) BlockQuote(out *bytes.Buffer, text []byte)                             {}
func (t *Tokenizer) BlockHtml(out *bytes.Buffer, text []byte)                              {}
func (t *Tokenizer) HRule(out *bytes.Buffer)                                               {}
func (t *Tokenizer) List(out *bytes.Buffer, text func() bool, flags int)                   {}
func (t *Tokenizer) ListItem(out *bytes.Buffer, text []byte, flags int)                    {}
func (t *Tokenizer) Paragraph(out *bytes.Buffer, text func() bool)                         {}
func (t *Tokenizer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (t *Tokenizer) TableRow(out *bytes.Buffer, text []byte)                               {}
func (t *Tokenizer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             {}
func (t *Tokenizer) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (t *Tokenizer) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (t *Tokenizer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}
func (t *Tokenizer) TitleBlock(out *bytes.Buffer, text []byte)                             {}

// Span-level callbacks
func (t *Tokenizer) AutoLink(out *bytes.Buffer, link []byte, kind int)                 {}
func (t *Tokenizer) CodeSpan(out *bytes.Buffer, text []byte)                           {}
func (t *Tokenizer) DoubleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (t *Tokenizer) Emphasis(out *bytes.Buffer, text []byte)                           {}
func (t *Tokenizer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)    {}
func (t *Tokenizer) LineBreak(out *bytes.Buffer)                                       {}
func (t *Tokenizer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {}
func (t *Tokenizer) RawHtmlTag(out *bytes.Buffer, tag []byte)                          {}
func (t *Tokenizer) TripleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (t *Tokenizer) StrikeThrough(out *bytes.Buffer, text []byte)                      {}
func (t *Tokenizer) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                 {}

// Low-level callbacks
func (t *Tokenizer) Entity(out *bytes.Buffer, entity []byte)   {}
func (t *Tokenizer) NormalText(out *bytes.Buffer, text []byte) {}

// func (t *Tokenizer) GetFlags() int { return int((*t).flags.Int64()) }
