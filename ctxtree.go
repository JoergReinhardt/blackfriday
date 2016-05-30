package blackfriday

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
