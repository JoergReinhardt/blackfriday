// Code generated by "stringer -type ValType types.go"; DO NOT EDIT

package agiledoc

import "fmt"

const (
	_ValType_name_0 = "NIL"
	_ValType_name_1 = "FLAG"
	_ValType_name_2 = "INTEGER"
	_ValType_name_3 = "FLOAT"
	_ValType_name_4 = "BYTE"
	_ValType_name_5 = "STRING"
)

var (
	_ValType_index_0 = [...]uint8{0, 3}
	_ValType_index_1 = [...]uint8{0, 4}
	_ValType_index_2 = [...]uint8{0, 7}
	_ValType_index_3 = [...]uint8{0, 5}
	_ValType_index_4 = [...]uint8{0, 4}
	_ValType_index_5 = [...]uint8{0, 6}
)

func (i ValType) String() string {
	switch {
	case i == 0:
		return _ValType_name_0
	case i == 2:
		return _ValType_name_1
	case i == 4:
		return _ValType_name_2
	case i == 8:
		return _ValType_name_3
	case i == 16:
		return _ValType_name_4
	case i == 32:
		return _ValType_name_5
	default:
		return fmt.Sprintf("ValType(%d)", i)
	}
}