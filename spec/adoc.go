package spec

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var fset *token.FileSet
var dirRoot string
var pad string
var buf bytes.Buffer
var out1 *os.File
var out2 *os.File

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
		fmt.Fprint(&buf, path)
		if !strings.Contains(path, pkg.Name) {
			fmt.Fprint(&buf, " ("+pkg.Name+")")
		}
		fmt.Fprint(&buf, ":")
		indent := buf.Len()
		fmt.Fprintln(&buf)
		buf.WriteTo(out1)
		for fn, file := range pkg.Files {
			fn := strings.Replace(fn, path+"/", "", 1)
			fmt.Fprint(&buf, pad[:indent], fn)
			indent := buf.Len()
			fmt.Fprintln(&buf)
			buf.WriteTo(out1)
			for _, decl := range file.Decls {
				interpret(decl, indent)
				buf.WriteTo(out1)
			}
		}
    if pkg.Name == "decls" {
		  p := New(pkg) // This is the important thing
		  fmt.Println(p.String())
    }
	}
	return nil
}

func interpret(decl ast.Decl, indent int) {
	switch decl.(type) {
	case *ast.GenDecl: //import, const, type, var
		reprGenDecl(decl.(*ast.GenDecl), indent)
	case *ast.FuncDecl: //function
		d := decl.(*ast.FuncDecl)
		reprFuncDecl(d, indent)
	}
}

func reprGenDecl(decl *ast.GenDecl, indent int) {
	switch decl.Tok {
	case token.IMPORT:
		fmt.Fprint(&buf, pad[:indent], "import")
	case token.CONST:
		fmt.Fprint(&buf, pad[:indent], "const")
	case token.TYPE:
		fmt.Fprint(&buf, pad[:indent], "type")
	case token.VAR:
		fmt.Fprint(&buf, pad[:indent], "var")
	}
	indent = buf.Len()
	if len(decl.Specs) == 1 {
		indent = 1
	} else {
		fmt.Fprintln(&buf)
		buf.WriteTo(out1)
	}
	for _, spec := range decl.Specs {
		reprSpec(spec, indent)
	}
}

func reprFuncDecl(decl *ast.FuncDecl, indent int) {
	fmt.Fprint(&buf, pad[:indent], "func "+decl.Name.Name)
	fmt.Fprintln(&buf)
	buf.WriteTo(out1)
	if decl.Recv == nil {
		return
	}
	fmt.Fprintln(out2, "> func "+decl.Name.Name)
	fmt.Fprintf(out2, "    Doc:  %#v\n", decl.Doc)
	fmt.Fprintf(out2, "    Recv: %#v\n", decl.Recv)
	fmt.Fprintf(out2, "    Type: %#v\n", decl.Type)
}

func reprSpec(spec ast.Spec, indent int) {
	switch spec.(type) {
	case *ast.ImportSpec:
		spec := spec.(*ast.ImportSpec)
		fmt.Fprint(&buf, pad[:indent], spec.Path.Value)
	case *ast.ValueSpec:
		spec := spec.(*ast.ValueSpec)
		for _, ident := range spec.Names {
			fmt.Fprint(&buf, pad[:indent], ident.Name)
		}
	case *ast.TypeSpec:
		spec := spec.(*ast.TypeSpec)
		fmt.Fprint(&buf, pad[:indent], spec.Name.Name)
	}
	fmt.Fprintln(&buf)
	buf.WriteTo(out1)
}

// Explore prints all the .go files (excluding tests...ish) that have the directory path root
// as a parent directory.
func Explore(root string) {
	defer out1.Close()
	defer out2.Close()
	dirRoot = root
	err := filepath.Walk(root, walker)
	check(err)
}

func init() {
	var err error
	out1, err = os.Create("./log/main.log")
	out2, err = os.Create("./log/funcs.log")
	check(err)

	fset = token.NewFileSet()
	b := make([]byte, 100)
	for i := 0; i < 100; i++ {
		b[i] = ' '
	}
	pad = string(b)
}

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
