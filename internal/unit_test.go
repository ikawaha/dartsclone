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
	"testing"
)

func TestUnit_SetOffset_GetOffset(t *testing.T) {
	t.Run("offset less than 1<<21", func(t *testing.T) {
		for i := uint32(0); i < 21; i++ {
			var u unit
			for _, x := range []int{-1, 0, +1} {
				expected := uint32(int(uint32(1)<<i) + x)
				if err := u.setOffset(expected); err != nil {
					t.Errorf("unexpected error, %v", err)
				}
				if got := u.offset(); got != expected {
					t.Errorf("expected %v, got %v (%v)", expected, got, i)
				}
			}
		}
	})
	t.Run("offset bigger or equal than 1<<21", func(t *testing.T) {
		for i := uint32(22); i < 29; i++ {
			var u unit
			for _, x := range []int{-1, 0, +1} {
				offset := uint32(int(uint32(1)<<i) + x)
				expected := (offset / 256) * 256 // granularity becomes rough. see http://d.hatena.ne.jp/s-yata/20100301/1267788256
				if err := u.setOffset(expected); err != nil {
					t.Errorf("unexpected error, %v", err)
				}
				if got := u.offset(); got != expected {
					t.Errorf("expected %v, got %v (%v)", expected, got, i)
				}
			}
		}
	})
	t.Run("too large offset", func(t *testing.T) {
		var u unit
		if err := u.setOffset(maxOffset); err == nil {
			t.Errorf("expected too large offset error")
		}
		if err := u.setOffset(maxOffset - 1); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	})
}

func TestUnit_SetLabel(t *testing.T) {
	var u unit
	for _, expected := range []byte{'a', 'b', 'c', '\a', '\n', 0, 0xFF} {
		u.setLabel(expected)
		if got := byte(uint32(u) & 0xFF); got != expected {
			t.Errorf("expected %c, got %c", expected, got)
		}
	}
}

func TestUnit_SetValue(t *testing.T) {
	var u unit
	for _, expected := range []uint32{0, 1, 2, 3, 1 << 10, 1 << 21, 1<<31 - 1} {
		u.setValue(expected)
		if got := u.value(); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}

func TestUnit_SetHasLeaf(t *testing.T) {
	var u unit
	for _, expected := range []bool{true, false, true, false, false, true} {
		u.setHasLeaf(expected)
		if got := u.hasLeaf(); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	}
}
