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

const (
	MmapedFileHeaderSize = 8
)

type MmapedDoubleArray struct {
	f   *os.File
	raw []byte
	r   *bytes.Reader
}

func OpenMmaped(name string) (TrieCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	var length int64
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		return nil, fmt.Errorf("broken header, %v", err)
	}
	return openMmap(f, 0, MmapedFileHeaderSize+int(length))
}

func openMmap(f *os.File, offset, length int) (*MmapedDoubleArray, error) {
	if int64(offset)%int64(os.Getpagesize()) != 0 {
		return nil, fmt.Errorf("offset parameter must be a multiple of the system's page size")
	}
	b, err := syscall.Mmap(int(f.Fd()), int64(offset), length, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		b, err = syscall.Mmap(int(f.Fd()), int64(offset), length, syscall.PROT_READ, syscall.MAP_PRIVATE)
		if err != nil {
			return nil, fmt.Errorf("mmap error, %v", err)
		}
	}
	return &MmapedDoubleArray{
		f:   f,
		raw: b,
		r:   bytes.NewReader(b[MmapedFileHeaderSize:]),
	}, nil
}

func (a MmapedDoubleArray) at(i uint32) (unit, error) {
	if _, err := a.r.Seek(int64(i*4), io.SeekStart); err != nil {
		return 0, fmt.Errorf("seek error, %v", err)
	}
	var ret uint32
	if err := binary.Read(a.r, binary.LittleEndian, &ret); err != nil {
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

func (a MmapedDoubleArray) Close() error {
	if err := syscall.Munmap(a.raw); err != nil {
		return fmt.Errorf("munmap error, %v", err)
	}
	return a.f.Close()
}
