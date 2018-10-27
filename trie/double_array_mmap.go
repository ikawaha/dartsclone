// +build mmap

package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"syscall"
)

type MmapedDoubleArray struct {
	raw *bytes.Reader
}

func OpenMmaped(name string) (Trie, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	var length int64
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		return nil, fmt.Errorf("broken header, %v", err)
	}
	offset := binary.Size(length)
	return openMmap(f, offset, int(length))
}

func openMmap(f *os.File, offset, length int) (*MmapedDoubleArray, error) {
	b, err := syscall.Mmap(int(f.Fd()), int64(offset), length, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap error, %v", err)
	}
	return &MmapedDoubleArray{
		raw: bytes.NewReader(b),
	}, nil
}

func (a MmapedDoubleArray) at(i uint32) (unit, error) {
	if _, err := a.raw.Seek(int64(i*4), io.SeekStart); err != nil {
		return 0, fmt.Errorf("seek error, %v", err)
	}
	var ret uint32
	if err := binary.Read(a.raw, binary.LittleEndian, &ret); err != nil {
		return 0, fmt.Errorf("read error, %v", err)
	}
	return unit(ret), nil
}

func (a MmapedDoubleArray) ExactMatchSearch(key string) (id, size int, err error) {
	return exactMatchSearch(a, key)
}

func (a MmapedDoubleArray) CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error) {
	return commonPrefixSearch(a, key, offset)
}

func (a MmapedDoubleArray) CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error {
	return commonPrefixSearchCallback(a, key, offset, callback)
}
