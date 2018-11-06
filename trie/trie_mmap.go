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

// OpenMmaped opens the named file of the double array and maps it on the memory.
func OpenMmaped(name string) (MmapedTrie, error) {
	return internal.OpenMmaped(name)
}
