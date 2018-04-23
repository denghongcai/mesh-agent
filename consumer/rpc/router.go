package rpc

import (
	"math/rand"
)

func selectRoute(chances []int64) int {
	chancesLength := len(chances)
	if chancesLength == 1 {
		return 0
	}
	var sumOfChances float64
	for _, chance := range chances {
		sumOfChances = sumOfChances + (1 / float64(chance))
	}

	randomFloat := rand.Float64() * sumOfChances

	for i := 0; i < chancesLength; i ++ {
		randomFloat = randomFloat - ( 1 / float64(chances[i]))
		if randomFloat < 0 {
			return i
		}
	}

	return chancesLength - 1
}