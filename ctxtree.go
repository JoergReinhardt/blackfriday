package blackfriday

import "sync"

//// NODE INTERFACE
///
// All node implementations need to implement the node interface to be
// applicable as element of a tree.
//
// each node has a private function to refere to it's parent and a Function to
// return it's contained elements type
type Node interface {
	parent() *Tree
	Type() uint
}

//// BASE NODE
///
// implements id, type and reference to a parent node
type BaseNode struct {
	id  string // since a value has no identity per default
	typ uint
	par func() *Tree // function returning a link to the containing node
}

// let each element implement variable interface per default
func (e BaseNode) Id() string        { return e.id }
func (e BaseNode) ElementType() uint { return e.typ }
func (e BaseNode) parent() *Tree     { return e.parent() }

func newBaseNode(id string, typ uint, parent *Tree) BaseNode {
	return BaseNode{
		id:  id,
		typ: typ,
		par: func() *Tree { return parent },
	}
}

//// ELEMENT NODE
///
// Element ist the most generic implementation of a node carrying a value. It
// also implements the Variable Interface by keeping id and position in its own
// fields.
type ElementNode struct {
	BaseNode
	Element // contains Document Element as bassed by parser value
	// implements Eval(), Type(), ToType(), so every ElementNode is also a Element
}

func newElementNode(base BaseNode, element Element) *ElementNode {
	return &ElementNode{
		base,
		element,
	}
}

//// VECTOR NODE (keeps sequence order)
///
// Can be one of unordered-, ordered-, linked-, or mapped (list of definitions)
// list. may be nested If unidentified values are part of the sequence. missing
// id's of anonymous elements of a list are genereated by using their place in
// the containing seqyuence
type VectorNode struct {
	BaseNode
	Vector // document elements, structural nodes, single values...
}

func newVectorNode(base BaseNode, vec Vector) *VectorNode {
	return &VectorNode{
		base,
		vec,
	}
}

//// MATRIX NODE
///
// A matrix node contains a slice of Node instances that contains the element
// of a matrix. The shape field contains the number of rows and columns those
// elements are spread about.
type MatrixNode struct {
	BaseNode
	shape  [2]int
	matrix Vector
}

func newMatrixNode(base BaseNode, mtx []Vector) *MatrixNode {

	var shape [2]int = [2]int{len(mtx), 0}

	// range over rows to determine shape
	for _, e := range mtx {
		// find longest row, update shape
		if e.Len() > shape[1] {
			v := e.Len()
			shape[1] = v
		}
	}

	// allocate array of size number of rows multiplied by number of
	// columns of longest row to get an indexed value for every possible
	// element in the matrix
	var matrix = make([]Vector, shape[0]*shape[1])

	// range over rows and generate matrix
	for o, r := range mtx {
		o := o
		r := r

		// range over elements of current row
		for i, e := range r {
			i := i
			e := e

			// allocate integer to hold position.
			var pos int

			// if on first line outer index is zero and cant be
			// further reduced.
			if o > 0 {
				// number of previous rows
				o = o - 1
			}
			// Some rows may be shorter than matrix wide, so
			// position is the sum of the fields of all previeous
			// rows as if they where full length. plus current rows
			// index.
			pos := (o * shape[1]) + i

			// assign current element at its position
			matrix[pos] = e
		}
	}

	return &MatrixNode{
		base,
		shape,
		matrix,
	}
}

// A tree is nothing more than a slice of nodes, contained in a value of type
// vector. The elements are interlinked to a tree by their parent() methods.
type Tree struct {
	*VectorNode
}

//// SYM TAB
///
// maps variable names to their corresponding elemnt values.
type symTab struct {
	*sync.RWMutex
	stab map[string]Element
}

func newSymTab() symTab {
	return symTab{
		&sync.RWMutex{},
		make(map[string]Element),
	}
}

//// NEWROOT
///
// The Whole document context is described as a tree, that needs to be rooted.
// NewRoot() returns a reference to a tree, that returns a reference to itself
// as it's parent. TODO: find out if that's actually such a good idea.
// Returning an empty tree would be another possibility.
func NewRoot(name string) *Tree {
	// allocate tree, define parent as reference to an empty Tree.
	t := Tree{
		&ElementNode{
			NewBaseNode(name, ROOT, nil), // parent left empty for now
		},
		nodes:  []Node{}, // contained nodes
		symtab: newSymTab(),
	}
	// replace parent function to return a reference to itself
	t.parent = func() *Tree { return &t }

	// return  reference to the root tree
	return &t
}
