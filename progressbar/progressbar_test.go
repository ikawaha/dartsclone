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

package progressbar

import (
	"testing"
)

func TestProgressBar(t *testing.T) {
	p := New()
	if p.ProgressBar != nil {
		t.Error("unexpected not nil")
	}
	t.Run("increments before SetMaximum()", func(t *testing.T) {
		p.Increment()
		p.Increment()
	})
	t.Run("set maximum", func(t *testing.T) {
		p.SetMaximum(100)
	})
	t.Run("increments", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			p.Increment()
		}
	})
	t.Run("increments over maximum", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			p.Increment()
		}
	})
}
