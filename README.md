# Gotokenizer
Simple keyword extraction

## Installation
```bash
go install github.com/lucasmenendez/gotagger
```

Then, set env var `STOPWORDS` to stopword lists path or store it into existing folder.

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