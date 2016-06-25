/*
CONTEXTUAL AGILE DOCUMENT

Ctx is supposed to implement a Representation of a tree of interlinked markdown
Documents and other data sources. Syntactic parts of the document are
represented as nodes in a abstract syntax tree, spanning all docyments and
other resolveable data sources, referenced by, or included into one or more
documents constituting the tree,

Lists and Tables will be parsed as vector and matrix datatypes, Their terminal
elements (list elements or table fields), get parsed as integer, float, string,
or one of various operators and keywords, specified by the ctx syntax
extension.

All elements are propagated to the symbol table and mapped by either their
identifyer, or a string containing of the type and a numeral representing of
the place in the sequence of elements of the same type, contained by the same
parent.

The symbol table is a data structure distributed over all nodes constituting
the global (forrest of) tree(s), Each node only containing variables defined
directly within the content, represented by this node,  All elements of the
symbol table can be referenced, accessed and manipulated by all nodes of the
tree, preventing naming conflicts by providing „full names” consisting of a
concatination of all the ancestors names for each variable.

The public interface of Blackfriday enforces the implementation of a
„Renderer”,

To implement an agile document, that keeps propagating values and
backreferences in realtime, an abstract representation of the whole document
tree is needed.

The abstract syntax tree would idealy be parsed by evalueating a stream of
tokens, generated and tagged by a lexer. Instead the containing parser struct
will call the callback functions defined by the renderer interface, pass raw
data in form of a byte slice and additional information in the form of
parameters, depending on the callback in question. An output buffer gets passed
to the callbacks with each call, expecting to get written to by the called
callback.

The parser generating functions of Blackfriday don't provide access to the
generated parser, but expect a byteslice containing the docyment, an instance
of a renderer and optionally an options instance and return a byteslice, that's
supposed to contain the rendered document, Parser and passed Renderer are
deallocated after each call.

A containing data structure to orchestrate the calls to the Parser, which get
nesccessary after each edit of the documents contents, is needed. It has to
provide the Renderer instance passed to the function generating the parser. The
Renderer has to contain a back reference to a part of the containing data
structure, representing the documents contents. That reference needs to be
mutable from within a call to the blackfriday parser,

The ast provides a data structure to hold the abstract syntax tree and to
coordinate all operations there up on, Each time, the input document changes, a
Renderer is instanciated and passed to the blackfriday function that parses it.
The sequence of calls to the renderers callback functions will be serialized to
and interpreted as a sequence of tokens, just like a lexer would provide them.
Those tokens get compared to and replace the last known sequence of tokens
representing the document, All tokens that changed will trigger the parser to
rebuild the tree of syntactic nodes representing the document, which will
replace the former syntactic representation held by the ast,

A tree walk is performed in turn, to propagate all new values to the symbol
table and resolve all references to their current values, evaluate all
calculations and functioncalls, defined in the documents elements, propagate
the results, rinse and repeat until tree is fully evaluated, finally walk tree
in sequential order of nodes reference to the token stream, to render the
document described by the syntax tree.
*/
package blackfriday
