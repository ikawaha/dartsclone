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
	"github.com/ikawaha/dartsclone/internal"
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
