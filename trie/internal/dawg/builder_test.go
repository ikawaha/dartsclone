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

func TestNewDAWGBuilder(t *testing.T) {
	b := NewBuilder()
	if got, expected := b.Root(), uint32(0); got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
	if got, expected := b.Size(), 1; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestDAWGBuilder_Init(t *testing.T) {
	b := Builder{}
	b.init()
	if len(b.nodes) != 1 {
		t.Errorf("nodes size: expected 1, got %v", len(b.nodes))
	} else if b.nodes[0].label != 0xFF {
		t.Errorf("nodes[0] label: expected 0xFF, got 0x%X", b.nodes[0].label)
	}
	if len(b.units) != 1 {
		t.Errorf("units size: expected 1, got %v", len(b.nodes))
	}
	if b.numStates != 1 {
		t.Errorf("num states: expected 1, got %v", b.numStates)
	}
	if len(b.nodeStack) != 1 {
		t.Errorf("node stack size: expected 1, got %v", len(b.nodeStack))
	}
}

func TestDAWGBuilder_AppendNode(t *testing.T) {
	b := NewBuilder() // build & initialize

	// first time
	id, err := b.appendNode()
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if id != 1 {
		t.Errorf("expected id=0, got %v", id)
	}
	if got, expected := len(b.nodes), 2; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}

	// second
	id, err = b.appendNode()
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if id != 2 {
		t.Errorf("expected id=0, got %v", id)
	}
	if got, expected := len(b.nodes), 3; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}

	// recycle
	b.nodes[id].label = 'a'
	b.freeNode(id)

	// reuse
	id, err = b.appendNode()
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if id != 2 {
		t.Errorf("expected id=0, got %v", id)
	}
	if b.nodes[id].label != byte(0) {
		t.Errorf("node crear error, %+v", b.nodes[id])
	}
	if len(b.recycleBin) != 0 {
		t.Errorf("node recycle error, %+v", b.recycleBin)
	}
}

func TestDAWGBuilder_Insert(t *testing.T) {
	t.Run("zero-length key", func(t *testing.T) {
		b := NewBuilder()
		b.init()
		if err := b.Insert("", uint32(0)); err == nil {
			t.Error("expected zero-length key error")
		}
	})

	t.Run("wrong key order", func(t *testing.T) {
		b := NewBuilder()
		b.init()
		if err := b.Insert("world", uint32(0)); err != nil {
			t.Errorf("insert error, %v", err)
		}
		if err := b.Insert("hello", uint32(1)); err == nil {
			t.Error("expected wrong key order error")
		}
	})

	t.Run("ok", func(t *testing.T) {
		keys := []string{
			"hello",
			"world",
		}
		b := NewBuilder()
		b.init()
		for i, v := range keys {
			if err := b.Insert(v, uint32(i)); err != nil {
				t.Errorf("unexpected insert error, %v", err)
			}
		}
	})
}

func TestDAWGBuilder_Hash(t *testing.T) {
	b := NewBuilder()
	for i, v := range []uint32{0, 1, 2, 1 << 10, 1<<21 - 1, 1 << 21, 1<<21 + 1, 1<<29 - 1, 1 << 31, 1<<32 - 1} {
		x := b.hash(v)
		y := b.hash(v)
		if x != y {
			t.Errorf("expected same hash value, input %v, %v<>%v (%v)", v, x, y, i)
		}
	}
}

func TestDAWGBuilder_Finish(t *testing.T) {
	b := NewBuilder()
	b.init()
	if err := b.Insert("a", uint32(0)); err != nil {
		t.Errorf("unexpected insert error, %v", err)
	}
	g, err := b.Finish()
	if err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if got, expected := g.Size(), 4; got != expected {
		t.Errorf("graph size: expected %v, got %v", expected, got)
	}
}

func TestBuilder_ExpandTable(t *testing.T) {
	t.Run("empty units", func(t *testing.T) {
		b := NewBuilder()
		b.expandTable()
		if got, expected := len(b.table), initialTableSize<<1; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
		for i, v := range b.table {
			if v != 0 {
				t.Errorf("unexpected value !=0, %v (%v)", v, i)
			}
		}
	})
	t.Run("has units", func(t *testing.T) {
		b := NewBuilder()
		value0UnitSize := 10
		for i := 0; i < value0UnitSize; i++ {
			b.units = append(b.units, unit(2)) // isState=true
			b.labels = append(b.labels, 0)
		}
		b.expandTable()
		if got, expected := len(b.table), initialTableSize<<1; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
		ii := []int{}
		for i, v := range b.table {
			if v != 0 {
				ii = append(ii, i)
			}
		}
		if got, expected := len(ii), value0UnitSize; got != expected {
			t.Errorf("expected %v, got %v, indexes %+v", expected, got, ii)
		}
	})

}
