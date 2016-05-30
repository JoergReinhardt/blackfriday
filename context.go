package blackfriday

import "bytes"

// Element Types as passed to the context renderer by the parser.
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

// All node implementations need to implement the node interface to be
// applicable as element of a tree.
type Node interface {
	parent() *Tree
	Type() int
	Content() []byte
}

// Element ist the broughtest implementation of the node interface. It is
// capable of representing every document, all nested Elements included. by
// representing the content in a private function.
type Element struct {
	typ     int           // node type
	content func() []byte //
	parent  func() *Tree  // function returning a link to the containing node
}

func (e *Element) Type() int { return (*e).typ }

func (n *Element) Content() []byte { return (*n).content() }

// Tree Node contains a list of Node interface instances. Since Tree itself
// implements the Node interface, this type is recursive.
type Tree struct {
	*Element        // identity and content of this tree of nodes
	content  []Node // slice of contained nodes
}

func (t *Tree) Content() (c []byte) {
	c = []byte{}
	// call content function of all contained nodes and concatenate output
	for _, n := range (*t).content {
		c = append(c, n.Content()...)
	}
	return c
}

// List Node contains a list of references to plain  elements
type List struct {
	*Element            // identity and content of this tree of nodes
	content  []*Element // slice of contained nodes
}

func (l *List) Content() (c []byte) {
	c = []byte{}
	// call content function of all contained nodes and concatenate output
	for _, n := range (*l).content {
		c = append(c, n.Content()...)
	}
	return c
}

func NewRoot() *Tree {
	// allocate tree, leave node reference empty
	t := Tree{
		nil,
		[]Node{},
	}
	// allocate node, set type to root, content to an empty byte slice and
	// a function to return a reference to the tree instance
	n := Element{
		typ: 0,
		content: func() (b []byte) {
			b = []byte{}
			for _, n := range (&t).content {
				b = append(b, n.Content()...)
			}
			return b
		},
		parent: func() *Tree { return &t },
	}

	// set trees identity node referencing its own tree as its parent.
	// !!! ENDLESS LOOP AHEAD, ALLWAYS CHECK NODE TYPE FOR ROOT !!!
	t.Element = &n

	// return  reference to the root tree
	return &t
}

////////////////////////////////////////////////////////////////////////////

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
type valt struct {
	typ int
	val []byte
}

//the context struct contains the symtab, that maps values to identity and
//type. It implements the variable type for all types of variables alike,
//Symtab key is a two field struct containing of the types identity and type.
//The Value contains the raw value in a byte slice and can return itself.
type Context struct {
	flags  int
	symtab map[string]valt
}

func (c *Context) Lookup(id string) (v Variable) {
	return v
}

///// IMPLEMENTING THE RENDERER API /////
// block-level callbacks
func (c *Context) BlockCode(out *bytes.Buffer, text []byte, lang string)                 {}
func (c *Context) BlockQuote(out *bytes.Buffer, text []byte)                             {}
func (c *Context) BlockHtml(out *bytes.Buffer, text []byte)                              {}
func (c *Context) Header(out *bytes.Buffer, text func() bool, level int, id string)      {}
func (c *Context) HRule(out *bytes.Buffer)                                               {}
func (c *Context) List(out *bytes.Buffer, text func() bool, flags int)                   {}
func (c *Context) ListItem(out *bytes.Buffer, text []byte, flags int)                    {}
func (c *Context) Paragraph(out *bytes.Buffer, text func() bool)                         {}
func (c *Context) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {}
func (c *Context) TableRow(out *bytes.Buffer, text []byte)                               {}
func (c *Context) TableHeaderCell(out *bytes.Buffer, text []byte, flags int)             {}
func (c *Context) TableCell(out *bytes.Buffer, text []byte, flags int)                   {}
func (c *Context) Footnotes(out *bytes.Buffer, text func() bool)                         {}
func (c *Context) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int)          {}
func (c *Context) TitleBlock(out *bytes.Buffer, text []byte)                             {}

// Span-level callbacks
func (c *Context) AutoLink(out *bytes.Buffer, link []byte, kind int)                 {}
func (c *Context) CodeSpan(out *bytes.Buffer, text []byte)                           {}
func (c *Context) DoubleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (c *Context) Emphasis(out *bytes.Buffer, text []byte)                           {}
func (c *Context) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte)    {}
func (c *Context) LineBreak(out *bytes.Buffer)                                       {}
func (c *Context) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {}
func (c *Context) RawHtmlTag(out *bytes.Buffer, tag []byte)                          {}
func (c *Context) TripleEmphasis(out *bytes.Buffer, text []byte)                     {}
func (c *Context) StrikeThrough(out *bytes.Buffer, text []byte)                      {}
func (c *Context) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                 {}

// Low-level callbacks
func (c *Context) Entity(out *bytes.Buffer, entity []byte)   {}
func (c *Context) NormalText(out *bytes.Buffer, text []byte) {}

// Header and footer
func (c *Context) DocumentHeader(out *bytes.Buffer) {}
func (c *Context) DocumentFooter(out *bytes.Buffer) {}

func (c *Context) GetFlags() int { return c.flags }

// INLINE CALLS ADDED TO PARSER
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
