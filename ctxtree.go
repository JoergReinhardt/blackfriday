package blackfriday

// All node implementations need to implement the node interface to be
// applicable as element of a tree.
type Node interface {
	parent() *Tree
	ElementType() int
}

//// ELEMENT NODE
///
// Element ist the most generic implementation of a node carrying a value. It
// also implements the Variable Interface by keeping id and position in its own
// fields.
type ElementNode struct {
	id       string // since a value has no identity per default
	typ      int
	position              // which part of document input is represented here
	parent   func() *Tree // function returning a link to the containing node
	// Eval(), Type(), ToType() -->
	Value // contains Document Element struct as bassed by parser
}

// let each element implement variable interface per default
func (e ElementNode) Id() string       { return e.id }
func (e ElementNode) ElementType() int { return e.typ }
func (e *ElementNode) Pos() position   { return e.position }

//// LIST NODE (keeps sequence order)
///
// Can be one of unordered, ordered, or definition list. may be nested If
// unidentified values are part of the sequence, generate missing id's
// sequentialy
type ListNode struct {
	id  string
	typ int
	position
	list []Node // document elements, structural nodes, single values...
}

func (e ListNode) Id() string       { return e.id }
func (e ListNode) ElementType() int { return e.typ }
func (e ListNode) Pos() position    { return e.position }

//// TREE DATA TYPE
///
// Tree Node contains a list of Node interface instances. Since Tree itself
// implements the Node interface, this type is recursive.
type Tree struct {
	id  string // identity and content of this tree of nodes
	typ int
	position
	parent func() *Tree
	nodes  []Node // slice of contained nodes
}

//// NEWROOT
///
// context is rooted in a tree, returning itself as its parent.
func NewRoot() *Tree {
	// allocate tree, leave node reference empty
	t := Tree{
		id:       "ROOT",
		typ:      ROOT,
		position: position{0, 0}, // start, end
		parent:   func() *Tree { return &Tree{} },
		nodes:    []Node{}, // contained nodes
	}

	// return  reference to the root tree
	return &t
}
