package blackfriday

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
	//////////////// DOC TYPES BLOCK LEVEL ///
	// complex types that represent semantic block level nodes of the document tree
	CTX_TYPE_SECTION   // reoresents all parts of document, contained by a section
	CTX_TYPE_PARAGRAPH // a block of text divided by newlines
	CTX_TYPE_LIST      // either unordered, ordered, or definition list
	CTX_TYPE_TABLE     // table is a list of lists
	CTX_TYPE_QUOTE     // represents quotet text
	CTX_TYPE_CODE      // represents a code block
	CTX_TYPE_FUNC      // type to store a function (possibly args & rets)
	CTX_TYPE_BLOCK     // handles all other blocks for now.
	//////////////// DOC TYPES INLINE ///
	// complex types that represent semantic inline nodes of the document tree
	CTX_TYPE_LINK   // type representing a link
	CTX_TYPE_FIGURE // type to provide figure inclusion
	CTX_FUNC_DEF    // define a function, possibly expecting parameters and/or returning a value
	CTX_TYPE_REF    // reference to a context variable (expands to value when evaluated)
	CTX_FUNC_CALL   // call of a context function (may evaluate to return value, or trigger side effect)
	CTX_TYPE_SPAN   // parse all other inline  elements for now.
	//////////////// OPERATORS
	CTX_TYPE_UNOP  // called when a unary operator is encountered
	CTX_TYPE_BINOP // called when a binary operator is encountered

	contextExtensions = 0 |
		EXTENSION_CTX_VAR_DEFINITIONS |
		EXTENSION_CTX_VAR_REFERENCES |
		EXTENSION_CTX_FNC_DEFINITIONS |
		EXTENSION_CTX_FNC_CALLS

	contextDefinitions = 0 |
		EXTENSION_CTX_VAR_DEFINITIONS |
		EXTENSION_CTX_FNC_DEFINITIONS

	contextReferences = 0 |
		EXTENSION_CTX_VAR_REFERENCES |
		EXTENSION_CTX_FNC_CALLS
)
