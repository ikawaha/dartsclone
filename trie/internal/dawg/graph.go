package dawg

import (
	"fmt"
)

type Graph struct {
	units           []unit
	labels          []byte
	isIntersections bitVector
}

func (g Graph) Root() uint32 {
	return 0
}

func (g Graph) Child(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.units[id].child(), nil
}

func (g Graph) Sibling(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	if g.units[id].hasSibling() {
		return id + 1, nil
	}
	return 0, nil
}

func (g Graph) Value(id uint32) (uint32, error) {
	if int(id) >= len(g.units) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.units[id].value(), nil
}

func (g Graph) IsLeaf(id uint32) (bool, error) {
	l, err := g.Label(id)
	return l == 0, err
}

func (g Graph) Label(id uint32) (byte, error) {
	if int(id) >= len(g.labels) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return g.labels[id], nil
}

func (g Graph) IsIntersection(id uint32) (bool, error) {
	return g.isIntersections.get(id)
}

func (g Graph) IntersectionID(id uint32) (int, error) {
	r, err := g.isIntersections.rank(id)
	return r - 1, err
}

func (g Graph) NumIntersections() int {
	return g.isIntersections.numOnes
}

func (g Graph) Size() int {
	return len(g.units)
}
