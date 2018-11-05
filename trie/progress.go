package trie

// ProgressFunction indicates progress bar of building a double array.
type ProgressFunction interface {
	// Increment with increase the current count on the progress bar.
	Increment(int)
}
