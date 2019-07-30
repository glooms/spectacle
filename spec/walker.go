package spec

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var fset *token.FileSet
var out *os.File

func filter(fi os.FileInfo) bool {
	return !strings.Contains(fi.Name(), "test")
}

func walker(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return nil
	}
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if pkg.Name == "decls" {
			p := NewPkg(pkg) // This is the important thing
			fmt.Println(p.String())
		}
	}
	return nil
}

// Explore prints all the .go files (excluding tests...ish) that have the directory path root
// as a parent directory.
func Explore(root string) {
	defer out.Close()
	err := filepath.Walk(root, walker)
	check(err)
}

func init() {
	var err error
  if _, err := os.Stat("./log"); os.IsNotExist(err) {
    os.Mkdir("./log", 0755)
  }
	out, err = os.Create("./log/debug.log")
	check(err)

	fset = token.NewFileSet()
}

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
