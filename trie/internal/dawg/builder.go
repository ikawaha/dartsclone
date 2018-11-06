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

package dawg

import (
	"fmt"
)

const initialTableSize = 1 << 10

// Builder represents the structure of Directed Acyclic Word Graph builder.
type Builder struct {
	Graph
	nodes      []node
	table      []int
	nodeStack  stack
	recycleBin stack
	numStates  int
}

// NewBuilder returns a DAWG builder.
func NewBuilder() *Builder {
	ret := Builder{}
	ret.init()
	return &ret
}

func (b *Builder) init() {
	b.table = make([]int, initialTableSize)
	b.appendNode()
	b.appendUnit()
	b.numStates = 1
	b.nodes[0].label = 0xFF
	b.nodeStack.push(0)
}

// Finish finishes a building DAWG.
func (b *Builder) Finish() (*Graph, error) {
	if err := b.flush(0); err != nil {
		return nil, fmt.Errorf("builder flush, %v", err)
	}
	b.units[0] = b.nodes[0].unit()
	b.labels[0] = b.nodes[0].label
	b.isIntersections.finish()
	b.nodes = nil
	b.table = nil
	b.nodeStack = nil
	b.recycleBin = nil

	return &b.Graph, nil
}

// Insert inserts the word and the value to a DAWG.
func (b *Builder) Insert(key string, value uint32) error {
	if len(key) == 0 {
		return fmt.Errorf("zero-length key")
	}
	id := 0
	keyPos := 0
	for ; keyPos <= len(key); keyPos++ {
		childID := b.nodes[id].child
		if childID == 0 {
			break
		}
		var keyLabel byte
		if keyPos < len(key) {
			keyLabel = key[keyPos]
		}
		if keyPos < len(key) && keyLabel == 0 {
			return fmt.Errorf("invalid null character")
		}

		unitLabel := b.nodes[childID].label
		if keyLabel < unitLabel {
			return fmt.Errorf("wrong key order")
		}
		if keyLabel > unitLabel {
			b.nodes[childID].hasSibling = true
			if err := b.flush(int(childID)); err != nil {
				return fmt.Errorf("builder flush, %v", err)
			}
			break
		}
		id = int(childID)
	}
	if keyPos > len(key) {
		return nil
	}
	for ; keyPos <= len(key); keyPos++ {
		var keyLabel byte
		if keyPos < len(key) {
			keyLabel = key[keyPos]
		}
		childID, err := b.appendNode()
		if err != nil {
			return fmt.Errorf("append node, %v", err)
		}
		if b.nodes[id].child == 0 {
			b.nodes[childID].isState = true
		}
		b.nodes[childID].sibling = int(b.nodes[id].child)
		b.nodes[childID].label = keyLabel
		b.nodes[id].child = uint32(childID)
		b.nodeStack.push(childID)
		id = childID
	}
	b.nodes[id].setValue(value)
	return nil
}

func (b *Builder) flush(id int) error {
	for {
		nodeID, err := b.nodeStack.top()
		if err != nil {
			return fmt.Errorf("node stack top, %v", err)
		}
		if nodeID == id {
			break
		}
		if err := b.nodeStack.pop(); err != nil {
			return fmt.Errorf("node stack pop, %v", err)
		}
		if b.numStates >= len(b.table)-(len(b.table)>>2) {
			b.expandTable()
		}
		numSiblings := 0
		for n := nodeID; n != 0; n = b.nodes[n].sibling {
			numSiblings++
		}
		matchID, hashID := b.findNode(nodeID)
		if matchID != 0 {
			if err := b.isIntersections.set(matchID, true); err != nil {
				return fmt.Errorf("intersections set, %v", err)
			}
		} else {
			unitID := 0
			for i := 0; i < numSiblings; i++ {
				unitID = b.appendUnit()
			}

			for n := nodeID; n != 0; n = b.nodes[n].sibling {
				b.units[unitID] = b.nodes[n].unit()
				b.labels[unitID] = b.nodes[n].label
				unitID--
			}
			matchID = unitID + 1
			b.table[hashID] = matchID
			b.numStates++
		}
		for n := nodeID; n != 0; {
			next := b.nodes[n].sibling
			b.freeNode(n)
			n = next
		}
		top, err := b.nodeStack.top()
		if err != nil {
			return fmt.Errorf("node stack top, %v", err)
		}
		b.nodes[top].child = uint32(matchID)
	}
	if err := b.nodeStack.pop(); err != nil {
		return fmt.Errorf("node stack pop, %v", err)
	}
	return nil
}

