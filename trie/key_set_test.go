package trie

import (
	"reflect"
	"sort"
	"testing"
)

func TestNewKeySet(t *testing.T) {
	t.Run("empty key set", func(t *testing.T) {
		if _, err := newSortedKeySet(nil, nil); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	})
	t.Run("keys", func(t *testing.T) {
		keys := []string{"hello", "world"}
		s, err := newSortedKeySet(keys, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if !reflect.DeepEqual(keys, s.keys) {
			t.Errorf("expected %+v, got %+v", keys, s.keys)
		}
	})
	t.Run("keys and values", func(t *testing.T) {
		keys := []string{"hi!", "hello", "world"}
		if _, err := newSortedKeySet(keys, []uint32{1, 2, 3, 4, 5}); err == nil {
			t.Errorf("expected invalid input error")
		}
		values := []uint32{100, 200, 300}
		s, err := newSortedKeySet(keys, values)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if !reflect.DeepEqual(keys, s.keys) {
			t.Errorf("expected %+v, got %+v", keys, s.keys)
		}
		if !reflect.DeepEqual(values, s.values) {
			t.Errorf("expected %+v, got %+v", keys, s.keys)
		}
	})
	t.Run("sort keys and values", func(t *testing.T) {
		keys := []string{"charlie", "bravo", "alpha"}
		values := []uint32{3, 2, 1}
		s, err := newSortedKeySet(keys, values)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if !sort.StringsAreSorted(s.keys) {
			t.Errorf("unexpected status, keys are not sorted, %v", keys)
		}
		if !sort.SliceIsSorted(s.values, func(i, j int) bool { return s.values[i] < s.values[j] }) {
			t.Errorf("unexpected status, values are not sorted, %v", values)
		}
	})
}

func TestKeySet_Size(t *testing.T) {
	t.Run("empty key set", func(t *testing.T) {
		var s keySet
		if got, expected := s.size(), 0; got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("keys", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		s, err := newSortedKeySet(keys, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if got, expected := s.size(), len(keys); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("keys and values", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		values := []uint32{1, 3, 5}
		s, err := newSortedKeySet(keys, values)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if got, expected := s.size(), len(keys); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("has duplicate key", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye", "aloha"}
		_, err := newSortedKeySet(keys, nil)
		if err == nil {
			t.Errorf("expected duplicate key error")
		}
	})
}

func TestKeySet_GetKey(t *testing.T) {
	t.Run("index out of bounds", func(t *testing.T) {
		var s keySet
		if _, err := s.getKey(0); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		if _, err := s.getKey(3); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		s.keys = []string{"hello", "goodbye"}
		if _, err := s.getKey(len(s.keys)); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		if _, err := s.getKey(-1); err == nil {
			t.Errorf("expected index out of bounds error")
		}
	})
	t.Run("get a key", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		s, err := newSortedKeySet(keys, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		for i := range keys {
			if got, err := s.getKey(i); err != nil {
			} else if expected := keys[i]; got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		}
	})
}

func TestKeySet_GetValue(t *testing.T) {
	t.Run("index out of bounds", func(t *testing.T) {
		var s keySet
		if _, err := s.getValue(-1); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		s.values = []uint32{3, 7}
		if _, err := s.getValue(len(s.values)); err == nil {
			t.Errorf("expected index out of bounds error")
		}
	})
	t.Run("key set has no values", func(t *testing.T) {
		var s keySet
		if got, err := s.getValue(3); err != nil {
			t.Errorf("unexpected error")
		} else if expected := uint32(3); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("get a value", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		values := []uint32{1, 3, 5}
		s, err := newSortedKeySet(keys, values)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		for i := range keys {
			if got, err := s.getValue(i); err != nil {
			} else if expected := values[i]; got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		}
	})
}

func TestKeySet_GetKeyByte(t *testing.T) {
	t.Run("index out of bounds", func(t *testing.T) {
		var s keySet
		if _, err := s.getKeyByte(0, 0); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		s.keys = []string{"hello", "world"}
		if _, err := s.getKeyByte(1, -1); err == nil {
			t.Errorf("expected index out of bounds error")
		}
		s.keys = []string{"hello", "world"}
		if _, err := s.getKeyByte(len(s.keys), 0); err == nil {
			t.Errorf("expected index out of bounds error")
		}
	})
	t.Run("get a byte of a key", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		s, err := newSortedKeySet(keys, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		for i, key := range keys {
			for j := 0; j < len(key); j++ {
				if got, err := s.getKeyByte(i, j); err != nil {
					t.Errorf("unexpected error, %v", err)
				} else if expected := keys[i][j]; got != expected {
					t.Errorf("expected %v, got %v", expected, got)
				}
			}
		}
	})
	t.Run("get null if index is over the key length", func(t *testing.T) {
		keys := []string{"aloha", "hello", "goodbye"}
		s, err := newSortedKeySet(keys, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if got, err := s.getKeyByte(0, 100); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := byte(0); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}

func TestKeySet_HasValues(t *testing.T) {
	var s keySet
	if got, expected := s.hasValues(), false; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
	s.values = []uint32{}
	if got, expected := s.hasValues(), false; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
	s.values = []uint32{1}
	if got, expected := s.hasValues(), true; got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
