// Code generated by "stringer -type ValueType ./"; DO NOT EDIT

package agiledoc

import "fmt"

const _ValueType_name = "EMPTYBOOLUINTINTEGERBYTESSTRINGFLOATRATIONALTUPLEFLAGLISTMATRIXKEYMAP"

var _ValueType_map = map[ValueType]string{
	0:    _ValueType_name[0:5],
	2:    _ValueType_name[5:9],
	4:    _ValueType_name[9:13],
	8:    _ValueType_name[13:20],
	16:   _ValueType_name[20:25],
	32:   _ValueType_name[25:31],
	64:   _ValueType_name[31:36],
	128:  _ValueType_name[36:44],
	256:  _ValueType_name[44:49],
	512:  _ValueType_name[49:53],
	1024: _ValueType_name[53:57],
	2048: _ValueType_name[57:63],
	4096: _ValueType_name[63:69],
}

func (i ValueType) String() string {
	if str, ok := _ValueType_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ValueType(%d)", i)
}
