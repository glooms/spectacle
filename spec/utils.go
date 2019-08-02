package spec

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/printer"
	"go/token"
	"os"
)

func log(a ...interface{}) {
	fmt.Fprintln(out, a...)
}

func vlog(i interface{}, prefix ...interface{}) {
	// Colored output doesn't work well with vim
	// fmt.Fprintf(out, "\x1b[38;2;%d;%d;%dm", 0xA0, 0xA0, 0x10)
	fmt.Fprint(out, prefix...)
	// fmt.Fprint(out, "\x1b[0m")
	fmt.Fprintf(out, "%#v\n", i)
}

// docPrint pretty prints all decls using "go/doc" and "go/printer".
// Use it as a reference when needed.
func docPrint(pkg *ast.Package, fset *token.FileSet, path string) {
	p := doc.New(pkg, path, doc.AllDecls)
	fmt.Println("doc:", p.Doc)
	fmt.Println("name:", p.Name)
	fmt.Println("imports:")
	for _, imp := range p.Imports {
		fmt.Println("  " + imp)
	}
	fmt.Println("files:")
	for _, file := range p.Filenames {
		fmt.Println("  " + file)
	}
	fmt.Println("consts:")
	for _, c := range p.Consts {
		printer.Fprint(os.Stdout, fset, c.Decl)
		fmt.Println()
	}
	fmt.Println("types:")
	for _, t := range p.Types {
		printer.Fprint(os.Stdout, fset, t.Decl)
		fmt.Println()
	}
	fmt.Println("vars:")
	for _, v := range p.Vars {
		printer.Fprint(os.Stdout, fset, v.Decl)
		fmt.Println()
	}
	fmt.Println("funcs:")
	for _, f := range p.Funcs {
		printer.Fprint(os.Stdout, fset, f.Decl)
		fmt.Println()
	}
}
