package internal

type ProgressFunction interface {
	Increment(int)
}
