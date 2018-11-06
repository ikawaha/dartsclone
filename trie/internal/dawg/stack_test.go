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

func TestStack_Top(t *testing.T) {
	t.Run("empty stack", func(t *testing.T) {
		var s stack
		if _, err := s.top(); err == nil {
			t.Errorf("expected empty stack error")
		}
	})
	t.Run("stack top", func(t *testing.T) {
		s := stack{1, 2, 3}
		if got, err := s.top(); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := 3; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}

func TestStack_Pop(t *testing.T) {
	t.Run("empty stack", func(t *testing.T) {
		var s stack
		if err := s.pop(); err == nil {
			t.Errorf("expected empty stack error")
		}
	})
	t.Run("pop stack", func(t *testing.T) {
		s := stack{1, 2, 3}
		if err := s.pop(); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		expected := stack{1, 2}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
}

func TestStack_Push(t *testing.T) {
	t.Run("empty stack", func(t *testing.T) {
		var s stack
		s.push(1)
		expected := stack{1}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
	t.Run("push stack", func(t *testing.T) {
		s := stack{1, 2, 3}
		s.push(4)
		expected := stack{1, 2, 3, 4}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
}
