package agiledoc

import ()

func (m TreeMap) Add(v ...Evaluable) Mapped     { return m }
func (m TreeBidiMap) Add(v ...Evaluable) Mapped { return m }
func (m HashMap) Add(v ...Evaluable) Mapped     { return m }
func (m HashBidiMap) Add(v ...Evaluable) Mapped { return m }
