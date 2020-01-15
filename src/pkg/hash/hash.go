package hash

// Wrapper Functions for Hashing.
// Wraps FNV Hash.

import (
	"hash/fnv"
)

// Sum computes the FNV hash of data.
func Sum(data []byte) uint64 {
	fnv := fnv.New64()
	fnv.Write(data)

	return fnv.Sum64()
}
