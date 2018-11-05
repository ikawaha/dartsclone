package dawg

import (
	"testing"
)

func TestGraph_Root(t *testing.T) {
	g := Graph{}
	if got, expected := g.Root(), uint32(0); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestGraph_Child(t *testing.T) {
	g := Graph{
		units: []unit{
			node{label: 1, hasSibling: false, isState: false, child: 1}.unit(),
			node{label: 2, hasSibling: true, isState: false, child: 2}.unit(),
			node{label: 3, hasSibling: false, isState: true, child: 3}.unit(),
			node{label: 4, hasSibling: true, isState: true, child: 4}.unit(),
		},
	}
	if _, err := g.Child(uint32(len(g.units))); err == nil {
		t.Errorf("expected index out of bounds error")
	}
	for i := range g.units {
		if got, err := g.Child(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := uint32(i) + 1; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestGraph_Sibling(t *testing.T) {
	t.Run("has sibling", func(t *testing.T) {
		g := Graph{
			units: []unit{
				node{label: 1, hasSibling: true, isState: false, child: 1}.unit(),
				node{label: 2, hasSibling: true, isState: false, child: 2}.unit(),
				node{label: 3, hasSibling: true, isState: true, child: 3}.unit(),
				node{label: 4, hasSibling: true, isState: true, child: 4}.unit(),
			},
		}
		for i := range g.units {
			if got, err := g.Sibling(uint32(i)); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if expected := uint32(i) + 1; got != expected {
				t.Errorf("expected %v, got %v, (%v)", expected, got, i)
			}
		}
	})
	t.Run("does not have sibling", func(t *testing.T) {
		g := Graph{
			units: []unit{
				node{label: 1, hasSibling: false, isState: false, child: 1}.unit(),
				node{label: 2, hasSibling: false, isState: false, child: 2}.unit(),
				node{label: 3, hasSibling: false, isState: true, child: 3}.unit(),
				node{label: 4, hasSibling: false, isState: true, child: 4}.unit(),
			},
		}
		if _, err := g.Sibling(uint32(len(g.units))); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		for i := range g.units {
			if got, err := g.Sibling(uint32(i)); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if expected := uint32(0); got != expected {
				t.Errorf("expected %v, got %v, (%v)", expected, got, i)
			}
		}
	})
}

func TestGraph_Value(t *testing.T) {
	g := Graph{
		units: []unit{
			0,
			1 << 1,
			2 << 1,
			3 << 1,
		},
	}
	if _, err := g.Value(uint32(len(g.units))); err == nil {
		t.Errorf("expected index out of bounds error")
	}
	for i := range g.units {
		if got, err := g.Value(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := uint32(i); got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestGraph_Label(t *testing.T) {
	g := Graph{
		labels: []byte{0, 2, 4, 6, 8},
	}
	for i := range g.labels {
		if got, err := g.Label(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := byte(i * 2); got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestGraph_IsLeaf(t *testing.T) {
	g := Graph{
		labels: []byte{0, 2, 0, 6, 0},
	}
	if _, err := g.IsLeaf(uint32(len(g.labels))); err == nil {
		t.Errorf("expected index out of bounds error")
	}
	for i := range g.labels {
		if got, err := g.IsLeaf(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := i%2 == 0; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestGraph_Size(t *testing.T) {
	g := Graph{
		units: []unit{0, 1, 2, 3, 4, 5, 6},
	}
	if got, expected := g.Size(), len(g.units); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestGraph_IntersectionID(t *testing.T) {
	var v bitVector
	for i := 0; i <= unitSize*2; i++ {
		v.append()
		if err := v.set(i, true); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}
	v.finish()

	g := Graph{isIntersections: v}
	if got, err := g.IntersectionID(uint32(60)); err != nil {
		t.Errorf("unexpected error, %v", err)
	} else if expected := 60; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestGraph_IsIntersection(t *testing.T) {
	var v bitVector
	for i := 0; i <= unitSize*2; i++ {
		v.append()
		if err := v.set(i, i%2 == 0); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}
	v.finish()

	g := Graph{isIntersections: v}

	t.Run("true if even", func(t *testing.T) {
		if got, err := g.IsIntersection(uint32(60)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := true; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("false if oddn", func(t *testing.T) {
		if got, err := g.IsIntersection(uint32(59)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := false; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}

func TestGraph_NumIntersections(t *testing.T) {
	var v bitVector
	for i := 0; i <= unitSize*2; i++ {
		v.append()
		if err := v.set(i, true); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}
	v.finish()

	g := Graph{isIntersections: v}

	if got, expected := g.NumIntersections(), unitSize*2+1; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
