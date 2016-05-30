package blackfriday

import "bytes"

// Document Element Types
const (
	ROOT     = 0
	AUTOLINK = 1 << iota
	BLOCKCODE
	BLOCKHTML
	BLOCKQUOTE
	CODESPAN
	DOCUMENTFOOTER
	DOCUMENTHEADER
	DOUBLEEMPHASIS
	EMPHASIS
	ENTITY
	FOOTNOTEITEM
	FOOTNOTEREF
	FOOTNOTES
	GETFLAGS
	HRULE
	HEADER
	IMAGE
	LINEBREAK
	LINK
	LIST
	LISTITEM
	NORMALTEXT
	PARAGRAPH
	RAWHTMLTAG
	STRIKETHROUGH
	TABLE
	TABLECELL
	TABLEHEADERCELL
	TABLEROW
	TITLEBLOCK
	TRIPLEEMPHASIS
)

// the value type contains the value as byte slice
type Value []byte

// when evaluated, a value will return a copy of itself
func (v Value) Eval() Value { return v }

// the Render Call returns the values content as raw byte slice
func (v Value) Render() []byte { return []byte(v) }

// the interface describing a variable by offering methods to fetch and
// (re)define it's value and perform type transformation, if nescessary
type Variable interface {
	// having context as first anonymous argument equals the signature of
	// a method of the context struct
	Eval(*Context) Value // return raw value
	Type(*Context) int   // return values type
	ToType(int) bool     // check if can be castet as passed tyoe
	CastAs(int) Variable // cast as other type
}
type tVal struct {
	typ int
	val []byte
}

// CONTEXT
//
//the context struct implements the symtab, that maps values to identity and
//type. It implements the variable type for all types of variables alike,
//Symtab key is a two field struct containing of the types identity and type.
//The Value contains the raw value in a byte slice and can return itself.
//
type Context map[string]tVal

func NewContext() Context { return make(map[string]tVal) }

func (c Context) Lookup(id string) (v Variable) {
	return v
}

// INLINE CALLS PASSED TO PARSER VIA OPTIONS
//
// context definition inline parser [:]
//
// gets called on the same byte encountered, as autolink [:]
// 1. check if context definithion
//  1.1.NO. call autolink, pass on orig params
//  1.2.YES --> check if:
//    1.2.1. value definition
//	1.2.1.1. integer
//	1.2.1.2. float
//	1.2.1.3. string
//	1.2.1.4. list (ul/ol)
//	1.2.1.5. dectionary (definition list)
//	1.2.1.6. matrix (table)
//    1.2.2. function definition
//	1.2.2.1. Parameter (Values)
//	1.2.2.2. Procedure (nested context code)
func contextDefinition(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return autoLink(p, out, data, offset)
}

// context reference inline parser [@]
// gets called on [@]. resolves refrences to variables and evaluates function
// calls
func contextReference(p *parser, out *bytes.Buffer, data []byte, offset int) int {
	return autoLink(p, out, data, offset)
}
