// This package tries to build some sort of api generator
package main

import (
	_ "./decls" // See if it breakes or not
	"./spec"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		spec.Explore(os.Args[1])
	} else {
		spec.Explore(".")
	}
}
