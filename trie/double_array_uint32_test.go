package trie

import (
	"bufio"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestDoubleArrayUint32_ExactMatchSearch(t *testing.T) {
	keys := []string{
		"a",
		"aa",
		"b",
		"cc",
		"hello",
		"world",
		"こんにちは",
	}
	t.Run("keys", func(t *testing.T) {
		a, err := BuildDoubleArray(keys, nil, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		for i, v := range keys {
			id, size, err := a.ExactMatchSearch(v)
			if err != nil {
				t.Errorf("unexpected error, %v", err)
			}
			if id != i || size != len(v) {
				t.Errorf("expected id=%v, size=%v, got id=%v, size=%v (%v)", i, len(v), id, size, string(v))
			}
		}
	})
	t.Run("keys and ids", func(t *testing.T) {
		ids := make([]uint32, len(keys))
		for i := range keys {
			ids[i] = uint32(i * 7)
		}
		a, err := BuildDoubleArray(keys, ids, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		for i, v := range keys {
			id, size, err := a.ExactMatchSearch(v)
			if err != nil {
				t.Errorf("unexpected error, %v", err)
			}
			if id != int(ids[i]) || size != len(v) {
				t.Errorf("expected id=%v, size=%v, got id=%v, size=%v (%v)", i, len(v), id, size, string(v))
			}
		}
	})
}

func TestDoubleArrayUint32_CommonPrefixSearch(t *testing.T) {
	keys := []string{
		"hello",
		"world",
		"電気",
		"電気通信",
		"電気通信大学",
		"電気通信大学大学院",
		"電気通信大学大学院大学",
	}
	t.Run("keys", func(t *testing.T) {
		a, err := BuildDoubleArray(keys, nil, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		ids, sizes, err := a.CommonPrefixSearch("電気通信大学大学院大学", 0)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if expected := []int{2, 3, 4, 5, 6}; !reflect.DeepEqual(expected, ids) {
			t.Errorf("ids: expected %v, got %v", expected, ids)
		}
		if expected := []int{6, 12, 18, 27, 33}; !reflect.DeepEqual(expected, sizes) {
			t.Errorf("sizes: expected %v, got %v", expected, sizes)
		}
	})
}

func TestDoubleArrayUint32_CommonPrefixSearchCallback(t *testing.T) {
	keys := []string{
		"hello",
		"world",
		"電気",
		"電気通信",
		"電気通信大学",
		"電気通信大学大学院",
		"電気通信大学大学院大学",
	}
	t.Run("keys", func(t *testing.T) {
		a, err := BuildDoubleArray(keys, nil, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		var ids, sizes []int
		a.CommonPrefixSearchCallback("電気通信大学大学院大学", 0, func(id, size int) {
			ids = append(ids, id)
			sizes = append(sizes, size)
		})
		if expected := []int{2, 3, 4, 5, 6}; !reflect.DeepEqual(expected, ids) {
			t.Errorf("ids: expected %v, got %v", expected, ids)
		}
		if expected := []int{6, 12, 18, 27, 33}; !reflect.DeepEqual(expected, sizes) {
			t.Errorf("sizes: expected %v, got %v", expected, sizes)
		}
	})
}

func TestOpen(t *testing.T) {
	f, err := os.Open("./_testdata/keys.txt")
	if err != nil {
		t.Errorf("unexpected open file error, %v", err)
	}
	var keys []string
	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		keys = append(keys, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Errorf("unexpected scanner error, %v", err)
	}
	sort.Strings(keys)

	da, err := Open("./_testdata/da_keys")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	for i, key := range keys {
		id, size, err := da.ExactMatchSearch(key)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		if expected := i; id != expected {
			t.Errorf("id: expected %v, got %v", i, id)
		}
		if expected := len(key); size != expected {
			t.Errorf("size: expected %v, got %v", i, id)
		}
	}
}
