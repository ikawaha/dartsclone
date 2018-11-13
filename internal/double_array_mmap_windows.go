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
	low, high := uint32(length), uint32(length>>32)
	fm, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, high, low, nil)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(fm)
	ptr, err := syscall.MapViewOfFile(fm, syscall.FILE_MAP_READ, 0, 0, uintptr(length))
	if err != nil {
		return nil, err
	}
	b := (*[maxBytes]byte)(unsafe.Pointer(ptr))[:size]

	ret := &MmapedDoubleArray{
		raw: b,
		r:   bytes.NewReader(b[MmapedFileHeaderSize:]),
	}
	runtime.SetFinalizer(ret, (*MmapedDoubleArray).Close)
	return ret, nil
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
	return exactMatchSearch(a, key)
}

// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the ids and it's lengths if found.
func (a MmapedDoubleArray) CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error) {
	return commonPrefixSearch(a, key, offset)
}

// CommonPrefixSearchCallback finds keywords sharing common prefix in an input and callback with id and it's length.
func (a MmapedDoubleArray) CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error {
	return commonPrefixSearchCallback(a, key, offset, callback)
}

// Close deletes the mapped memory and closes the opened file.
func (a *MmapedDoubleArray) Close() error {
	if a.raw == nil {
		return nil
	}
	data := a.raw
	a.raw = nil
	runtime.SetFinalizer(a, nil)
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&data[0])))
}
