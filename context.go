package blackfriday

import "bytes"

// markdown interface to be implemented

// autoLink(p *parser, out *bytes.Buffer, data []byte, offset int) : int
// context specific INLINE callbacks are methods of the parser struct
var (
	ctxBraces          = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxDef             = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxRef             = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxColonOrAutoLink = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxProdOrEmphasis  = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpAdd        = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpSubstract  = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpDivide     = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpDot        = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpCross      = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	ctxBinOpEqual      = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	// common INLINE CALLBACKS
	CodeSpan       = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	DoubleEmphasis = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	Emphasis       = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	Image          = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	LineBreak      = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	Link           = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	RawHtmlTag     = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	TripleEmphasis = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	StrikeThrough  = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	Entity         = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }
	NormalText     = func(p *parser, out *bytes.Buffer, data []byte, offset int) int { return offset }

	// BLOCK level elements are implemended, kind of mis-using the rendere
	// interface by using it to implemet the context parser, simply because there
	// is no other public api to imolement block level parsing.
	//
	// to get output rendered, once the AST is evaluated, the renderer calls html
	// render by default, or latex if set in the options.
	//
	//It doesn't generate any output by itselfe, but populates the AST, propagates
	//Values and backtracks references.
	//
	// DOCUMENT LEVEL BLOCKS
	//
	// TITLE BLOCK
	TitleBlock = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	// DOCUMENT HEADER AND FOOTER
	DocumentHeader = func(p *parser, out *bytes.Buffer, offset int) int { return offset }
	DocumentFooter = func(p *parser, out *bytes.Buffer, offset int) int { return offset }
	// SECTION (identifyed by header tag, or slug of title)
	Header    = func(p *parser, out *bytes.Buffer, text []byte, id string, offset int) int { return offset }
	Paragraph = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	// LIST
	List     = func(p *parser, out *bytes.Buffer, text []byte, flags int, offset int) int { return offset }
	ListItem = func(p *parser, out *bytes.Buffer, text []byte, flags int, offset int) int { return offset }
	// TABLE
	Table = func(p *parser, out *bytes.Buffer, header []byte, body []byte, columnData []int, offset int) int {
		return offset
	}
	TableRow        = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	TableHeaderCell = func(p *parser, out *bytes.Buffer, text []byte, flags int, offset int) int { return offset }
	TableCell       = func(p *parser, out *bytes.Buffer, text []byte, flags int, offset int) int { return offset }
	// BLOCKS OF UNIFORM CONTENT
	BlockCode  = func(p *parser, out *bytes.Buffer, text []byte, lang string, offset int) int { return offset }
	BlockQuote = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	BlockHtml  = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	HRule      = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	// FOOTNOTES
	Footnotes    = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	FootnoteItem = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	FootnoteRef  = func(p *parser, out *bytes.Buffer, text []byte, offset int) int { return offset }
	// DOCUMENT-AUTOLINKER (basicly tunnel through
	AutoLink = func(p *parser, out *bytes.Buffer, link []byte, offset int) int { return offset }
	GetFlags = func(p *parser, out *bytes.Buffer, link []byte, offset int) int { return offset }
)
