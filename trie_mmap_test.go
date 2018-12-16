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

// +build mmap

package dartsclone

import (
	"bufio"
	"os"
	"testing"
)

func BenchmarkTRIEMmaped(b *testing.B) {
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

	trie, err := OpenMmaped("./internal/_testdata/da_keys")
	defer trie.Close()
	if err != nil {
		b.Fatalf("unexpected error, dartsclone open, %v", err)
	}
	b.Run("exact match search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range keys {
				if id, _, err := trie.ExactMatchSearch(v); id < 0 || err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("common prefix match search", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range keys {
				if ret, err := trie.CommonPrefixSearch(v, 0); len(ret) == 0 || err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
				}
			}
		}
	})
}

func TestOpenMmaped(t *testing.T) {
	f, err := os.Open("./internal/_testdata/keys.txt")
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

	trie, err := OpenMmaped("./internal/_testdata/da_keys")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer func() {
		if err := trie.Close(); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}()
	for i, key := range keys {
		id, size, err := trie.ExactMatchSearch(key)
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
