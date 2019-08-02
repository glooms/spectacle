package spec

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var out *os.File
var debug_out *os.File

var specs map[string]*Spec

func Build(root string) map[string]*Spec {
  specs = map[string]*Spec{}
	err := filepath.Walk(root, walker)
	check(err)
  return specs
}

func init() {
  if _, err := os.Stat("./log"); os.IsNotExist(err) {
    os.Mkdir("./log", 0755)
  }
}

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
  fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, filter, parser.ParseComments)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		if pkg.Name != "main" {
      createLogs(pkg.Name)
      spec := New(fset, pkg)
      log(spec.String())
      specs[pkg.Name] = spec
      closeLogs()
		}
	}
	return nil
}

func createLogs(name string) {
  var err error
  out, err = os.Create("./log/" + name + ".log")
  check(err)
	debug_out, err = os.Create("./log/" + name + ".debug")
  check(err)
}

func closeLogs() {
  out.Close()
  debug_out.Close()
}

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}

func log(a ...interface{}) {
	fmt.Fprintln(out, a...)
}

func vlog(i interface{}, prefix ...interface{}) {
	fmt.Fprint(out, prefix...)
	fmt.Fprintf(out, "%#v\n", i)
}

func debug(a ...interface{}) {
	fmt.Fprintln(debug_out, a...)
}

func vdebug(i interface{}, prefix ...interface{}) {
	fmt.Fprint(debug_out, prefix...)
	fmt.Fprintf(debug_out, "%#v\n", i)
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
