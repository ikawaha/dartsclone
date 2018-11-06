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
)

type node struct {
	label byte

	child      uint32
	sibling    int
	isState    bool
	hasSibling bool
}

func (n *node) reset() {
	n.child = 0
	n.sibling = 0
	n.label = 0
	n.isState = false
	n.hasSibling = false
}

func (n node) value() uint32 {
	return n.child
}

func (n *node) setValue(value uint32) {
	n.child = value
}

func (n node) unit() unit {
	hasSibling := uint32(0)
	if n.hasSibling {
		hasSibling = 1
	}
	if n.label == 0 {
		return unit((n.child << 1) | hasSibling)
	}
	isState := uint32(0)
	if n.isState {
		isState = 2
	}
	return unit((n.child << 2) | isState | hasSibling)
}

func (n *node) String() string {
	return fmt.Sprintf("child: %d, sibling: %d, label: %c, is_state: %v, has_sibling: %v",
		n.child, n.sibling, n.label, n.isState, n.hasSibling)
}
