package dawg

import (
	"fmt"
)

type stack []int

func (s stack) top() (int, error) {
	if len(s) == 0 {
		return 0, fmt.Errorf("empty stack error")
	}
	return s[len(s)-1], nil
}

func (s *stack) pop() error {
	if s == nil || len(*s) == 0 {
		return fmt.Errorf("empty stack error")
	}
	*s = (*s)[0 : len(*s)-1]
	return nil
}

func (s *stack) push(item int) {
	*s = append(*s, item)
}
