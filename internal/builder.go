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
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ikawaha/dartsclone/internal/dawg"
)

const (
	blockSize      = 256
	numExtraBlocks = 16
	numExtras      = blockSize * numExtraBlocks

	upperMask = 0xFF << 21
	lowerMask = 0xFF
)

type extraUnit struct {
	prev    int
	next    int
	isFixed bool
	isUsed  bool
}

// DoubleArrayBuilder represents the builder of the double array.
type DoubleArrayBuilder struct {
	units      []unit
	extras     []extraUnit
	labels     []byte
	table      []int
	extrasHead int

	progress ProgressFunction
}

// BuildDoubleArray constructs a double array from given keywords and values.
// The parameter values sets nil if no values.
func BuildDoubleArray(keys []string, values []uint32, progress ProgressFunction) (*DoubleArrayUint32, error) {
	b := NewDoubleArrayBuilder(progress)
	if err := b.Build(keys, values); err != nil {
		return nil, fmt.Errorf("build error, %v", err)
	}
	return &DoubleArrayUint32{array: b.toArray()}, nil
}

// NewDoubleArrayBuilder returns a builder of the double array with progress function.
// The parameter progress sets nil if no progress bar.
func NewDoubleArrayBuilder(progress ProgressFunction) *DoubleArrayBuilder {
	return &DoubleArrayBuilder{
		progress: progress,
	}
}

// Build constructs a double array from given keys and values.
func (b *DoubleArrayBuilder) Build(keys []string, values []uint32) error {
	keySet, err := newSortedKeySet(keys, values)
	if err != nil {
		return fmt.Errorf("build key set, %v", err)
	}
	if !keySet.hasValues() {
		if err := b.buildFromKeySetHeader(keySet); err != nil {
			return fmt.Errorf("build from key set header, %v", err)
		}
		return nil
	}
	g, err := b.buildDAWG(keySet)
	if err != nil {
		return fmt.Errorf("build DAWG, %v", err)
	}
	if err := b.buildFromDAWGHeader(g); err != nil {
		return fmt.Errorf("build from DAWG header, %v", err)
	}
	return nil
}

func (b DoubleArrayBuilder) toArray() []uint32 {
	var ret []uint32
	for _, u := range b.units {
		ret = append(ret, uint32(u))
	}
	return ret
}

// WriteTo write to the serialize data of the double array.
func (b DoubleArrayBuilder) WriteTo(w io.Writer) (int64, error) {
	var size int64
	x := int64(len(b.units) * 4) // size * uint32
	if err := binary.Write(w, binary.LittleEndian, x); err != nil {
		return size, fmt.Errorf("header write error, %v", err)
	}
	size += int64(binary.Size(x))
	for _, v := range b.units {
		if err := binary.Write(w, binary.LittleEndian, uint32(v)); err != nil {
			return size, err
		}
		size += 4
	}
	return size, nil
}

func (b DoubleArrayBuilder) numBlocks() int {
	return len(b.units) / blockSize
}

func (b DoubleArrayBuilder) getExtras(id int) *extraUnit {
	return &b.extras[id%numExtras]
}

func (b DoubleArrayBuilder) buildDAWG(keySet *keySet) (*dawg.Graph, error) {
	dawgBuilder := dawg.NewBuilder()
	for i := 0; i < keySet.size(); i++ {
		k, err := keySet.getKey(i)
		if err != nil {
			return nil, fmt.Errorf("key set get key, %v", err)
		}
		v, err := keySet.getValue(i)
		if err != nil {
			return nil, fmt.Errorf("key set get value, %v", err)
		}
		if err := dawgBuilder.Insert(k, v); err != nil {
			return nil, fmt.Errorf("DAWG builder insert, %v", err)
		}

		// progress bar
		if b.progress != nil {
			b.progress.Increment(1)
		}
	}
	g, err := dawgBuilder.Finish()
	if err != nil {
		return nil, fmt.Errorf("DAWG builder finish, %v", err)
	}
	return g, nil
}

