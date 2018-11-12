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

func BenchmarkTrieMmaped(b *testing.B) {
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
}
