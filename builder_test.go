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

package dartsclone

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewDoubleArrayBuilder(t *testing.T) {
	t.Run("w/ values", func(t *testing.T) {
		f, err := os.Open("./internal/_testdata/keys.txt")
		if err != nil {
			t.Errorf("unexpected open file error, %v", err)
		}
		var (
			keys   []string
			values []uint32
		)
		scanner := bufio.NewScanner(f)
		for i := 0; scanner.Scan(); i++ {
			keys = append(keys, scanner.Text())
			values = append(values, uint32(i))
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("unexpected scanner error, %v", err)
		}
		b := NewDoubleArrayBuilder(nil)
		if err := b.Build(keys, values); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	})
	t.Run("w/o values", func(t *testing.T) {
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
		b := NewDoubleArrayBuilder(nil)
		if err := b.Build(keys, nil); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	})
}

func TestDoubleArrayBuilder_WriteTo(t *testing.T) {
	t.Run("build from keys.txt and write to file", func(t *testing.T) {
		f, err := os.Open("./internal/_testdata/keys.txt")
		if err != nil {
			t.Errorf("unexpected open file error, %v", err)
		}
		var (
			keys   []string
			values []uint32
		)
		scanner := bufio.NewScanner(f)
		for i := 0; scanner.Scan(); i++ {
			keys = append(keys, scanner.Text())
			values = append(values, uint32(i))
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("unexpected scanner error, %v", err)
		}
		b := NewDoubleArrayBuilder(nil)
		if err := b.Build(keys, values); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		fp, err := ioutil.TempFile("", "da_write_to_test")
		if err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		defer os.Remove(fp.Name())
		if _, err := b.WriteTo(fp); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	})
}

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
