package trie

import (
	"io"

	"github.com/ikawaha/dartsclone/trie/internal"
)

// BuildDoubleArray returns a double array trie for keys and values.
func BuildDoubleArray(keys []string, values []uint32, progress ProgressFunction) (Trie, error) {
	return internal.BuildDoubleArray(keys, values, progress)
}

// DoubleArrayBuilder represents builder of the double array.
type DoubleArrayBuilder struct {
	*internal.DoubleArrayBuilder
}

// NewDoubleArrayBuilder creates a builder of the double array.
func NewDoubleArrayBuilder(progress ProgressFunction) *DoubleArrayBuilder {
	return &DoubleArrayBuilder{
		DoubleArrayBuilder: internal.NewDoubleArrayBuilder(progress),
	}
}

// WriteTo write to the serialize data of the double array.
func (b DoubleArrayBuilder) WriteTo(w io.Writer) (int64, error) {
	return b.DoubleArrayBuilder.WriteTo(w)
}
