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
	"encoding/binary"
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/unix"
)

const (
	MmapedFileHeaderSize = 8
)

// MmapedDoubleArray represents the TRIE data structure mapped on the virtual memory address.
type MmapedDoubleArray struct {
	raw []byte
}

// OpenMmaped opens the named file of double array and maps it on the memory.
func OpenMmaped(name string) (*MmapedDoubleArray, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()
	if size != int64(int(size)) {
		return nil, fmt.Errorf("too large file")
	}
	return openMmap(f, 0, int(size))
}

func openMmap(f *os.File, offset, size int) (*MmapedDoubleArray, error) {
	if int64(offset)%int64(os.Getpagesize()) != 0 {
		return nil, fmt.Errorf("offset parameter must be a multiple of the system's page size")
	}
	if size%unitSize != 0 {
		return nil, fmt.Errorf("invalid file size, %v", size)
	}
	b, err := unix.Mmap(int(f.Fd()), int64(offset), size, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap error, %v", err)
	}
	ret := &MmapedDoubleArray{
		raw: b,
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
	return unix.Munmap(data)
}

func (a MmapedDoubleArray) at(i uint32) (unit, error) {
	if int(i+1)*unitSize > len(a.raw) {
		return 0, fmt.Errorf("index out of bounds")
	}
	ret := binary.LittleEndian.Uint32(a.raw[i*unitSize : (i+1)*unitSize])
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

// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the array of pairs (id and it's length) if found.
func (a MmapedDoubleArray) CommonPrefixSearch(key string, offset int) ([][2]int, error) {
	var ret [][2]int
	nodePos := uint32(0)
	unit, err := a.at(nodePos)
	if err != nil {
		return ret, err
	}
	nodePos ^= unit.offset()
	for i := offset; i < len(key); i++ {
		k := key[i]
		nodePos ^= uint32(k)
		unit, err := a.at(nodePos)
		if err != nil {
			return ret, err
		}
		if unit.label() != k {
			break
		}
		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			u, err := a.at(nodePos)
			if err != nil {
				return ret, err
			}
			ret = append(ret, [2]int{int(u.value()), i + 1})
		}
	}
	return ret, nil
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
