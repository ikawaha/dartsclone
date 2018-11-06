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

type seekerAt interface {
	at(uint32) (unit, error)
}

func exactMatchSearch(a seekerAt, key string) (id, size int, err error) {
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

func commonPrefixSearch(a seekerAt, key string, offset int) (ids, sizes []int, err error) {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return ids, sizes, err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return ids, sizes, err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return ids, sizes, err
			}
			ids = append(ids, int(u.value()))
			sizes = append(sizes, i+1)
		}
	}
	return ids, sizes, nil
}

func commonPrefixSearchCallback(a seekerAt, key string, offset int, callback func(id, size int)) error {
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
