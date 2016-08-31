package agiledoc

import (
	"strconv"
	"strings"
)

// PARSE STRING OR BYTE SLICE REPRESENTATIONS OF NUMBERS TO NUMERALS
func parseInt(v Value) (Value, error) {
	i, err := strconv.Atoi(v.String())
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}
func parseUint(v Value) (Value, error) {
	i, err := strconv.ParseUint(v.String(), 2, 8)
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}
func parseFloat(v Value) (Value, error) {
	f, err := strconv.ParseFloat(v.String(), 10)
	if err != nil {
		return nil, err
	}
	return NativeToValue(f), nil
}
func parseBool(v Value) (Value, error) {
	i, err := strconv.ParseBool(v.String())
	if err != nil {
		return nil, err
	}
	return NativeToValue(i), nil
}

type Language int

const (
	ENGLISH Language = 0
	DEUTSCH Language = 1 + iota
)

type Literal interface {
	Text(l Language) string
}
type Numeral interface {
	Literal // Text(l Language) string
	Number() int
}

//go:generate -command stringer -type Number
type Number int

const (
	zero, null Number = 0, 0
	one, eins  Number = 1 + iota, 1 + iota
	two, zwei
	three, drei
	four, vier
	five, fünf
	six, sechs
	seven, sieben
	eight, acht
	nine, neun
	ten, zehn
	eleven, elf
	twelve, zwölf
	thirteen, dreizehn

	teen, zig           Number = 10, teen
	twen, zwan          Number = 20, twen
	hundred, hundert    Number = 100, hundred
	thousend, tausend   Number = 1000, thousend
	million, millionen  Number = 1000000, million
	billion, milliarde  Number = 1000000000, billion
	trillion, billiarde Number = 1000000000000, trillion
)

func parseLiteral(v Value) Value {
	var r Value
	for i := 0; i < 13; i++ {
		n := Number(i)
		if strings.Compare( // compare decapitalized value…
			strings.ToLower(
				v.String(),
			), // with string function of number equal to current
			//loop index (allready decapitalized)
			n.String(),
		) == 0 { // if identical, current index is the parsed number
			r = NativeToValue(i)
		}
	}
	return r
}
