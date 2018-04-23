package rpc

import (
	"math/rand"
)

func selectRoute(chances []int64) int {
	chancesLength := len(chances)
	if chancesLength == 1 {
		return 0
	}
	var sumOfChances int64
	for _, chance := range chances {
		sumOfChances = sumOfChances + chance
	}

	randomFloat := rand.Float64() * float64(sumOfChances)

	for float64(sumOfChances) > randomFloat {
		sumOfChances = sumOfChances - chances[chancesLength -1]
		chancesLength = chancesLength - 1
	}

	return chancesLength - 1
}