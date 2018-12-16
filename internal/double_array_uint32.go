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
	"encoding/binary"
	"fmt"
	"os"
)

// DoubleArrayUint32 represents the TRIE data structure.
type DoubleArrayUint32 struct {
	array []uint32
}

// Open opens the named file of the double array.
func Open(name string) (*DoubleArrayUint32, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size != int64(int(size)) {
		return nil, fmt.Errorf("too large file")
	}
	var ret DoubleArrayUint32
	ret.array = make([]uint32, 0, size/4)
	for i := int64(0); i < size; i += 4 {
		var u uint32
		if err := binary.Read(f, binary.LittleEndian, &u); err != nil {
			return nil, fmt.Errorf("broken array, %v", err)
		}
		ret.array = append(ret.array, u)
	}
	return &ret, nil
}

func (a DoubleArrayUint32) at(i uint32) (unit, error) {
	if int(i) >= len(a.array) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return unit(a.array[i]), nil
}

// ExactMatchSearch searches TRIE by a given keyword and returns the id and it's length if found.
func (a DoubleArrayUint32) ExactMatchSearch(key string) (id, size int, err error) {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return -1, -1, err
	}
	for i := 0; i < len(key); i++ {
		nodePos ^= unit.offset() ^ uint32(key[i])
		unit, err = a.at(nodePos)
		if err != nil {
			return -1, -1, err
		}
		if unit.label() != key[i] {
			return -1, 0, nil
		}
	}
	if !unit.hasLeaf() {
		return -1, 0, nil
	}
	unit, err = a.at(nodePos ^ unit.offset())
	if err != nil {
		return -1, -1, err
	}
	return int(unit.value()), len(key), nil
}

// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the array of pairs (id and it's length) if found.
func (a DoubleArrayUint32) CommonPrefixSearch(key string, offset int) ([][2]int, error) {
	var ret [][2]int
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return ret, err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return ret, err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return ret, err
			}
			ret = append(ret, [2]int{int(u.value()), i + 1})
		}
	}
	return ret, nil
}

// CommonPrefixSearchCallback finds keywords sharing common prefix in an input and callback with id and it's length.
func (a DoubleArrayUint32) CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return err
			}
			callback(int(u.value()), i+1)
		}
	}
	return nil
}
