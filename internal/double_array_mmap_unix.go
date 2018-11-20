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

// +build mmap,linux mmap,darwin

package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
)

const (
	MmapedFileHeaderSize = 8
)

// MmapedDoubleArray represents the TRIE data structure mapped on the virtual memory address.
type MmapedDoubleArray struct {
	raw []byte
	r   *bytes.Reader
}

// OpenMmaped opens the named file of double array and maps it on the memory.
func OpenMmaped(name string) (*MmapedDoubleArray, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
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
		return nil, fmt.Errorf("mmap error, %v", err)
	}

	ret := &MmapedDoubleArray{
		raw: b,
		r:   bytes.NewReader(b[MmapedFileHeaderSize:]),
	}
	runtime.SetFinalizer(ret, (*MmapedDoubleArray).Close)
	return ret, nil
}

// Close deletes the mapped memory and closes the opened file.
func (a *MmapedDoubleArray) Close() error {
	if a.raw == nil {
		return nil
	}
	data := a.raw
	a.raw = nil
	runtime.SetFinalizer(a, nil)
	return syscall.Munmap(data)
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

// ExactMatchSearch searches TRIE by a given keyword and returns the id and it's length if found.
func (a MmapedDoubleArray) ExactMatchSearch(key string) (id, size int, err error) {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return -1, -1, err
	}
	for i := 0; i < len(key); i++ {
		nodePos ^= unit.offset() ^ uint32(key[i])
		unit, err = a.at(nodePos)
		if err != nil {
			return -1, -1, err
		}
		if unit.label() != key[i] {
			return -1, 0, nil
		}
	}
	if !unit.hasLeaf() {
		return -1, 0, nil
	}
	unit, err = a.at(nodePos ^ unit.offset())
	if err != nil {
		return -1, -1, err
	}
	return int(unit.value()), len(key), nil
}

// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the ids and it's lengths if found.
func (a MmapedDoubleArray) CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error) {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return ids, sizes, err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return ids, sizes, err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return ids, sizes, err
			}
			ids = append(ids, int(u.value()))
			sizes = append(sizes, i+1)
		}
	}
	return ids, sizes, nil
}

// CommonPrefixSearchCallback finds keywords sharing common prefix in an input and callback with id and it's length.
func (a MmapedDoubleArray) CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error {
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return err
			}
			callback(int(u.value()), i+1)
		}
	}
	return nil
}
