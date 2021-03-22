package helpers

import (
	"math/rand"
)

func rangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

type NumRange struct {
	low int
	hi  int
}

func GenerateVerCode() int {
	var code int
	num := []NumRange{
		{1000, 10000},
	}

	for _, test := range num {
		code = rangeIn(test.low, test.hi)
		// fmt.Printf("Num %d <= %d <= %d\n", test.low, rangeIn(test.low, test.hi), test.hi)
	}

	return code
}
