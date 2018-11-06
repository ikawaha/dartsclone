// Copyright 2018 ikawaha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
