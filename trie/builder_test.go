package trie

import (
	"testing"
)

func TestBuildDoubleArray(t *testing.T) {
	keys := []string{
		"電気",
		"電気通信",
		"電気通信大学",
		"電気通信大学大学院",
		"電気通信大学大学院大学",
	}
	t.Run("build", func(t *testing.T) {
		a, err := BuildDoubleArray(keys, nil, nil)
		if err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		t.Run("check", func(t *testing.T) {
			ids, sizes, err := a.CommonPrefixSearch("電気通信大学大学院大学", 0)
			if err != nil {
				t.Errorf("unexpected error, %v", err)
			}
			for i := 0; i < len(ids); i++ {
				if got, expected := ids[i], i; got != expected {
					t.Errorf("got %v, expected %v", got, expected)
				}
				if got, expected := "電気通信大学大学院大学"[0:sizes[i]], keys[i]; got != expected {
					t.Errorf("got %v, expected %v", got, expected)
				}
			}
		})
	})
}
