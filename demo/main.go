package main

import (
	"fmt"
	"io/ioutil"
	"github.com/lucasmenendez/gotagger"
)

func main() {
	var c string
	if r, e := ioutil.ReadFile("demo/input"); e != nil {
		fmt.Println(e)
		return
	} else {
		c = string(r)
	}


	if tags, err := gotagger.Tag("es", c); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(tags)
	}
}