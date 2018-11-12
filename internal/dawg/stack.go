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
