/*
collection types are a wee bit harder to auto-create, since there are more
different principle roles, a collection can take depending on the contained
type, as well, as on the way the collection structures the data:

  - flat lists, of values of one or several types, indexed by position,
  - maps, with values mapped on to keys of a certain type.
  - semantic lists, that feature rows, of possibly different length.
  - semantic tables with colums and rows, possibly named.
  - numeric vectors with indexed fields,
  - numeric matrices featuring two or more dimensions,

Many of those roles can. should, or even must be implemented in a certain way,
to provide the features needed. The performance those features can achieve is
highly dependent on the implementation and the type, they are intended to deal
with. Unnessccessary features can generate considerable overhead and so on…

To represent all those possibylities, the collection interface is layered. The
base level is represented by the collection interface, basically wraping
gods/container. The second layer provides differentiation between Lists,
Stacks, Sets, Maps and Trees, again wrapping gods second level interfaces. In
its eternally gods represents all features, collections can come with,
implemented in most efficient ways. On top of that, a third level is provided,
that features division between semantic role, a collections is supposed to take
(see above).

The first layers are determinable by the type(s) of values passed to the
generator. the Constructor returns a sensible default instance, of a second
level type. This is achieved in a two step process, creating an array based
list of all passed values at first, using its enumerables, to analyze which
type(s) are given and convert to the appropriate second level type, by unifying
the contained value types if nesccessary. Which type to convert contained
values to, is decided so that no information is lost. In case of a simple list
containing values of identical type (including tuples with numeric key), the
created array list is returned as the default type. All pairs, with symbolic
key types, get wrapped in a hash map. All second level types come with methods
to convert between one another, and to all third level types that can be
implemented, by the method set of that particular second level type,

Its is up to the user of the library, to determine which role instances are
supposed to take in her, or his application and to choose an appropriate second
level type directly, or one of the existing, specialized implementations of
third level typs and possinly define and implement additional ones.
implementations of third level types. conversion to third level may again
involve two steps, to first generate the appropriate second level type, that
comes with the nesccessary features to implement- and therefore comes with the
methods to convert to the wanted third level type.

         ┌──────┬────────────┐
   Big	 │ Rat	│ Lst,Map,Set│
   ──────┼──────┼────────────┤
   Simple│ Tuple│ Collection │
  	 └──────┴────┬───────┘
                     │
		     ▼
  		   paired
*/
package agiledoc

import (
	//al "github.com/emirpasic/gods/lists/arraylist"
	"math/big"
)

func nativeToValue(i interface{}) (r Evaluable) {

	switch i.(type) {
	case bool: // a boolean returns a flag with the first bit set
		if i.(bool) {
			r = wrap(intPool.Get().(*big.Int).SetInt64(1)).(BitFlag)
		} else {
			r = wrap(intPool.Get().(*big.Int).SetInt64(0)).(BitFlag)
		}
	case uint, uint8, uint16, uint32, uint64, ValueType: // a uint is assumed to be a single byte
		r = divideUints(i)
	case int, int16, int32, int64: // integers are integer
		r = divideInts(i)
	case float32: // floating point values get assigned to rationals
		r = wrap(ratPool.Get().(*big.Rat).SetFloat64(float64(i.(float32))))
	case float64: // floating point values get assigned to rationals
		r = wrap(ratPool.Get().(*big.Rat).SetFloat64(i.(float64)))
	case []byte: // == uint8
		r = wrap(intPool.Get().(*big.Int).SetBytes(i.([]byte)))
	case string: // a string gets assigned by its bislice as well
		val, ok := intPool.Get().(*big.Int).SetString(i.(string), 10)
		if ok {
			r = wrap(val)
		}
	}
	return r
}
func divideUints(i interface{}) (r Evaluable) {
	switch i.(type) {
	case ValueType:
		r = wrap(intPool.Get().(*big.Int).SetUint64(uint64(i.(ValueType))))
	case byte:
		r = wrap(intPool.Get().(*big.Int).SetBytes([]byte(i.([]byte))))
	case rune:
		r = wrap(intPool.Get().(*big.Int).SetUint64(uint64(i.(uint32))))
	case uint:
		r = wrap(intPool.Get().(*big.Int).SetUint64(uint64(i.(uint))))
	case uint16:
		r = wrap(intPool.Get().(*big.Int).SetUint64(uint64(i.(uint16))))
	case uint64:
		r = wrap(intPool.Get().(*big.Int).SetUint64(i.(uint64)))
	}
	return r
}
func divideInts(i interface{}) (r Evaluable) {
	switch i.(type) {
	case int:
		r = wrap(intPool.Get().(*big.Int).SetInt64(int64(i.(int))))
	case int16:
		r = wrap(intPool.Get().(*big.Int).SetInt64(int64(i.(int16))))
	case int32:
		r = wrap(intPool.Get().(*big.Int).SetInt64(int64(i.(int32))))
	case int64:
		r = wrap(intPool.Get().(*big.Int).SetInt64(int64(i.(int64))))
	}
	return r
}

// collection generator determines which type of collection to allocate, based
// on the parameters it gets passed
func Collect(v ...Evaluable) (r Collected) {
	return r
}
