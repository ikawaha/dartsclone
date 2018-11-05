package dawg

import (
	"testing"
)

func TestUnit_Child(t *testing.T) {
	type testdata struct {
		unit     unit
		expected uint32
	}
	for i, v := range []testdata{
		{unit: node{label: 1, hasSibling: false, isState: false, child: 1}.unit(), expected: 1},
		{unit: node{label: 2, hasSibling: true, isState: false, child: 2}.unit(), expected: 2},
		{unit: node{label: 3, hasSibling: false, isState: true, child: 3}.unit(), expected: 3},
		{unit: node{label: 4, hasSibling: true, isState: true, child: 4}.unit(), expected: 4},
	} {
		if got, expected := v.unit.child(), v.expected; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestUnit_HasSibling(t *testing.T) {
	type testdata struct {
		unit     unit
		expected bool
	}
	for i, v := range []testdata{
		{unit: node{label: 1, hasSibling: false, isState: false, child: 1}.unit(), expected: false},
		{unit: node{label: 2, hasSibling: true, isState: false, child: 2}.unit(), expected: true},
		{unit: node{label: 3, hasSibling: false, isState: true, child: 3}.unit(), expected: false},
		{unit: node{label: 4, hasSibling: true, isState: true, child: 4}.unit(), expected: true},
	} {
		if got, expected := v.unit.hasSibling(), v.expected; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestUnit_IsState(t *testing.T) {
	type testdata struct {
		unit     unit
		expected bool
	}
	for i, v := range []testdata{
		{unit: node{label: 1, hasSibling: false, isState: false, child: 1}.unit(), expected: false},
		{unit: node{label: 2, hasSibling: true, isState: false, child: 2}.unit(), expected: false},
		{unit: node{label: 3, hasSibling: false, isState: true, child: 3}.unit(), expected: true},
		{unit: node{label: 4, hasSibling: true, isState: true, child: 4}.unit(), expected: true},
	} {
		if got, expected := v.unit.isState(), v.expected; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}

func TestUnit_Value(t *testing.T) {
	for i, expected := range []uint32{0, 1, 3, 5, 7, 9} {
		var n node
		n.setValue(expected)
		u := n.unit()
		if got, expected := u.value(), expected; got != expected {
			t.Errorf("expected %v, got %v, (%v)", expected, got, i)
		}
	}
}