func (b *DoubleArrayBuilder) buildFromDAWGHeader(g *dawg.Graph) error {
	numUnits := 1
	for numUnits < g.Size() {
		numUnits <<= 1
	}
	b.table = make([]int, g.NumIntersections())
	b.extras = make([]extraUnit, numExtras)

	b.reserveID(0)
	b.extras[0].isUsed = true
	if err := b.units[0].setOffset(1); err != nil {
		return fmt.Errorf("set offset, %v", err)
	}
	b.units[0].setLabel(0)

	if id, err := g.Child(g.Root()); err != nil {
		return fmt.Errorf("invalid root child, %v", err)
	} else if id != 0 {
		if err := b.buildFromDAWGInsert(g, g.Root(), 0); err != nil {
			return fmt.Errorf("insert from DAWG, %v", err)
		}
	}
	b.fixAllBlocks()
	b.extras = nil
	b.labels = nil
	b.table = nil
	return nil
}

func (b *DoubleArrayBuilder) buildFromDAWGInsert(g *dawg.Graph, dawgID uint32, dicID int) error {
	dawgChildID, _ := g.Child(dawgID)
	if ok, err := g.IsIntersection(dawgChildID); err != nil {
		return fmt.Errorf("invalid intersection, %v", err)
	} else if ok {
		intersectionID, err := g.IntersectionID(dawgChildID)
		if err != nil {
			return fmt.Errorf("invalid intersection ID, %v", err)
		}
		offset := b.table[intersectionID]
		if offset != 0 {
			offset ^= dicID
			if (offset&upperMask) == 0 || (offset&lowerMask) == 0 {
				ok, err := g.IsLeaf(dawgChildID)
				if err != nil {
					return fmt.Errorf("invalid leaf, %v", err)
				}
				if ok {
					b.units[dicID].setHasLeaf(true)
				}
				if err := b.units[dicID].setOffset(uint32(offset)); err != nil {
					return fmt.Errorf("set offset, %v", err)
				}
				return nil
			}
		}
	}

	offset, err := b.arrangeFromDAWG(g, dawgID, dicID)
	if err != nil {
		return fmt.Errorf("arrange from DAWG, %v", err)
	}
	if ok, err := g.IsIntersection(dawgChildID); err != nil {
		return fmt.Errorf("invalid intersection, %v", err)
	} else if ok {
		iid, err := g.IntersectionID(dawgChildID)
		if err != nil {
			return fmt.Errorf("invalid intersection ID, %v", err)
		}
		b.table[iid] = offset
	}

	for {
		childLabel, _ := g.Label(dawgChildID)
		dicChildID := offset ^ int(childLabel)
		if childLabel != 0 {
			if err := b.buildFromDAWGInsert(g, dawgChildID, dicChildID); err != nil {
				return fmt.Errorf("insert from DAWG, %v", err)
			}
		}
		var err error
		dawgChildID, err = g.Sibling(dawgChildID)
		if err != nil {
			return fmt.Errorf("sibling, %v", err)
		}
		if dawgChildID == 0 {
			break
		}
	}
	return nil
}

func (b *DoubleArrayBuilder) arrangeFromDAWG(g *dawg.Graph, dawgID uint32, dicID int) (int, error) {
	if dicID < 0 || dicID >= len(b.units) {
		return -1, fmt.Errorf("dicID, index out of bounds, %v", dicID)
	}
	b.labels = []byte{}
	dawgChildID, _ := g.Child(dawgID)

	for dawgChildID != 0 {
		label, err := g.Label(dawgChildID)
		if err != nil {
			return -1, fmt.Errorf("label, %v", err)
		}
		b.labels = append(b.labels, label)
		dawgChildID, err = g.Sibling(dawgChildID)
		if err != nil {
			return -1, fmt.Errorf("sibling, %v", err)
		}
	}
	offset := b.findValidOffset(dicID)
	if err := b.units[dicID].setOffset(uint32(dicID ^ offset)); err != nil {
		return -1, fmt.Errorf("set offset, %v", err)
	}

	var err error
	dawgChildID, err = g.Child(dawgID)
	if err != nil {
		return -1, fmt.Errorf("child, %v", err)
	}
	for _, l := range b.labels {
		dicChildID := offset ^ int(l)
		b.reserveID(dicChildID)
		if ok, err := g.IsLeaf(dawgChildID); err != nil {
			return -1, fmt.Errorf("invalid leaf, %v", err)
		} else if !ok {
			b.units[dicChildID].setLabel(l)
		} else {
			b.units[dicID].setHasLeaf(true)
			v, err := g.Value(dawgChildID)
			if err != nil {
				return -1, fmt.Errorf("invalid value, %v", err)
			}
			b.units[dicChildID].setValue(v)
		}
		dawgChildID, err = g.Sibling(dawgChildID)
		if err != nil {
			return -1, fmt.Errorf("sibling, %v", err)
		}
	}
	b.getExtras(offset).isUsed = true
	return offset, nil
}

