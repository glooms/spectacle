// This package tries to build some sort of api generator
package main

import (
	"./spec"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		spec.Build(os.Args[1])
	} else {
		spec.Build(".")
	}
}
