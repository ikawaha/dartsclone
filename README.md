# dartsclone : Double Array TRIE liblary

[![Build Status](https://travis-ci.org/ikawaha/dartsclone.svg?branch=master)](https://travis-ci.org/ikawaha/dartsclone)
[![Build status](https://ci.appveyor.com/api/projects/status/2ku3oes7oe7nlw2x/branch/master?svg=true)](https://ci.appveyor.com/project/ikawaha/dartsclone/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/ikawaha/dartsclone/badge.svg)](https://coveralls.io/github/ikawaha/dartsclone)

Port of [Sudachi's dartsclone library](https://github.com/WorksApplications/Sudachi/tree/develop/src/main/java/com/worksap/nlp/dartsclone) to Go. 


## Build & Save

```Go:
package main

import (
	"os"

	"github.com/ikawaha/dartsclone"
)

func main() {
	keys := []string{
		"電気",
		"電気通信",
		"電気通信大学",
		"電気通信大学大学院",
		"電気通信大学大学院大学",
	}

	// Build
	builder := dartsclone.NewBuilder(nil)
	if err := builder.Build(keys, nil); err != nil {
		panic(err)
	}
	// Save
	f, err := os.Create("my-double-array-file")
	if err != nil {
		panic(err)
	}
	builder.WriteTo(f)
	f.Close()
}
```

## Load & Search

```Go:
package main

import (
	"fmt"
	"github.com/ikawaha/dartsclone"
)

func main() {
	trie, err := dartsclone.Open("my-double-array-file")
	if err != nil {
		panic(err)
	}
	ids, sizes, err := trie.CommonPrefixSearch("電気通信大学大学院大学", 0)
	for i := 0; i < len(ids); i++ {
		fmt.Printf("id=%d, common prefix=%s\n", ids[i], "電気通信大学大学院大学"[0:sizes[i]])
	}
}
```

### outputs

```
id=0, common prefix=電気
id=1, common prefix=電気通信
id=2, common prefix=電気通信大学
id=3, common prefix=電気通信大学大学院
id=4, common prefix=電気通信大学大学院大学
```


## Use memory mapping

* Build Tags : mmap
* Support OS : linux, osx


```Go:
package main

import (
	"fmt"
	"github.com/ikawaha/dartsclone"
)

func main() {
	trie, err := dartsclone.OpenMmaped("my-double-array-file") // ← ★
	if err != nil {
		panic(err)
	}
	defer trie.Close() // ← ★

	ids, sizes, err := trie.CommonPrefixSearch("電気通信大学大学院大学", 0)
	for i := 0; i < len(ids); i++ {
		fmt.Printf("id=%d, common prefix=%s\n", ids[i], "電気通信大学大学院大学"[0:sizes[i]])
	}
}
```

