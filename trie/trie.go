package trie

import (
	"github.com/ikawaha/dartsclone/trie/internal"
)

// Trie is the TRIE interface.
type Trie interface {
	// ExactMatchSearch searches TRIE by a given keyword and returns the id and it's length if found.
	ExactMatchSearch(key string) (id, size int, err error)
	// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the ids and it's lengths if found.
	CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error)
	// CommonPrefixSearchCallback finds keywords sharing common prefix in an input and callback with id and it's length.
	CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error
}

// Open opens the named file of the double array.
func Open(name string) (Trie, error) {
	return internal.Open(name)
}
