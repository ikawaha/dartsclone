package trie

type ProgressFunction interface {
	Increment(int)
}
