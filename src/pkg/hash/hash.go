package hash

// Wrapper Functions for Hashing.
// Wraps SHA1 Hash.

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
)

// Sum computes the SHA1 hash of data and extracts a 64 bit integer out of it.
func Sum(data []byte) uint64 {
	var id uint64

	hasher := sha1.New()
	hasher.Write(data)

	data = hasher.Sum(nil)
	binary.Read(bytes.NewBuffer(data[:8]), binary.LittleEndian, &id)

	return id
}
