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

package dartsclone

import (
	"io"

	"github.com/ikawaha/dartsclone/internal"
)

// BuildTRIE returns a dartsclone TRIE for keys and values.
func BuildTRIE(keys []string, values []uint32, progress ProgressFunction) (Trie, error) {
	return internal.BuildDoubleArray(keys, values, progress)
}

// Builder represents builder of the dartsclone TRIE.
type Builder struct {
	*internal.DoubleArrayBuilder
}

// NewBuilder creates a builder of the dartsclone TRIE.
func NewBuilder(progress ProgressFunction) *Builder {
	return &Builder{
		DoubleArrayBuilder: internal.NewDoubleArrayBuilder(progress),
	}
}

// WriteTo write to the serialize data of the dartsclone TRIE.
func (b Builder) WriteTo(w io.Writer) (int64, error) {
	return b.DoubleArrayBuilder.WriteTo(w)
}
