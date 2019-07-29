// This package tries to build some sort of api generator
package main

import (
	"./gdoc"
	"./spec"
	"fmt"
	"os"
)

func GdocEx() bool {
	gdoc.Generate("./animator")
	return true
}

func main() {
	//adoc.Explore("/c/Users/dxa/wsl/git/go/src/go")
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		spec.Explore(os.Args[1])
	} else {
		spec.Explore(".")
	}
}
