// Code generated by "stringer -type Number"; DO NOT EDIT

package types

import "fmt"

const (
	_Number_name_0 = "NEGATIVEZERO"
	_Number_name_1 = "ONETWOTHREEFOURFIVESIXSEVENEIGHTNINETENELEVENTWELVETHIRTEEN"
)

var (
	_Number_index_0 = [...]uint8{0, 8, 12}
	_Number_index_1 = [...]uint8{0, 3, 6, 11, 15, 19, 22, 27, 32, 36, 39, 45, 51, 59}
)

func (i Number) String() string {
	switch {
	case -1 <= i && i <= 0:
		i -= -1
		return _Number_name_0[_Number_index_0[i]:_Number_index_0[i+1]]
	case 3 <= i && i <= 15:
		i -= 3
		return _Number_name_1[_Number_index_1[i]:_Number_index_1[i+1]]
	default:
		return fmt.Sprintf("Number(%d)", i)
	}
}
