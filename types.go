package agiledoc

import (
	"math/big"
	"strings"
	//bf "github.com/russross/blackfriday"
)

type (
	TokenType uint32
	ValueType uint8
)

//go:generate -command stringer -type ValueType ./
const (
	// BASIC VALUE TYPES
	EMPTY ValueType = 0
	BOOL  ValueType = 1 << iota
	INTEGER
	FLOAT
	BYTE
	BYTESLICE
	STRING
	CHAR
)

//// VALUE INTERFACE
///
type Value interface {
	Type() ValueType
}
type (
	/// VALUE TYPES
	emptyVar  struct{}
	boolVal   bool
	intVal    struct{ *big.Int }
	floatVal  struct{ *big.Float }
	byteVal   byte
	byteSlice []byte
	strVal    string
)

func (emptyVar) Type() ValueType  { return EMPTY }
func (boolVal) Type() ValueType   { return BOOL }
func (intVal) Type() ValueType    { return INTEGER }
func (floatVal) Type() ValueType  { return FLOAT }
func (byteVal) Type() ValueType   { return BYTE }
func (byteSlice) Type() ValueType { return BYTESLICE }
func (strVal) Type() ValueType    { return STRING }

type ( // FUNCTION TYPES (first passed type replaces receiver)
	parseFunc       func(Value, interface{}) Value // parse arbitrary typed value to value type
	convertToFunc   func(Value, ValueType) Value   // convert passed value instance to its receiver type
	convertFromFunc func(Value, Value) Value       // convert passed value instance to its receiver type
)

// generate value from arbitrary value
func Parse(v Value, i interface{}) (r Value) {
	switch i.(type) {
	case bool:
	case int, int8, int16, int32, int64, *big.Int:
	case float32, float64, *big.Float:
	case byte:
	case []byte:
	case string:
	}
	return r
}

// convert to value of a certain type
func ConvertToFunc(v Value, t ValueType) func() Value {

	var fnc func() Value

	// switch on receiver type
	switch v.Type() {
	case BOOL:
	case INTEGER:
	case FLOAT:
	case BYTE:
	case BYTESLICE:
	case STRING:
	}
	return fnc
}

var ConversionFunctions = map[ValueType]map[ValueType]func(Value) Value{
	// return bool
	BOOL: map[ValueType]func(Value) Value{
		INTEGER: func(v Value) Value {
			i := v.(intVal).Int // assert integer
			if v.(intVal).Cmp(i) > 0 {
				return boolVal(true)
			} else {
				return boolVal(false)
			}
		},
		FLOAT: func(v Value) Value {
			i := v.(floatVal).Float // assert float
			if v.(floatVal).Cmp(i) > 0 {
				return boolVal(true)
			} else {
				return boolVal(false)
			}
		},
		BYTE: func(v Value) Value {
			u := uint8(v.(byteVal)) // assert uint8
			if u > 0 {
				return boolVal(true)
			} else {
				return boolVal(false)
			}
		},
		BYTESLICE: func(v Value) Value {
			s := v.(byteSlice) // assert byte slice
			if len(s) > 0 {
				return boolVal(true)
			} else {
				return boolVal(false)
			}
		},
		STRING: func(v Value) Value {
			s := string(v.(strVal))
			if strings.Compare(s, "true") > 0 {
				return boolVal(true)
			} else {
				return boolVal(false)
			}
		},
	},
	// return *big.Int
	INTEGER: map[ValueType]func(Value) Value{
		BOOL: func(v Value) Value {
			if v.(boolVal) {
				return intVal{big.NewInt(1)}
			} else {
				return intVal{big.NewInt(-1)}
			}
		},
		FLOAT: func(v Value) Value {
			f, _ := v.(floatVal).Int64()
			if f > -1 {
				return intVal{big.NewInt(1)}
			} else {
				return intVal{big.NewInt(-1)}
			}
		},
		BYTE: func(v Value) Value {
			if uint8(v.(byteVal)) > 0 {
				return intVal{big.NewInt(1)}
			} else {
				return intVal{big.NewInt(0)}
			}
		},
		BYTESLICE: func(v Value) Value {
			if len(v.(byteSlice)) > -1 {
				return intVal{big.NewInt(1)}
			} else {
				return intVal{big.NewInt(-1)}
			}
		},
		STRING: func(v Value) Value {
			if strings.Compare(string(v.(strVal)), "true") > -1 {
				return intVal{big.NewInt(1)}
			} else {
				return intVal{big.NewInt(-1)}
			}
		},
	},
	// return *big.Float
	FLOAT: map[ValueType]func(Value) Value{
		BOOL: func(v Value) Value {
			if v.(boolVal) {
				return floatVal{big.NewFloat(1)}
			} else {
				return floatVal{big.NewFloat(-1)}
			}
		},
		INTEGER: func(v Value) Value {
			f, _ := v.(floatVal).Int64()
			if f > -1 {
				return floatVal{big.NewFloat(1)}
			} else {
				return floatVal{big.NewFloat(-1)}
			}
		},
		BYTE: func(v Value) Value {
			if uint8(v.(byteVal)) > 0 {
				return floatVal{big.NewFloat(1)}
			} else {
				return floatVal{big.NewFloat(0)}
			}
		},
		BYTESLICE: func(v Value) Value {
			if len(v.(byteSlice)) > -1 {
				return floatVal{big.NewFloat(1)}
			} else {
				return floatVal{big.NewFloat(-1)}
			}
		},
		STRING: func(v Value) Value {
			if strings.Compare(string(v.(strVal)), "true") > -1 {
				return floatVal{big.NewFloat(1)}
			} else {
				return floatVal{big.NewFloat(-1)}
			}
		},
	},
	// return byte
	BYTE: map[ValueType]func(Value) Value{
		BOOL: func(v Value) Value {
			var r byteVal = 0
			return r
		},
		INTEGER: func(v Value) Value {
			var r byteVal = 0
			return r
		},
		FLOAT: func(v Value) Value {
			var r byteVal = 0
			return r
		},
		BYTESLICE: func(v Value) Value {
			var r byteVal = 0
			return r
		},
		STRING: func(v Value) Value {
			var r byteVal = 0
			return r
		},
	},
	// return byte slice
	BYTESLICE: map[ValueType]func(Value) Value{
		BOOL:    func(v Value) (r Value) { return r },
		INTEGER: func(v Value) (r Value) { return r },
		FLOAT:   func(v Value) (r Value) { return r },
		BYTE:    func(v Value) (r Value) { return r },
		STRING:  func(v Value) (r Value) { return r },
	},
	// return string
	STRING: map[ValueType]func(Value) Value{
		BOOL:      func(v Value) (r Value) { return r },
		INTEGER:   func(v Value) (r Value) { return r },
		FLOAT:     func(v Value) (r Value) { return r },
		BYTE:      func(v Value) (r Value) { return r },
		BYTESLICE: func(v Value) (r Value) { return r },
	},
}
