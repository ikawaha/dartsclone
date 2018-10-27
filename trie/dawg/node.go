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
