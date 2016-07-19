package agiledoc

import (
	// "github.com/emirpasic/gods/containers"
	"github.com/emirpasic/gods/lists"
	// 	"github.com/emirpasic/gods/lists/arraylist"
	// 	"github.com/emirpasic/gods/lists/doublylinkedlist"
	// 	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps"
	// 	 "github.com/emirpasic/gods/maps/hashbidimap"
	// 	 "github.com/emirpasic/gods/maps/hashmap"
	// 	"github.com/emirpasic/gods/maps/treebidimap"
	// 	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets"
	// 	"github.com/emirpasic/gods/sets/hashset"
	// 	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/stacks"
	// 	"github.com/emirpasic/gods/stacks/arraystack"
	// 	"github.com/emirpasic/gods/stacks/linkedliststack"
	// "github.com/emirpasic/gods/trees"
	// 	"github.com/emirpasic/gods/trees/binaryheap"
	// 	"github.com/emirpasic/gods/trees/redblacktree"
	// 	"github.com/emirpasic/gods/utils"
)

func (v keyVal) Key() string { return v.id }

func (v cntVal) ContType() CntType { return v.CntType }
func (v *cntVal) Values() (r []Val) {
	r = []Val{}
	for _, val := range (*v).Values() {
		val := val.Value()
		r = append(r, val)
	}
	return r
}

func (v *cntVal) Add(vals ...Val) {
	switch {
	case v.ContType()&LISTS != 0:
		for _, val := range vals {
			val := val.Value()
			(*v).Container.(lists.List).Add(val)
		}
	case v.ContType()&SETS != 0:
		for _, val := range vals {
			val := val.Value()
			(*v).Container.(sets.Set).Add(val)
		}
	case v.ContType()&STACKS != 0:
		for _, val := range vals {
			val := val.Value()
			(*v).Container.(stacks.Stack).Push(val)
		}
	case v.ContType()&MAPS != 0:
		for i, kv := range vals {
			i := i
			kv := kv.Value()
			var key string
			// if the value is of type keyVal, it will implement
			// the Var interface and have a key method, which is
			// used to set the map key, otherwise take string
			// representation of the intefer index of current
			// element.
			if kv.Type()&KEYVAL != 0 {
				key = kv.(Var).Key()
			} else {
				key = string(i)
			}
			(*v).Container.(maps.Map).Put(key, kv.Value())
		}
		// 	case v.ContType()&REDBLACK != 0:
		// 	case v.ContType()&BINHEAP != 0:
	}
}
