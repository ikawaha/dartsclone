package trie

type Trie interface {
	ExactMatchSearch(key string) (id, size int, err error)
	CommonPrefixSearch(key string, offset int) (ids, sizes []int, err error)
	CommonPrefixSearchCallback(key string, offset int, callback func(id, size int)) error
}