func (b *DoubleArrayBuilder) buildFromKeySetHeader(keySet *keySet) error {
	numUnits := 1
	for numUnits < keySet.size() {
		numUnits <<= 1
	}
	b.extras = make([]extraUnit, numExtras)

	b.reserveID(0)
	b.extras[0].isUsed = true
	if err := b.units[0].setOffset(1); err != nil {
		return fmt.Errorf("set offset, %v", err)
	}
	b.units[0].setLabel(0)

	if keySet.size() > 0 {
		b.buildFromKeySetInsert(keySet, 0, keySet.size(), 0, 0)
	}

	b.fixAllBlocks()

	b.extras = nil
	b.labels = nil

	return nil
}

func (b *DoubleArrayBuilder) buildFromKeySetInsert(keySet *keySet, begin, end, depth, dicID int) error {
	offset, err := b.arrangeFromKeySet(keySet, begin, end, depth, dicID)
	if err != nil {
		return fmt.Errorf("arrange from key set, %v", err)
	}
	for begin < end {
		if b, err := keySet.getKeyByte(begin, depth); err != nil {
			return fmt.Errorf("get key byte, %v", err)
		} else if b != 0 {
			break
		}
		begin++
	}
	if begin == end {
		return nil
	}
	lastBegin := begin
	lastLabel, err := keySet.getKeyByte(begin, depth)
	if err != nil {
		return fmt.Errorf("get key byte, %v", err)
	}
	for {
		begin++
		if begin >= end {
			break
		}
		label, err := keySet.getKeyByte(begin, depth)
		if err != nil {
			return fmt.Errorf("get key byte, %v", err)
		}
		if label != lastLabel {
			b.buildFromKeySetInsert(keySet, lastBegin, begin, depth+1, offset^int(lastLabel))
			lastBegin = begin
			lastLabel, err = keySet.getKeyByte(begin, depth)
			if err != nil {
				return fmt.Errorf("get key byte, %v", err)
			}
		}
	}
	b.buildFromKeySetInsert(keySet, lastBegin, end, depth+1, int(uint32(offset)^uint32(lastLabel)))
	return nil
}

func (b *DoubleArrayBuilder) arrangeFromKeySet(keySet *keySet, begin, end, depth, dicID int) (int, error) {
	b.labels = []byte{}
	value := -1
	for i := begin; i < end; i++ {
		label, err := keySet.getKeyByte(i, depth)
		if err != nil {
			return -1, fmt.Errorf("%+v, i=%v, depth=%v, %v", keySet, i, depth, err)
		}
		if label == 0 {
			key, err := keySet.getKey(i)
			if err != nil {
				return -1, fmt.Errorf("get key (%v), %v", i, err)
			}
			if depth < len(key) {
				return -1, fmt.Errorf("invalid null character, %v", key)
			}
			if value == -1 {
				val, err := keySet.getValue(i)
				if err != nil {
					panic(err)
				}
				value = int(val)
			}
			// progress bar
			if b.progress != nil {
				b.progress.Increment(1)
			}
		}
		if len(b.labels) == 0 {
			b.labels = append(b.labels, label)
		} else if label != b.labels[len(b.labels)-1] {
			if label < b.labels[len(b.labels)-1] {
				return -1, fmt.Errorf("wrong key order")
			}
			b.labels = append(b.labels, label)
		}
	}

	offset := b.findValidOffset(dicID)
	if dicID < 0 || dicID >= len(b.units) {
		return -1, fmt.Errorf("dicID, index out of bounds, %v", dicID)
	}
	if err := b.units[dicID].setOffset(uint32(dicID ^ offset)); err != nil {
		return -1, fmt.Errorf("set offset, %v", err)
	}

	for _, l := range b.labels {
		dicChildID := offset ^ int(l)
		b.reserveID(dicChildID)
		if l != 0 {
			b.units[dicChildID].setLabel(l)
		} else {
			b.units[dicID].setHasLeaf(true)
			b.units[dicChildID].setValue(uint32(value))
		}
	}
	b.getExtras(offset).isUsed = true

	return offset, nil
}