func (b *Builder) expandTable() {
	tableSize := len(b.table) << 1
	b.table = make([]int, tableSize)
	for id := 1; id < len(b.units); id++ {
		if b.labels[id] == 0 || b.units[id].isState() {
			findResult := b.findUnit(id)
			hashID := findResult[1]
			b.table[hashID] = id
		}
	}
}

func (b *Builder) findUnit(id int) [2]int {
	result := [2]int{}
	hashID := b.hashUnit(id) % len(b.table)
	for {
		unitID := b.table[hashID]
		if unitID == 0 {
			break
		}
		hashID = (hashID + 1) % len(b.table)
	}
	result[1] = hashID
	return result
}

func (b *Builder) findNode(nodeID int) (matchID, hashID int) {
	hashID = b.hashNode(nodeID) % len(b.table)
	for {
		unitID := b.table[hashID]
		if unitID == 0 {
			break
		}
		if b.areEqual(nodeID, unitID) {
			return unitID, hashID
		}
		hashID = (hashID + 1) % len(b.table)
	}
	return 0, hashID
}

func (b *Builder) areEqual(nodeID, unitID int) bool {
	for n := b.nodes[nodeID].sibling; n != 0; n = b.nodes[n].sibling {
		if !b.units[unitID].hasSibling() {
			return false
		}
		unitID++
	}
	if b.units[unitID].hasSibling() {
		return false
	}
	for n := nodeID; n != 0; n = b.nodes[n].sibling {
		if (b.nodes[n].unit() != b.units[unitID]) || (b.nodes[n].label != b.labels[unitID]) {
			return false
		}
		unitID--
	}
	return true
}

func (b *Builder) hashUnit(id int) int {
	hashValue := 0
	for ; id != 0; id++ {
		u := b.units[id]
		label := b.labels[id]
		hashValue ^= b.hash((uint32(label) << 24) ^ uint32(u))
		if !b.units[id].hasSibling() {
			break
		}
	}
	return hashValue

}

func (b *Builder) hashNode(id int) int {
	hashValue := 0
	for ; id != 0; id = b.nodes[id].sibling {
		u := b.nodes[id].unit()
		label := b.nodes[id].label
		hashValue ^= b.hash((uint32(label) << 24) ^ uint32(u))
	}
	return hashValue
}

func (b *Builder) appendUnit() int {
	b.isIntersections.append()
	b.units = append(b.units, unit(0))
	b.labels = append(b.labels, 0)
	return b.isIntersections.size - 1
}

func (b *Builder) appendNode() (id int, err error) {
	if len(b.recycleBin) > 0 {
		top, err := b.recycleBin.top()
		if err != nil {
			return -1, fmt.Errorf("recycle bin top, %v", err)
		}
		id = top
		b.nodes[id].reset()
		if err := b.recycleBin.pop(); err != nil {
			return -1, fmt.Errorf("recycle bin pop, %v", err)
		}
		return id, nil
	}
	id = len(b.nodes)
	b.nodes = append(b.nodes, node{})
	return id, nil
}

func (b *Builder) freeNode(id int) {
	b.recycleBin.push(id)
}

func (b Builder) hash(key uint32) int {
	key = ^key + (key << 15)
	key = key ^ (key >> 12)
	key = key + (key << 2)
	key = key ^ (key >> 4)
	key = key * 2057
	key = key ^ (key >> 16)
	return int(key)
}
