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

// Graph represents the Directed Acyclic Word Graph.
type Graph struct {
	units           []unit
	labels          []byte
	isIntersections bitVector
}

// Root returns the root ID.
func (g Graph) Root() uint32 {
	return 0
}

// Child returns the child unit ID of a unit.
func (g Graph) Child(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.units[id].child(), nil
}

// Sibling returns the sibling unit ID of a unit.
func (g Graph) Sibling(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	if g.units[id].hasSibling() {
		return id + 1, nil
	}
	return 0, nil
}

// Value returns the value of a unit.
func (g Graph) Value(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.units[id].value(), nil
}

// IsLeaf is true if unit has no sibling.
func (g Graph) IsLeaf(id uint32) (bool, error) {
	l, err := g.Label(id)
	return l == 0, err
}

// Label returns the character byte of a unit.
func (g Graph) Label(id uint32) (byte, error) {
	if int(id) >= len(g.labels) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.labels[id], nil
}

// IsIntersection is true if a unit is intersection of other units.
func (g Graph) IsIntersection(id uint32) (bool, error) {
	return g.isIntersections.get(id)
}

// IntersectionID returns the intersection unit ID.
func (g Graph) IntersectionID(id uint32) (int, error) {
	r, err := g.isIntersections.rank(id)
	return r - 1, err
}

// NumIntersections returns the number of intersection units.
func (g Graph) NumIntersections() int {
	return g.isIntersections.numOnes
}

// Size returns the graph size.
func (g Graph) Size() int {
	return len(g.units)
}
