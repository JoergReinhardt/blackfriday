// Code generated by "stringer -type ValueType"; DO NOT EDIT

package agiledoc

import "fmt"

const _ValueType_name = "EMPTYBOOLUINTINTEGERBYTESTEXTFLOATRATIONALPAIRFLAGLISTSTACKTABLEMATRIXSETMAP"

var _ValueType_map = map[ValueType]string{
	0:     _ValueType_name[0:5],
	2:     _ValueType_name[5:9],
	4:     _ValueType_name[9:13],
	8:     _ValueType_name[13:20],
	16:    _ValueType_name[20:25],
	32:    _ValueType_name[25:29],
	64:    _ValueType_name[29:34],
	128:   _ValueType_name[34:42],
	256:   _ValueType_name[42:46],
	512:   _ValueType_name[46:50],
	1024:  _ValueType_name[50:54],
	2048:  _ValueType_name[54:59],
	4096:  _ValueType_name[59:64],
	8192:  _ValueType_name[64:70],
	16384: _ValueType_name[70:73],
	32768: _ValueType_name[73:76],
}

func (i ValueType) String() string {
	if str, ok := _ValueType_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ValueType(%d)", i)
}
