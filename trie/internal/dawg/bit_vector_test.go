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
	"reflect"
	"testing"
)

func TestBitVector_Append(t *testing.T) {
	var v bitVector
	for i := 1; i <= unitSize; i++ {
		v.append()
		if expected, got := 1, len(v.units); expected != got {
			t.Errorf("size of units: expected %v, got %v", expected, got)
		}
		if expected, got := i, v.size; expected != got {
			t.Errorf("size: expected %v, got %v", expected, got)
		}
	}
	v.append()
	if expected, got := 2, len(v.units); expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
	if expected, got := unitSize+1, v.size; expected != got {
		t.Errorf("size: expected %v, got %v", expected, got)
	}
}

func TestBitVector_Set(t *testing.T) {
	var v bitVector

	if err := v.set(1, true); err == nil {
		t.Errorf("expected index out of bounds error")
	}

	v.units = []uint32{0, 0}
	if err := v.set(0, true); err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if expected, got := 1, v.units[0]; expected != int(got) {
		t.Errorf("[0...01]: expected %v, got %v", expected, got)
	}
	if err := v.set(1, true); err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if expected, got := 3, v.units[0]; expected != int(got) {
		t.Errorf("[0...11]: expected %v, got %v", expected, got)
	}
	if err := v.set(0, false); err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if expected, got := 2, v.units[0]; expected != int(got) {
		t.Errorf("[0...10]: expected %v, got %v", expected, got)
	}
	if err := v.set(1, false); err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if expected, got := 0, v.units[0]; expected != int(got) {
		t.Errorf("[0...00]: expected %v, got %v", expected, got)
	}

	if err := v.set(unitSize, true); err != nil {
		t.Errorf("unexpected error, %v", err)
	}
	if expected, got := 1, v.units[1]; expected != int(got) {
		t.Errorf("[...][0...01]: expected %v, got %v", expected, got)
	}
}

func TestBitVector_Empty(t *testing.T) {
	var v bitVector
	if !v.empty() {
		t.Errorf("expected empty")
	}
	v.append()
	if v.empty() {
		t.Errorf("expected not empty")
	}
}

func TestBitVector_popCount(t *testing.T) {
	if expected, got := 1, popCount(2); expected != got {
		t.Errorf("popCount(2): expected %v, got %v", expected, got)
	}
	if expected, got := 2, popCount(3); expected != got {
		t.Errorf("popCount(3): expected %v, got %v", expected, got)
	}
	if expected, got := 8, popCount(0xFF); expected != got {
		t.Errorf("popCount(0xFF): expected %v, got %v", expected, got)
	}
	if expected, got := 16, popCount(0xFFFF); expected != got {
		t.Errorf("popCount(0xFFFF): expected %v, got %v", expected, got)
	}
	if expected, got := 24, popCount(0xFFFFFF); expected != got {
		t.Errorf("popCount(0xFFFFFF): expected %v, got %v", expected, got)
	}
	if expected, got := 32, popCount(0xFFFFFFFF); expected != got {
		t.Errorf("popCount(0xFFFFFFFF): expected %v, got %v", expected, got)
	}
	if expected, got := 16, popCount(0x55555555); expected != got {
		t.Errorf("popCount(0x55555555): expected %v, got %v", expected, got)
	}
}

func TestBitVector_Get(t *testing.T) {
	var v bitVector
	v.units = []uint32{0, 0}

	if _, err := v.get(unitSize * 2); err == nil {
		t.Error("expected index out of bounds error")
	}

	for i := range []uint32{1, 3, 5, 7, 9, unitSize, unitSize + 1, unitSize + 2} {
		if v, err := v.get(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if v {
			t.Errorf("expected false, but true")
		}
	}

	for i := range []int{1, 3, 5, 7, 9, unitSize, unitSize + 1, unitSize + 2} {
		if err := v.set(i, true); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if v, err := v.get(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if !v {
			t.Errorf("expected true, but false")
		}
	}
}

func TestBitVector_Build(t *testing.T) {
	var v bitVector
	for i := 0; i <= unitSize*2; i++ {
		v.append()
		if err := v.set(i, true); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}
	v.finish()
	if expected := 3; expected != len(v.units) {
		t.Errorf("expected %v, got %v", expected, len(v.units))
	}
	if expected := []int{0, 32, 64}; !reflect.DeepEqual([]int{0, 32, 64}, v.ranks) {
		t.Errorf("expected %v, got %v", expected, v.ranks)
	}
	if expected := 65; expected != v.numOnes {
		t.Errorf("expected %v, got %v", expected, v.numOnes)
	}
}

func TestBitVector_Rank(t *testing.T) {
	var v bitVector
	for i := 0; i <= unitSize*2; i++ {
		v.append()
		if err := v.set(i, true); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}
	v.finish()

	if _, err := v.rank(unitSize * 3); err == nil {
		t.Error("expected index out of bounds error")
	}
	for i := 0; i <= unitSize*2; i++ {
		if v, err := v.rank(uint32(i)); err != nil {
			t.Errorf("unexpected error, %v", v)
		} else if i+1 != v {
			t.Errorf("expected %v, got %v", i+1, v)
		}
	}
}
