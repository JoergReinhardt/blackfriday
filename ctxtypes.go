package blackfriday

import ()

const (
	//// CTX BASE TYPES
	///
	// terminal values to store actual data in an appropriate type
	DOC_ELEMENT = 0         // zero value representing part of input as a byteslice
	CTX_INT     = 1 << iota // single integer
	CTX_FLOAT               // single float
	CTX_STRING              // arbitrary piece of string
	CTX_VECTOR              // list, deflist, or map of values
	CTX_MATRIX              // twodimensional array of vectors

	//// OPERATORS
	///
	// Operators represent arithmetric, logic, binary operators and so on.
	// A unary operation can also represent a user defined function.
	// Parameters can be passed as a vector.
	//
	// When instanciated, operator and operands are passed to the
	// constructing fucntion. it determines the appropriate type, operands
	// need to be casted as, casts and passes them to the appropriate
	// function implementing the operation referenced by the passed
	// operator.
	// CTX_UNA_OP
	// CTX_BIN_OP

	//// TOKEN TYPES (DOCUMENT ELEMENT TYPES)
	///
	// an integer for every type of element that can be contained by a
	// markdown document.
	EMPTY = 0
	ROOT  = 1 << iota
	AUTOLINK
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

//// ELEMENT INTERFACE
///
// An anomynous element of a known type can be evaluated to a byte slice of it's
// contents. It can also return it's own type and be casted as value of another
// type, by using it's ToType() method
type Element interface {
	Eval() []byte       // return raw content as a byte slice
	Type() int          // return the type of the contained value
	ToType(int) Element // cast as element of another type
}

type Vector interface {
	Element                // vector type, evaluates to. row of its elements,
	Len() int              // return length of vector
	Idx(index int) Element // return a single element by its reference
}

type Matrix interface {
	Element        // matrix type, evaluates to list of rows of column elements
	Shape() [2]int // rows χ columns
	Row(ridx int) Vector
	Col(cidx int) Vector
}

//// VARIABLE INTERFACE
///
// A variable is a strig identifyer, that references some continuous Element of
// the input stream. The element is defined by it's position in the input stream.
type Variable interface {
	Element
	Id() string
	Pos() *position
}

// the position of an element in the input stream is defined by it's starting
// and ending byte.
type position struct {
	start int // == offset
	end   int // start - end = len
}

// the end of a position can be moved forward, if the following bytes turn out
// to belong to the same element.
func (p *position) Fwd() { (*p).end++ }

// the end can also be moved backward, if a consumed element turns out not to
// be part of this element.
func (p *position) Rewd() { (*p).end-- }

// returns the length of the element
func (p position) Len() int { return p.end - p.start }

// A Position instance can generate a new position instance whith a starting
// point set to the next byte after this elements ending byte.
func (p position) NextPos() *position {
	return &position{
		p.end + 1,
		p.end + 1,
	}
}
