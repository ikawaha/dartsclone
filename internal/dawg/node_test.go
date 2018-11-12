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
	"testing"
)

func TestNode_SetValue(t *testing.T) {
	var n node
	for i, expected := range []uint32{0, 1, 2, 3} {
		n.setValue(expected)
		if n.child != expected {
			t.Errorf("expected %v, got %v, (%d)", expected, n.child, i)
		}
	}
}

func TestNode_GetValue(t *testing.T) {
	var n node
	for i, expected := range []uint32{0, 1, 2, 3} {
		n.setValue(expected)
		if got := n.value(); got != expected {
			t.Errorf("expected %v, got %v, (%d)", expected, got, i)
		}
	}
}

func TestNode_Unit(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var n node
		if got, expected := uint32(n.unit()), uint32(0); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("no label", func(t *testing.T) {
		n := node{
			label:      0,
			hasSibling: true,
			isState:    true,
			child:      1,
		}
		if got, expected := uint32(n.unit()), uint32(1<<1+1); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})

	t.Run("has label", func(t *testing.T) {
		type testdata struct {
			n        node
			expected uint32
		}
		for i, v := range []testdata{
			{n: node{label: 1, hasSibling: false, isState: false, child: 1}, expected: 1 << 2},
			{n: node{label: 2, hasSibling: true, isState: false, child: 1}, expected: 1<<2 + 1},
			{n: node{label: 3, hasSibling: false, isState: true, child: 1}, expected: 1<<2 + 2},
			{n: node{label: 4, hasSibling: true, isState: true, child: 1}, expected: 1<<2 + 3},
		} {
			if got, expected := uint32(v.n.unit()), v.expected; got != expected {
				t.Errorf("expected %v, got %v, (%v)", expected, got, i)
			}
		}
	})
}

func TestNode_Reset(t *testing.T) {
	n := node{
		child:      3,
		sibling:    4,
		label:      'a',
		hasSibling: true,
		isState:    true,
	}
	n.reset()
	if got, expected := n.child, uint32(0); got != expected {
		t.Errorf("child: expected %v, got %v", expected, got)
	}
	if got, expected := n.sibling, 0; got != expected {
		t.Errorf("sibling: expected %v, got %v", expected, got)
	}
	if got, expected := n.label, byte(0); got != expected {
		t.Errorf("label: expected %v, got %v", expected, got)
	}
	if got, expected := n.hasSibling, false; got != expected {
		t.Errorf("hasSibling: expected %v, got %v", expected, got)
	}
	if got, expected := n.isState, false; got != expected {
		t.Errorf("isState: expected %v, got %v", expected, got)
	}
}

func TestNode_String(t *testing.T) {
	n := node{
		child:      3,
		sibling:    4,
		label:      'a',
		hasSibling: true,
		isState:    true,
	}
	expected := "child: 3, sibling: 4, label: a, is_state: true, has_sibling: true"
	if got := n.String(); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
