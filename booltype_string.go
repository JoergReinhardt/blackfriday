// Code generated by "stringer -type BoolType"; DO NOT EDIT

package agiledoc

import "fmt"

const (
	_BoolType_name_0 = "NATIVE"
	_BoolType_name_1 = "LISTED"
	_BoolType_name_2 = "SIGNED"
	_BoolType_name_3 = "BIT_FLAG"
	_BoolType_name_4 = "UINT_FLAG"
	_BoolType_name_5 = "VAL_TYPE"
	_BoolType_name_6 = "TOKEN_TYPE"
	_BoolType_name_7 = "NODE_TYPE"
)

var (
	_BoolType_index_0 = [...]uint8{0, 6}
	_BoolType_index_1 = [...]uint8{0, 6}
	_BoolType_index_2 = [...]uint8{0, 6}
	_BoolType_index_3 = [...]uint8{0, 8}
	_BoolType_index_4 = [...]uint8{0, 9}
	_BoolType_index_5 = [...]uint8{0, 8}
	_BoolType_index_6 = [...]uint8{0, 10}
	_BoolType_index_7 = [...]uint8{0, 9}
)

func (i BoolType) String() string {
	switch {
	case i == 0:
		return _BoolType_name_0
	case i == 2:
		return _BoolType_name_1
	case i == 4:
		return _BoolType_name_2
	case i == 8:
		return _BoolType_name_3
	case i == 16:
		return _BoolType_name_4
	case i == 32:
		return _BoolType_name_5
	case i == 64:
		return _BoolType_name_6
	case i == 128:
		return _BoolType_name_7
	default:
		return fmt.Sprintf("BoolType(%d)", i)
	}
}
