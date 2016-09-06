// Big	 | Rat	| Lst,Mao,Set|
// ──────|──────|────────────|
// Simple| Tuple| Collection |
//	 └──────────┬────────┘
//		  Complex
package agiledoc

import (
	"math/big"
)

func ComplexToSimple(b Evaluable) Val {
	// define a new closure
	return Val(func() *big.Int {
		// allocate a new big Int reference
		return new(big.Int).SetBytes(
			// converts string representation of a complex value to
			// a byte slice and stores it in a big Int
			[]byte(b.String()))
	})
}
func SimpleToTuple(b Val) Rat {
	// leaves value uninterpreted but dereferences it
	var val = *b()
	// define a new closure
	return Rat(func() *big.Rat {
		// allocate a new big Rat reference
		return new(big.Rat).SetFrac(
			// since no key is known yet, set numerable to one and
			// let denominator reference value
			new(big.Int).SetInt64(1),
			&val,
		)
	})
}
func RankedValueToTuple(i int, v Val) Rat {
	var idx = *big.NewInt(int64(i))
	var val = *v()
	return Rat(func() *big.Rat {
		// allocate a new big Rat reference
		return new(big.Rat).SetFrac(
			// since no key is known yet, set numerable to one and
			// let denominator reference value
			&idx,
			&val,
		)
	})
}
func ValuePairToTuple(k Val, v Val) Rat {
	var key = *k()
	var val = *v()
	return Rat(func() *big.Rat {
		// allocate a new big Rat reference
		return new(big.Rat).SetFrac(
			// since no key is known yet, set numerable to one and
			// let denominator reference value
			&key,
			&val,
		)
	})
}
