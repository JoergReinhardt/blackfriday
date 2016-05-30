package blackfriday

import (
	"sync"
)

// These are the possible flag Values for the context renderer
//
// Definitions of variables and functions, are propagated to the context,
// references to context variables and calls of context functions are
// backtracked to the corresponding definitions, values parameter and return
// values, expandet and evaluated.
//
// Sections are addressable as value by their ID and may contain, paragraphs,
// span type elements, each of the blocks specialy treated by context renderer
// and every arbitrary other block, as well as other (Sub-)Sections, The
// „document-root” describes a virtual section containing a complete document,
// with all its contained variables, references, links images and figures,
// includes and so on.
//
// Context variables and functions have an ID, depending on the scope chosen.
// That ID has either to be unique for the whole document, or to be prefixed by
// the dot concatenated chane of its ancestor elements names, beginning with
// the upmost level of section headings, right below the „document-root”.
const (
	/////////////// BASE TYPES ///
	// terminal values to actually store data in an appropriate type
	CTX_TYPE_INT    = 1 << iota // single integer
	CTX_TYPE_FLOAT              // single float
	CTX_TYPE_STRING             // arbitrary piece of string
	CTX_TYPE_VECTOR             // list, deflist, or map of values
	CTX_TYPE_MATRIX             // twodimensional array of vectors
	//////////////// OPERATORS
	CTX_TYPE_UNARYOP  // called when a unary operator is encountered
	CTX_TYPE_BINARYOP // called when a binary operator is encountered

	// Document Element Types
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

type Value interface {
	Eval() Value
	Type() int
	ToType(int) Variable
}
type position struct {
	start int // == offset
	end   int // start - end = len
}
type Variable interface {
	Id() string
	Pos() *position
	Value
}

// CONTEXT
//
//the context struct implements the symtab, that maps values to identity and
//type. It implements the variable type for all types of variables alike,
//Symtab key is a two field struct containing of the types identity and type.
//The Value contains the raw value in a byte slice and can return itself.
//
type Context struct {
	flags  int
	lock   sync.RWMutex
	symtab map[string]Value
	ast    *Tree
}

func NewContext() *Context {
	return &Context{
		lock:   sync.RWMutex{},
		symtab: make(map[string]Val),
		ast:    NewRoot(),
	}
}

func (c *Context) Lookup(id string) Value {
	return (*c).symtab[id]
}

func (c *Context) Declare(id string, v Value) {
	(*c).symtab[id] = v
}
