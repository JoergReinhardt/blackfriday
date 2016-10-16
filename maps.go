package agiledoc

////////////////////////////////////////////////////////////////////////////////////
//// MAPS ////
//////////////
func (m HashMap) Add(v ...Evaluable) (r Mapped) {
	for i, v := range v {
		if v, ok := v.(Tupled); ok {
			k := v.Key()
			v := v.Value()
			r = putToMap(m, k, v)
		}
		r = putToMap(m, Value(i), v)
	}
	return r
}
func (m HashMap) Eval() Evaluable                     { return evalCollection(m()) }
func (m HashMap) Type() ValueType                     { return MAP }
func (m HashMap) Size() int                           { return collectionSize(m()) }
func (m HashMap) Empty() bool                         { return emptyCollection(m()) }
func (m HashMap) Clear() Collected                    { return clearCollection(m()) }
func (m HashMap) Put(k Evaluable, v Evaluable) Mapped { return putToMap(m, k, v) }
func (m HashMap) Get(v Evaluable) (Evaluable, bool)   { return getFromMap(m, v) }
func (m HashMap) Keys() []Evaluable                   { return keysOfMap(m) }
func (m HashMap) Values() []Evaluable                 { return collectionValues(m()) }
func (m HashMap) Remove(v Evaluable) Mapped           { return removeFromMap(m, v) }
func (m HashMap) Serialize() []byte                   { return serializeMap(m) }
func (m HashMap) Interfaces() []interface{}           { return interfacesFromMap(m) }
func (m HashMap) String() string                      { return mapToString(m) }

func (m HashBidiMap) Add(v ...Evaluable) (r Mapped) {
	for i, v := range v {
		if v, ok := v.(Tupled); ok {
			k := v.Key()
			v := v.Value()
			r = putToMap(m, k, v)
		}
		r = putToMap(m, Value(i), v)
	}
	return r
}
func (m HashBidiMap) Eval() Evaluable                     { return evalCollection(m()) }
func (m HashBidiMap) Type() ValueType                     { return MAP }
func (m HashBidiMap) Size() int                           { return collectionSize(m()) }
func (m HashBidiMap) Empty() bool                         { return emptyCollection(m()) }
func (m HashBidiMap) Clear() Collected                    { return clearCollection(m()) }
func (m HashBidiMap) Put(k Evaluable, v Evaluable) Mapped { return putToMap(m, k, v) }
func (m HashBidiMap) Get(v Evaluable) (Evaluable, bool)   { return getFromMap(m, v) }
func (m HashBidiMap) Keys() []Evaluable                   { return keysOfMap(m) }
func (m HashBidiMap) Values() []Evaluable                 { return collectionValues(m()) }
func (m HashBidiMap) Remove(v Evaluable) Mapped           { return removeFromMap(m, v) }
func (m HashBidiMap) Serialize() []byte                   { return serializeMap(m) }
func (m HashBidiMap) Interfaces() []interface{}           { return interfacesFromMap(m) }
func (m HashBidiMap) String() string                      { return mapToString(m) }

func (m TreeMap) Add(v ...Evaluable) (r Mapped) {
	for i, v := range v {
		if v, ok := v.(Tupled); ok {
			k := v.Key()
			v := v.Value()
			r = putToMap(m, k, v)
		}
		r = putToMap(m, Value(i), v)
	}
	return r
}
func (m TreeMap) Eval() Evaluable                     { return evalCollection(m()) }
func (m TreeMap) Type() ValueType                     { return MAP }
func (m TreeMap) Size() int                           { return collectionSize(m()) }
func (m TreeMap) Empty() bool                         { return emptyCollection(m()) }
func (m TreeMap) Clear() Collected                    { return clearCollection(m()) }
func (m TreeMap) Put(k Evaluable, v Evaluable) Mapped { return putToMap(m, k, v) }
func (m TreeMap) Get(v Evaluable) (Evaluable, bool)   { return getFromMap(m, v) }
func (m TreeMap) Keys() []Evaluable                   { return keysOfMap(m) }
func (m TreeMap) Values() []Evaluable                 { return collectionValues(m()) }
func (m TreeMap) Remove(v Evaluable) Mapped           { return removeFromMap(m, v) }
func (m TreeMap) Serialize() []byte                   { return serializeMap(m) }
func (m TreeMap) Interfaces() []interface{}           { return interfacesFromMap(m) }
func (m TreeMap) String() string                      { return mapToString(m) }

func (m TreeBidiMap) Add(v ...Evaluable) (r Mapped) {
	for i, v := range v {
		if v, ok := v.(Tupled); ok {
			k := v.Key()
			v := v.Value()
			r = putToMap(m, k, v)
		}
		r = putToMap(m, Value(i), v)
	}
	return r
}
func (m TreeBidiMap) Eval() Evaluable                     { return evalCollection(m()) }
func (m TreeBidiMap) Type() ValueType                     { return MAP }
func (m TreeBidiMap) Size() int                           { return collectionSize(m()) }
func (m TreeBidiMap) Empty() bool                         { return emptyCollection(m()) }
func (m TreeBidiMap) Clear() Collected                    { return clearCollection(m()) }
func (m TreeBidiMap) Put(k Evaluable, v Evaluable) Mapped { return putToMap(m, k, v) }
func (m TreeBidiMap) Get(v Evaluable) (Evaluable, bool)   { return getFromMap(m, v) }
func (m TreeBidiMap) Keys() []Evaluable                   { return keysOfMap(m) }
func (m TreeBidiMap) Values() []Evaluable                 { return collectionValues(m()) }
func (m TreeBidiMap) Remove(v Evaluable) Mapped           { return removeFromMap(m, v) }
func (m TreeBidiMap) Serialize() []byte                   { return serializeMap(m) }
func (m TreeBidiMap) Interfaces() []interface{}           { return interfacesFromMap(m) }
func (m TreeBidiMap) String() string                      { return mapToString(m) }
