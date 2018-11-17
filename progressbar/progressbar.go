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
	ansi "github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v2"
)

// ProgressBar represents a progress bar that implements ProgressFunction interface.
type ProgressBar struct {
	*progressbar.ProgressBar
}

// New create a progress bar.
func New() *ProgressBar {
	return &ProgressBar{}
}

// SetMaximum sets the maximum of the progress bar.
func (p *ProgressBar) SetMaximum(max int) {
	p.ProgressBar = progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetBytes(10000),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("[cyan][1/3][reset] Writing moshable file..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

// Increment with increase the current count on the progress bar.
func (p *ProgressBar) Increment() {
	if p.ProgressBar != nil {
		p.Add(1)
	}
}
