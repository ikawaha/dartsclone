package trie

import (
	"fmt"
)

const MaxOffset = 1 << 29

type unit uint32

func (u *unit) setHasLeaf(hasLeaf bool) {
	if hasLeaf {
		*u = unit(int32(*u) | 1<<8)
		return
	}
	*u = unit(uint32(*u) & ^(uint32(1) << 8))
}

func (u *unit) setValue(value uint32) {
	*u = unit(value | (1 << 31))
}

func (u *unit) setLabel(label byte) {
	*u = unit(uint32(*u) & ^uint32(0xFF) | uint32(label))
}

func (u *unit) setOffset(offset uint32) error {
	if offset >= MaxOffset {
		return fmt.Errorf("failed to modify unit, too large offset")
	}
	*u = unit(uint32(*u) & ((1 << 31) | (1 << 8) | 0xFF))
	if offset < 1<<21 {
		*u = unit(uint32(*u) | (offset << 10))
		return nil
	}
	*u = unit(uint32(*u) | ((offset << 2) | (1 << 9)))
	return nil
}

func (u unit) label() byte {
	return byte(uint32(u) & ((1 << 31) | 0xFF))
}

func (u unit) offset() uint32 {
	return (uint32(u) >> 10) << ((uint32(u) & (1 << 9)) >> 6)
}

func (u unit) hasLeaf() bool {
	return ((uint32(u) >> 8) & 1) == 1
}

func (u unit) value() uint32 {
	return uint32(u) & ((1 << 31) - 1)
}
