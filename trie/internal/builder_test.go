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
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestDoubleArrayBuilder_Build(t *testing.T) {
	t.Run("small test w/ values", func(t *testing.T) {
		b := NewDoubleArrayBuilder(nil)
		keys := []string{"aaa", "bbb"}
		values := []uint32{7, 5}
		if err := b.Build(keys, values); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		da := DoubleArrayUint32{array: b.toArray()}
		for i, key := range keys {
			if id, size, err := da.ExactMatchSearch(key); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if id != int(values[i]) {
				t.Errorf("unexpected id, expected %v, got %v", values[i], id)
			} else if size != len(key) {
				t.Errorf("unexpected size, expected %v, got %v", len(key), size)
			}
		}
	})
	t.Run("small test w/o values", func(t *testing.T) {
		b := NewDoubleArrayBuilder(nil)
		keys := []string{"aaa", "bbb"}
		if err := b.Build(keys, nil); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		da := DoubleArrayUint32{array: b.toArray()}
		for i, key := range keys {
			if id, size, err := da.ExactMatchSearch(key); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if id != i {
				t.Errorf("unexpected id, expected %v, got %v", i, id)
			} else if size != len(key) {
				t.Errorf("unexpected size, expected %v, got %v", len(key), size)
			}
		}
	})
	t.Run("large test w/ values", func(t *testing.T) {
		f, err := os.Open("./_testdata/keys.txt")
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
		da := DoubleArrayUint32{array: b.toArray()}
		for i, key := range keys {
			if id, size, err := da.ExactMatchSearch(key); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if id != int(values[i]) {
				t.Errorf("unexpected id, expected %v, got %v", values[i], id)
			} else if size != len(key) {
				t.Errorf("unexpected size, expected %v, got %v", len(key), size)
			}
		}

	})
	t.Run("large test w/o values", func(t *testing.T) {
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
		b := NewDoubleArrayBuilder(nil)
		if err := b.Build(keys, nil); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
		da := DoubleArrayUint32{array: b.toArray()}
		for i, key := range keys {
			if id, size, err := da.ExactMatchSearch(key); err != nil {
				t.Errorf("unexpected error, %v", err)
			} else if id != i {
				t.Errorf("unexpected id, expected %v, got %v", i, id)
			} else if size != len(key) {
				t.Errorf("unexpected size, expected %v, got %v", len(key), size)
			}
		}
	})
}

func TestDoubleArrayBuilder_WriteTo(t *testing.T) {
	t.Run("small test", func(t *testing.T) {

		builder := DoubleArrayBuilder{
			units: []unit{1, 2, 3, 4, 5},
		}
		var b bytes.Buffer
		if size, err := builder.WriteTo(&b); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := int64(8 + 4*5); size != expected {
			t.Errorf("expected %v, got %v", expected, size)
		}
		got := b.Bytes()
		expected := []byte{
			20, 0, 0, 0, 0, 0, 0, 0, //size int64(20) <<little endian>>
			1, 0, 0, 0, // uint32(1)
			2, 0, 0, 0, // uint32(2)
			3, 0, 0, 0, // uint32(3)
			4, 0, 0, 0, // uint32(4)
			5, 0, 0, 0, // uint32(5)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
	t.Run("build from keys.txt and write to file", func(t *testing.T) {
		f, err := os.Open("./_testdata/keys.txt")
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
