package blackfriday

import "bytes"

// markdown interface to be implemented

// autoLink(p *parser, out *bytes.Buffer, data []byte, offset int) : int
// context specific INLINE callbacks are methods of the parser struct
func ctxBraces(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['(']
}
func ctxDef(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['{']
}
func ctxRef(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['$']
}
func ctxColonOrAutoLink(p *parser, out *bytes.Buffer, data []byte, offset int) { // [':']
}
func ctxProdOrEmphasis(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['*']
}
func ctxBinOpAdd(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['+']
}
func ctxBinOpSubstract(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['-']
}
func ctxBinOpDivide(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['÷']
}
func ctxBinOpDot(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['·']
}
func ctxBinOpCross(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['×']
}
func ctxBinOpEqual(p *parser, out *bytes.Buffer, data []byte, offset int) { // ['=']
}

// common INLINE CALLBACKS
func (c *Context) CodeSpan(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) DoubleEmphasis(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) Emphasis(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) Image(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) LineBreak(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) Link(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) RawHtmlTag(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) TripleEmphasis(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) StrikeThrough(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) Entity(p *parser, out *bytes.Buffer, data []byte, offset int) {
}
func (c *Context) NormalText(p *parser, out *bytes.Buffer, data []byte, offset int) {
}

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
func (c *Context) TitleBlock(out *bytes.Buffer, text []byte) {
}

//
// DOCUMENT HEADER AND FOOTER
func (c *Context) DocumentHeader(out *bytes.Buffer) {
}
func (c *Context) DocumentFooter(out *bytes.Buffer) {
}

// SECTION (identifyed by header tag, or slug of title)
func (c *Context) Header(out *bytes.Buffer, text func() bool, level int, id string) {
}
func (c *Context) Paragraph(out *bytes.Buffer, text func() bool) {
}

// LIST
func (c *Context) List(out *bytes.Buffer, text func() bool, flags int) {
}
func (c *Context) ListItem(out *bytes.Buffer, text []byte, flags int) {
}

// TABLE
func (c *Context) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
}
func (c *Context) TableRow(out *bytes.Buffer, text []byte) {
}
func (c *Context) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
}
func (c *Context) TableCell(out *bytes.Buffer, text []byte, flags int) {
}

// BLOCKS OF UNIFORM CONTENT
func (c *Context) BlockCode(out *bytes.Buffer, text []byte, lang string) {
}
func (c *Context) BlockQuote(out *bytes.Buffer, text []byte) {
}
func (c *Context) BlockHtml(out *bytes.Buffer, text []byte) {
}
func (c *Context) HRule(out *bytes.Buffer) {
}

// FOOTNOTES
func (c *Context) Footnotes(out *bytes.Buffer, text func() bool) {
}
func (c *Context) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
}
func (c *Context) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
}

// DOCUMENT-AUTOLINKER (basicly tunnel through
func (c *Context) AutoLink(out *bytes.Buffer, link []byte, kind int) {
}
func (c *Context) GetFlags() int {
	return c.flags
}
