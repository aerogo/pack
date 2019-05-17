package pack

import "github.com/OneOfOne/xxhash"

// HashString hashes a string.
func HashString(data string) uint64 {
	h := xxhash.NewS64(0)
	_, _ = h.WriteString(data)
	return h.Sum64()
}
