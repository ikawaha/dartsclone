package dawg

const (
	unitSize = 32
)

type unit uint32

func (u unit) child() uint32 {
	return uint32(u) >> 2
}

func (u unit) hasSibling() bool {
	return (uint32(u) & 1) == 1
}

func (u unit) value() uint32 {
	return uint32(u) >> 1
}

func (u unit) isState() bool {
	return (uint32(u) & 2) == 2
}
