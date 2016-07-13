package agiledoc

type (
	tokenType uint32
	ValueType uint16
)

//go:generate -command stringer -type ValueType ./
const (
	// BASIC VALUE TYPES// {{{
	EMPTY ValueType = 0
	BOOL  ValueType = 1 << iota
	INTEGER
	FLOAT
	BYTE
	BYTES
	STRING
	VECTOR
	MATRIX
	// }}}
	// SINGLE & MULTI XOR MAP// {{{
	singleValueTypes   = BOOL | INTEGER | FLOAT | STRING
	multipleValueTypes = VECTOR | MATRIX // }}}
	// BLOCK LEVEL NODES// {{{
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
	F_REF // footnote reference// }}}
	// LOW LEVEL NODES// {{{
	ENTITY
	TEXT // }}}
	// XORED SETS OF COMMON NODES// {{{
	blockElements = DOCUMENT | D_HEADER | D_FOOTER | SECTION | TITLE | PARAGRAPH | CODE | QUOTE | HTML | HRULE | LIST | L_ITEM | TABLE | T_HEADER_CELL | T_ROW | T_CELL | FOOTNOTES | F_ITEM

	spanElements = AUTO_LINK | CODE_SPAN | LINE_BREAK | EMPHASIS | DOUBLE_EMPHASIS | TRIPLE_EMPHASIS | STRIKE_THROUGHT | RAW_HTML_TAG | LINK | IMAGE | F_REF

	lowLevelElements = ENTITY | TEXT // }}}
)

//// VALUE INTERFACE
///
// a value has a type and can be evaluated to a byte slice representation of
// it's content
type Value interface {
	Type() ValueType
	ToType(ValueType) Value
	Eval() []byte
}
