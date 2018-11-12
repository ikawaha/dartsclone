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

package dawg

import (
	"fmt"
	"math/bits"
)

type bitVector struct {
	units   []uint32
	ranks   []int
	numOnes int
	size    int
}

func (v bitVector) get(id uint32) (bool, error) {
	unitID := id / unitSize
	if int(unitID) > len(v.units)-1 {
		return false, fmt.Errorf("index out of bounds")
	}
	return (v.units[unitID] >> uint(id%unitSize) & 1) == 1, nil
}

func (v bitVector) rank(id uint32) (int, error) {
	unitID := id / unitSize
	if int(unitID) > len(v.units)-1 {
		return -1, fmt.Errorf("index out of bounds")
	}
	return v.ranks[unitID] + popCount(v.units[unitID]&(^uint32(0)>>uint(unitSize-(id%unitSize)-1))), nil

}

func (v *bitVector) set(id int, bit bool) error {
	index := id / unitSize
	if index < 0 || index > len(v.units)-1 {
		return fmt.Errorf("index out of bounds")
	}
	if bit {
		v.units[index] = v.units[index] | 1<<uint(id%unitSize)
		return nil
	}
	v.units[index] = v.units[index] & ^(1 << uint(id%unitSize))
	return nil
}

func (v bitVector) empty() bool {
	return len(v.units) == 0
}

func (v *bitVector) append() {
	if (v.size % unitSize) == 0 {
		v.units = append(v.units, 0)
	}
	v.size++
}

func (v *bitVector) finish() {
	v.ranks = make([]int, len(v.units))
	v.numOnes = 0
	for i := 0; i < len(v.units); i++ {
		v.ranks[i] = v.numOnes
		v.numOnes += popCount(v.units[i])
	}
}

func popCount(v uint32) int {
	return bits.OnesCount32(v)
}
