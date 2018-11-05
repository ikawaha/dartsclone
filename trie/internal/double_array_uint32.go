package internal

import (
	"encoding/binary"
	"fmt"
	"os"
)

type DoubleArrayUint32 struct {
	array []uint32
}

func Open(name string) (*DoubleArrayUint32, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var length int64
	if err := binary.Read(f, binary.LittleEndian, &length); err != nil {
		return nil, fmt.Errorf("broken header, %v", err)
	}
	var ret DoubleArrayUint32
	ret.array = make([]uint32, 0, length/4)
	for i := int64(0); i < length; i += 4 {
		var u uint32
		if err := binary.Read(f, binary.LittleEndian, &u); err != nil {
			return nil, fmt.Errorf("broken array, %v", err)
		}
		ret.array = append(ret.array, u)
	}
	return &ret, nil
}

func (a DoubleArrayUint32) at(i uint32) (unit, error) {
	if int(i) >= len(a.array) {
		return 0, fmt.Errorf("index out of bounds")
	}
	return unit(a.array[i]), nil
}

// ExactMatchSearch searches TRIE by a given keyword and returns the id and it's length if found.
func (a DoubleArrayUint32) ExactMatchSearch(key string) (id, size int, err error) {
	return exactMatchSearch(a, key)
}

// CommonPrefixSearch finds keywords sharing common prefix in an input and returns the ids and it's lengths if found.
func (a DoubleArrayUint32) CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error) {
	return commonPrefixSearch(a, key, offset)
}

// CommonPrefixSearchCallback finds keywords sharing common prefix in an input and callback with id and it's length.
func (a DoubleArrayUint32) CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error {
	return commonPrefixSearchCallback(a, key, offset, callback)
}
