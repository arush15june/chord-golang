package util

import (
	"math/rand"
)

// IsBetweenID checks if the cmp lies between low and high (exclusive of low, inclusive of high).
func IsBetweenID(cmp uint64, low uint64, high uint64) bool {
	if cmp > low && cmp <= high || (cmp == high) {
		return true
	}
	return false
}

// GetRandomBetween returns a random integer value between low and high.
func GetRandomBetween(low int, high int) int {
	return rand.Intn(high-low) + low
}
