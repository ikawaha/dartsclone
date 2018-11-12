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

package internal

import (
	"fmt"
	"sort"
)

type keySet struct {
	keys   []string
	values []uint32
}

func (s keySet) Len() int           { return len(s.keys) }
func (s keySet) Less(i, j int) bool { return s.keys[i] < s.keys[j] }
func (s keySet) Swap(i, j int) {
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
	if s.hasValues() && len(s.keys) == len(s.values) {
		s.values[i], s.values[j] = s.values[j], s.values[i]
	}
}

func newSortedKeySet(keys []string, values []uint32) (*keySet, error) {
	if len(values) != 0 && len(keys) != len(values) {
		return nil, fmt.Errorf("invalid input, keys=%v, values=%v", len(keys), len(values))
	}
	s := keySet{
		keys:   keys,
		values: values,
	}
	if !sort.StringsAreSorted(keys) {
		sort.Sort(s)
	}
	prev := ""
	for i, v := range keys {
		if i != 0 && prev == v {
			return nil, fmt.Errorf("duplicate key error, %v", v)
		}
		prev = v
	}
	return &s, nil
}

func (s keySet) size() int {
	return len(s.keys)
}

func (s keySet) getKey(id int) (string, error) {
	if id < 0 || id >= len(s.keys) {
		return "", fmt.Errorf("index out of bounds")
	}
	return s.keys[id], nil
}

func (s keySet) getKeyByte(keyID, byteID int) (byte, error) {
	if keyID < 0 || keyID >= len(s.keys) {
		return 0, fmt.Errorf("index out of bounds")
	}
	if byteID < 0 {
		return 0, fmt.Errorf("index out of bounds")
	} else if byteID >= len(s.keys[keyID]) {
		return 0, nil // THIS IS A SPEC!
	}
	return s.keys[keyID][byteID], nil
}

func (s keySet) hasValues() bool {
	return len(s.values) > 0
}

func (s keySet) getValue(id int) (uint32, error) {
	if id < 0 {
		return 0, fmt.Errorf("index out of bounds")
	}
	if !s.hasValues() {
		return uint32(id), nil
	}
	if int(id) >= len(s.values) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return s.values[int(id)], nil
}
