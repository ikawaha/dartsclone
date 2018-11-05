// +build mmap

package internal

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMmapedDoubleArray_ExactMatchSearch(t *testing.T) {
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
		builder := DoubleArrayBuilder{}
		builder.Build(keys, nil)
		var b bytes.Buffer
		builder.WriteTo(&b)
		mmaped := MmapedDoubleArray{r: bytes.NewReader(b.Bytes()[8:])} // header size is 8
		for i, v := range keys {
			id, size, err := mmaped.ExactMatchSearch(v)
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
		builder := DoubleArrayBuilder{}
		builder.Build(keys, ids)
		var b bytes.Buffer
		builder.WriteTo(&b)
		mmaped := MmapedDoubleArray{r: bytes.NewReader(b.Bytes()[8:])} // header size is 8
		for i, v := range keys {
			id, size, err := mmaped.ExactMatchSearch(v)
			if err != nil {
				t.Errorf("unexpected error, %v", err)
			}
			if id != int(ids[i]) || size != len(v) {
				t.Errorf("expected id=%v, size=%v, got id=%v, size=%v (%v)", i, len(v), id, size, string(v))
			}
		}
	})
}

func TestMmapedDoubleArray_CommonPrefixSearch(t *testing.T) {
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
		builder := DoubleArrayBuilder{}
		builder.Build(keys, nil)
		var b bytes.Buffer
		builder.WriteTo(&b)
		mmaped := MmapedDoubleArray{r: bytes.NewReader(b.Bytes()[8:])} // header size is 8
		ids, sizes, err := mmaped.CommonPrefixSearch("電気通信大学大学院大学", 0)
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

func TestMmapedDoubleArray_CommonPrefixSearchCallback(t *testing.T) {
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
		builder := DoubleArrayBuilder{}
		builder.Build(keys, nil)
		var b bytes.Buffer
		builder.WriteTo(&b)
		mmaped := MmapedDoubleArray{r: bytes.NewReader(b.Bytes()[8:])} // header size is 8
		var ids, sizes []int
		mmaped.CommonPrefixSearchCallback("電気通信大学大学院大学", 0, func(id, size int) {
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

func TestOpenMmap(t *testing.T) {
	t.Run("open error", func(t *testing.T) {
		_, err := OpenMmaped("./_testdata/not-found-file-error")
		if err == nil {
			t.Fatalf("expected file not found error")
		}
	})
	t.Run("open sample binary of mmaped double array", func(t *testing.T) {
		da, err := OpenMmaped("./_testdata/mmapbin_20_1_2_3_4_5")
		if err != nil {
			t.Fatalf("unexpected error, %v", err)
		}
		defer func() {
			if err := da.Close(); err != nil {
				t.Errorf("unexpected error, %v", err)
			}
		}()
		for i := 0; i < 5; i++ {
			if got, err := da.at(uint32(i)); err != nil {
				t.Errorf("unexpected error, %v (%v)", err, i)
			} else if expected := unit(i + 1); got != expected {
				t.Errorf("expected, %v, got %v", expected, got)
			}
		}
	})
}

func TestMmapedDoubleArray_At(t *testing.T) {
	da, err := OpenMmaped("./_testdata/mmapbin_20_1_2_3_4_5")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	defer func() {
		if err := da.Close(); err != nil {
			t.Errorf("unexpected error, %v", err)
		}
	}()
	t.Run("in range", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			if got, err := da.at(uint32(i)); err != nil {
				t.Errorf("unexpected error, %v (%v)", err, i)
			} else if expected := unit(i + 1); got != expected {
				t.Errorf("expected, %v, got %v", expected, got)
			}
		}
	})
	t.Run("out of range", func(t *testing.T) {
		if _, err := da.at(uint32(5)); err == nil {
			t.Errorf("expected read error")
		}
		// recover ok
		if got, err := da.at(uint32(4)); err != nil {
			t.Errorf("unexpected error, %v", err)
		} else if expected := unit(5); got != expected {
			t.Errorf("expected %v, got %v", expected, got)
		}
	})
}
