// +build mmap

package trie

import (
	"github.com/ikawaha/dartsclone/trie/internal"
)

// MmapedTrie is the TRIE interface.
type MmapedTrie interface {
	Trie
	// Close deletes mapped memory and closes mapped file.
	Close() error
}

// OpenMmaped opens the named file of double array and maps it on the memory.
func OpenMmaped(name string) (MmapedTrie, error) {
	return internal.OpenMmaped(name)
}