func (b DoubleArrayBuilder) findValidOffset(id int) int {
	if b.extrasHead >= len(b.units) {
		return len(b.units) | (id & lowerMask)
	}
	unfixedID := b.extrasHead
	memo := map[int]struct{}{}
	for {
		if _, ok := memo[unfixedID]; !ok {
			memo[unfixedID] = struct{}{}
		} else {
			panic(fmt.Sprintf("runtime error, unfixedID=%v, memo=%+v", unfixedID, memo))
		}
		offset := unfixedID ^ int(b.labels[0])
		if b.isValidOffset(id, offset) {
			return offset
		}
		unfixedID = b.getExtras(unfixedID).next
		if unfixedID == b.extrasHead {
			break
		}
	}
	return len(b.units) | (id & lowerMask)
}

func (b DoubleArrayBuilder) isValidOffset(id, offset int) bool {
	if b.getExtras(offset).isUsed {
		return false
	}
	relOffset := id ^ offset
	if (relOffset&lowerMask) != 0 && (relOffset&upperMask) != 0 {
		return false
	}
	for i := 1; i < len(b.labels); i++ {
		if b.getExtras(offset ^ int(b.labels[i])).isFixed {
			return false
		}
	}
	return true
}

func (b *DoubleArrayBuilder) reserveID(id int) {
	if id >= len(b.units) {
		b.expandUnits()
	}
	if id == b.extrasHead {
		b.extrasHead = b.getExtras(id).next
		if b.extrasHead == id {
			b.extrasHead = len(b.units)
		}
	}
	b.getExtras(b.getExtras(id).prev).next = b.getExtras(id).next
	b.getExtras(b.getExtras(id).next).prev = b.getExtras(id).prev
	b.getExtras(id).isFixed = true
}

func (b *DoubleArrayBuilder) expandUnits() {
	srcNumUnits := len(b.units)
	srcNumBlocks := b.numBlocks()

	destNumUnits := srcNumUnits + blockSize
	destNumBlocks := srcNumBlocks + 1

	if destNumBlocks > numExtraBlocks {
		b.fixBlock(srcNumBlocks - numExtraBlocks)
	}
	for i := srcNumUnits; i < destNumUnits; i++ {
		b.units = append(b.units, unit(0))
	}
	if destNumBlocks > numExtraBlocks {
		for id := srcNumUnits; id < destNumUnits; id++ {
			b.getExtras(id).isUsed = false
			b.getExtras(id).isFixed = false
		}
	}
	for i := srcNumUnits + 1; i < destNumUnits; i++ {
		b.getExtras(i - 1).next = i
		b.getExtras(i).prev = i - 1
	}
	b.getExtras(srcNumUnits).prev = destNumUnits
	b.getExtras(destNumUnits - 1).next = srcNumUnits

	b.getExtras(srcNumUnits).prev = b.getExtras(b.extrasHead).prev
	b.getExtras(destNumUnits - 1).next = b.extrasHead

	b.getExtras(b.getExtras(b.extrasHead).prev).next = srcNumUnits
	b.getExtras(b.extrasHead).prev = destNumUnits - 1
}

func (b *DoubleArrayBuilder) fixAllBlocks() {
	begin := 0

	if b.numBlocks() > numExtraBlocks {
		begin = b.numBlocks() - numExtraBlocks
	}
	end := b.numBlocks()

	for blockID := begin; blockID < end; blockID++ {
		b.fixBlock(blockID)
	}
}

func (b *DoubleArrayBuilder) fixBlock(blockID int) {
	begin := blockID * blockSize
	end := begin + blockSize

	unusedOffset := 0
	for offset := begin; offset < end; offset++ {
		if !b.getExtras(offset).isUsed {
			unusedOffset = offset
			break
		}
	}
	for id := begin; id < end; id++ {
		if !b.getExtras(id).isFixed {
			b.reserveID(id)
			b.units[id].setLabel(byte(id ^ unusedOffset))
		}
	}
}
