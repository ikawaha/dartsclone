package internal

type arrayAt interface {
	at(uint32) (unit, error)
}

func exactMatchSearch(a arrayAt, key string) (id, size int, err error) {
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

func commonPrefixSearch(a arrayAt, key string, offset int) (ids, sizes []int, err error) {
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

func commonPrefixSearchCallback(a arrayAt, key string, offset int, callback func(id, size int)) error {
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
