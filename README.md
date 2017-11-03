# Gotokenizer
Simple keyword extraction

## Installation
```bash
go install github.com/lucasmenendez/gotagger
```

## Demo
````go
package main

import (
	"fmt"
	"github.com/lucasmenendez/gotagger"
)

func main() {
	var text, lang string = "<input-text>", "<input-lang>"
	if tags, err := gotagger.Tag(lang, text); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(tags)
	}
}
````