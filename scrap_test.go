package agiledoc

import (
	"math/big"
	"testing"
)

func TestRatQuo(t *testing.T) {
	n, d, v := new(big.Rat), new(big.Rat), new(big.Rat)
	n.SetFrac64(13, 17)
	d.SetFrac64(11, 3)
	v.Quo(n, d)

	t.Log("values: ", v, n, d, "numerator: ", v.Num(), "denominator: ", v.Denom())

	f := new(big.Rat)
	f.SetFrac64(2, 4)
	t.Log(f)
}
