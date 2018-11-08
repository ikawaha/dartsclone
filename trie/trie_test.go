package trie

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/ikawaha/da"
)

func BenchmarkTrie(b *testing.B) {
	f, err := os.Open("./internal/_testdata/keys.txt")
	if err != nil {
		b.Fatalf("unexpected open file error, %v", err)
	}
	var keys []string
	scanner := bufio.NewScanner(f)
	for i := 0; scanner.Scan(); i++ {
		keys = append(keys, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		b.Fatalf("unexpected scanner error, %v", err)
	}

	b.Run("github.com/ikawaha/da", func(b *testing.B) {
		da := da.DoubleArray{}
		err = da.Build(keys)
		if err := da.Build(keys); err != nil {
			b.Fatalf("unexpected error, da build, %v", err)
		}
		b.Run("exact match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if _, ok := da.Find(strings.NewReader(v)); !ok {
						b.Fatalf("unexpected error, missing the keyword %v", v)
					}
				}
			}
		})
		b.Run("common prefix match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if ids, _ := da.CommonPrefixSearch(strings.NewReader(v)); len(ids) == 0 {
						b.Fatalf("unexpected error, missing a keyword %v", v)
					}
				}
			}
		})
	})
	b.Run("dartsclone", func(b *testing.B) {
		da, err := Open("./internal/_testdata/da_keys")
		if err != nil {
			b.Fatalf("unexpected error, dartsclone open, %v", err)
		}
		b.Run("exact match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if id, _, err := da.ExactMatchSearch(v); id < 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("common prefix match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if ids, _, err := da.CommonPrefixSearch(v, 0); len(ids) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("dartsclone mmap", func(b *testing.B) {
		da, err := OpenMmaped("./internal/_testdata/da_keys")
		defer da.Close()
		if err != nil {
			b.Fatalf("unexpected error, dartsclone open, %v", err)
		}
		b.Run("exact match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if id, _, err := da.ExactMatchSearch(v); id < 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("common prefix match search", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range keys {
					if ids, _, err := da.CommonPrefixSearch(v, 0); len(ids) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
}
