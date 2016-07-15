package agiledoc

import (
	"bytes"
	"sync"
)

// the position type provides byte indices of start and end of a token in the
// input
type pos [2]int

// a parameter has a name and a value
// flags are considered a parameter
type keyVal struct {
	key string
	val Value
}

// token is of token type, references a start and end position in the input
// stream and a slice of optional parameters
type token struct {
	ttype    tokenType
	position pos
	params   []keyVal
}

func (t token) len() int { return t.position[1] - t.position[0] }

//// BUFFER
///
// a lockable buffer that contains the document in form of a byte slice
type buffer struct {
	*sync.RWMutex
	*bytes.Buffer
}

func newBuffer() *buffer {
	return &buffer{
		&sync.RWMutex{},
		bytes.NewBuffer([]byte{}),
	}
}

//// QUEUE OF TOKENS
///
// the queue is concurrently write- and readable and provides the current
// position in the token slice
type queue struct {
	*sync.RWMutex
	curPos int
	head   int
	tail   int
	queue  []token
}

// add token at current position
func (q *queue) add(tok token) {
	(*q).Lock()
	defer (*q).Unlock()

	// split head from tail
	h, t := (*q).queue[:q.curPos], (*q).queue[q.curPos+1:]
	(*q).queue = append(h, tok)
	(*q).queue = append(q.queue, t...)
	(*q).tail = q.tail + 1
}

// add token at current position
func (q *queue) del() {
	(*q).Lock()
	defer (*q).Unlock()

	// split head from tail and ommit current token
	h, t := (*q).queue[:q.curPos-1], (*q).queue[q.curPos+1:]
	(*q).queue = append(h, t...)
	(*q).tail = q.tail - 1
}

// append token at the end of the queue
func (q *queue) append(t token) {
	(*q).Lock()
	defer (*q).Unlock()

	// advance position counter
	(*q).curPos = (*q).curPos + t.len()
	// append the token to the queue
	(*q).queue = append(q.queue, t)
}

func newQueue() *queue {
	return &queue{
		&sync.RWMutex{},
		0,
		0,
		0,
		[]token{},
	}
}

//// PARSER
///
type parser struct {
	*buffer // contains document as bute slice
	*queue  // contains tokens to be parsed
}

//// TOKENIZER
///
// tokenizer implements blackfriday renderer. token content is written to
// buffer, token instance is appended to queue
type tokenizer struct {
	flags   int
	curPos  int
	*buffer // references parsers buffer
	*queue  // references parsers queue
}

// newtoken is called by all methods implementing the blackfriday renderer
// interface, to generate a token and propagate it to the caller
func (t *tokenizer) newtoken(typ tokenType, content []byte, params ...keyVal) {

	// instanciate a new token
	tok := token{
		ttype: typ, // determined by each callback
		position: pos{ // instanciate pos from last position to lp plus elements length
			(*t).curPos + 1,                // starts one byte after end of last tag
			(*t).curPos + 1 + len(content), // ends at (start + length of content)
		},
		params: params, // slice of name/value pairs
	}

	// update tokeniers last position of to end of new token
	(*t).curPos = tok.position[1]

	// send token to queue
	(*t).queue.append(tok)

	// write content to buffer
	(*t).Read(content)
}

// Header and footer
func (t *tokenizer) DocumentHeader(out *bytes.Buffer) {
	// (*t).newtoken(D_HEADER, out.Bytes(), keyVal{"", emptyVal{}})
}

func (t *tokenizer) DocumentFooter(out *bytes.Buffer) {
	// (*t).newtoken(D_FOOTER, out.Bytes(), keyVal{"", emptyVal{}})
}

// Document Blocks
func (t *tokenizer) Header(out *bytes.Buffer, text func() bool, level int, id string) { // header as in headline as in section
	// 	var params []keyVal = []keyVal{}
	// 	params = append(params, keyVal{"level", newVal(INTEGER, level)})
	// 	params = append(params, keyVal{"id", newVal(STRING, id)})
	//
	// 	if text() { // call callback to generate byteslice and writes it to the buffer
	// 		(*t).newtoken( // generates a new token and sends it to the callers queue
	// 			SECTION,
	// 			out.Bytes(),
	// 			params...,
	// 		)
	// 	}
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

func newtokenizer(queue chan token, name string, flags ...int) *tokenizer {

	var nf int = 0

	for _, f := range flags {
		f := f
		nf = nf | f
	}

	return &tokenizer{
		flags:  nf,
		curPos: 0,
		buffer: newBuffer(),
	}
}
