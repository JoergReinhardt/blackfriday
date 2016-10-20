package agiledoc

import (
	"testing"
)

const (
	TRUE  bool = true
	FALSE bool = false
)

func TestBool(t *testing.T) {
	v := Value(TRUE)
	(*t).Log(v)

}
