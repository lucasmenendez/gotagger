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
	"github.com/lucasmenendez/gotokenizer"
	"github.com/lucasmenendez/gotagger"
)

func main() {
	var limit int = 15
	var text, lang string = "<input-text>", "<input-lang>"
	
	sentences := gotokenizer.Sentences(text)
    
    var words [][]string
    for _, s := range sentences {
        words = append(words, gotokenizer.Words(s))
    }

    if tags, err := gotagger.GetTags(words, lang, limit); err != nil {
        fmt.Println(err)
    } else {
    	fmt.Printf("%q\n", tags)
    }
}
````